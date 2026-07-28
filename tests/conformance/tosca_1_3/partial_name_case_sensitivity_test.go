package tosca_1_3_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tliron/go-puccini/tests/conformance/internal/testsupport"
)

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 5.2, TOSCA normative type names; 5.2.1, Additional requirements
// Requirement: TOSCA13-5.2.1-001
// Expected: the exact Type URI, Shorthand, and Type Qualified names resolve.
// Category: positive, namespace resolution, three name forms, direct
func TestNormativeTypeNamesExactCaseAccepted(t *testing.T) {
	for _, name := range []string{"tosca.datatypes.json", "json", "tosca:json"} {
		t.Run(name, func(t *testing.T) {
			if _, parseProblems, err := testsupport.ParseSource(t, normativeDataTypeNameTemplate(name)); err != nil {
				t.Fatalf("exact normative type name %q was rejected: %v\n%s", name, err, parseProblems)
			}
		})
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 5.2, TOSCA normative type names; 5.2.1, Additional requirements
// Requirement: TOSCA13-5.2.1-001
// Expected: case variants of every normative name form are rejected.
// Category: negative, namespace resolution, three name forms, direct
func TestNormativeTypeNameCaseMismatchRejected(t *testing.T) {
	for _, name := range []string{"tosca.datatypes.Json", "JSON", "tosca:Json"} {
		t.Run(name, func(t *testing.T) {
			_, parseProblems, err := testsupport.ParseSource(t, normativeDataTypeNameTemplate(name))
			if err == nil {
				t.Fatalf("case-mismatched normative type name %q was accepted", name)
			}
			for _, fragment := range []string{"unknown data type", name} {
				if !strings.Contains(parseProblems, fragment) {
					t.Fatalf("diagnostic for %q does not contain %q:\n%s", name, fragment, parseProblems)
				}
			}
		})
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 3.6.8, Import definition; 5.2.1, Additional requirements
// Requirement: TOSCA13-5.2.1-001
// Expected: an imported qualified name resolves with its exact declared case.
// Category: positive, import, qualified name, direct
func TestImportedQualifiedTypeNameExactCaseAccepted(t *testing.T) {
	if _, parseProblems, err := testsupport.ParseFile(t, importedCaseTemplate(t, "ext:CaseNode")); err != nil {
		t.Fatalf("exact imported qualified type name was rejected: %v\n%s", err, parseProblems)
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 3.6.8, Import definition; 5.2.1, Additional requirements
// Requirement: TOSCA13-5.2.1-001
// Expected: an imported qualified name with different case does not resolve.
// Category: negative, import, qualified name, direct
func TestImportedQualifiedTypeNameCaseMismatchRejected(t *testing.T) {
	path := importedCaseTemplate(t, "ext:casenode")
	_, parseProblems, err := testsupport.ParseFile(t, path)
	if err == nil {
		t.Fatal("case-mismatched imported qualified type name was accepted")
	}
	for _, fragment := range []string{"ext:casenode", "reference to unknown node type"} {
		if !strings.Contains(parseProblems, fragment) {
			t.Fatalf("imported-name diagnostic does not contain %q:\n%s", fragment, parseProblems)
		}
	}
}

// Specification: TOSCA Version 2.0
// Section: namespace name resolution (cross-version regression)
// Expected: exact custom type-name matching remains case-sensitive.
// Category: positive, negative, cross-version regression
func TestTosca20TypeNameCaseSensitivityUnchanged(t *testing.T) {
	valid := `tosca_definitions_version: tosca_2_0
node_types:
  CaseNode: {}
service_template:
  node_templates:
    node:
      type: CaseNode
`
	if _, parseProblems, err := testsupport.ParseSource(t, valid); err != nil {
		t.Fatalf("TOSCA 2.0 exact type name failed: %v\n%s", err, parseProblems)
	}
	invalid := strings.Replace(valid, "type: CaseNode", "type: casenode", 1)
	if _, _, err := testsupport.ParseSource(t, invalid); err == nil {
		t.Fatal("TOSCA 2.0 case-mismatched type name was accepted")
	}
}

func normativeDataTypeNameTemplate(name string) string {
	return `tosca_definitions_version: tosca_simple_yaml_1_3
node_types:
  TestNode:
    derived_from: tosca.nodes.Root
    properties:
      document:
        type: ` + name + `
`
}

func importedCaseTemplate(t *testing.T, parent string) string {
	t.Helper()
	directory := t.TempDir()
	importedPath := filepath.Join(directory, "imported.yaml")
	mainPath := filepath.Join(directory, "main.yaml")
	imported := `tosca_definitions_version: tosca_simple_yaml_1_3
node_types:
  CaseNode:
    derived_from: tosca.nodes.Root
`
	main := `tosca_definitions_version: tosca_simple_yaml_1_3
imports:
  - file: imported.yaml
    namespace_prefix: ext
node_types:
  LocalNode:
    derived_from: ` + parent + `
`
	if err := os.WriteFile(importedPath, []byte(imported), 0o600); err != nil {
		t.Fatalf("write imported fixture: %v", err)
	}
	if err := os.WriteFile(mainPath, []byte(main), 0o600); err != nil {
		t.Fatalf("write main fixture: %v", err)
	}
	return mainPath
}
