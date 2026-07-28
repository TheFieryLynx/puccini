package tosca_1_3_test

import (
	"archive/zip"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/tliron/go-puccini/tests/conformance/internal/testsupport"
)

const csarTosca13WithoutMetadata = `tosca_definitions_version: tosca_simple_yaml_1_3
topology_template: {}
`

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 6.1 Overall Structure of a CSAR; 6.3 Archive without TOSCA-Metadata
// Expected: accepted when the sole root definition has metadata, template_name, and template_version
// Category: positive, CSAR read
func TestCSARWithoutMetadataAcceptsRequiredRootMetadata(t *testing.T) {
	path := writeCSAR(t, map[string]string{
		"service.yaml": `tosca_definitions_version: tosca_simple_yaml_1_3
metadata:
  template_name: root-entry
  template_version: 1.0
topology_template: {}
`,
	})

	if _, parseProblems, err := testsupport.ParseFile(t, path); err != nil {
		t.Fatalf("valid metadata-less TOSCA 1.3 CSAR failed: %v\n%s", err, parseProblems)
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 6.1 Overall Structure of a CSAR; 6.3 Archive without TOSCA-Metadata
// Expected: rejected when root metadata or either required metadata key is absent
// Category: negative, CSAR read
func TestCSARWithoutMetadataRequiresRootMetadata(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		expected []string
	}{
		{
			name:     "metadata section",
			source:   csarTosca13WithoutMetadata,
			expected: []string{"metadata", "missing required keyname"},
		},
		{
			name: "template_name",
			source: `tosca_definitions_version: tosca_simple_yaml_1_3
metadata:
  template_version: 1.0
topology_template: {}
`,
			expected: []string{"metadata.template_name", "missing required keyname"},
		},
		{
			name: "template_version",
			source: `tosca_definitions_version: tosca_simple_yaml_1_3
metadata:
  template_name: root-entry
topology_template: {}
`,
			expected: []string{"metadata.template_version", "missing required keyname"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			path := writeCSAR(t, map[string]string{"service.yaml": test.source})
			_, parseProblems, err := testsupport.ParseFile(t, path)
			if err == nil {
				t.Fatal("invalid metadata-less TOSCA 1.3 CSAR was accepted")
			}
			for _, expected := range test.expected {
				if !strings.Contains(parseProblems, expected) {
					t.Fatalf("diagnostic does not contain %q:\n%s", expected, parseProblems)
				}
			}
		})
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 6.1 Overall Structure of a CSAR; 6.3 Archive without TOSCA-Metadata
// Expected: missing required metadata key diagnostics have stable order
// Category: boundary, determinism, CSAR read
func TestCSARWithoutMetadataReportsMissingNamesDeterministically(t *testing.T) {
	path := writeCSAR(t, map[string]string{
		"service.yaml": `tosca_definitions_version: tosca_simple_yaml_1_3
metadata: {}
topology_template: {}
`,
	})

	_, first, firstErr := testsupport.ParseFile(t, path)
	_, second, secondErr := testsupport.ParseFile(t, path)
	if firstErr == nil || secondErr == nil {
		t.Fatal("root metadata without required names was accepted")
	}
	if first != second {
		t.Fatalf("CSAR diagnostics are not deterministic:\nfirst:\n%s\nsecond:\n%s", first, second)
	}
	nameIndex := strings.Index(first, "metadata.template_name")
	versionIndex := strings.Index(first, "metadata.template_version")
	if nameIndex < 0 || versionIndex < 0 || nameIndex > versionIndex {
		t.Fatalf("required metadata diagnostics are absent or unstable:\n%s", first)
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 6.3 Archive without TOSCA-Metadata
// Expected: rejected when more than one root YAML file could be the entry definition
// Category: boundary, archive structure
func TestCSARWithoutMetadataRejectsMultipleRootDefinitions(t *testing.T) {
	path := writeCSAR(t, map[string]string{
		"first.yaml":  csarTosca13WithoutMetadata,
		"second.yaml": csarTosca13WithoutMetadata,
	})

	_, parseProblems, err := testsupport.ParseFile(t, path)
	if err == nil {
		t.Fatal("CSAR with multiple potential root definitions was accepted")
	}
	if !strings.Contains(parseProblems, "more than one potential service template at the root") {
		t.Fatalf("wrong multiple-root rejection reason:\n%s", parseProblems)
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 6.2 TOSCA Meta File
// Expected: version 1.1 block_0 with Entry-Definitions selects and accepts the entry
// Category: positive, CSAR metadata
func TestCSARWithMetadataAcceptsVersion11AndEntryDefinitions(t *testing.T) {
	path := writeCSAR(t, map[string]string{
		"TOSCA-Metadata/TOSCA.meta": validTOSCAMeta("1.1", "definitions/service.yaml"),
		"definitions/service.yaml":  csarTosca13WithoutMetadata,
	})

	if _, parseProblems, err := testsupport.ParseFile(t, path); err != nil {
		t.Fatalf("valid TOSCA.meta 1.1 CSAR failed: %v\n%s", err, parseProblems)
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 6.2 TOSCA Meta File
// Expected: block_0 without Entry-Definitions is rejected for the selected TOSCA 1.3 grammar
// Category: negative, CSAR metadata
func TestCSARMetaRequiresEntryDefinitions(t *testing.T) {
	path := writeCSAR(t, map[string]string{
		"TOSCA-Metadata/TOSCA.meta": `TOSCA-Meta-File-Version: 1.1
CSAR-Version: 1.1
Created-By: conformance-test
`,
		"service.yaml": csarTosca13WithoutMetadata,
	})

	_, parseProblems, err := testsupport.ParseFile(t, path)
	if err == nil {
		t.Fatal("TOSCA 1.3 CSAR without Entry-Definitions was accepted")
	}
	for _, expected := range []string{"invalid TOSCA 1.3 CSAR metadata", "Entry-Definitions"} {
		if !strings.Contains(parseProblems, expected) {
			t.Fatalf("diagnostic does not contain %q:\n%s", expected, parseProblems)
		}
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 6.2 TOSCA Meta File
// Expected: TOSCA-Meta-File-Version must denote 1.1
// Category: negative, version boundary
func TestCSARMetaRequiresVersion11(t *testing.T) {
	path := writeCSAR(t, map[string]string{
		"TOSCA-Metadata/TOSCA.meta": validTOSCAMeta("1.0", "service.yaml"),
		"service.yaml":              csarTosca13WithoutMetadata,
	})

	_, parseProblems, err := testsupport.ParseFile(t, path)
	if err == nil {
		t.Fatal("TOSCA 1.3 CSAR with TOSCA.meta version 1.0 was accepted")
	}
	for _, expected := range []string{"invalid TOSCA 1.3 CSAR metadata", "TOSCA-Meta-File-Version", "1.0", "1.1"} {
		if !strings.Contains(parseProblems, expected) {
			t.Fatalf("diagnostic does not contain %q:\n%s", expected, parseProblems)
		}
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 6.2 TOSCA Meta File
// Expected: Entry-Definitions remains a path within the CSAR root
// Category: negative, path traversal
func TestCSARMetaRejectsUnsafeEntryDefinitions(t *testing.T) {
	path := writeCSAR(t, map[string]string{
		"TOSCA-Metadata/TOSCA.meta": validTOSCAMeta("1.1", "../service.yaml"),
		"../service.yaml":           csarTosca13WithoutMetadata,
	})

	_, parseProblems, err := testsupport.ParseFile(t, path)
	if err == nil {
		t.Fatal("TOSCA 1.3 CSAR accepted a traversing Entry-Definitions path")
	}
	for _, expected := range []string{"invalid TOSCA 1.3 CSAR metadata", "Entry-Definitions", "relative to the CSAR root"} {
		if !strings.Contains(parseProblems, expected) {
			t.Fatalf("diagnostic does not contain %q:\n%s", expected, parseProblems)
		}
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 3.10.1.1 Metadata keynames; 6.1 Overall Structure of a CSAR
// Expected: metadata remains optional outside the archive-without-TOSCA-Metadata context
// Category: regression, scope
func TestOrdinaryServiceTemplateMetadataRemainsOptional(t *testing.T) {
	if _, parseProblems, err := testsupport.ParseSource(t, csarTosca13WithoutMetadata); err != nil {
		t.Fatalf("ordinary TOSCA 1.3 service template unexpectedly requires metadata: %v\n%s", err, parseProblems)
	}
}

// Specification: TOSCA 1.3 sections 6.1 and 6.2 only
// Expected: the TOSCA 1.3 CSAR policy does not alter TOSCA 2.0 grammar behavior
// Category: regression, version isolation
func TestCSARValidationDoesNotChangeTosca20(t *testing.T) {
	path := writeCSAR(t, map[string]string{
		"service.yaml": `tosca_definitions_version: tosca_2_0
service_template: {}
`,
	})

	if _, parseProblems, err := testsupport.ParseFile(t, path); err != nil {
		t.Fatalf("TOSCA 1.3 CSAR validation changed TOSCA 2.0 behavior: %v\n%s", err, parseProblems)
	}
}

func validTOSCAMeta(version string, entryDefinitions string) string {
	return "TOSCA-Meta-File-Version: " + version + "\n" +
		"CSAR-Version: 1.1\n" +
		"Created-By: conformance-test\n" +
		"Entry-Definitions: " + entryDefinitions + "\n"
}

func writeCSAR(t *testing.T, entries map[string]string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "service.csar")
	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("create CSAR fixture: %v", err)
	}

	writer := zip.NewWriter(file)
	names := make([]string, 0, len(entries))
	for name := range entries {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		entry, createErr := writer.Create(name)
		if createErr != nil {
			t.Fatalf("create CSAR entry %q: %v", name, createErr)
		}
		if _, writeErr := entry.Write([]byte(entries[name])); writeErr != nil {
			t.Fatalf("write CSAR entry %q: %v", name, writeErr)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close CSAR writer: %v", err)
	}
	if err := file.Close(); err != nil {
		t.Fatalf("close CSAR file: %v", err)
	}
	return path
}
