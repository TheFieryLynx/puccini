package tosca_1_3_test

import (
	"strings"
	"testing"

	"github.com/tliron/go-puccini/tests/conformance/internal/testsupport"
)

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.5.1, Required Keynames
// Requirement: TOSCA13-3.5.1-002
// Expected: a derived artifact definition may inherit both required keynames.
// Category: positive, inheritance, direct
func TestDerivedArtifactInheritsRequiredKeynames(t *testing.T) {
	source := inheritedArtifactSource(`
  ParentNode:
    derived_from: tosca.nodes.Root
    artifacts:
      payload:
        type: TestArtifact
        file: parent.bin
  ChildNode:
    derived_from: ParentNode
    artifacts:
      payload:
        description: inherited required keynames
`)
	if _, problems, err := testsupport.ParseSource(t, source); err != nil {
		t.Fatalf("derived artifact did not inherit type and file: %v\n%s", err, problems)
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.5.1, Required Keynames
// Requirement: TOSCA13-3.5.1-002
// Expected: an explicitly replaced type is combined with the inherited file.
// Category: positive, boundary, inheritance, direct
func TestDerivedArtifactMayOverrideOneInheritedKey(t *testing.T) {
	source := `tosca_definitions_version: tosca_simple_yaml_1_3
artifact_types:
  ParentArtifact:
    derived_from: tosca.artifacts.File
  ChildArtifact:
    derived_from: ParentArtifact
node_types:
  ParentNode:
    derived_from: tosca.nodes.Root
    artifacts:
      payload:
        type: ParentArtifact
        file: parent.bin
  ChildNode:
    derived_from: ParentNode
    artifacts:
      payload:
        type: ChildArtifact
        description: inherited file
`
	if _, problems, err := testsupport.ParseSource(t, source); err != nil {
		t.Fatalf("derived artifact did not combine explicit type with inherited file: %v\n%s", err, problems)
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 3.5.1, Required Keynames; 3.6.7.1, Artifact definition
// Requirement: TOSCA13-3.5.1-002
// Expected: a base long-form artifact definition cannot omit either required key.
// Category: negative, effective definition, direct
func TestBaseArtifactStillRequiresTypeAndFile(t *testing.T) {
	for _, test := range []struct {
		name       string
		definition string
		missing    string
	}{
		{"type", "file: payload.bin", "type"},
		{"file", "type: TestArtifact", "file"},
	} {
		t.Run(test.name, func(t *testing.T) {
			assertArtifactRequiredKeyRejected(t, inheritedArtifactSource(`
  BaseNode:
    derived_from: tosca.nodes.Root
    artifacts:
      payload:
        `+test.definition+`
`), test.missing)
		})
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 3.5.1, Required Keynames; 3.6.7.1, Artifact definition
// Requirement: TOSCA13-3.5.1-002
// Expected: a new child artifact has no parent definition to supply required keys.
// Category: negative, inheritance boundary, direct
func TestNewDerivedArtifactStillRequiresTypeAndFile(t *testing.T) {
	source := inheritedArtifactSource(`
  ParentNode:
    derived_from: tosca.nodes.Root
  ChildNode:
    derived_from: ParentNode
    artifacts:
      payload:
        description: no parent artifact
`)
	assertArtifactRequiredKeyRejected(t, source, "type")
	assertArtifactRequiredKeyRejected(t, source, "file")
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 3.5.1, Required Keynames; 3.6.7.2.1, short notation
// Requirement: TOSCA13-3.5.1-002
// Expected: scalar short notation retains its specified type-inference semantics.
// Category: positive, notation boundary, regression
func TestArtifactShortNotationInferenceRemainsValid(t *testing.T) {
	source := inheritedArtifactSource(`
  BaseNode:
    derived_from: tosca.nodes.Root
    artifacts:
      payload: payload.bin
`)
	if _, problems, err := testsupport.ParseSource(t, source); err != nil {
		t.Fatalf("artifact short notation was rejected: %v\n%s", err, problems)
	}
}

// Specification: TOSCA Version 2.0
// Category: cross-version regression
// Expected: the TOSCA 1.3 post-inheritance policy does not alter TOSCA 2.0.
func TestTosca20ArtifactRequirednessUnchanged(t *testing.T) {
	source := `tosca_definitions_version: tosca_2_0
artifact_types:
  TestArtifact: {}
node_types:
  ParentNode:
    artifacts:
      payload:
        type: TestArtifact
        file: parent.bin
  ChildNode:
    derived_from: ParentNode
    artifacts:
      payload:
        description: inherited required keynames
`
	if _, problems, err := testsupport.ParseSource(t, source); err != nil {
		t.Fatalf("TOSCA 2.0 artifact behavior changed: %v\n%s", err, problems)
	}
}

func inheritedArtifactSource(nodeTypes string) string {
	return `tosca_definitions_version: tosca_simple_yaml_1_3
artifact_types:
  TestArtifact:
    derived_from: tosca.artifacts.File
node_types:
` + nodeTypes
}

func assertArtifactRequiredKeyRejected(t *testing.T, source string, missing string) {
	t.Helper()
	_, problems, err := testsupport.ParseSource(t, source)
	if err == nil {
		t.Fatalf("artifact definition missing %s was accepted", missing)
	}
	for _, fragment := range []string{missing, "missing"} {
		if !strings.Contains(problems, fragment) {
			t.Fatalf("diagnostic for missing %s does not contain %q:\n%s", missing, fragment, problems)
		}
	}
}
