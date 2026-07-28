package tosca_v1_3

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/lestrrat-go/helium"
	"github.com/lestrrat-go/helium/xsd"
	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/tliron/go-ard"
	"github.com/tliron/go-kutil/reflection"
	"github.com/tliron/go-puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/go-puccini/tosca/parsing"
	"github.com/tliron/yamlkeys"
)

const (
	maxExternalSchemaBytes      = 1 << 20
	maxExternalSchemaDepth      = 64
	maxExternalSchemaNodes      = 100_000
	maxExternalSchemaReferences = 128
	externalSchemaTimeout       = time.Second
)

const (
	externalSchemaInvalidDeclaration  = "invalid external schema declaration"
	externalSchemaUnsupportedLanguage = "unsupported external schema language"
	externalSchemaTooLarge            = "external schema too large"
	externalSchemaAccessDenied        = "external schema access denied"
	externalSchemaReferenceCycle      = "external schema reference cycle"
	externalSchemaCompilationFailure  = "external schema compilation failure"
)

const (
	jsonSchemaDraft4HTTP  = "http://json-schema.org/draft-04/schema#"
	jsonSchemaDraft4HTTPS = "https://json-schema.org/draft-04/schema#"
)

type externalSchemaDeclaration struct {
	property   *tosca_v2_0.PropertyDefinition
	schema     string
	identifier string
	typeName   string
	schemaType string
	path       string
}

type externalSchemaCacheKey struct {
	schemaType string
	digest     [sha256.Size]byte
}

type externalSchemaCompileFunc func(string, string) string

type externalSchemaValidator struct {
	cache   map[externalSchemaCacheKey]string
	compile externalSchemaCompileFunc
}

func newExternalSchemaValidator() *externalSchemaValidator {
	return &externalSchemaValidator{
		cache:   make(map[externalSchemaCacheKey]string),
		compile: compileExternalSchema,
	}
}

func validateExternalPropertySchemas(file *tosca_v2_0.File) {
	roots := []any{file}
	file.Context.Namespace.Range(func(entity parsing.EntityPtr) bool {
		roots = append(roots, entity)
		return true
	})
	declarations := collectExternalSchemaDeclarations(roots...)
	validator := newExternalSchemaValidator()
	for _, declaration := range declarations {
		if category := validator.compileCached(declaration.schemaType, declaration.schema); category != "" {
			declaration.property.Context.ReportError(fmt.Errorf(
				"%s for property %q (type %q, schema %q)",
				category,
				declaration.property.Name,
				declaration.typeName,
				declaration.identifier,
			))
		}
	}
}

func (self *externalSchemaValidator) compileCached(schemaType string, schema string) string {
	key := externalSchemaCacheKey{
		schemaType: schemaType,
		digest:     sha256.Sum256([]byte(schema)),
	}
	if category, ok := self.cache[key]; ok {
		return category
	}
	category := self.compile(schemaType, schema)
	self.cache[key] = category
	return category
}

func collectExternalSchemaDeclarations(roots ...any) []externalSchemaDeclaration {
	var declarations []externalSchemaDeclaration
	seen := make(map[*tosca_v2_0.PropertyDefinition]struct{})

	for _, root := range roots {
		reflection.TraverseEntities(root, false, func(entity any) bool {
			property, ok := entity.(*tosca_v2_0.PropertyDefinition)
			if !ok {
				return true
			}
			if _, ok := seen[property]; ok {
				return false
			}
			seen[property] = struct{}{}
			// Name lookup already reports an unresolved property type. Avoid a
			// second, misleading schema-language diagnostic in that case.
			if property.DataType == nil {
				return true
			}

			schemas := schemaConstraintStrings(property.ValidationClause)
			if len(schemas) == 0 {
				return true
			}

			typeName, schemaType := externalPropertySchemaType(property)
			for _, schema := range schemas {
				declarations = append(declarations, externalSchemaDeclaration{
					property:   property,
					schema:     schema,
					identifier: externalSchemaIdentifier(schema),
					typeName:   typeName,
					schemaType: schemaType,
					path:       property.Context.Path.String(),
				})
			}
			return true
		})
	}

	sort.SliceStable(declarations, func(left int, right int) bool {
		if declarations[left].path == declarations[right].path {
			if declarations[left].typeName == declarations[right].typeName {
				return declarations[left].schema < declarations[right].schema
			}
			return declarations[left].typeName < declarations[right].typeName
		}
		return declarations[left].path < declarations[right].path
	})
	return declarations
}

func schemaConstraintStrings(clause *tosca_v2_0.ValidationClause) []string {
	if clause == nil {
		return nil
	}
	var schemas []string
	if clause.Operator == "schema" {
		for _, argument := range clause.Arguments {
			if schema, ok := argument.(string); ok {
				schemas = append(schemas, schema)
			}
		}
		return schemas
	}
	for _, argument := range clause.Arguments {
		collectSchemaConstraintStrings(argument, &schemas)
	}
	return schemas
}

func collectSchemaConstraintStrings(value any, schemas *[]string) {
	switch value := value.(type) {
	case ard.List:
		for _, item := range value {
			collectSchemaConstraintStrings(item, schemas)
		}
	case ard.Map:
		for key, operand := range value {
			operator := strings.TrimPrefix(yamlkeys.KeyString(key), "$")
			if operator == "schema" {
				if schema, ok := operand.(string); ok {
					*schemas = append(*schemas, schema)
				}
				continue
			}
			if operator == "and" || operator == "or" || operator == "not" || operator == "xor" {
				collectSchemaConstraintStrings(operand, schemas)
			}
		}
	case *tosca_v2_0.ValidationClause:
		*schemas = append(*schemas, schemaConstraintStrings(value)...)
	}
}

func externalPropertySchemaType(property *tosca_v2_0.PropertyDefinition) (string, string) {
	typeName := ""
	if property.DataTypeName != nil {
		typeName = *property.DataTypeName
	}
	for dataType := property.DataType; dataType != nil; dataType = dataType.Parent {
		switch dataType.Name {
		case "json", "tosca.datatypes.json":
			return firstNonEmpty(typeName, dataType.Name), "json"
		case "xml", "tosca.datatypes.xml":
			return firstNonEmpty(typeName, dataType.Name), "xml"
		}
	}
	return firstNonEmpty(typeName, "<unresolved>"), ""
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func externalSchemaIdentifier(schema string) string {
	digest := sha256.Sum256([]byte(schema))
	return "sha256:" + hex.EncodeToString(digest[:8])
}

func compileExternalSchema(schemaType string, schema string) string {
	if len(schema) > maxExternalSchemaBytes {
		return externalSchemaTooLarge
	}
	if strings.TrimSpace(schema) == "" {
		return externalSchemaInvalidDeclaration
	}

	switch schemaType {
	case "json":
		return compileJSONSchema(schema)
	case "xml":
		return compileXMLSchema(schema)
	default:
		return externalSchemaUnsupportedLanguage
	}
}

func compileJSONSchema(schema string) string {
	document, err := jsonschema.UnmarshalJSON(strings.NewReader(schema))
	if err != nil {
		return externalSchemaCompilationFailure
	}
	if category := inspectJSONSchema(document); category != "" {
		return category
	}

	digest := sha256.Sum256([]byte(schema))
	resource := "urn:puccini:tosca:1.3:external-schema:" + hex.EncodeToString(digest[:])
	compiler := jsonschema.NewCompiler()
	compiler.DefaultDraft(jsonschema.Draft4)
	compiler.UseLoader(denyExternalSchemaURLLoader{})
	if err := compiler.AddResource(resource, document); err != nil {
		return externalSchemaCompilationFailure
	}
	if _, err := compiler.Compile(resource); err != nil {
		if errors.Is(err, errExternalSchemaAccessDenied) {
			return externalSchemaAccessDenied
		}
		return externalSchemaCompilationFailure
	}
	return ""
}

func inspectJSONSchema(document any) string {
	type valueAtDepth struct {
		value any
		depth int
	}

	stack := []valueAtDepth{{value: document, depth: 1}}
	nodes := 0
	references := 0

	for len(stack) > 0 {
		last := len(stack) - 1
		current := stack[last]
		stack = stack[:last]
		nodes++
		if nodes > maxExternalSchemaNodes || current.depth > maxExternalSchemaDepth {
			return externalSchemaCompilationFailure
		}

		switch value := current.value.(type) {
		case map[string]any:
			if dialect, ok := value["$schema"]; ok {
				dialect, ok := dialect.(string)
				if !ok {
					return externalSchemaCompilationFailure
				}
				if !isJSONSchemaDraft4(dialect) {
					return externalSchemaUnsupportedLanguage
				}
			}
			if reference, ok := value["$ref"]; ok {
				reference, ok := reference.(string)
				if !ok {
					return externalSchemaCompilationFailure
				}
				references++
				if references > maxExternalSchemaReferences {
					return externalSchemaCompilationFailure
				}
				if reference != "" && !strings.HasPrefix(reference, "#") {
					return externalSchemaAccessDenied
				}
			}
			for _, child := range value {
				stack = append(stack, valueAtDepth{value: child, depth: current.depth + 1})
			}
		case []any:
			for _, child := range value {
				stack = append(stack, valueAtDepth{value: child, depth: current.depth + 1})
			}
		}
	}
	return ""
}

func isJSONSchemaDraft4(dialect string) bool {
	dialect = strings.TrimSuffix(dialect, "#")
	return dialect == strings.TrimSuffix(jsonSchemaDraft4HTTP, "#") ||
		dialect == strings.TrimSuffix(jsonSchemaDraft4HTTPS, "#")
}

func compileXMLSchema(schema string) string {
	if category := inspectXMLSchema(schema); category != "" {
		return category
	}

	compileContext, cancel := context.WithTimeout(context.Background(), externalSchemaTimeout)
	defer cancel()

	parser := helium.NewParser().
		AllowNetwork(false).
		BlockXXE(true).
		MaxDepth(maxExternalSchemaDepth).
		MaxNodeContentSize(maxExternalSchemaBytes)
	document, err := parser.Parse(compileContext, []byte(schema))
	if err != nil {
		if compileContext.Err() != nil {
			return externalSchemaCompilationFailure
		}
		return externalSchemaCompilationFailure
	}
	compiler := xsd.NewCompiler().
		FS(nil).
		Parser(parser).
		Version(xsd.Version10).
		Label("inline TOSCA external schema")
	if _, err := compiler.Compile(compileContext, document); err != nil {
		if compileContext.Err() != nil {
			return externalSchemaCompilationFailure
		}
		return externalSchemaCompilationFailure
	}
	return ""
}

func inspectXMLSchema(schema string) string {
	decoder := xml.NewDecoder(strings.NewReader(schema))
	depth := 0
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			return ""
		}
		if err != nil {
			return externalSchemaCompilationFailure
		}
		switch token := token.(type) {
		case xml.Directive:
			if strings.HasPrefix(strings.ToUpper(strings.TrimSpace(string(token))), "DOCTYPE") {
				return externalSchemaAccessDenied
			}
		case xml.StartElement:
			depth++
			if depth > maxExternalSchemaDepth {
				return externalSchemaCompilationFailure
			}
			if token.Name.Space == "http://www.w3.org/2001/XMLSchema" &&
				(token.Name.Local == "include" || token.Name.Local == "import" ||
					token.Name.Local == "redefine" || token.Name.Local == "override") {
				for _, attribute := range token.Attr {
					if attribute.Name.Local == "schemaLocation" && strings.TrimSpace(attribute.Value) != "" {
						return externalSchemaAccessDenied
					}
				}
			}
		case xml.EndElement:
			depth--
		}
	}
}

var errExternalSchemaAccessDenied = errors.New("external schema access denied")

type denyExternalSchemaURLLoader struct{}

func (denyExternalSchemaURLLoader) Load(string) (any, error) {
	return nil, errExternalSchemaAccessDenied
}
