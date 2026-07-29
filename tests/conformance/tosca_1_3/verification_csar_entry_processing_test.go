package tosca_1_3_test

import (
	"strings"
	"testing"

	"github.com/tliron/go-puccini/normal"
	"github.com/tliron/go-puccini/tests/conformance/internal/testsupport"
)

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 6.1 Overall Structure of a CSAR; 6.2 TOSCA Meta File;
// 6.3 Archive without TOSCA-Metadata
// Requirements: TOSCA13-6.1-005, TOSCA13-6.1-008,
// TOSCA13-6.2-001, TOSCA13-6.3-002, TOSCA13-6.3-005
// Expected: the real CSAR processor selects the normative entry definition,
// accepts every permitted placement/shape, rejects invalid root definitions
// and ambiguous roots, and preserves the selected template's semantic result.
// Category: positive, negative, boundary, CSAR selection, grammar, direct
func TestVerificationCSAREntryProcessing(t *testing.T) {
	t.Run("TOSCA13-6.1-005", func(t *testing.T) {
		valid := writeCSAR(t, map[string]string{
			"service.yaml": rootCSAREntry("valid-root"),
		})
		serviceTemplate := parseVerifiedCSAR(t, valid)
		assertSelectedCSARNode(t, serviceTemplate, "valid-root")

		invalid := writeCSAR(t, map[string]string{
			"service.yaml": `tosca_definitions_version: tosca_simple_yaml_1_3
metadata:
  template_name: invalid-root
  template_version: 1.0
service_template: {}
`,
		})
		assertRejectedCSAR(t, invalid, "unsupported keyname", "service_template")
	})

	t.Run("TOSCA13-6.1-008", func(t *testing.T) {
		path := writeCSAR(t, map[string]string{
			"TOSCA-Metadata/TOSCA.meta":   validTOSCAMeta("1.1", "arbitrary/deep/service.yaml"),
			"arbitrary/deep/service.yaml": metadataCSAREntry("arbitrary-directory"),
		})
		serviceTemplate := parseVerifiedCSAR(t, path)
		assertSelectedCSARNode(t, serviceTemplate, "arbitrary-directory")
	})

	t.Run("TOSCA13-6.2-001", func(t *testing.T) {
		path := writeCSAR(t, map[string]string{
			"TOSCA-Metadata/TOSCA.meta": validTOSCAMeta("1.1", "entry.yaml"),
			"entry.yaml":                metadataCSAREntry("block-zero-only"),
		})
		serviceTemplate := parseVerifiedCSAR(t, path)
		assertSelectedCSARNode(t, serviceTemplate, "block-zero-only")

		missingEntry := writeCSAR(t, map[string]string{
			"TOSCA-Metadata/TOSCA.meta": `TOSCA-Meta-File-Version: 1.1
CSAR-Version: 1.1
Created-By: conformance-test
`,
			"entry.yaml": metadataCSAREntry("not-selected"),
		})
		assertRejectedCSAR(t, missingEntry,
			"invalid TOSCA 1.3 CSAR metadata",
			"Entry-Definitions",
		)
	})

	t.Run("TOSCA13-6.3-002", func(t *testing.T) {
		oneRoot := writeCSAR(t, map[string]string{
			"service.yaml":      rootCSAREntry("single-root"),
			"nested/other.yaml": metadataCSAREntry("nested-not-entry"),
		})
		serviceTemplate := parseVerifiedCSAR(t, oneRoot)
		assertSelectedCSARNode(t, serviceTemplate, "single-root")

		multipleRoots := writeCSAR(t, map[string]string{
			"first.yaml":  rootCSAREntry("first"),
			"second.yaml": rootCSAREntry("second"),
		})
		assertRejectedCSAR(t, multipleRoots,
			"more than one potential service template at the root",
		)
	})

	t.Run("TOSCA13-6.3-005", func(t *testing.T) {
		valid := writeCSAR(t, map[string]string{
			"service.yml": rootCSAREntry("valid-definitions"),
		})
		serviceTemplate := parseVerifiedCSAR(t, valid)
		assertSelectedCSARNode(t, serviceTemplate, "valid-definitions")

		missingMetadata := writeCSAR(t, map[string]string{
			"service.yaml": csarTosca13WithoutMetadata,
		})
		assertRejectedCSAR(t, missingMetadata,
			"metadata",
			"missing required keyname",
		)
	})
}

func parseVerifiedCSAR(t *testing.T, path string) *normal.ServiceTemplate {
	t.Helper()
	serviceTemplate, problems, err := testsupport.ParseFile(t, path)
	if err != nil {
		t.Fatalf("valid TOSCA 1.3 CSAR failed: %v\n%s", err, problems)
	}
	return serviceTemplate
}

func assertRejectedCSAR(t *testing.T, path string, expected ...string) {
	t.Helper()
	_, problems, err := testsupport.ParseFile(t, path)
	if err == nil {
		t.Fatal("invalid TOSCA 1.3 CSAR was accepted")
	}
	for _, fragment := range expected {
		if !strings.Contains(problems, fragment) {
			t.Fatalf("CSAR diagnostic does not contain %q:\n%s", fragment, problems)
		}
	}
}

func assertSelectedCSARNode(t *testing.T, serviceTemplate *normal.ServiceTemplate, name string) {
	t.Helper()
	if serviceTemplate == nil || serviceTemplate.NodeTemplates[name] == nil {
		t.Fatalf("selected CSAR entry nodes = %#v, want %q",
			serviceTemplate, name)
	}
}

func rootCSAREntry(nodeName string) string {
	return `tosca_definitions_version: tosca_simple_yaml_1_3
metadata:
  template_name: root-entry
  template_version: 1.0
topology_template:
  node_templates:
    ` + nodeName + `:
      type: tosca.nodes.Root
`
}

func metadataCSAREntry(nodeName string) string {
	return `tosca_definitions_version: tosca_simple_yaml_1_3
topology_template:
  node_templates:
    ` + nodeName + `:
      type: tosca.nodes.Root
`
}
