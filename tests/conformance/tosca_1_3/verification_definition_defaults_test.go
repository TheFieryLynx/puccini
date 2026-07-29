package tosca_1_3_test

import (
	"strings"
	"testing"

	"github.com/tliron/go-puccini/normal"
	"github.com/tliron/go-puccini/tests/conformance/internal/testsupport"
)

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 3.6.10.4 Grammar; 3.6.10.5 Additional Requirements;
// 3.6.12.2 Grammar; 3.6.14.2 Grammar; 3.6.14.3 Additional Requirements
// Requirements: TOSCA13-3.6.10.4-008, TOSCA13-3.6.10.5-003,
// TOSCA13-3.6.12.2-010, TOSCA13-3.6.14.2-015,
// TOSCA13-3.6.14.3-002
// Expected: omitted required flags default to true; compatible property,
// attribute, and parameter defaults render through their effective data types;
// incompatible literal defaults are rejected during rendering.
// Category: positive, negative, default, rendering, normalization, direct
func TestVerificationDefinitionDefaults(t *testing.T) {
	t.Run("TOSCA13-3.6.10.4-008", func(t *testing.T) {
		valid := definitionDefaultsTemplate(`
    properties:
      configured:
        type: string
`, "", `
      properties:
        configured: ready
`)
		serviceTemplate, problems, err := testsupport.ParseSource(t, valid)
		if err != nil {
			t.Fatalf("assigned default-required property failed: %v\n%s", err, problems)
		}
		assertNormalizedPrimitive(t, serviceTemplate.NodeTemplates["node"].Properties["configured"], "ready")

		source := definitionDefaultsTemplate(`
    properties:
      configured:
        type: string
`, "", "")
		assertDefinitionDefaultRejected(
			t,
			source,
			`node_templates["node"].properties["configured"]`,
			"required",
		)
	})

	t.Run("TOSCA13-3.6.10.5-003", func(t *testing.T) {
		source := definitionDefaultsTemplate(`
    properties:
      count:
        type: integer
        default: 7
`, "", "")
		serviceTemplate, problems, err := testsupport.ParseSource(t, source)
		if err != nil {
			t.Fatalf("compatible property default failed: %v\n%s", err, problems)
		}
		assertNormalizedPrimitive(t, serviceTemplate.NodeTemplates["node"].Properties["count"], 7)

		assertDefinitionDefaultRejected(
			t,
			definitionDefaultsTemplate(`
    properties:
      count:
        type: integer
        default: invalid
`, "", ""),
			`node_types["example.Node"].properties["count"].default`,
			"integer",
		)
	})

	t.Run("TOSCA13-3.6.12.2-010", func(t *testing.T) {
		source := definitionDefaultsTemplate(`
    attributes:
      observed:
        type: integer
      derived:
        type: integer
        default: { get_attribute: [SELF, observed] }
`, "", "")
		serviceTemplate, problems, err := testsupport.ParseSource(t, source)
		if err != nil {
			t.Fatalf("compatible attribute default failed: %v\n%s", err, problems)
		}
		value := serviceTemplate.NodeTemplates["node"].Attributes["derived"]
		if _, ok := value.(*normal.FunctionCall); !ok {
			t.Fatalf("normalized attribute default has type %T, want function call", value)
		}
	})

	t.Run("TOSCA13-3.6.14.2-015", func(t *testing.T) {
		valid := definitionDefaultsTemplate("", `
    configured:
      type: string
      default: ready
`, "")
		serviceTemplate, problems, err := testsupport.ParseSource(t, valid)
		if err != nil {
			t.Fatalf("satisfied default-required input failed: %v\n%s", err, problems)
		}
		assertNormalizedPrimitive(t, serviceTemplate.Inputs["configured"], "ready")

		assertDefinitionDefaultRejected(
			t,
			definitionDefaultsTemplate("", `
    configured:
      type: string
`, ""),
			`inputs["configured"]`,
			"required",
		)
	})

	t.Run("TOSCA13-3.6.14.3-002", func(t *testing.T) {
		source := definitionDefaultsTemplate("", `
    count:
      type: integer
      default: 7
`, "")
		serviceTemplate, problems, err := testsupport.ParseSource(t, source)
		if err != nil {
			t.Fatalf("compatible parameter default failed: %v\n%s", err, problems)
		}
		assertNormalizedPrimitive(t, serviceTemplate.Inputs["count"], 7)

		assertDefinitionDefaultRejected(
			t,
			definitionDefaultsTemplate("", `
    count:
      type: integer
      default: invalid
`, ""),
			`inputs["count"].default`,
			"integer",
		)
	})
}

func definitionDefaultsTemplate(nodeDefinition, inputs, nodeProperties string) string {
	inputDefinitions := ""
	if inputs != "" {
		inputDefinitions = "\n  inputs:" + inputs
	}
	return `tosca_definitions_version: tosca_simple_yaml_1_3
description: TOSCA 1.3 definition-default verification
node_types:
  example.Node:
    derived_from: tosca.nodes.Root` + nodeDefinition + `
topology_template:` + inputDefinitions + `
  node_templates:
    node:
      type: example.Node` + nodeProperties + `
`
}

func assertDefinitionDefaultRejected(t *testing.T, source string, expected ...string) {
	t.Helper()
	_, problems, err := testsupport.ParseSource(t, source)
	if err == nil {
		t.Fatal("invalid definition default input was accepted")
	}
	for _, fragment := range expected {
		if !strings.Contains(problems, fragment) {
			t.Fatalf("definition-default diagnostic does not contain %q:\n%s", fragment, problems)
		}
	}
}
