package tosca_1_3_test

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/tliron/go-puccini/normal"
	"github.com/tliron/go-puccini/tests/conformance/internal/testsupport"
)

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 4.3.2.2, Parameters
// Requirement: TOSCA13-4.3.2.2-003
// Expected: join accepts the required non-empty list and optional delimiter.
// Category: positive, boundary, normalization, direct
func TestJoinRequiredListAccepted(t *testing.T) {
	tests := []struct {
		name string
		call string
	}{
		{"without-delimiter", "{ join: [[alpha, beta]] }"},
		{"with-delimiter", `{ join: [[alpha, beta], "-"] }`},
		{"single-element", "{ join: [[alpha]] }"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			serviceTemplate, parseProblems, err := testsupport.ParseSource(t, functionTemplate(test.call))
			if err != nil {
				t.Fatalf("valid join call was rejected: %v\n%s", err, parseProblems)
			}
			call := normalizedOutputFunction(t, serviceTemplate)
			if call.Name != "tosca.function.join" {
				t.Fatalf("normalized function name = %q", call.Name)
			}
			if len(call.Arguments) < 1 || len(call.Arguments) > 2 {
				t.Fatalf("normalized join arguments = %#v", call.Arguments)
			}
		})
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 4.3.2.2, Parameters
// Requirement: TOSCA13-4.3.2.2-003
// Expected: string value expressions are valid list elements and list producers.
// Category: positive, nested expression, direct
func TestJoinStringExpressionListAccepted(t *testing.T) {
	source := `tosca_definitions_version: tosca_simple_yaml_1_3
topology_template:
  inputs:
    source:
      type: string
      default: beta
    parts:
      type: list
      entry_schema:
        type: string
      default: [alpha, beta]
  outputs:
    nested:
      value: { join: [[alpha, { get_input: source }], "-"] }
    produced:
      value: { join: [{ get_input: parts }, "-"] }
`
	serviceTemplate, parseProblems, err := testsupport.ParseSource(t, source)
	if err != nil {
		t.Fatalf("join string expressions were rejected: %v\n%s", err, parseProblems)
	}
	for _, name := range []string{"nested", "produced"} {
		value, ok := serviceTemplate.Outputs[name].(*normal.FunctionCall)
		if !ok || value.FunctionCall.Name != "tosca.function.join" {
			t.Fatalf("output %s did not normalize as join: %T", name, serviceTemplate.Outputs[name])
		}
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 4.3.2.2, Parameters
// Requirement: TOSCA13-4.3.2.2-003
// Expected: every violation of the required list/delimiter shape is rejected in read.
// Category: negative, argument grammar, boundary, direct
func TestJoinRequiredListRejected(t *testing.T) {
	tests := []struct {
		name   string
		call   string
		reason string
	}{
		{"missing", "{ join: [] }", "requires list"},
		{"scalar-first", "{ join: [alpha] }", "first argument"},
		{"empty-list", "{ join: [[]] }", "one or more"},
		{"non-string-element", "{ join: [[alpha, 1]] }", "list element 2"},
		{"bad-delimiter", "{ join: [[alpha], 1] }", "delimiter"},
		{"excess", `{ join: [[alpha], "-", extra] }`, "at most 2"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assertFunctionRejected(t, test.call, "join", test.reason)
		})
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 4.3.2.2, Parameters
// Requirement: TOSCA13-4.3.2.2-003
// Expected: one list element and either absence or presence of a string delimiter are valid boundaries.
// Category: positive, boundary, direct
func TestJoinSingleElementAndOptionalDelimiter(t *testing.T) {
	for _, call := range []string{"{ join: [[only]] }", `{ join: [[only], ":"] }`} {
		if _, parseProblems, err := testsupport.ParseSource(t, functionTemplate(call)); err != nil {
			t.Fatalf("valid join boundary failed: %v\n%s", err, parseProblems)
		}
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 4.7.1.2, Parameters
// Requirement: TOSCA13-4.7.1.2-003
// Expected: one known Node Type name is accepted and normalized.
// Category: positive, resolution, normalization, direct
func TestGetNodesOfTypeRequiredNameAccepted(t *testing.T) {
	serviceTemplate, parseProblems, err := testsupport.ParseSource(t, getNodesOfTypeTemplate("ParentNode"))
	if err != nil {
		t.Fatalf("known node type name was rejected: %v\n%s", err, parseProblems)
	}
	call := normalizedOutputFunction(t, serviceTemplate)
	if call.Name != "tosca.function.get_nodes_of_type" {
		t.Fatalf("normalized function name = %q", call.Name)
	}
	if got := normalizedArguments(call.Arguments); !reflect.DeepEqual(got, []any{"ParentNode"}) {
		t.Fatalf("normalized arguments = %#v", got)
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 4.7.1.2, Parameters
// Requirement: TOSCA13-4.7.1.2-003
// Expected: a missing or non-string Node Type name is rejected in read.
// Category: negative, argument grammar, direct
func TestGetNodesOfTypeRequiredNameRejected(t *testing.T) {
	for _, test := range []struct {
		call   string
		reason string
	}{
		{"{ get_nodes_of_type: [] }", "requires node_type_name"},
		{"{ get_nodes_of_type: 1 }", "must be a string"},
		{"{ get_nodes_of_type: [ParentNode, extra] }", "exactly 1"},
	} {
		assertFunctionRejected(t, test.call, "get_nodes_of_type", test.reason)
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 4.7.1.2, Parameters
// Requirement: TOSCA13-4.7.1.2-003
// Expected: the required name resolves specifically to a Node Type.
// Category: negative, name resolution, kind, direct
func TestGetNodesOfTypeUnknownOrWrongKindRejected(t *testing.T) {
	for _, test := range []struct {
		name string
		want string
	}{
		{"AbsentNode", "not found"},
		{"string", "not a Node Type"},
	} {
		_, parseProblems, err := testsupport.ParseSource(t, getNodesOfTypeTemplate(test.name))
		if err == nil {
			t.Fatalf("invalid node type name %q was accepted", test.name)
		}
		for _, fragment := range []string{"get_nodes_of_type", test.name, test.want} {
			if !strings.Contains(parseProblems, fragment) {
				t.Fatalf("diagnostic for %q does not contain %q:\n%s", test.name, fragment, parseProblems)
			}
		}
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 4.7.1.2, Parameters
// Requirement: TOSCA13-4.7.1.2-003
// Expected: a derived Node Type is a valid target name.
// Category: positive, hierarchy, boundary, direct
func TestGetNodesOfDerivedTypeAccepted(t *testing.T) {
	if _, parseProblems, err := testsupport.ParseSource(t, getNodesOfTypeTemplate("ChildNode")); err != nil {
		t.Fatalf("derived node type was rejected: %v\n%s", err, parseProblems)
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 3.6.8, Import definition; 4.7.1.2, Parameters
// Requirement: TOSCA13-4.7.1.2-003
// Expected: a namespace-qualified imported Node Type satisfies the required name.
// Category: positive, import, namespace, resolution, direct
func TestGetNodesOfTypeResolvesImportedType(t *testing.T) {
	directory := t.TempDir()
	importedPath := filepath.Join(directory, "imported.yaml")
	mainPath := filepath.Join(directory, "main.yaml")
	imported := `tosca_definitions_version: tosca_simple_yaml_1_3
node_types:
  ImportedNode:
    derived_from: tosca.nodes.Root
`
	main := `tosca_definitions_version: tosca_simple_yaml_1_3
imports:
  - file: imported.yaml
    namespace_prefix: imported
topology_template:
  outputs:
    result:
      value: { get_nodes_of_type: imported:ImportedNode }
`
	if err := os.WriteFile(importedPath, []byte(imported), 0o600); err != nil {
		t.Fatalf("write imported fixture: %v", err)
	}
	if err := os.WriteFile(mainPath, []byte(main), 0o600); err != nil {
		t.Fatalf("write main fixture: %v", err)
	}
	if _, parseProblems, err := testsupport.ParseFile(t, mainPath); err != nil {
		t.Fatalf("imported node type name was rejected: %v\n%s", err, parseProblems)
	}
}

// Specification: TOSCA Version 2.0
// Section: function invocation (cross-version isolation)
// Expected: TOSCA 1.3 argument and resolution checks do not alter TOSCA 2.0 calls.
// Category: positive, cross-version regression
func TestTosca20JoinAndGetNodesOfTypeBehaviorUnchanged(t *testing.T) {
	source := `tosca_definitions_version: tosca_2_0
node_types:
  TestNode: {}
service_template:
  outputs:
    joined:
      value: { $join: [[alpha, beta], "-"] }
    nodes:
      value: { $get_nodes_of_type: TestNode }
`
	serviceTemplate, parseProblems, err := testsupport.ParseSource(t, source)
	if err != nil {
		t.Fatalf("TOSCA 2.0 function behavior changed: %v\n%s", err, parseProblems)
	}
	if _, ok := serviceTemplate.Outputs["joined"].(*normal.FunctionCall); !ok {
		t.Fatalf("TOSCA 2.0 join normalized as %T", serviceTemplate.Outputs["joined"])
	}
	if _, ok := serviceTemplate.Outputs["nodes"].(*normal.FunctionCall); !ok {
		t.Fatalf("TOSCA 2.0 get_nodes_of_type normalized as %T", serviceTemplate.Outputs["nodes"])
	}
}

func getNodesOfTypeTemplate(typeName string) string {
	return `tosca_definitions_version: tosca_simple_yaml_1_3
node_types:
  ParentNode:
    derived_from: tosca.nodes.Root
  ChildNode:
    derived_from: ParentNode
topology_template:
  outputs:
    result:
      value: { get_nodes_of_type: ` + typeName + ` }
`
}
