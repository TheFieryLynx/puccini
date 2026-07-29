package tosca_1_3_test

import (
	"strings"
	"testing"

	"github.com/tliron/go-puccini/normal"
	"github.com/tliron/go-puccini/tests/conformance/internal/testsupport"
)

// TOSCA 1.3 §3.6.12.4, "Additional Requirements".
// Category: positive, negative, boundary, inheritance, function, provenance,
// normalization, and cross-version regression.
// Expected: an attribute-definition default is derived from actual-state
// attributes or operation outputs, never a literal, property, or input.
func TestPartialAttributeDefaultProvenanceAcceptsActualStateSources(t *testing.T) {
	tests := []struct {
		name       string
		expression string
	}{
		{name: "attribute", expression: "{ get_attribute: [SELF, observed] }"},
		{name: "operation-output", expression: "{ get_operation_output: [SELF, Standard, create, result] }"},
		{name: "calculated", expression: "{ concat: [{ get_attribute: [SELF, observed] }, -suffix] }"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			serviceTemplate, parseProblems, err := testsupport.ParseSource(
				t,
				attributeDefaultProvenanceTemplate(test.expression, "", ""),
			)
			if err != nil {
				t.Fatalf("actual-state attribute default failed: %v\n%s", err, parseProblems)
			}

			value := serviceTemplate.NodeTemplates["node"].Attributes["derived"]
			if _, ok := value.(*normal.FunctionCall); !ok {
				t.Fatalf("normalized attribute default has type %T, want function call", value)
			}
		})
	}
}

func TestPartialAttributeDefaultProvenanceRejectsForbiddenSources(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		reason     string
	}{
		{name: "literal-string", expression: "configured", reason: "hard-coded"},
		{name: "numeric-looking-string", expression: `"7"`, reason: "hard-coded"},
		{name: "input", expression: "{ get_input: configured }", reason: "get_input"},
		{name: "property", expression: "{ get_property: [SELF, configured] }", reason: "get_property"},
		{
			name:       "mixed-calculation",
			expression: "{ concat: [{ get_attribute: [SELF, observed] }, { get_input: configured }] }",
			reason:     "get_input",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			problems := rejectAttributeDefaultProvenance(
				t,
				attributeDefaultProvenanceTemplate(test.expression, "", ""),
			)
			for _, fragment := range []string{
				`node_types["example.Node"].attributes["derived"].default`,
				"actual-state attribute or operation output",
				test.reason,
			} {
				if !strings.Contains(problems, fragment) {
					t.Fatalf("provenance diagnostic does not contain %q:\n%s", fragment, problems)
				}
			}
		})
	}
}

func TestPartialAttributeDefaultProvenanceRejectsHardCodedStructuredLeaf(t *testing.T) {
	source := attributeDefaultProvenanceTemplate(
		"",
		`      derived:
        type: map
        entry_schema:
          type: string
        default:
          actual: { get_attribute: [SELF, observed] }
          literal: configured
`,
		"",
	)
	problems := rejectAttributeDefaultProvenance(t, source)
	if !strings.Contains(problems, `default["literal"]`) || !strings.Contains(problems, "hard-coded") {
		t.Fatalf("structured hard-coded leaf failed for the wrong reason:\n%s", problems)
	}
}

func TestPartialAttributeDefaultProvenanceInheritance(t *testing.T) {
	validInherited := attributeDefaultProvenanceTemplate(
		"{ get_attribute: [SELF, observed] }",
		"",
		`  example.Derived:
    derived_from: example.Node
`,
	)
	if _, parseProblems, err := testsupport.ParseSource(t, validInherited); err != nil {
		t.Fatalf("valid inherited attribute default failed: %v\n%s", err, parseProblems)
	}

	invalidOverride := attributeDefaultProvenanceTemplate(
		"{ get_attribute: [SELF, observed] }",
		"",
		`  example.Derived:
    derived_from: example.Node
    attributes:
      derived:
        type: string
        default: configured
`,
	)
	problems := rejectAttributeDefaultProvenance(t, invalidOverride)
	if !strings.Contains(problems, `node_types["example.Derived"].attributes["derived"].default`) {
		t.Fatalf("invalid refined default failed at the wrong path:\n%s", problems)
	}
}

func TestPartialAttributeDefaultProvenanceDiagnosticIsDeterministic(t *testing.T) {
	source := attributeDefaultProvenanceTemplate("{ get_input: configured }", "", "")
	first := rejectAttributeDefaultProvenance(t, source)
	second := rejectAttributeDefaultProvenance(t, source)
	if attributeDefaultDiagnosticBody(first) != attributeDefaultDiagnosticBody(second) {
		t.Fatalf("attribute default diagnostic is nondeterministic:\nfirst:\n%s\nsecond:\n%s", first, second)
	}
}

func TestPartialAttributeDefaultProvenanceDoesNotChangeTOSCA20(t *testing.T) {
	source := `tosca_definitions_version: tosca_2_0
description: TOSCA 2.0 attribute-default provenance regression
imports:
  - profile: org.oasis-open.simple:2.0
    namespace: tosca
node_types:
  example.Node:
    derived_from: tosca:Root
    attributes:
      configured:
        type: string
        default: literal
service_template:
  node_templates:
    node:
      type: example.Node
`
	if _, parseProblems, err := testsupport.ParseSource(t, source); err != nil {
		t.Fatalf("TOSCA 2.0 attribute-default behavior changed: %v\n%s", err, parseProblems)
	}
}

func TestPartialAttributeDefaultProvenancePreservesPinnedRootStateDefault(t *testing.T) {
	source := `tosca_definitions_version: tosca_simple_yaml_1_3
description: pinned normative Root attribute-default regression
topology_template:
  node_templates:
    node:
      type: tosca.nodes.Root
`
	serviceTemplate, parseProblems, err := testsupport.ParseSource(t, source)
	if err != nil {
		t.Fatalf("pinned Root state default failed: %v\n%s", err, parseProblems)
	}
	assertNormalizedPrimitive(t, serviceTemplate.NodeTemplates["node"].Attributes["state"], "initial")
}

func attributeDefaultProvenanceTemplate(expression string, attributeDefinition string, extraNodeTypes string) string {
	defaultDefinition := ""
	if expression != "" {
		defaultDefinition = `
      derived:
        type: string
        default: ` + expression
	}
	if attributeDefinition != "" {
		defaultDefinition = "\n" + attributeDefinition
	}

	return `tosca_definitions_version: tosca_simple_yaml_1_3
description: TOSCA 1.3 attribute-default provenance conformance
node_types:
  example.Node:
    derived_from: tosca.nodes.Root
    properties:
      configured:
        type: string
        required: false
    attributes:
      observed:
        type: string` + defaultDefinition + `
` + extraNodeTypes + `
topology_template:
  inputs:
    configured:
      type: string
      default: configured
  node_templates:
    node:
      type: example.Node
`
}

func rejectAttributeDefaultProvenance(t *testing.T, source string) string {
	t.Helper()
	_, parseProblems, err := testsupport.ParseSource(t, source)
	if err == nil {
		t.Fatalf("invalid attribute-definition default provenance was accepted")
	}
	return parseProblems
}

func attributeDefaultDiagnosticBody(problems string) string {
	if index := strings.Index(problems, "@"); index >= 0 {
		return problems[index:]
	}
	return problems
}
