package tosca_v1_3

import (
	"fmt"
	"net/url"
	"reflect"
	"sort"
	"strings"

	"github.com/tliron/go-ard"
	"github.com/tliron/go-puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/go-puccini/tosca/parsing"
	"github.com/tliron/yamlkeys"
)

const (
	reservedToscaNamespacePrefix = "tosca"
	reservedToscaNamespaceHost   = "docs.oasis-open.org"
	reservedToscaNamespacePath   = "/tosca"
)

func validateNamespaceDeclarations(context *parsing.Context, namespace *string, imports tosca_v2_0.Imports) {
	validateDefinitionsVersionFirstLine(context)
	validateDefaultNamespace(context, namespace)
	validateImportPrefixes(imports)
}

func validateDefinitionsVersionFirstLine(context *parsing.Context) {
	versionContext, ok := context.GetFieldChild("tosca_definitions_version")
	if !ok {
		return
	}
	versionRow, _ := versionContext.GetLocation()
	for _, fieldContext := range context.FieldChildren() {
		row, _ := fieldContext.GetLocation()
		if row > 0 && row < versionRow {
			versionContext.ReportValueMalformed(
				"tosca_definitions_version",
				"must be the first line/key of YAML",
			)
			return
		}
	}
}

func validateDefaultNamespace(context *parsing.Context, namespace *string) {
	if namespace == nil {
		return
	}

	parsed, err := url.Parse(*namespace)
	if err != nil {
		return
	}
	path := strings.TrimSuffix(parsed.EscapedPath(), "/")
	if (parsed.Scheme == "http") &&
		strings.EqualFold(parsed.Hostname(), reservedToscaNamespaceHost) &&
		((path == reservedToscaNamespacePath) || strings.HasPrefix(path, reservedToscaNamespacePath+"/")) {
		context.FieldChild("namespace", *namespace).ReportValueInvalid(
			"namespace",
			"reserved for TOSCA-approved specifications and work",
		)
	}
}

func validateImportPrefixes(imports tosca_v2_0.Imports) {
	seen := make(map[string]*tosca_v2_0.Import)
	for _, import_ := range imports {
		if import_.Namespace == nil || *import_.Namespace == "" {
			continue
		}

		prefix := *import_.Namespace
		prefixContext := import_.Context.FieldChild("namespace_prefix", prefix)
		if prefix == reservedToscaNamespacePrefix {
			prefixContext.ReportValueInvalid(
				"namespace prefix",
				fmt.Sprintf("%q is reserved and cannot be used on import", prefix),
			)
			continue
		}

		if previous, ok := seen[prefix]; ok {
			prefixContext.ReportValueInvalid(
				"namespace prefix",
				fmt.Sprintf("duplicate namespace_prefix %q; first declared at %s", prefix, previous.Context.Path.String()),
			)
			continue
		}
		seen[prefix] = import_
	}
}

func validateImportedDefinitionIdentities(context *parsing.Context) {
	type importedDefinition struct {
		namespace string
		localName string
		version   string
		typeName  string
		url       string
		context   *parsing.Context
	}

	var definitions []importedDefinition
	seenEntities := make(map[parsing.EntityPtr]struct{})
	context.Namespace.Range(func(entityPtr parsing.EntityPtr) bool {
		if _, ok := seenEntities[entityPtr]; ok {
			return true
		}
		seenEntities[entityPtr] = struct{}{}

		entityContext := parsing.GetContext(entityPtr)
		if entityContext.URL.String() == context.URL.String() {
			return true
		}
		namespace := entityContext.GetCanonicalNamespace()
		if namespace == nil || *namespace == "" {
			return true
		}

		definitions = append(definitions, importedDefinition{
			namespace: *namespace,
			localName: entityContext.Name,
			version:   definitionVersion(entityContext.Data),
			typeName:  reflect.TypeOf(entityPtr).String(),
			url:       entityContext.URL.String(),
			context:   entityContext,
		})
		return true
	})

	sort.Slice(definitions, func(i int, j int) bool {
		left, right := definitions[i], definitions[j]
		if left.namespace != right.namespace {
			return left.namespace < right.namespace
		}
		if left.localName != right.localName {
			return left.localName < right.localName
		}
		if left.version != right.version {
			return left.version < right.version
		}
		if left.typeName != right.typeName {
			return left.typeName < right.typeName
		}
		return left.url < right.url
	})

	for index := 1; index < len(definitions); index++ {
		previous, current := definitions[index-1], definitions[index]
		if previous.namespace == current.namespace &&
			previous.localName == current.localName &&
			previous.version == current.version &&
			previous.typeName == current.typeName &&
			previous.url != current.url {
			current.context.ReportError(fmt.Errorf(
				"imported definitions with namespace URI %q, local name %q, and version %q must be equivalent; conflicting definitions in %q and %q",
				current.namespace,
				current.localName,
				current.version,
				previous.url,
				current.url,
			))
		}
	}
}

func definitionVersion(data ard.Value) string {
	map_, ok := data.(ard.Map)
	if !ok {
		return ""
	}
	for key, value := range map_ {
		if yamlkeys.KeyString(key) == "version" {
			return yamlkeys.KeyString(value)
		}
	}
	return ""
}
