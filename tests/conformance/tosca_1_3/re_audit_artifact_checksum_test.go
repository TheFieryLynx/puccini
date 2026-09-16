package tosca_1_3_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TOSCA 1.3 §3.6.7.1 Keynames, §3.5.1 Required keynames, F06.
// Checksum requires an effective algorithm. Positive, negative, inheritance,
// short/long notation and normalized metadata assertions.
func TestReAuditArtifactChecksum(t *testing.T) {
	for _, tc := range []struct{ name, parent, child, assignment, phase, checksum, algorithm string }{
		{"pair", "{ type: tosca.artifacts.File, file: file.txt, checksum: abc, checksum_algorithm: SHA-256 }", "{}", "", "", "abc", "SHA-256"},
		{"missing-definition-algorithm", "{ type: tosca.artifacts.File, file: file.txt, checksum: abc }", "{}", "", "inheritance", "", ""},
		{"inherited-algorithm", "{ type: tosca.artifacts.File, file: file.txt, checksum_algorithm: SHA-256 }", "{ checksum: abc }", "", "", "abc", "SHA-256"},
		{"missing-refinement-algorithm", "{ type: tosca.artifacts.File, file: file.txt }", "{ checksum: abc }", "", "inheritance", "", ""},
		{"inherited-pair", "{ type: tosca.artifacts.File, file: file.txt, checksum: abc, checksum_algorithm: SHA-256 }", "{}", "{ file: other.txt }", "", "abc", "SHA-256"},
		{"assignment-inherits-algorithm", "{ type: tosca.artifacts.File, file: file.txt, checksum_algorithm: SHA-256 }", "{}", "{ checksum: abc }", "", "abc", "SHA-256"},
		{"assignment-missing-algorithm", "{ type: tosca.artifacts.File, file: file.txt }", "{}", "{ checksum: abc }", "inheritance", "", ""},
		{"algorithm-alone", "{ type: tosca.artifacts.File, file: file.txt, checksum_algorithm: SHA-256 }", "{}", "", "", "", "SHA-256"},
		{"neither", "{ type: tosca.artifacts.File, file: file.txt }", "{}", "", "", "", ""},
		{"short-refinement", "{ type: tosca.artifacts.File, file: file.txt, checksum: abc, checksum_algorithm: SHA-256 }", "other.txt", "", "", "abc", "SHA-256"},
		{"short-assignment", "{ type: tosca.artifacts.File, file: file.txt, checksum: abc, checksum_algorithm: SHA-256 }", "{}", "other.txt", "", "abc", "SHA-256"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			source := reAuditArtifactSource(t, tc.parent, tc.child, tc.assignment)
			if tc.phase != "" {
				reAuditReject(t, source, tc.phase, "checksum_algorithm", "a")
				return
			}
			r := reAuditParse(t, source)
			if r.Phase != "" {
				t.Fatalf("%s: %s", r.Phase, r.Problems)
			}
			a := r.Template.NodeTemplates["n"].Artifacts["a"]
			if a == nil || a.Checksum != tc.checksum || a.ChecksumAlgorithm != tc.algorithm {
				t.Fatalf("normalized artifact = %#v", a)
			}
		})
	}

	for _, tc := range []struct {
		name, fields string
		reject       bool
	}{
		{"standalone-pair", "checksum: abc, checksum_algorithm: SHA-256", false},
		{"standalone-missing", "checksum: abc", true},
		{"empty-checksum-still-requires-algorithm", "checksum: ''", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			source := reAuditArtifactSource(t, "{}", "{}", "{ type: tosca.artifacts.File, file: file.txt, "+tc.fields+" }")
			source = strings.ReplaceAll(source, "    artifacts: { a: {} }\n", "")
			if tc.reject {
				reAuditReject(t, source, "inheritance", "checksum_algorithm")
				return
			}
			r := reAuditParse(t, source)
			if r.Phase != "" {
				t.Fatalf("%s: %s", r.Phase, r.Problems)
			}
			if a := r.Template.NodeTemplates["n"].Artifacts["a"]; a.Checksum != "abc" || a.ChecksumAlgorithm != "SHA-256" {
				t.Fatalf("%#v", a)
			}
		})
	}
	for _, field := range []string{"checksum: []", "checksum_algorithm: 7", "checksum: abc, checksum_algorithm: null"} {
		t.Run(field, func(t *testing.T) {
			reAuditReject(t, reAuditArtifactSource(t, "{ type: tosca.artifacts.File, file: file.txt, "+field+" }", "{}", ""), "read", "checksum")
		})
	}
}
func reAuditArtifactSource(t *testing.T, parent, child, assignment string) string {
	// An explicit empty assignment instantiates the inherited artifact.
	if assignment == "" {
		assignment = "{}"
	}
	source := fmt.Sprintf(`tosca_definitions_version: tosca_simple_yaml_1_3
node_types:
  Parent:
    derived_from: tosca.nodes.Root
    artifacts: { a: %s }
  Child:
    derived_from: Parent
    artifacts: { a: %s }
topology_template:
  node_templates:
    n:
      type: Child
      artifacts: { a: %s }
`, parent, child, assignment)
	for _, name := range []string{"file.txt", "other.txt"} {
		path := filepath.Join(t.TempDir(), name)
		if err := os.WriteFile(path, []byte("fixture"), 0600); err != nil {
			t.Fatal(err)
		}
		source = strings.ReplaceAll(source, name, path)
	}
	return source
}
