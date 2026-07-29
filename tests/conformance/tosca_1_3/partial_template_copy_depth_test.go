package tosca_1_3_test

import (
	"strings"
	"testing"

	"github.com/tliron/go-puccini/normal"
	"github.com/tliron/go-puccini/tests/conformance/internal/testsupport"
)

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.8.3.3, Additional requirements
// Requirement: TOSCA13-3.8.3.3-001
// Expected: a node template may copy a complete source template and override
// copied values.
// Category: positive, boundary, normalization, direct
func TestPartialNodeTemplateCopyAcceptsCompleteSource(t *testing.T) {
	serviceTemplate, problems, err := testsupport.ParseSource(t, nodeCopyTemplate("direct", "base"))
	if err != nil {
		t.Fatalf("one-level node copy was rejected: %v\n%s", err, problems)
	}

	node := serviceTemplate.NodeTemplates["direct"]
	if node == nil {
		t.Fatal("copied node template is absent from normalized output")
	}
	if node.Description != "direct node" {
		t.Fatalf("copied node description = %q, want local override", node.Description)
	}
	value, ok := node.Properties["marker"].(*normal.Primitive)
	if !ok || value.Primitive != "base value" {
		t.Fatalf("copied node property = %#v, want inherited base value", node.Properties["marker"])
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.8.3.3, Additional requirements
// Requirement: TOSCA13-3.8.3.3-001
// Expected: the source named by copy must not itself use copy.
// Category: negative, copy depth, parser phase, direct
func TestPartialNodeTemplateCopyRejectsCopiedSource(t *testing.T) {
	problems := rejectTemplateCopy(t, nodeCopyTemplate("leaf", "middle"))
	for _, fragment := range []string{
		`node_templates["leaf"].copy`,
		"middle",
		"source template",
		"must not itself use copy",
	} {
		if !strings.Contains(problems, fragment) {
			t.Fatalf("node copy-depth diagnostic does not contain %q:\n%s", fragment, problems)
		}
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.8.4.3, Additional requirements
// Requirement: TOSCA13-3.8.4.3-001
// Expected: a relationship template may copy a complete source template.
// Category: positive, boundary, direct
func TestPartialRelationshipTemplateCopyAcceptsCompleteSource(t *testing.T) {
	if _, problems, err := testsupport.ParseSource(t, relationshipCopyTemplate("direct", "base")); err != nil {
		t.Fatalf("one-level relationship copy was rejected: %v\n%s", err, problems)
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.8.4.3, Additional requirements
// Requirement: TOSCA13-3.8.4.3-001
// Expected: the source relationship template named by copy must not itself
// use copy.
// Category: negative, copy depth, parser phase, direct
func TestPartialRelationshipTemplateCopyRejectsCopiedSource(t *testing.T) {
	problems := rejectTemplateCopy(t, relationshipCopyTemplate("leaf", "middle"))
	for _, fragment := range []string{
		`relationship_templates["leaf"].copy`,
		"middle",
		"source template",
		"must not itself use copy",
	} {
		if !strings.Contains(problems, fragment) {
			t.Fatalf("relationship copy-depth diagnostic does not contain %q:\n%s", fragment, problems)
		}
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 3.8.3.3 and 3.8.4.3, Additional requirements
// Requirements: TOSCA13-3.8.3.3-001, TOSCA13-3.8.4.3-001
// Expected: copy-depth diagnostics are stable across processing operations.
// Category: negative, determinism, direct
func TestPartialTemplateCopyDepthDiagnosticIsDeterministic(t *testing.T) {
	source := nodeCopyTemplate("leaf", "middle")
	first := templateCopyDiagnosticBody(rejectTemplateCopy(t, source))
	second := templateCopyDiagnosticBody(rejectTemplateCopy(t, source))
	if first != second {
		t.Fatalf("copy-depth diagnostic is nondeterministic:\nfirst:\n%s\nsecond:\n%s", first, second)
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.8.3.3, Additional requirements
// Requirement: TOSCA13-3.8.3.3-001
// Expected: the existing copy-loop diagnostic remains distinct from the
// source-completeness diagnostic.
// Category: negative, cycle regression, direct
func TestPartialTemplateCopyLoopStillRejected(t *testing.T) {
	source := `tosca_definitions_version: tosca_simple_yaml_1_3
topology_template:
  node_templates:
    node:
      copy: node
`
	problems := rejectTemplateCopy(t, source)
	if !strings.Contains(problems, "endless loop") {
		t.Fatalf("copy loop failed for the wrong reason:\n%s", problems)
	}
}

// Specification: TOSCA Version 2.0
// Cross-version isolation: the TOSCA 1.3 source-completeness rule is not used
// as normative evidence for, or silently imposed on, TOSCA 2.0.
// Category: positive, cross-version regression
func TestPartialTemplateCopyDepthDoesNotChangeTOSCA20(t *testing.T) {
	source := `tosca_definitions_version: tosca_2_0
node_types:
  CopyNode:
    properties:
      marker:
        type: string
service_template:
  node_templates:
    base:
      type: CopyNode
      properties:
        marker: base
    middle:
      copy: base
    leaf:
      copy: middle
`
	serviceTemplate, problems, err := testsupport.ParseSource(t, source)
	if err != nil {
		t.Fatalf("TOSCA 2.0 recursive-copy behavior changed: %v\n%s", err, problems)
	}
	if serviceTemplate.NodeTemplates["leaf"] == nil {
		t.Fatal("TOSCA 2.0 recursively copied node is absent")
	}
}

func nodeCopyTemplate(target string, source string) string {
	return `tosca_definitions_version: tosca_simple_yaml_1_3
node_types:
  CopyNode:
    derived_from: tosca.nodes.Root
    properties:
      marker:
        type: string
topology_template:
  node_templates:
    base:
      type: CopyNode
      description: base node
      properties:
        marker: base value
    middle:
      copy: base
      description: middle node
    ` + target + `:
      copy: ` + source + `
      description: direct node
`
}

func relationshipCopyTemplate(target string, source string) string {
	return `tosca_definitions_version: tosca_simple_yaml_1_3
relationship_types:
  CopyRelationship:
    derived_from: tosca.relationships.Root
    properties:
      marker:
        type: string
topology_template:
  relationship_templates:
    base:
      type: CopyRelationship
      description: base relationship
      properties:
        marker: base value
    middle:
      copy: base
      description: middle relationship
    ` + target + `:
      copy: ` + source + `
      description: direct relationship
`
}

func rejectTemplateCopy(t *testing.T, source string) string {
	t.Helper()
	_, problems, err := testsupport.ParseSource(t, source)
	if err == nil {
		t.Fatal("template whose copy source also uses copy was accepted")
	}
	return problems
}

func templateCopyDiagnosticBody(problems string) string {
	if index := strings.Index(problems, "topology_template"); index >= 0 {
		return problems[index:]
	}
	return problems
}
