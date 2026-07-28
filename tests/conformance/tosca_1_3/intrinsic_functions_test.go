package tosca_1_3_test

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/tliron/exturl"
	"github.com/tliron/go-ard"
	"github.com/tliron/go-kutil/problems"
	"github.com/tliron/go-kutil/terminal"
	cloutjs "github.com/tliron/go-puccini/clout/js"
	"github.com/tliron/go-puccini/normal"
	"github.com/tliron/go-puccini/tests/conformance/internal/testsupport"
	"github.com/tliron/go-puccini/tosca/parsing"
)

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 4.3.1.2, 4.3.3.2, 4.4.1.2, 4.4.2.2, 4.5.1.2, 4.6.1.2, 4.8.1.2
// Requirements: TOSCA13-4.3.1.2-003, TOSCA13-4.3.3.2-003,
// TOSCA13-4.3.3.2-006, TOSCA13-4.3.3.2-010,
// TOSCA13-4.4.1.2-003, TOSCA13-4.4.2.2-004,
// TOSCA13-4.4.2.2-010, TOSCA13-4.5.1.2-004,
// TOSCA13-4.5.1.2-010, TOSCA13-4.6.1.2-003,
// TOSCA13-4.6.1.2-006, TOSCA13-4.6.1.2-009,
// TOSCA13-4.6.1.2-012, TOSCA13-4.8.1.2-004,
// TOSCA13-4.8.1.2-007
// Expected: rejected in the read phase as a malformed intrinsic function call
// Category: negative, required argument, argument count
func TestIntrinsicFunctionRequiredArguments(t *testing.T) {
	tests := []struct {
		name     string
		call     string
		function string
		argument string
	}{
		{"concat-empty", "{ concat: [] }", "concat", "string value expression"},
		{"token-missing-string", "{ token: [] }", "token", "string_with_tokens"},
		{"token-missing-separators", "{ token: [abc] }", "token", "string_of_token_chars"},
		{"token-missing-index", `{ token: [abc, ":"] }`, "token", "substring_index"},
		{"get_input-empty", "{ get_input: [] }", "get_input", "input_property_name"},
		{"get_property-missing-entity", "{ get_property: [] }", "get_property", "modelable_entity_name"},
		{"get_property-missing-property", "{ get_property: [SELF] }", "get_property", "property_name"},
		{"get_attribute-missing-entity", "{ get_attribute: [] }", "get_attribute", "modelable_entity_name"},
		{"get_attribute-missing-attribute", "{ get_attribute: [SELF] }", "get_attribute", "attribute_name"},
		{"get_operation_output-missing-entity", "{ get_operation_output: [] }", "get_operation_output", "modelable_entity_name"},
		{"get_operation_output-missing-interface", "{ get_operation_output: [SELF] }", "get_operation_output", "interface_name"},
		{"get_operation_output-missing-operation", "{ get_operation_output: [SELF, Standard] }", "get_operation_output", "operation_name"},
		{"get_operation_output-missing-output", "{ get_operation_output: [SELF, Standard, create] }", "get_operation_output", "output_variable_name"},
		{"get_artifact-missing-entity", "{ get_artifact: [] }", "get_artifact", "modelable_entity_name"},
		{"get_artifact-missing-artifact", "{ get_artifact: [SELF] }", "get_artifact", "artifact_name"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assertFunctionRejected(t, test.call, test.function, test.argument)
		})
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 4.3.1.1-4.3.3.2, 4.4.1.1-4.6.1.2, 4.8.1.1-4.8.1.2
// Requirements: all IDs listed in intrinsic-functions.yaml
// Expected: invalid YAML argument types, fixed-arity overflow, and unknown
// additional entries are rejected in the read phase for the intended function
// Category: negative, YAML type, argument grammar, argument count
func TestIntrinsicFunctionArgumentGrammar(t *testing.T) {
	tests := []struct {
		name     string
		call     string
		function string
		reason   string
	}{
		{"concat-non-string", "{ concat: [prefix, 1] }", "concat", "argument 2"},
		{"token-string-type", `{ token: [1, ":", 0] }`, "token", "string_with_tokens"},
		{"token-separator-type", "{ token: [a:b, 1, 0] }", "token", "string_of_token_chars"},
		{"token-index-type", `{ token: [a:b, ":", zero] }`, "token", "substring_index"},
		{"token-excess", `{ token: [a:b, ":", 0, extra] }`, "token", "exactly 3"},
		{"get_input-scalar-type", "{ get_input: 1 }", "get_input", "input_property_name"},
		{"get_input-path-type", "{ get_input: [source, true] }", "get_input", "nested input"},
		{"get_property-entity-type", "{ get_property: [1, name] }", "get_property", "modelable_entity_name"},
		{"get_property-path-type", "{ get_property: [SELF, name, true] }", "get_property", "property path"},
		{"get_attribute-entity-type", "{ get_attribute: [1, state] }", "get_attribute", "modelable_entity_name"},
		{"get_attribute-path-type", "{ get_attribute: [SELF, state, true] }", "get_attribute", "attribute path"},
		{"get_operation_output-type", "{ get_operation_output: [SELF, Standard, create, 1] }", "get_operation_output", "output_variable_name"},
		{"get_operation_output-excess", "{ get_operation_output: [SELF, Standard, create, result, extra] }", "get_operation_output", "exactly 4"},
		{"get_artifact-entity-type", "{ get_artifact: [1, config] }", "get_artifact", "modelable_entity_name"},
		{"get_artifact-location-type", "{ get_artifact: [SELF, config, 1] }", "get_artifact", "location"},
		{"get_artifact-remove-type", "{ get_artifact: [SELF, config, LOCAL_FILE, no] }", "get_artifact", "remove"},
		{"get_artifact-excess", "{ get_artifact: [SELF, config, LOCAL_FILE, true, extra] }", "get_artifact", "at most 4"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assertFunctionRejected(t, test.call, test.function, test.reason)
		})
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 4.3.1-4.6.1, 4.8.1
// Requirements: all IDs listed in intrinsic-functions.yaml
// Expected: valid short/list forms are accepted and preserve the exact function
// name and arguments in deterministic normalized output
// Category: positive, boundary, nesting, normalization
func TestIntrinsicFunctionValidFormsAndNormalization(t *testing.T) {
	tests := []struct {
		name      string
		call      string
		function  string
		arguments []any
	}{
		{"concat-one", "{ concat: [prefix] }", "concat", []any{"prefix"}},
		{"concat-many", "{ concat: [prefix, suffix] }", "concat", []any{"prefix", "suffix"}},
		{"token", `{ token: [a:b, ":", 1] }`, "token", []any{"a:b", ":", 1}},
		{"get_input-scalar", "{ get_input: source }", "get_input", []any{"source"}},
		{"get_input-list", "{ get_input: [source, nested, 0] }", "get_input", []any{"source", "nested", 0}},
		{"get_property", "{ get_property: [node, inherited_name] }", "get_property", []any{"node", "inherited_name"}},
		{"get_attribute", "{ get_attribute: [node, inherited_state] }", "get_attribute", []any{"node", "inherited_state"}},
		{"get_operation_output", "{ get_operation_output: [node, Test, create, result] }", "get_operation_output", []any{"node", "Test", "create", "result"}},
		{"get_artifact-minimum", "{ get_artifact: [node, config] }", "get_artifact", []any{"node", "config"}},
		{"get_artifact-complete", "{ get_artifact: [node, config, LOCAL_FILE, true] }", "get_artifact", []any{"node", "config", "LOCAL_FILE", true}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			source := functionTemplate(test.call)
			if test.function == "get_property" || test.function == "get_attribute" ||
				test.function == "get_operation_output" || test.function == "get_artifact" {
				source = resolutionTemplate(test.call)
			}
			serviceTemplate, problems, err := testsupport.ParseSource(t, source)
			if err != nil {
				t.Fatalf("valid %s call failed: %v\n%s", test.function, err, problems)
			}
			call := normalizedOutputFunction(t, serviceTemplate)
			if call.Name != "tosca.function."+test.function {
				t.Fatalf("normalized function name = %q, want %q", call.Name, "tosca.function."+test.function)
			}
			arguments := normalizedArguments(call.Arguments)
			if !reflect.DeepEqual(arguments, test.arguments) {
				t.Fatalf("normalized arguments = %#v, want %#v", arguments, test.arguments)
			}
		})
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 4.3.1.2, 4.4.1.2
// Requirements: TOSCA13-4.3.1.2-003, TOSCA13-4.4.1.2-003
// Expected: nested intrinsic functions remain function calls, without semantic loss
// Category: positive, nesting, normalization
func TestIntrinsicFunctionNesting(t *testing.T) {
	serviceTemplate, problems, err := testsupport.ParseSource(t, functionTemplate(
		"{ concat: [prefix-, { get_input: source }] }",
	))
	if err != nil {
		t.Fatalf("nested function call failed: %v\n%s", err, problems)
	}

	outer := normalizedOutputFunction(t, serviceTemplate)
	if outer.Name != "tosca.function.concat" || len(outer.Arguments) != 2 {
		t.Fatalf("unexpected outer function: %#v", outer)
	}
	innerValue, ok := outer.Arguments[1].(*normal.FunctionCall)
	if !ok {
		t.Fatalf("nested argument was not normalized as a function call: %T", outer.Arguments[1])
	}
	if innerValue.FunctionCall.Name != "tosca.function.get_input" {
		t.Fatalf("nested function name = %q", innerValue.FunctionCall.Name)
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 4.4.1.2, 4.4.2.2, 4.5.1.2, 4.6.1.2, 4.8.1.2
// Requirements: TOSCA13-4.4.1.2-003, TOSCA13-4.4.2.2-002,
// TOSCA13-4.5.1.2-002, TOSCA13-4.6.1.2-001,
// TOSCA13-4.6.1.2-004, TOSCA13-4.6.1.2-007,
// TOSCA13-4.6.1.2-010, TOSCA13-4.8.1.2-002
// Expected: existing local and inherited definitions resolve; absent references
// are rejected in the rendering phase and identify the failing reference
// Category: positive, negative, inheritance, reference resolution
func TestIntrinsicFunctionResolution(t *testing.T) {
	valid := []struct {
		name string
		call string
	}{
		{"get_input", "{ get_input: source }"},
		{"get_property-inherited", "{ get_property: [node, inherited_name] }"},
		{"get_attribute-inherited", "{ get_attribute: [node, inherited_state] }"},
		{"get_operation_output-inherited", "{ get_operation_output: [node, Test, create, result] }"},
		{"get_artifact-inherited", "{ get_artifact: [node, config] }"},
	}
	for _, test := range valid {
		t.Run(test.name, func(t *testing.T) {
			_, problems, err := testsupport.ParseSource(t, resolutionTemplate(test.call))
			if err != nil {
				t.Fatalf("valid resolved call failed: %v\n%s", err, problems)
			}
		})
	}

	invalid := []struct {
		name      string
		call      string
		function  string
		reference string
	}{
		{"get_input-missing", "{ get_input: absent }", "get_input", "absent"},
		{"get_property-unknown-entity", "{ get_property: [absent, inherited_name] }", "get_property", "absent"},
		{"get_property-wrong-entity-kind", "{ get_property: [source, inherited_name] }", "get_property", "source"},
		{"get_property-unknown-property", "{ get_property: [node, absent] }", "get_property", "absent"},
		{"get_attribute-unknown-entity", "{ get_attribute: [absent, inherited_state] }", "get_attribute", "absent"},
		{"get_attribute-unknown-attribute", "{ get_attribute: [node, absent] }", "get_attribute", "absent"},
		{"get_operation_output-unknown-interface", "{ get_operation_output: [node, Absent, create, result] }", "get_operation_output", "Absent"},
		{"get_operation_output-unknown-operation", "{ get_operation_output: [node, Test, absent, result] }", "get_operation_output", "absent"},
		{"get_operation_output-unknown-output", "{ get_operation_output: [node, Test, create, absent] }", "get_operation_output", "absent"},
		{"get_artifact-unknown-artifact", "{ get_artifact: [node, absent] }", "get_artifact", "absent"},
	}
	for _, test := range invalid {
		t.Run(test.name, func(t *testing.T) {
			_, problems, err := testsupport.ParseSource(t, resolutionTemplate(test.call))
			if err == nil {
				t.Fatalf("unresolved %s reference was accepted", test.function)
			}
			if !strings.Contains(problems, "malformed "+test.function+" function") ||
				!strings.Contains(problems, test.reference) ||
				!strings.Contains(problems, `outputs["result"].value`) {
				t.Fatalf("wrong unresolved-reference diagnostic:\n%s", problems)
			}
			if strings.Contains(problems, "unknown data type") || strings.Contains(problems, "unsupported keyname") {
				t.Fatalf("unrelated earlier failure masked reference validation:\n%s", problems)
			}
		})
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 3.6.8 Import definition; 4.4.2.2 get_property parameters
// Requirements: TOSCA13-4.4.2.2-002, TOSCA13-4.4.2.2-010
// Expected: a namespace-qualified imported type contributes its inherited
// property definition to function reference resolution
// Category: positive, import, namespace, inheritance, resolution
func TestIntrinsicFunctionResolutionThroughNamespacedImport(t *testing.T) {
	directory := t.TempDir()
	importedPath := filepath.Join(directory, "imported.yaml")
	mainPath := filepath.Join(directory, "main.yaml")

	imported := `tosca_definitions_version: tosca_simple_yaml_1_3
node_types:
  ImportedNode:
    derived_from: tosca.nodes.Root
    properties:
      imported_name:
        type: string
        default: imported
`
	main := `tosca_definitions_version: tosca_simple_yaml_1_3
imports:
  - file: imported.yaml
    namespace_prefix: imported
node_types:
  LocalNode:
    derived_from: imported:ImportedNode
topology_template:
  node_templates:
    node:
      type: LocalNode
  outputs:
    result:
      value: { get_property: [node, imported_name] }
`
	if err := os.WriteFile(importedPath, []byte(imported), 0o600); err != nil {
		t.Fatalf("write imported fixture: %v", err)
	}
	if err := os.WriteFile(mainPath, []byte(main), 0o600); err != nil {
		t.Fatalf("write main fixture: %v", err)
	}

	serviceTemplate, parseProblems, err := testsupport.ParseFile(t, mainPath)
	if err != nil {
		t.Fatalf("namespaced imported function reference failed: %v\n%s", err, parseProblems)
	}
	call := normalizedOutputFunction(t, serviceTemplate)
	if call.Name != "tosca.function.get_property" {
		t.Fatalf("normalized function name = %q", call.Name)
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 4.1, Reserved Function Keywords; 4.4-4.8
// Requirements: TOSCA13-4.4.2.2-002, TOSCA13-4.5.1.2-002,
// TOSCA13-4.6.1.2-001, TOSCA13-4.8.1.2-002
// Expected: contextual keywords resolve only in a modelable-entity context
// Category: positive, negative, context
func TestIntrinsicFunctionContexts(t *testing.T) {
	validSource := `tosca_definitions_version: tosca_simple_yaml_1_3
node_types:
  ContextNode:
    derived_from: tosca.nodes.Root
    properties:
      source:
        type: string
      copy:
        type: string
topology_template:
  inputs:
    source:
      type: string
      default: value
  node_templates:
    node:
      type: ContextNode
      properties:
        source: { get_input: source }
        copy: { get_property: [SELF, source] }
  outputs:
    result:
      value: { get_property: [node, source] }
`
	if _, problems, err := testsupport.ParseSource(t, validSource); err != nil {
		t.Fatalf("valid property/output contexts failed: %v\n%s", err, problems)
	}

	relationshipSource := `tosca_definitions_version: tosca_simple_yaml_1_3
relationship_types:
  ContextRelationship:
    derived_from: tosca.relationships.Root
    properties:
      endpoint_name:
        type: string
topology_template:
  relationship_templates:
    rel:
      type: ContextRelationship
      properties:
        endpoint_name: { get_property: [SOURCE, name] }
`
	if _, problems, err := testsupport.ParseSource(t, relationshipSource); err != nil {
		t.Fatalf("valid SOURCE relationship context failed: %v\n%s", err, problems)
	}

	invalid := []struct {
		name string
		call string
		key  string
	}{
		{"self-in-topology-output", "{ get_property: [SELF, inherited_name] }", "SELF"},
		{"source-in-topology-output", "{ get_attribute: [SOURCE, inherited_state] }", "SOURCE"},
		{"target-in-topology-output", "{ get_artifact: [TARGET, config] }", "TARGET"},
	}
	for _, test := range invalid {
		t.Run(test.name, func(t *testing.T) {
			_, problems, err := testsupport.ParseSource(t, resolutionTemplate(test.call))
			if err == nil {
				t.Fatalf("context-invalid keyword %s was accepted", test.key)
			}
			if !strings.Contains(problems, "invalid context") || !strings.Contains(problems, test.key) {
				t.Fatalf("wrong context diagnostic:\n%s", problems)
			}
		})
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 3.8.3 Node Template; 3.6.20 Interface Assignment;
// 3.6.23.3 Call operation activity; 4.4-4.6 Intrinsic functions
// Requirements: TOSCA13-4.4.1.2-003, TOSCA13-4.4.2.2-002,
// TOSCA13-4.5.1.2-002, TOSCA13-4.6.1.2-001
// Expected: functions are accepted in property, attribute, output, interface
// operation input, and workflow call-operation input assignments
// Category: positive, context, property, attribute, output, interface, workflow
func TestIntrinsicFunctionAssignmentContexts(t *testing.T) {
	source := `tosca_definitions_version: tosca_simple_yaml_1_3
node_types:
  ContextNode:
    derived_from: tosca.nodes.Root
    properties:
      configured:
        type: string
    attributes:
      state:
        type: string
    interfaces:
      Test:
        type: tosca.interfaces.Root
        operations:
          create:
            inputs:
              source:
                type: string
            outputs:
              result:
                type: string
topology_template:
  inputs:
    source:
      type: string
      default: value
  node_templates:
    node:
      type: ContextNode
      properties:
        configured: { get_input: source }
      attributes:
        state: { get_operation_output: [SELF, Test, create, result] }
      interfaces:
        Test:
          operations:
            create:
              inputs:
                source: { get_input: source }
  workflows:
    deploy:
      steps:
        configure:
          target: node
          activities:
            - call_operation:
                operation: Test.create
                inputs:
                  source: { get_input: source }
  outputs:
    result:
      value: { get_property: [node, configured] }
`
	serviceTemplate, parseProblems, err := testsupport.ParseSource(t, source)
	if err != nil {
		t.Fatalf("valid assignment contexts failed: %v\n%s", err, parseProblems)
	}
	if serviceTemplate.Outputs["result"] == nil || serviceTemplate.Workflows["deploy"] == nil {
		t.Fatal("assignment contexts were not preserved through normalization")
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 4.3.1 concat; 4.3.3 token; 4.4.1 get_input
// Requirements: TOSCA13-4.3.1.2-003, TOSCA13-4.3.3.2-003,
// TOSCA13-4.3.3.2-006, TOSCA13-4.3.3.2-010,
// TOSCA13-4.4.1.2-003
// Expected: normalized calls evaluate deterministically to the specified value
// Category: positive, evaluation, result type, result value
func TestIntrinsicFunctionEvaluation(t *testing.T) {
	tests := []struct {
		name       string
		call       string
		want       any
		resolution bool
	}{
		{"concat", "{ concat: [prefix-, suffix] }", "prefix-suffix", false},
		{"token", `{ token: [a:b:c, ":", 1] }`, "b", false},
		{"get_input", "{ get_input: source }", "value", false},
		{"get_property", "{ get_property: [node, inherited_name] }", "inherited", true},
		{"get_attribute-unresolved-runtime", "{ get_attribute: [node, inherited_state] }", nil, true},
		{"get_artifact", "{ get_artifact: [node, config] }", "file:///etc/hosts", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			source := functionTemplate(test.call)
			if test.resolution {
				source = resolutionTemplate(test.call)
			}
			serviceTemplate, parseProblems, err := testsupport.ParseSource(t, source)
			if err != nil {
				t.Fatalf("parse failed: %v\n%s", err, parseProblems)
			}
			clout, err := serviceTemplate.Compile()
			if err != nil {
				t.Fatalf("compile failed: %v", err)
			}

			urlContext := exturl.NewContext()
			defer urlContext.Release()
			evaluationProblems := problems.NewProblems(terminal.NewStylist(false))
			execContext := cloutjs.ExecContext{
				Clout:      clout,
				Problems:   evaluationProblems,
				URLContext: urlContext,
				Format:     "yaml",
			}
			execContext.Coerce()
			if !evaluationProblems.Empty() {
				t.Fatalf("evaluation failed:\n%s", evaluationProblems.ToString(false))
			}

			tosca, ok := clout.Properties["tosca"].(ard.StringMap)
			if !ok {
				t.Fatalf("compiled TOSCA properties have type %T", clout.Properties["tosca"])
			}
			outputs, ok := tosca["outputs"].(ard.StringMap)
			if !ok {
				t.Fatalf("compiled outputs have type %T: %#v", tosca["outputs"], tosca["outputs"])
			}
			if got := outputs["result"]; !reflect.DeepEqual(got, test.want) {
				t.Fatalf("evaluated result = %#v (%T), want %#v (%T)", got, got, test.want, test.want)
			}
		})
	}
}

func assertFunctionRejected(t *testing.T, call string, function string, reason string) {
	t.Helper()
	_, problems, err := testsupport.ParseSource(t, functionTemplate(call))
	if err == nil {
		t.Fatalf("invalid %s call was accepted", function)
	}
	if !strings.Contains(problems, "malformed "+function+" function") {
		t.Fatalf("wrong error category for %s:\n%s", function, problems)
	}
	if !strings.Contains(problems, reason) {
		t.Fatalf("diagnostic does not identify %q for %s:\n%s", reason, function, problems)
	}
	if !strings.Contains(problems, `outputs["result"].value`) {
		t.Fatalf("diagnostic does not identify the function path:\n%s", problems)
	}
	if strings.Contains(problems, "unknown data type") || strings.Contains(problems, "unsupported keyname") {
		t.Fatalf("unrelated earlier failure masked function validation:\n%s", problems)
	}
}

func resolutionTemplate(call string) string {
	return fmt.Sprintf(`tosca_definitions_version: tosca_simple_yaml_1_3
node_types:
  ParentNode:
    derived_from: tosca.nodes.Root
    properties:
      inherited_name:
        type: string
        default: inherited
    attributes:
      inherited_state:
        type: string
    artifacts:
      config:
        type: tosca.artifacts.File
        file: config.txt
    interfaces:
      Test:
        type: tosca.interfaces.Root
        operations:
          create:
            outputs:
              result:
                type: string
  ChildNode:
    derived_from: ParentNode
topology_template:
  inputs:
    source:
      type: string
      default: value
  node_templates:
    node:
      type: ChildNode
      artifacts:
        config:
          type: tosca.artifacts.File
          file: /etc/hosts
  outputs:
    result:
      value: %s
`, call)
}

func functionTemplate(call string) string {
	return fmt.Sprintf(`tosca_definitions_version: tosca_simple_yaml_1_3
topology_template:
  inputs:
    source:
      type: string
      default: value
  node_templates:
    node:
      type: tosca.nodes.Root
  outputs:
    result:
      value: %s
`, call)
}

func normalizedOutputFunction(t *testing.T, serviceTemplate *normal.ServiceTemplate) *parsing.FunctionCall {
	t.Helper()
	value, ok := serviceTemplate.Outputs["result"]
	if !ok {
		t.Fatal("normalized output result is absent")
	}
	function, ok := value.(*normal.FunctionCall)
	if !ok {
		t.Fatalf("normalized output is %T, want *normal.FunctionCall", value)
	}
	return function.FunctionCall
}

func normalizedArguments(arguments []any) []any {
	normalized := make([]any, len(arguments))
	for index, argument := range arguments {
		switch value := argument.(type) {
		case *normal.Primitive:
			normalized[index] = value.Primitive
		default:
			normalized[index] = value
		}
	}
	return normalized
}
