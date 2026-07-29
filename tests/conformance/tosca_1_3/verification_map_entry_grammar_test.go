package tosca_1_3_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/tliron/go-puccini/normal"
	"github.com/tliron/go-puccini/tests/conformance/internal/testsupport"
)

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.3.5.1.2 Multi-line grammar
// Requirement: TOSCA13-3.3.5.1.2-002
// Expected: every multi-line map entry has a key; keyed entries are retained
// and an unkeyed YAML sequence is rejected as a non-map at the assigned field.
// Category: positive, negative, boundary, rendering, normalization, direct
func TestVerificationMultilineMapEntryKey(t *testing.T) {
	t.Run("TOSCA13-3.3.5.1.2-002", func(t *testing.T) {
		serviceTemplate, problems, err := testsupport.ParseSource(t, mapVerificationTemplate(`
          first: one
          second: two`))
		if err != nil {
			t.Fatalf("valid keyed map failed: %v\n%s", err, problems)
		}
		value, ok := serviceTemplate.NodeTemplates["node"].Properties["payload"].(*normal.Map)
		if !ok {
			t.Fatalf("normalized payload has type %T", serviceTemplate.NodeTemplates["node"].Properties["payload"])
		}
		for key, expected := range map[string]string{"first": "one", "second": "two"} {
			assertNormalizedMapPrimitive(t, value, key, expected)
		}

		_, problems, err = testsupport.ParseSource(t, mapVerificationTemplate(`
          - one
          - two`))
		if err == nil {
			t.Fatal("unkeyed sequence was accepted as a TOSCA map")
		}
		for _, expected := range []string{`properties["payload"]`, "map"} {
			if !strings.Contains(problems, expected) {
				t.Fatalf("map diagnostic does not contain %q:\n%s", expected, problems)
			}
		}
	})
}

func mapVerificationTemplate(value string) string {
	return fmt.Sprintf(`tosca_definitions_version: tosca_simple_yaml_1_3
node_types:
  MapHolder:
    derived_from: tosca.nodes.Root
    properties:
      payload:
        type: map
        entry_schema:
          type: string
topology_template:
  node_templates:
    node:
      type: MapHolder
      properties:
        payload:%s
`, value)
}
