package tosca_1_3_test

import (
	"strings"
	"testing"

	"github.com/tliron/go-puccini/tests/conformance/internal/testsupport"
)

// TOSCA 1.3 §3.8.2.2.3, "Extended grammar with Property Assignments for
// the relationship's Interfaces".
// Category: positive, negative, boundary, resolution, normalization, and
// cross-version regression.
// Expected: node_filter is valid only when the assignment's explicit node
// keyname resolves to a Node Type, never a Node Template.
func TestPartialRequirementNodeFilterAcceptsNodeType(t *testing.T) {
	serviceTemplate, parseProblems, err := testsupport.ParseSource(
		t,
		requirementNodeFilterTemplate("example.Target", true),
	)
	if err != nil {
		t.Fatalf("node type with node_filter failed: %v\n%s", err, parseProblems)
	}

	requirements := serviceTemplate.NodeTemplates["source"].Requirements
	if len(requirements) != 1 {
		t.Fatalf("normalized requirement count = %d, want 1", len(requirements))
	}
	requirement := requirements[0]
	if requirement.NodeTypeName == nil ||
		(*requirement.NodeTypeName != "example.Target" && !strings.HasSuffix(*requirement.NodeTypeName, "::example.Target")) {
		t.Fatalf("normalized node type = %v, want example.Target", requirementNodeTypeName(requirement.NodeTypeName))
	}
	if len(requirement.NodeTemplatePropertyValidation) == 0 {
		t.Fatalf("normalized node_filter was discarded")
	}
}

func TestPartialRequirementNodeFilterRejectsNodeTemplate(t *testing.T) {
	problems := rejectRequirementNodeFilter(
		t,
		requirementNodeFilterTemplate("target", true),
	)
	for _, fragment := range []string{
		`node_templates["source"].requirements{0}.node`,
		"target",
		"node_filter",
		"Node Type",
		"Node Template",
	} {
		if !strings.Contains(problems, fragment) {
			t.Fatalf("node-template diagnostic does not contain %q:\n%s", fragment, problems)
		}
	}
}

func TestPartialRequirementNodeFilterRejectsMissingNode(t *testing.T) {
	problems := rejectRequirementNodeFilter(
		t,
		requirementNodeFilterTemplate("", true),
	)
	if !strings.Contains(problems, "node_filter") ||
		!strings.Contains(problems, "requires an explicit node keyname") {
		t.Fatalf("missing-node diagnostic has the wrong reason:\n%s", problems)
	}
}

func TestPartialRequirementNodeFilterWithoutFilterMayTargetTemplate(t *testing.T) {
	if _, parseProblems, err := testsupport.ParseSource(
		t,
		requirementNodeFilterTemplate("target", false),
	); err != nil {
		t.Fatalf("node-template assignment without node_filter failed: %v\n%s", err, parseProblems)
	}
}

func TestPartialRequirementNodeFilterUnknownNodeStillFailsLookup(t *testing.T) {
	problems := rejectRequirementNodeFilter(
		t,
		requirementNodeFilterTemplate("example.Missing", true),
	)
	if !strings.Contains(problems, "node") || !strings.Contains(problems, "unknown") {
		t.Fatalf("unknown node failed for the wrong reason:\n%s", problems)
	}
	if strings.Contains(problems, "requires an explicit node keyname") {
		t.Fatalf("unknown node also produced a misleading missing-node diagnostic:\n%s", problems)
	}
}

func TestPartialRequirementNodeFilterDiagnosticIsDeterministic(t *testing.T) {
	source := requirementNodeFilterTemplate("target", true)
	first := rejectRequirementNodeFilter(t, source)
	second := rejectRequirementNodeFilter(t, source)
	if requirementNodeFilterDiagnosticBody(first) != requirementNodeFilterDiagnosticBody(second) {
		t.Fatalf("node_filter diagnostic is nondeterministic:\nfirst:\n%s\nsecond:\n%s", first, second)
	}
}

func TestPartialRequirementNodeFilterDoesNotChangeTOSCA20(t *testing.T) {
	source := `tosca_definitions_version: tosca_2_0
description: TOSCA 2.0 requirement node_filter regression
imports:
  - profile: org.oasis-open.simple:2.0
    namespace: tosca
capability_types:
  example.Capability:
    derived_from: tosca:Node
node_types:
  example.Source:
    derived_from: tosca:Root
    requirements:
      - target:
          capability: example.Capability
          node: example.Target
  example.Target:
    derived_from: tosca:Root
    properties:
      size:
        type: integer
    capabilities:
      feature: example.Capability
service_template:
  node_templates:
    target:
      type: example.Target
      properties:
        size: 2
    source:
      type: example.Source
      requirements:
        - target:
            node: target
            node_filter:
              $equal:
                - 1
                - 1
`
	if _, parseProblems, err := testsupport.ParseSource(t, source); err != nil {
		t.Fatalf("TOSCA 2.0 requirement node_filter behavior changed: %v\n%s", err, parseProblems)
	}
}

func requirementNodeFilterTemplate(node string, includeFilter bool) string {
	assignment := ""
	if node != "" {
		assignment += "\n            node: " + node
	}
	if includeFilter {
		assignment += `
            node_filter:
              properties:
                - size:
                    equal: 1`
	}

	return `tosca_definitions_version: tosca_simple_yaml_1_3
description: TOSCA 1.3 requirement node_filter conformance
capability_types:
  example.Capability:
    derived_from: tosca.capabilities.Node
node_types:
  example.Source:
    derived_from: tosca.nodes.Root
    requirements:
      - target:
          capability: example.Capability
          node: example.Target
  example.Target:
    derived_from: tosca.nodes.Root
    properties:
      size:
        type: integer
    capabilities:
      feature:
        type: example.Capability
topology_template:
  node_templates:
    target:
      type: example.Target
      properties:
        size: 2
    source:
      type: example.Source
      requirements:
        - target:` + assignment + `
`
}

func rejectRequirementNodeFilter(t *testing.T, source string) string {
	t.Helper()
	_, parseProblems, err := testsupport.ParseSource(t, source)
	if err == nil {
		t.Fatalf("invalid requirement node_filter assignment was accepted")
	}
	return parseProblems
}

func requirementNodeFilterDiagnosticBody(problems string) string {
	if index := strings.Index(problems, "@"); index >= 0 {
		return problems[index:]
	}
	return problems
}

func requirementNodeTypeName(value *string) string {
	if value == nil {
		return "<nil>"
	}
	return *value
}
