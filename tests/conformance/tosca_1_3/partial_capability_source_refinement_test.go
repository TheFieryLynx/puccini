package tosca_1_3_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tliron/go-puccini/tests/conformance/internal/testsupport"
)

// TOSCA 1.3 §3.7.2.4, "Additional requirements".
// Category: positive, negative, boundary, inheritance, resolution, and
// cross-version regression.
// Expected: an explicitly refined valid_source_types list contains only the
// same node types or types derived from the parent capability definition's set.
func TestPartialCapabilitySourceRefinementAcceptsCompatibleTypes(t *testing.T) {
	tests := []struct {
		name        string
		childSource string
	}{
		{name: "same-type", childSource: "example.Source"},
		{name: "derived-type", childSource: "example.DerivedSource"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			source := capabilitySourceRefinementTemplate(test.childSource)
			if _, parseProblems, err := testsupport.ParseSource(t, source); err != nil {
				t.Fatalf("compatible valid_source_types refinement failed: %v\n%s", err, parseProblems)
			}
		})
	}
}

func TestPartialCapabilitySourceRefinementRejectsUnrelatedType(t *testing.T) {
	problems := rejectCapabilitySourceRefinement(t, capabilitySourceRefinementTemplate("example.Unrelated"))
	for _, fragment := range []string{
		`node_types["example.DerivedHost"].capabilities["custom_feature"].valid_source_types`,
		"example.Unrelated",
		"derived from one of the types in the parent set",
	} {
		if !strings.Contains(problems, fragment) {
			t.Fatalf("incompatible source diagnostic does not contain %q:\n%s", fragment, problems)
		}
	}
}

func TestPartialCapabilitySourceRefinementInheritsOmittedList(t *testing.T) {
	source := capabilitySourceRefinementTemplate("")
	if _, parseProblems, err := testsupport.ParseSource(t, source); err != nil {
		t.Fatalf("omitted valid_source_types did not inherit: %v\n%s", err, parseProblems)
	}
}

func TestPartialCapabilitySourceRefinementRejectsMixedList(t *testing.T) {
	problems := rejectCapabilitySourceRefinement(t, capabilitySourceRefinementTemplate(
		"example.DerivedSource, example.Unrelated",
	))
	if !strings.Contains(problems, "example.Unrelated") {
		t.Fatalf("mixed source list failed for the wrong reason:\n%s", problems)
	}
	if strings.Contains(problems, "example.DerivedSource") &&
		strings.Contains(problems, `type "example.DerivedSource" must be derived`) {
		t.Fatalf("compatible list member was rejected:\n%s", problems)
	}
}

func TestPartialCapabilitySourceRefinementUnknownTypeStillFailsLookup(t *testing.T) {
	problems := rejectCapabilitySourceRefinement(t, capabilitySourceRefinementTemplate("example.Missing"))
	if !strings.Contains(problems, "valid_source") || !strings.Contains(problems, "unknown node type") {
		t.Fatalf("unknown source type failed for the wrong reason:\n%s", problems)
	}
}

func TestPartialCapabilitySourceRefinementResolvesImportedTypes(t *testing.T) {
	directory := t.TempDir()
	importedPath := filepath.Join(directory, "sources.yaml")
	imported := `tosca_definitions_version: tosca_simple_yaml_1_3
namespace: urn:example:capability-sources
node_types:
  Source:
    derived_from: tosca.nodes.Root
  DerivedSource:
    derived_from: Source
`
	if err := os.WriteFile(importedPath, []byte(imported), 0o600); err != nil {
		t.Fatalf("write imported source types: %v", err)
	}

	mainPath := filepath.Join(directory, "main.yaml")
	main := `tosca_definitions_version: tosca_simple_yaml_1_3
imports:
  - file: sources.yaml
    namespace_prefix: ext
node_types:
  BaseHost:
    derived_from: tosca.nodes.Root
    capabilities:
      custom_feature:
        type: tosca.capabilities.Node
        valid_source_types: [ext:Source]
  DerivedHost:
    derived_from: BaseHost
    capabilities:
      custom_feature:
        type: tosca.capabilities.Node
        valid_source_types: [ext:DerivedSource]
`
	if err := os.WriteFile(mainPath, []byte(main), 0o600); err != nil {
		t.Fatalf("write importing template: %v", err)
	}
	if _, parseProblems, err := testsupport.ParseFile(t, mainPath); err != nil {
		t.Fatalf("qualified imported source types failed: %v\n%s", err, parseProblems)
	}
}

func TestPartialCapabilitySourceRefinementDiagnosticIsDeterministic(t *testing.T) {
	source := capabilitySourceRefinementTemplate("example.Unrelated")
	first := rejectCapabilitySourceRefinement(t, source)
	second := rejectCapabilitySourceRefinement(t, source)
	if capabilitySourceDiagnosticBody(first) != capabilitySourceDiagnosticBody(second) {
		t.Fatalf("source refinement diagnostic is nondeterministic:\nfirst:\n%s\nsecond:\n%s", first, second)
	}
}

func TestPartialCapabilitySourceRefinementDoesNotChangeTOSCA20(t *testing.T) {
	source := `tosca_definitions_version: tosca_2_0
description: TOSCA 2.0 capability source refinement regression
capability_types:
  example.Capability: {}
node_types:
  example.Source: {}
  example.Unrelated: {}
  example.BaseHost:
    capabilities:
      custom_feature:
        type: example.Capability
        valid_source_node_types: [example.Source]
  example.DerivedHost:
    derived_from: example.BaseHost
    capabilities:
      custom_feature:
        type: example.Capability
        valid_source_node_types: [example.Unrelated]
`
	if _, parseProblems, err := testsupport.ParseSource(t, source); err != nil {
		t.Fatalf("TOSCA 2.0 source refinement behavior changed: %v\n%s", err, parseProblems)
	}
}

func capabilitySourceRefinementTemplate(childSources string) string {
	childDefinition := `
        type: tosca.capabilities.Node`
	if childSources != "" {
		childDefinition += `
        valid_source_types: [` + childSources + `]`
	}

	return `tosca_definitions_version: tosca_simple_yaml_1_3
description: TOSCA 1.3 capability source refinement conformance
node_types:
  example.Source:
    derived_from: tosca.nodes.Root
  example.DerivedSource:
    derived_from: example.Source
  example.Unrelated:
    derived_from: tosca.nodes.Root
  example.BaseHost:
    derived_from: tosca.nodes.Root
    capabilities:
      custom_feature:
        type: tosca.capabilities.Node
        valid_source_types: [example.Source]
  example.DerivedHost:
    derived_from: example.BaseHost
    capabilities:
      custom_feature:` + childDefinition + `
`
}

func rejectCapabilitySourceRefinement(t *testing.T, source string) string {
	t.Helper()
	_, parseProblems, err := testsupport.ParseSource(t, source)
	if err == nil {
		t.Fatalf("incompatible valid_source_types refinement was accepted")
	}
	return parseProblems
}

func capabilitySourceDiagnosticBody(problems string) string {
	if index := strings.Index(problems, "@"); index >= 0 {
		return problems[index:]
	}
	return problems
}
