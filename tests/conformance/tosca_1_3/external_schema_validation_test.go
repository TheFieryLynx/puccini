package tosca_1_3_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tliron/go-puccini/normal"
	"github.com/tliron/go-puccini/tests/conformance/internal/testsupport"
)

const validJSONSchema = `{"$schema":"http://json-schema.org/draft-04/schema#","type":"object","required":["name"],"properties":{"name":{"type":"string","pattern":"^[a-z]+$"},"count":{"type":"integer","minimum":1,"maximum":3}},"additionalProperties":false}`

const validXMLSchema = `<xs:schema xmlns:xs="http://www.w3.org/2001/XMLSchema"><xs:element name="payload"><xs:complexType><xs:sequence><xs:element name="name" type="xs:string"/><xs:element name="count" type="xs:positiveInteger" minOccurs="0"/></xs:sequence></xs:complexType></xs:element></xs:schema>`

const maxExternalSchemaReferenceTestCount = 129

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 3.6.3.2, 3.6.3.4, 3.6.10.5, 5.3.2, 5.3.4
// Requirement: TOSCA13-3.6.10.5-005
// Expected: inline JSON Schema Draft 4 and XML Schema 1.0 strings matching the
// resolved json/xml external property type compile successfully.
// Category: positive, declaration, language, scalar, mapping, sequence, nested
func TestExternalSchemaValidDeclarations(t *testing.T) {
	source := externalSchemaTemplate(`
      json_payload:
        type: json
        constraints:
          - schema: >-
              `+validJSONSchema+`
      xml_payload:
        type: xml
        constraints:
          - schema: >-
              `+validXMLSchema+`
      json_array:
        type: json
        constraints:
          - schema: >-
              {"$schema":"http://json-schema.org/draft-04/schema#","type":"array","items":{"type":"integer"},"minItems":1,"maxItems":3}
`, `
        json_payload: '{"name":"valid","count":2}'
        xml_payload: '<payload><name>valid</name><count>2</count></payload>'
        json_array: '[1,2,3]'
`)

	serviceTemplate, parseProblems, err := testsupport.ParseSource(t, source)
	if err != nil {
		t.Fatalf("valid external schemas failed: %v\n%s", err, parseProblems)
	}
	if serviceTemplate.NodeTemplates["node"] == nil {
		t.Fatal("valid external-schema template did not normalize")
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 3.6.3.1, 3.6.3.4, 3.6.10.4, 3.6.10.5
// Requirement: TOSCA13-3.6.10.5-005
// Expected: malformed, empty, wrong-language, wrong-type, and unsupported
// external-schema declarations fail in the read or render phase with a stable
// category, property path, and external type.
// Category: negative, declaration, compilation, diagnostic
func TestExternalSchemaInvalidDeclarations(t *testing.T) {
	tests := []struct {
		name       string
		properties string
		category   string
		fragments  []string
	}{
		{
			name: "malformed JSON",
			properties: `
      malformed_json:
        type: json
        constraints:
          - schema: '{"type":'
`,
			category:  "external schema compilation failure",
			fragments: []string{"malformed_json", "json"},
		},
		{
			name: "semantically invalid JSON schema",
			properties: `
      invalid_json:
        type: json
        constraints:
          - schema: '{"$schema":"http://json-schema.org/draft-04/schema#","type":7}'
`,
			category:  "external schema compilation failure",
			fragments: []string{"invalid_json", "json"},
		},
		{
			name: "empty schema",
			properties: `
      empty_schema:
        type: json
        constraints:
          - schema: "   "
`,
			category:  "invalid external schema declaration",
			fragments: []string{"empty_schema", "json"},
		},
		{
			name: "unsupported JSON dialect",
			properties: `
      future_json:
        type: json
        constraints:
          - schema: '{"$schema":"https://json-schema.org/draft/2020-12/schema","type":"object"}'
`,
			category:  "unsupported external schema language",
			fragments: []string{"future_json", "json"},
		},
		{
			name: "malformed XML schema",
			properties: `
      malformed_xml:
        type: xml
        constraints:
          - schema: '<xs:schema xmlns:xs="http://www.w3.org/2001/XMLSchema">'
`,
			category:  "external schema compilation failure",
			fragments: []string{"malformed_xml", "xml"},
		},
		{
			name: "semantically invalid XML schema",
			properties: `
      invalid_xml:
        type: xml
        constraints:
          - schema: '<xs:schema xmlns:xs="http://www.w3.org/2001/XMLSchema"><xs:element name="payload" type="xs:notAType"/></xs:schema>'
`,
			category:  "external schema compilation failure",
			fragments: []string{"invalid_xml", "xml"},
		},
		{
			name: "schema on non-external type",
			properties: `
      ordinary:
        type: string
        constraints:
          - schema: '{"type":"string"}'
`,
			category:  "unsupported external schema language",
			fragments: []string{"ordinary", "string"},
		},
		{
			name: "path is not an inline schema form",
			properties: `
      path_form:
        type: json
        constraints:
          - schema: ./schemas/payload.json
`,
			category:  "external schema compilation failure",
			fragments: []string{"path_form", "json"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assertExternalSchemaRejected(t, externalSchemaTemplate(test.properties, ""), test.category, test.fragments...)
		})
	}

	t.Run("wrong YAML operand type", func(t *testing.T) {
		problems := parseExternalSchemaFailure(t, externalSchemaTemplate(`
      wrong_operand:
        type: json
        constraints:
          - schema: {type: object}
`, ""))
		for _, fragment := range []string{"wrong_operand", "schema", "string"} {
			if !strings.Contains(problems, fragment) {
				t.Fatalf("wrong-operand diagnostic does not contain %q:\n%s", fragment, problems)
			}
		}
		if strings.Contains(problems, "external schema compilation failure") {
			t.Fatalf("wrong YAML type reached schema compilation instead of read validation:\n%s", problems)
		}
	})

	t.Run("unknown declaration keyname", func(t *testing.T) {
		problems := parseExternalSchemaFailure(t, externalSchemaTemplate(`
      unknown_key:
        type: json
        constraints:
          - schema_language: json-schema
`, ""))
		if !strings.Contains(problems, "schema_language") || !strings.Contains(problems, "unsupported operator") {
			t.Fatalf("unknown schema declaration operator was not rejected in the read phase:\n%s", problems)
		}
	})

	t.Run("direct external-schema key is not the normative grammar form", func(t *testing.T) {
		problems := parseExternalSchemaFailure(t, externalSchemaTemplate(`
      direct_key:
        type: json
        external-schema: '{"type":"object"}'
`, ""))
		if !strings.Contains(problems, "external-schema") || !strings.Contains(problems, "unsupported") {
			t.Fatalf("ambiguous direct external-schema key was accepted:\n%s", problems)
		}
	})

	t.Run("TOSCA 2.0 validation form is rejected by 1.3", func(t *testing.T) {
		problems := parseExternalSchemaFailure(t, externalSchemaTemplate(`
      v2_form:
        type: json
        validation:
          $schema: '{"type":"object"}'
`, ""))
		if !strings.Contains(problems, "validation") || !strings.Contains(problems, "unsupported") {
			t.Fatalf("TOSCA 2.0 validation form was accepted by the 1.3 grammar:\n%s", problems)
		}
	})
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 3.6.10.5, 3.6.10.6, 3.7.9
// Requirement: TOSCA13-3.6.10.5-005
// Expected: schema constraints retained through type inheritance are compiled,
// and an unrelated same-named property without a schema is unaffected.
// Category: positive, negative, inheritance, refinement, scope
func TestExternalSchemaInheritance(t *testing.T) {
	valid := `tosca_definitions_version: tosca_simple_yaml_1_3
node_types:
  Base:
    derived_from: tosca.nodes.Root
    properties:
      payload:
        type: json
        required: false
        constraints:
          - schema: >-
              ` + validJSONSchema + `
  Derived:
    derived_from: Base
    properties:
      payload:
        required: true
        default: '{"name":"default"}'
      unrelated:
        type: string
        required: false
topology_template:
  node_templates:
    node:
      type: Derived
`
	if _, parseProblems, err := testsupport.ParseSource(t, valid); err != nil {
		t.Fatalf("inherited external schema failed: %v\n%s", err, parseProblems)
	}

	invalid := strings.Replace(validJSONSchema, `"type":"object"`, `"type":7`, 1)
	assertExternalSchemaRejected(t, strings.Replace(valid, validJSONSchema, invalid, 1),
		"external schema compilation failure", "payload", "json")
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 3.6.3.4, 3.6.10.5, 3.10.2
// Requirement: TOSCA13-3.6.10.5-005
// Expected: property definitions in an imported type file are validated in
// that file's post-inheritance render operation; the inline schema is not
// resolved relative to the importer, repository, or CSAR root.
// Category: positive, negative, import, URI-context isolation
func TestExternalSchemaImportedDefinition(t *testing.T) {
	directory := t.TempDir()
	importPath := filepath.Join(directory, "types.yaml")
	mainPath := filepath.Join(directory, "main.yaml")
	write := func(path string, source string) {
		t.Helper()
		if err := os.WriteFile(path, []byte(source), 0o600); err != nil {
			t.Fatalf("write %s: %v", path, err)
		}
	}

	imported := `tosca_definitions_version: tosca_simple_yaml_1_3
node_types:
  Imported:
    derived_from: tosca.nodes.Root
    properties:
      payload:
        type: json
        required: false
        constraints:
          - schema: >-
              ` + validJSONSchema + `
`
	main := `tosca_definitions_version: tosca_simple_yaml_1_3
imports:
  - types.yaml
topology_template: {}
`
	write(importPath, imported)
	write(mainPath, main)
	if _, parseProblems, err := testsupport.ParseFile(t, mainPath); err != nil {
		t.Fatalf("valid schema in imported definition failed: %v\n%s", err, parseProblems)
	}

	write(importPath, strings.Replace(imported, `"type":"object"`, `"type":7`, 1))
	_, parseProblems, err := testsupport.ParseFile(t, mainPath)
	if err == nil {
		t.Fatal("invalid schema in imported definition was accepted")
	}
	for _, fragment := range []string{"external schema compilation failure", "payload", "json", "types.yaml"} {
		if !strings.Contains(parseProblems, fragment) {
			t.Fatalf("imported-schema diagnostic does not contain %q:\n%s", fragment, parseProblems)
		}
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 3.6.3.2, 3.6.10.5, 3.6.11.2
// Requirement: TOSCA13-3.6.10.5-005
// Expected: the mandatory processor check compiles the schema declaration; it
// does not implement the separate MAY for validating defaults, assignments, or
// unevaluated/evaluated intrinsic-function results against that schema.
// Category: positive, optional-boundary, default, assignment, function, null
func TestExternalSchemaValueValidationIsOptional(t *testing.T) {
	source := externalSchemaTemplate(`
      default_mismatch:
        type: json
        default: '{"count":0}'
        constraints:
          - schema: >-
              `+validJSONSchema+`
      assignment_mismatch:
        type: json
        constraints:
          - schema: >-
              `+validJSONSchema+`
      optional_absent:
        type: json
        required: false
        constraints:
          - schema: >-
              `+validJSONSchema+`
      function_value:
        type: json
        constraints:
          - schema: >-
              `+validJSONSchema+`
`, `
        assignment_mismatch: '{"extra":true}'
        function_value: { concat: ['{"count":', '0}'] }
`)
	serviceTemplate, parseProblems, err := testsupport.ParseSource(t, source)
	if err != nil {
		t.Fatalf("optional value-schema validation was applied as a MUST: %v\n%s", err, parseProblems)
	}
	node := serviceTemplate.NodeTemplates["node"]
	if node == nil {
		t.Fatal("value-boundary template did not normalize")
	}
	if _, ok := node.Properties["function_value"].(*normal.FunctionCall); !ok {
		t.Fatalf("intrinsic function was not preserved for later evaluation: %T", node.Properties["function_value"])
	}
	if value := node.Properties["optional_absent"]; value != nil {
		t.Fatalf("absent optional value became %#v", value)
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 3.6.3.2, 3.6.10.5
// Requirement: TOSCA13-3.6.10.5-005
// Expected: external references are denied without reading files or reaching
// loopback/link-local/private networks, including through XML entities.
// Category: negative, security, URI resolution, offline
func TestExternalSchemaExternalReferencesAreDenied(t *testing.T) {
	tests := []struct {
		name       string
		typeName   string
		schema     string
		identifier string
	}{
		{
			name:       "JSON file traversal",
			typeName:   "json",
			schema:     `{"$schema":"http://json-schema.org/draft-04/schema#","$ref":"file:///etc/passwd"}`,
			identifier: "file:///etc/passwd",
		},
		{
			name:       "JSON relative traversal",
			typeName:   "json",
			schema:     `{"$schema":"http://json-schema.org/draft-04/schema#","$ref":"../../secret.json"}`,
			identifier: "../../secret.json",
		},
		{
			name:       "JSON loopback",
			typeName:   "json",
			schema:     `{"$schema":"http://json-schema.org/draft-04/schema#","$ref":"http://127.0.0.1/schema"}`,
			identifier: "http://127.0.0.1/schema",
		},
		{
			name:       "JSON cloud metadata",
			typeName:   "json",
			schema:     `{"$schema":"http://json-schema.org/draft-04/schema#","$ref":"http://169.254.169.254/latest/meta-data/"}`,
			identifier: "http://169.254.169.254/latest/meta-data/",
		},
		{
			name:       "unsupported process scheme",
			typeName:   "json",
			schema:     `{"$schema":"http://json-schema.org/draft-04/schema#","$ref":"exec:///bin/sh"}`,
			identifier: "exec:///bin/sh",
		},
		{
			name:       "XSD relative traversal",
			typeName:   "xml",
			schema:     `<xs:schema xmlns:xs="http://www.w3.org/2001/XMLSchema"><xs:include schemaLocation="../../secret.xsd"/></xs:schema>`,
			identifier: "../../secret.xsd",
		},
		{
			name:       "XSD loopback import",
			typeName:   "xml",
			schema:     `<xs:schema xmlns:xs="http://www.w3.org/2001/XMLSchema"><xs:import namespace="urn:test" schemaLocation="http://127.0.0.1/schema.xsd"/></xs:schema>`,
			identifier: "http://127.0.0.1/schema.xsd",
		},
		{
			name:       "XML external entity",
			typeName:   "xml",
			schema:     `<!DOCTYPE xs:schema [<!ENTITY secret SYSTEM "file:///etc/passwd">]><xs:schema xmlns:xs="http://www.w3.org/2001/XMLSchema">&secret;</xs:schema>`,
			identifier: "external entity",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			source := externalSchemaTemplate(fmt.Sprintf(`
      protected:
        type: %s
        constraints:
          - schema: >-
              %s
`, test.typeName, test.schema), "")
			problems := parseExternalSchemaFailure(t, source)
			for _, fragment := range []string{"external schema access denied", "protected", test.typeName} {
				if !strings.Contains(problems, fragment) {
					t.Fatalf("access-denied diagnostic does not contain %q:\n%s", fragment, problems)
				}
			}
			if strings.Contains(problems, "root:x:") {
				t.Fatalf("diagnostic disclosed local-file content:\n%s", problems)
			}
		})
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 3.6.3.4, 3.6.10.5
// Requirement: TOSCA13-3.6.10.5-005
// Expected: bounded schema input and reference-cycle handling terminate with
// deterministic diagnostics; inline references do not trigger retrieval.
// Category: positive, negative, security, limit, cycle
func TestExternalSchemaSecurityLimits(t *testing.T) {
	oversized := strings.Repeat(" ", (1<<20)+1)
	assertExternalSchemaRejected(t, externalSchemaTemplate(`
      oversized:
        type: json
        constraints:
          - schema: "`+oversized+`"
`, ""), "external schema too large", "oversized", "json")

	selfReference := `{"$schema":"http://json-schema.org/draft-04/schema#","definitions":{"node":{"type":"object","properties":{"next":{"$ref":"#/definitions/node"}}}},"$ref":"#/definitions/node"}`
	if _, parseProblems, err := testsupport.ParseSource(t, externalSchemaTemplate(`
      recursive:
        type: json
        constraints:
          - schema: >-
              `+selfReference+`
`, "")); err != nil {
		t.Fatalf("bounded in-document JSON reference failed: %v\n%s", err, parseProblems)
	}

	deep := strings.Repeat("<xs:annotation>", 70) + strings.Repeat("</xs:annotation>", 70)
	assertExternalSchemaRejected(t, externalSchemaTemplate(`
      deeply_nested:
        type: xml
        constraints:
          - schema: >-
              <xs:schema xmlns:xs="http://www.w3.org/2001/XMLSchema">`+deep+`</xs:schema>
`, ""), "external schema compilation failure", "deeply_nested", "xml")

	deepJSON := `{"$schema":"http://json-schema.org/draft-04/schema#"`
	for range 70 {
		deepJSON += `,"allOf":[{"type":"object"`
	}
	deepJSON += strings.Repeat("}]", 70) + "}"
	assertExternalSchemaRejected(t, externalSchemaTemplate(`
      deeply_nested_json:
        type: json
        constraints:
          - schema: >-
              `+deepJSON+`
`, ""), "external schema compilation failure", "deeply_nested_json", "json")

	var manyReferences strings.Builder
	manyReferences.WriteString(`{"$schema":"http://json-schema.org/draft-04/schema#","allOf":[`)
	for index := range maxExternalSchemaReferenceTestCount {
		if index > 0 {
			manyReferences.WriteByte(',')
		}
		fmt.Fprintf(&manyReferences, `{"$ref":"#/definitions/value%d"}`, index)
	}
	manyReferences.WriteString(`],"definitions":{`)
	for index := range maxExternalSchemaReferenceTestCount {
		if index > 0 {
			manyReferences.WriteByte(',')
		}
		fmt.Fprintf(&manyReferences, `"value%d":{"type":"string"}`, index)
	}
	manyReferences.WriteString("}}")
	assertExternalSchemaRejected(t, externalSchemaTemplate(`
      excessive_references:
        type: json
        constraints:
          - schema: >-
              `+manyReferences.String()+`
`, ""), "external schema compilation failure", "excessive_references", "json")
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 3.6.10.5, 14.3
// Requirement: TOSCA13-3.6.10.5-005
// Expected: repeated processing emits the same stable category order, and
// duplicate schema declarations use one operation-local compilation result.
// Category: positive, negative, determinism, cache
func TestExternalSchemaDeterminismAndCache(t *testing.T) {
	source := externalSchemaTemplate(`
      first:
        type: json
        constraints:
          - schema: '{"type":7}'
      second:
        type: json
        constraints:
          - schema: '{"type":7}'
`, "")
	first := parseExternalSchemaFailure(t, source)
	second := parseExternalSchemaFailure(t, source)
	for _, problems := range []string{first, second} {
		if count := strings.Count(problems, "external schema compilation failure"); count != 2 {
			t.Fatalf("diagnostic count = %d, want 2:\n%s", count, problems)
		}
		if strings.Index(problems, "first") > strings.Index(problems, "second") {
			t.Fatalf("diagnostic order is unstable:\n%s", problems)
		}
	}

	validSource := externalSchemaTemplate(`
      first:
        type: json
        constraints:
          - schema: >-
              `+validJSONSchema+`
      second:
        type: json
        constraints:
          - schema: >-
              `+validJSONSchema+`
`, `
        first: '{"name":"one"}'
        second: '{"name":"two"}'
`)
	var normalized [][]byte
	validPath := filepath.Join(t.TempDir(), "deterministic.yaml")
	if err := os.WriteFile(validPath, []byte(validSource), 0o600); err != nil {
		t.Fatalf("write deterministic fixture: %v", err)
	}
	for range 2 {
		serviceTemplate, parseProblems, err := testsupport.ParseFile(t, validPath)
		if err != nil {
			t.Fatalf("deterministic valid schema failed: %v\n%s", err, parseProblems)
		}
		encoded, err := json.Marshal(serviceTemplate)
		if err != nil {
			t.Fatalf("marshal normalized template: %v", err)
		}
		normalized = append(normalized, encoded)
	}
	if string(normalized[0]) != string(normalized[1]) {
		t.Fatal("normalized output changed across identical processing operations")
	}
}

// Specification: TOSCA Version 2.0 (regression only; not a normative source
// for TOSCA13-3.6.10.5-005)
// Expected: activating TOSCA 1.3 external-schema compilation does not change
// TOSCA 2.0 grammar or ordinary property rendering.
// Category: cross-version regression, version isolation
func TestExternalSchemaTosca20Isolation(t *testing.T) {
	source := `tosca_definitions_version: tosca_2_0
node_types:
  Example:
    properties:
      ordinary:
        type: string
        required: false
service_template:
  node_templates:
    node:
      type: Example
      properties:
        ordinary: unchanged
`
	serviceTemplate, parseProblems, err := testsupport.ParseSource(t, source)
	if err != nil {
		t.Fatalf("TOSCA 2.0 regression template failed: %v\n%s", err, parseProblems)
	}
	assertNormalizedPrimitive(t, serviceTemplate.NodeTemplates["node"].Properties["ordinary"], "unchanged")
}

func externalSchemaTemplate(properties string, assignments string) string {
	source := `tosca_definitions_version: tosca_simple_yaml_1_3
node_types:
  ExternalSchemaNode:
    derived_from: tosca.nodes.Root
    properties:
` + properties
	if strings.TrimSpace(assignments) == "" {
		return source + "\ntopology_template: {}\n"
	}
	source += `
topology_template:
  node_templates:
    node:
      type: ExternalSchemaNode
`
	source += "      properties:\n" + assignments
	return source
}

func parseExternalSchemaFailure(t *testing.T, source string) string {
	t.Helper()
	_, problems, err := testsupport.ParseSource(t, source)
	if err == nil {
		t.Fatal("invalid external-schema template was accepted")
	}
	return problems
}

func assertExternalSchemaRejected(t *testing.T, source string, category string, fragments ...string) {
	t.Helper()
	problems := parseExternalSchemaFailure(t, source)
	for _, fragment := range append([]string{category}, fragments...) {
		if !strings.Contains(problems, fragment) {
			t.Fatalf("diagnostic does not contain %q:\n%s", fragment, problems)
		}
	}
}
