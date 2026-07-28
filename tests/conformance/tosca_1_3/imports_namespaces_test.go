package tosca_1_3_test

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tliron/go-puccini/tests/conformance/internal/testsupport"
)

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 3.6.8.1 Keynames; 3.6.8.2.1 Single-line grammar;
// 3.6.8.2.2 Multi-line grammar; 3.6.8.2.4 Import URI processing requirements
// Requirements: TOSCA13-3.6.8.1-002, TOSCA13-3.6.8.1-003,
// TOSCA13-3.6.8.2.2-003
// Expected: short, long, relative, repository, and namespaced imports resolve
// Category: positive, import grammar, repository, relative path, resolution
func TestImportNotationAndResolution(t *testing.T) {
	directory := t.TempDir()
	writeToscaFixture(t, filepath.Join(directory, "imported.yaml"), importedNodeType("ImportedNode", ""))
	writeToscaFixture(t, filepath.Join(directory, "types", "relative.yaml"), importedNodeType("RelativeNode", ""))

	repositoryDirectory := filepath.Join(directory, "repository")
	writeToscaFixture(t, filepath.Join(repositoryDirectory, "repository.yaml"), importedNodeType("RepositoryNode", ""))
	repositoryURL := (&url.URL{Scheme: "file", Path: filepath.ToSlash(repositoryDirectory) + "/"}).String()

	tests := []struct {
		name       string
		imports    string
		parentName string
		extra      string
	}{
		{"short", "  - imported.yaml", "ImportedNode", ""},
		{"long", "  - file: imported.yaml\n    namespace_prefix: ext", "ext:ImportedNode", ""},
		{"relative", "  - types/relative.yaml", "RelativeNode", ""},
		{
			"repository",
			"  - file: repository.yaml\n    repository: definitions",
			"RepositoryNode",
			fmt.Sprintf("repositories:\n  definitions: %s\n", repositoryURL),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mainPath := filepath.Join(directory, "main-"+test.name+".yaml")
			source := fmt.Sprintf(`tosca_definitions_version: tosca_simple_yaml_1_3
%simports:
%s
node_types:
  LocalNode:
    derived_from: %s
`, test.extra, test.imports, test.parentName)
			writeToscaFixture(t, mainPath, source)
			if _, problems, err := testsupport.ParseFile(t, mainPath); err != nil {
				t.Fatalf("valid %s import failed: %v\n%s", test.name, err, problems)
			}
		})
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 3.1 TOSCA Namespace URI and alias; 3.1.1 TOSCA Namespace prefix;
// 3.1.2 TOSCA Namespacing in TOSCA Service Templates
// Requirements: TOSCA13-3.1-001, TOSCA13-3.1.1-001,
// TOSCA13-3.1.2-002, TOSCA13-3.1.2-003, TOSCA13-3.1.2-004,
// TOSCA13-3.1.2-005
// Expected: the selected 1.3 grammar establishes normative types and permits
// a non-reserved default namespace
// Category: positive, default namespace, implicit import, namespace resolution
func TestNamespaceDeclarations(t *testing.T) {
	source := `tosca_definitions_version: tosca_simple_yaml_1_3
namespace: urn:example:application
node_types:
  LocalNode:
    derived_from: tosca.nodes.Root
topology_template:
  node_templates:
    node:
      type: LocalNode
`
	if _, problems, err := testsupport.ParseSource(t, source); err != nil {
		t.Fatalf("valid namespace and normative type lookup failed: %v\n%s", err, problems)
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.1.2 TOSCA Namespacing in TOSCA Service Templates
// Requirement: TOSCA13-3.1.2-001
// Expected: tosca_definitions_version is physically the first YAML key line
// Category: positive, negative, read phase, first-line grammar
func TestDefinitionsVersionFirstLine(t *testing.T) {
	valid := `tosca_definitions_version: tosca_simple_yaml_1_3
description: valid
`
	if _, problems, err := testsupport.ParseSource(t, valid); err != nil {
		t.Fatalf("first-line definitions version failed: %v\n%s", err, problems)
	}

	withLeadingComments := `# comments are not YAML mapping keys

tosca_definitions_version: tosca_simple_yaml_1_3
description: valid
`
	if _, problems, err := testsupport.ParseSource(t, withLeadingComments); err != nil {
		t.Fatalf("definitions version as first YAML key failed: %v\n%s", err, problems)
	}

	invalid := `description: invalid ordering
tosca_definitions_version: tosca_simple_yaml_1_3
`
	_, problems, err := testsupport.ParseSource(t, invalid)
	if err == nil {
		t.Fatal("non-first-line tosca_definitions_version was accepted")
	}
	assertDiagnostic(t, problems, "tosca_definitions_version", "first line")
	assertNoUnrelatedImportFailure(t, problems)
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.1.3.1 Additional Requirements
// Requirements: TOSCA13-3.1.3.1-001, TOSCA13-3.1.3.1-002
// Expected: non-OASIS templates cannot declare the reserved TOSCA URI or a
// path within its reserved hierarchy as their default namespace
// Category: negative, read phase, reserved namespace
func TestReservedNamespacePolicy(t *testing.T) {
	for _, namespace := range []string{
		"http://docs.oasis-open.org/tosca",
		"http://docs.oasis-open.org/tosca/vendor/profile",
	} {
		t.Run(strings.TrimPrefix(namespace, "http://"), func(t *testing.T) {
			source := fmt.Sprintf(`tosca_definitions_version: tosca_simple_yaml_1_3
namespace: %s
`, namespace)
			_, problems, err := testsupport.ParseSource(t, source)
			if err == nil {
				t.Fatalf("reserved namespace %q was accepted", namespace)
			}
			assertDiagnostic(t, problems, "namespace", namespace, "reserved")
			assertNoUnrelatedImportFailure(t, problems)
		})
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 3.6.8.1 Keynames; 3.6.8.2.2 Multi-line grammar;
// 3.6.8.2.3 Requirements
// Requirements: TOSCA13-3.6.8.1-002, TOSCA13-3.6.8.1-003,
// TOSCA13-3.6.8.2.2-003, TOSCA13-3.6.8.2.3-002
// Expected: malformed imports and reserved namespace_prefix are rejected in
// the read phase at the responsible argument
// Category: negative, YAML type, required key, reserved prefix
func TestImportGrammarErrors(t *testing.T) {
	tests := []struct {
		name       string
		importItem string
		diagnostic []string
	}{
		{"wrong-import-type", "  - 1", []string{"imports", "map", "string"}},
		{"missing-file", "  - namespace_prefix: ext", []string{"imports", "file", "required"}},
		{"empty-file", `  - file: ""`, []string{"invalid import file", "must not be empty"}},
		{"wrong-file-type", "  - file: 1", []string{"imports", "file", "string"}},
		{"reserved-prefix", "  - file: imported.yaml\n    namespace_prefix: tosca", []string{"namespace_prefix", "tosca", "reserved"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			directory := t.TempDir()
			writeToscaFixture(t, filepath.Join(directory, "imported.yaml"), importedNodeType("ImportedNode", ""))
			mainPath := filepath.Join(directory, "main.yaml")
			writeToscaFixture(t, mainPath, "tosca_definitions_version: tosca_simple_yaml_1_3\nimports:\n"+test.importItem+"\n")

			_, problems, err := testsupport.ParseFile(t, mainPath)
			if err == nil {
				t.Fatalf("invalid import %s was accepted", test.name)
			}
			assertDiagnostic(t, problems, test.diagnostic...)
			assertNoUnrelatedImportFailure(t, problems)
		})
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 3.6.8 Import definition; 3.6.8.1 Keynames
// Requirements: TOSCA13-3.6.8.1-002, TOSCA13-3.6.8.1-003
// Expected: missing resources and unknown repositories fail with a stable
// diagnostic naming the import or repository
// Category: negative, import resolution, repository reference
func TestImportResolutionErrors(t *testing.T) {
	t.Run("missing-file", func(t *testing.T) {
		source := `tosca_definitions_version: tosca_simple_yaml_1_3
imports:
  - absent.yaml
`
		_, problems, err := testsupport.ParseSource(t, source)
		if err == nil {
			t.Fatal("missing imported file was accepted")
		}
		assertDiagnostic(t, problems, "absent.yaml")
	})

	t.Run("unknown-repository", func(t *testing.T) {
		directory := t.TempDir()
		writeToscaFixture(t, filepath.Join(directory, "imported.yaml"), importedNodeType("ImportedNode", ""))
		mainPath := filepath.Join(directory, "main.yaml")
		writeToscaFixture(t, mainPath, `tosca_definitions_version: tosca_simple_yaml_1_3
imports:
  - file: imported.yaml
    repository: absent
`)
		_, problems, err := testsupport.ParseFile(t, mainPath)
		if err == nil {
			t.Fatal("unknown import repository was accepted")
		}
		assertDiagnostic(t, problems, "repository", "absent")
		assertNoUnrelatedImportFailure(t, problems)
	})
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 3.1.3 Rules to avoid namespace collisions;
// 3.1.3.1 Additional Requirements
// Requirements: TOSCA13-3.1.3.1-004, TOSCA13-3.1.3.1-005
// Expected: conflicting prefixes fail deterministically; distinct prefixes
// may expose types from the same declared namespace
// Category: positive, negative, prefix collision, namespace URI
func TestNamespacePrefixCollisions(t *testing.T) {
	t.Run("duplicate-prefix", func(t *testing.T) {
		directory := t.TempDir()
		writeToscaFixture(t, filepath.Join(directory, "one.yaml"), importedNodeType("One", "urn:example:one"))
		writeToscaFixture(t, filepath.Join(directory, "two.yaml"), importedNodeType("Two", "urn:example:two"))
		mainPath := filepath.Join(directory, "main.yaml")
		writeToscaFixture(t, mainPath, `tosca_definitions_version: tosca_simple_yaml_1_3
imports:
  - file: one.yaml
    namespace_prefix: ext
  - file: two.yaml
    namespace_prefix: ext
`)
		_, problems, err := testsupport.ParseFile(t, mainPath)
		if err == nil {
			t.Fatal("duplicate namespace_prefix was accepted")
		}
		assertDiagnostic(t, problems, "namespace_prefix", "ext", "duplicate")
	})

	t.Run("multiple-prefixes-one-uri", func(t *testing.T) {
		directory := t.TempDir()
		writeToscaFixture(t, filepath.Join(directory, "one.yaml"), importedNodeType("One", "urn:example:shared"))
		writeToscaFixture(t, filepath.Join(directory, "two.yaml"), importedNodeType("Two", "urn:example:shared"))
		mainPath := filepath.Join(directory, "main.yaml")
		writeToscaFixture(t, mainPath, `tosca_definitions_version: tosca_simple_yaml_1_3
imports:
  - file: one.yaml
    namespace_prefix: one
  - file: two.yaml
    namespace_prefix: two
node_types:
  LocalOne:
    derived_from: one:One
  LocalTwo:
    derived_from: two:Two
`)
		if _, problems, err := testsupport.ParseFile(t, mainPath); err != nil {
			t.Fatalf("distinct prefixes for one namespace failed: %v\n%s", err, problems)
		}
	})
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 3.1.3.1 Additional Requirements; 3.6.8 Import definition
// Requirements: TOSCA13-3.1.3.1-004, TOSCA13-3.1.3.1-005
// Expected: qualified imported parents and function targets resolve, while
// unknown namespaces and qualified names fail in the namespace phase
// Category: positive, negative, qualified name, inheritance, function target
func TestQualifiedImportedTypeResolution(t *testing.T) {
	directory := t.TempDir()
	writeToscaFixture(t, filepath.Join(directory, "imported.yaml"), `tosca_definitions_version: tosca_simple_yaml_1_3
namespace: urn:example:imported
node_types:
  Parent:
    derived_from: tosca.nodes.Root
    properties:
      inherited:
        type: string
        default: inherited
`)

	validPath := filepath.Join(directory, "valid.yaml")
	writeToscaFixture(t, validPath, `tosca_definitions_version: tosca_simple_yaml_1_3
imports:
  - file: imported.yaml
    namespace_prefix: ext
node_types:
  Child:
    derived_from: ext:Parent
topology_template:
  node_templates:
    node:
      type: Child
  outputs:
    result:
      value: { get_property: [node, inherited] }
`)
	if _, problems, err := testsupport.ParseFile(t, validPath); err != nil {
		t.Fatalf("qualified parent/function target failed: %v\n%s", err, problems)
	}

	for _, test := range []struct {
		name   string
		parent string
	}{
		{"missing-namespace", "absent:Parent"},
		{"unknown-qualified-name", "ext:Absent"},
	} {
		t.Run(test.name, func(t *testing.T) {
			path := filepath.Join(directory, test.name+".yaml")
			writeToscaFixture(t, path, fmt.Sprintf(`tosca_definitions_version: tosca_simple_yaml_1_3
imports:
  - file: imported.yaml
    namespace_prefix: ext
node_types:
  Child:
    derived_from: %s
`, test.parent))
			_, problems, err := testsupport.ParseFile(t, path)
			if err == nil {
				t.Fatalf("unknown qualified parent %q was accepted", test.parent)
			}
			assertDiagnostic(t, problems, "derived_from", test.parent, "unknown")
			assertNoUnrelatedImportFailure(t, problems)
		})
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 3.1.3 Rules to avoid namespace collisions; 3.6.8 Import definition
// Requirements: TOSCA13-3.1.3.1-004, TOSCA13-3.1.3.1-005
// Expected: transitive imports are visible through the parent import prefix
// and cyclic imports fail naming the responsible resource
// Category: positive, negative, transitive import, import graph, cycle
func TestImportGraphResolution(t *testing.T) {
	t.Run("transitive", func(t *testing.T) {
		directory := t.TempDir()
		writeToscaFixture(t, filepath.Join(directory, "base.yaml"), importedNodeType("Base", "urn:example:base"))
		writeToscaFixture(t, filepath.Join(directory, "middle.yaml"), `tosca_definitions_version: tosca_simple_yaml_1_3
imports:
  - base.yaml
node_types:
  Middle:
    derived_from: Base
`)
		mainPath := filepath.Join(directory, "main.yaml")
		writeToscaFixture(t, mainPath, `tosca_definitions_version: tosca_simple_yaml_1_3
imports:
  - file: middle.yaml
    namespace_prefix: graph
node_types:
  Local:
    derived_from: graph:Middle
  LocalBase:
    derived_from: graph:Base
`)
		if _, problems, err := testsupport.ParseFile(t, mainPath); err != nil {
			t.Fatalf("transitive import lookup failed: %v\n%s", err, problems)
		}
	})

	t.Run("cycle", func(t *testing.T) {
		directory := t.TempDir()
		aPath := filepath.Join(directory, "a.yaml")
		writeToscaFixture(t, aPath, "tosca_definitions_version: tosca_simple_yaml_1_3\nimports:\n  - b.yaml\n")
		writeToscaFixture(t, filepath.Join(directory, "b.yaml"), "tosca_definitions_version: tosca_simple_yaml_1_3\nimports:\n  - a.yaml\n")
		_, problems, err := testsupport.ParseFile(t, aPath)
		if err == nil {
			t.Fatal("cyclic import graph was accepted")
		}
		assertDiagnostic(t, problems, "endless loop", "a.yaml")
	})
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.1.3.1 Additional Requirements
// Requirement: TOSCA13-3.1.3.1-005
// Expected: an import whose declared definitions version selects a different
// grammar is rejected as an incompatible import before namespace lookup
// Category: negative, version isolation, import resolution
func TestImportDifferentDefinitionsVersionIsRejected(t *testing.T) {
	directory := t.TempDir()
	writeToscaFixture(t, filepath.Join(directory, "different-version.yaml"), `tosca_definitions_version: tosca_2_0
node_types:
  Foreign: {}
`)
	mainPath := filepath.Join(directory, "main.yaml")
	writeToscaFixture(t, mainPath, `tosca_definitions_version: tosca_simple_yaml_1_3
imports:
  - file: different-version.yaml
    namespace_prefix: foreign
`)

	_, problems, err := testsupport.ParseFile(t, mainPath)
	if err == nil {
		t.Fatal("TOSCA 1.3 accepted an import governed by TOSCA 2.0")
	}
	assertDiagnostic(t, problems, "incompatible import", "different-version.yaml")
	assertNoUnrelatedImportFailure(t, problems)
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 3.1.3 Rules to avoid namespace collisions; 3.6.8 Import definition
// Requirements: TOSCA13-3.1.3.1-004, TOSCA13-3.1.3.1-005
// Expected: importing the same document twice is deterministic and idempotent
// Category: positive, duplicate import, import graph
func TestDuplicateImportIsIdempotent(t *testing.T) {
	directory := t.TempDir()
	writeToscaFixture(t, filepath.Join(directory, "imported.yaml"), importedNodeType("Imported", "urn:example:imported"))
	mainPath := filepath.Join(directory, "main.yaml")
	writeToscaFixture(t, mainPath, `tosca_definitions_version: tosca_simple_yaml_1_3
imports:
  - imported.yaml
  - imported.yaml
node_types:
  Local:
    derived_from: Imported
`)
	if _, problems, err := testsupport.ParseFile(t, mainPath); err != nil {
		t.Fatalf("idempotent duplicate import failed: %v\n%s", err, problems)
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.1.3.1 Additional Requirements
// Requirement: TOSCA13-3.1.3.1-004
// Expected: distinct imported definitions with the same namespace URI, local
// name, and version cannot coexist when their definitions are not equivalent
// Category: negative, namespace identity, collision, resolution
func TestConflictingImportedDefinitionIdentity(t *testing.T) {
	directory := t.TempDir()
	writeToscaFixture(t, filepath.Join(directory, "one.yaml"), `tosca_definitions_version: tosca_simple_yaml_1_3
namespace: urn:example:shared
node_types:
  Shared:
    version: 1.0
    derived_from: tosca.nodes.Root
    properties:
      value:
        type: string
`)
	writeToscaFixture(t, filepath.Join(directory, "two.yaml"), `tosca_definitions_version: tosca_simple_yaml_1_3
namespace: urn:example:shared
node_types:
  Shared:
    version: 1.0
    derived_from: tosca.nodes.Root
    properties:
      value:
        type: integer
`)
	mainPath := filepath.Join(directory, "main.yaml")
	writeToscaFixture(t, mainPath, `tosca_definitions_version: tosca_simple_yaml_1_3
imports:
  - file: one.yaml
    namespace_prefix: one
  - file: two.yaml
    namespace_prefix: two
`)

	_, problems, err := testsupport.ParseFile(t, mainPath)
	if err == nil {
		t.Fatal("non-equivalent imported definitions with the same identity were accepted")
	}
	assertDiagnostic(t, problems, "equivalent", "urn:example:shared", "Shared", "one.yaml", "two.yaml")
	assertNoUnrelatedImportFailure(t, problems)
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.1.3.1 Additional Requirements
// Requirements: TOSCA13-3.1.3.1-004, TOSCA13-3.1.3.1-005
// Expected: local/imported collisions and their diagnostics are deterministic
// Category: negative, collision, deterministic resolution order
func TestImportCollisionDeterminism(t *testing.T) {
	directory := t.TempDir()
	writeToscaFixture(t, filepath.Join(directory, "one.yaml"), importedNodeType("Collision", "urn:example:one"))
	writeToscaFixture(t, filepath.Join(directory, "two.yaml"), importedNodeType("Collision", "urn:example:two"))
	mainPath := filepath.Join(directory, "main.yaml")
	writeToscaFixture(t, mainPath, `tosca_definitions_version: tosca_simple_yaml_1_3
imports:
  - one.yaml
  - two.yaml
node_types:
  Collision:
    derived_from: tosca.nodes.Root
`)

	var first string
	for iteration := 0; iteration < 10; iteration++ {
		_, problems, err := testsupport.ParseFile(t, mainPath)
		if err == nil {
			t.Fatal("local/imported name collision was accepted")
		}
		assertDiagnostic(t, problems, "ambiguous", "Collision")
		if iteration == 0 {
			first = problems
		} else if problems != first {
			t.Fatalf("non-deterministic collision diagnostic on iteration %d:\nfirst:\n%s\ncurrent:\n%s", iteration, first, problems)
		}
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.1.3.1 Additional Requirements
// Requirements: TOSCA13-3.1.3.1-007 through -013, -015 through -018,
// and -020 through -027 as listed in imports-namespaces.yaml
// Expected: duplicate local definition, template, and nested key names fail
// during YAML decoding or structural read with the duplicate name identified
// Category: negative, local-name collision, read phase
func TestLocalNameCollisions(t *testing.T) {
	tests := []struct {
		name   string
		source string
		key    string
	}{
		{"repositories", duplicateRootMap("repositories", "repo", "https://example.invalid/one", "https://example.invalid/two"), "repo"},
		{"data-types", duplicateRootMap("data_types", "Duplicate", "{ derived_from: string }", "{ derived_from: string }"), "Duplicate"},
		{"node-types", duplicateRootMap("node_types", "Duplicate", "{ derived_from: tosca.nodes.Root }", "{ derived_from: tosca.nodes.Root }"), "Duplicate"},
		{"relationship-types", duplicateRootMap("relationship_types", "Duplicate", "{ derived_from: tosca.relationships.Root }", "{ derived_from: tosca.relationships.Root }"), "Duplicate"},
		{"capability-types", duplicateRootMap("capability_types", "Duplicate", "{ derived_from: tosca.capabilities.Root }", "{ derived_from: tosca.capabilities.Root }"), "Duplicate"},
		{"artifact-types", duplicateRootMap("artifact_types", "Duplicate", "{ derived_from: tosca.artifacts.Root }", "{ derived_from: tosca.artifacts.Root }"), "Duplicate"},
		{"interface-types", duplicateRootMap("interface_types", "Duplicate", "{ derived_from: tosca.interfaces.Root }", "{ derived_from: tosca.interfaces.Root }"), "Duplicate"},
		{"node-templates", duplicateTopologyMap("node_templates", "duplicate", "{ type: tosca.nodes.Root }"), "duplicate"},
		{"relationship-templates", duplicateTopologyMap("relationship_templates", "duplicate", "{ type: tosca.relationships.Root }"), "duplicate"},
		{"inputs", duplicateTopologyMap("inputs", "duplicate", "{ type: string }"), "duplicate"},
		{"outputs", duplicateTopologyMap("outputs", "duplicate", "{ value: one }"), "duplicate"},
		{"properties", duplicateNodeTypeField("properties", "duplicate", "{ type: string }"), "duplicate"},
		{"attributes", duplicateNodeTypeField("attributes", "duplicate", "{ type: string }"), "duplicate"},
		{"artifacts", duplicateNodeTypeField("artifacts", "duplicate", "{ type: tosca.artifacts.File, file: one }"), "duplicate"},
		{"capabilities", duplicateNodeTypeField("capabilities", "duplicate", "{ type: tosca.capabilities.Root }"), "duplicate"},
		{"interfaces", duplicateNodeTypeField("interfaces", "duplicate", "{ type: tosca.interfaces.Root }"), "duplicate"},
		{
			"requirements",
			`tosca_definitions_version: tosca_simple_yaml_1_3
node_types:
  DuplicateContainer:
    derived_from: tosca.nodes.Root
    requirements:
      - duplicate:
          capability: tosca.capabilities.Node
      - duplicate:
          capability: tosca.capabilities.Node
`,
			"duplicate",
		},
		{
			"policies",
			`tosca_definitions_version: tosca_simple_yaml_1_3
topology_template:
  policies:
    - duplicate:
        type: tosca.policies.Root
    - duplicate:
        type: tosca.policies.Root
`,
			"duplicate",
		},
		{"groups", duplicateTopologyMap("groups", "duplicate", "{ type: tosca.groups.Root, members: [] }"), "duplicate"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, problems, err := testsupport.ParseSource(t, test.source)
			if err == nil {
				t.Fatalf("duplicate %s name was accepted", test.name)
			}
			assertDiagnostic(t, diagnosticText(problems, err), "duplicate map key: "+test.key)
		})
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.1.3.1 Additional Requirements
// Requirements: TOSCA13-3.1.3.1-001, TOSCA13-3.1.3.1-002
// Expected: a service template cannot contain competing default namespace
// declarations
// Category: negative, duplicate namespace declaration, YAML decoding
func TestDuplicateNamespaceDeclaration(t *testing.T) {
	source := `tosca_definitions_version: tosca_simple_yaml_1_3
namespace: urn:example:one
namespace: urn:example:two
`
	_, problems, err := testsupport.ParseSource(t, source)
	if err == nil {
		t.Fatal("duplicate namespace declarations were accepted")
	}
	assertDiagnostic(t, diagnosticText(problems, err), "duplicate", "namespace")
}

func writeToscaFixture(t *testing.T, path string, source string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatalf("create fixture directory: %v", err)
	}
	if err := os.WriteFile(path, []byte(source), 0o600); err != nil {
		t.Fatalf("write fixture %s: %v", path, err)
	}
}

func importedNodeType(name string, namespace string) string {
	namespaceLine := ""
	if namespace != "" {
		namespaceLine = "namespace: " + namespace + "\n"
	}
	return fmt.Sprintf(`tosca_definitions_version: tosca_simple_yaml_1_3
%snode_types:
  %s:
    derived_from: tosca.nodes.Root
`, namespaceLine, name)
}

func assertDiagnostic(t *testing.T, problems string, fragments ...string) {
	t.Helper()
	lower := strings.ToLower(problems)
	for _, fragment := range fragments {
		if !strings.Contains(lower, strings.ToLower(fragment)) {
			t.Fatalf("diagnostic does not contain %q:\n%s", fragment, problems)
		}
	}
}

func assertNoUnrelatedImportFailure(t *testing.T, problems string) {
	t.Helper()
	for _, unrelated := range []string{"unknown data type", "unsupported TOSCA definitions version"} {
		if strings.Contains(problems, unrelated) {
			t.Fatalf("unrelated earlier failure masked import/namespace validation:\n%s", problems)
		}
	}
}

func diagnosticText(problems string, err error) string {
	if err == nil {
		return problems
	}
	return problems + "\n" + err.Error()
}

func duplicateRootMap(section string, key string, first string, second string) string {
	return fmt.Sprintf(`tosca_definitions_version: tosca_simple_yaml_1_3
%s:
  %s: %s
  %s: %s
`, section, key, first, key, second)
}

func duplicateTopologyMap(section string, key string, value string) string {
	return fmt.Sprintf(`tosca_definitions_version: tosca_simple_yaml_1_3
topology_template:
  %s:
    %s: %s
    %s: %s
`, section, key, value, key, value)
}

func duplicateNodeTypeField(field string, key string, value string) string {
	return fmt.Sprintf(`tosca_definitions_version: tosca_simple_yaml_1_3
node_types:
  DuplicateContainer:
    derived_from: tosca.nodes.Root
    %s:
      %s: %s
      %s: %s
`, field, key, value, key, value)
}
