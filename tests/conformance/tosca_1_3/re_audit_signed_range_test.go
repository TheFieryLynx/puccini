package tosca_1_3_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

// TOSCA 1.3 §3.3.3.1 Grammar, §3.3.3.2 Keywords, F10;
// §3.6.3 constraint in_range and §5.3.11 PortSpec domain.
// Positive/negative signed bounds, exact boundaries, UNBOUNDED, normalization,
// constrained and inherited data types. Invalid native values fail rendering.
func TestReAuditSignedRange(t *testing.T) {
	for _, tc := range []struct {
		value string
		lower int
		upper any
	}{
		{"[-2, -1]", -2, float64(-1)}, {"[0, 0]", 0, float64(0)},
		{"[-1, 0]", -1, float64(0)}, {"[1, 4]", 1, float64(4)},
		{"[-2, 4]", -2, float64(4)}, {"[-2, UNBOUNDED]", -2, nil},
	} {
		t.Run(tc.value, func(t *testing.T) {
			r := reAuditParse(t, reAuditSignedRangeSource(tc.value, ""))
			if r.Phase != "" {
				t.Fatalf("%s: %s", r.Phase, r.Problems)
			}
			raw, err := json.Marshal(r.Template.NodeTemplates["n"].Properties["p"])
			if err != nil {
				t.Fatal(err)
			}
			var p map[string]any
			if err = json.Unmarshal(raw, &p); err != nil {
				t.Fatal(err)
			}
			value, ok := p["$primitive"].(map[string]any)
			if !ok || value["lower"] != float64(tc.lower) || value["upper"] != tc.upper {
				t.Fatalf("signed normalized bounds: %s", raw)
			}
		})
	}
	for _, value := range []string{"[-1, -2]", "[1, -1]", "[0]", "[0, 1, 2]", "{lower: -2, upper: -1}", "['-2', -1]", "[-2, '-1']", "[UNBOUNDED, 2]", "[-2, unbounded]", "[-2, null]", "[-2.5, -1]"} {
		t.Run("invalid-"+value, func(t *testing.T) { reAuditReject(t, reAuditSignedRangeSource(value, ""), "rendering", "p") })
	}
	for _, tc := range []struct {
		value string
		valid bool
	}{
		{"[-2, -1]", true}, {"[-3, 0]", true}, {"[-4, -1]", false}, {"[-2, 1]", false}, {"[-2, UNBOUNDED]", false},
	} {
		t.Run("constraint-"+tc.value, func(t *testing.T) {
			source := reAuditSignedRangeSource(tc.value, ", constraints: [{ in_range: [-3, 0] }]")
			if !tc.valid {
				reAuditReject(t, source, "rendering", "in_range constraint not satisfied")
				return
			}
			r := reAuditParse(t, source)
			if r.Phase != "" {
				t.Fatalf("%s: %s", r.Phase, r.Problems)
			}
		})
	}

	t.Run("unbounded-constraint", func(t *testing.T) {
		r := reAuditParse(t, reAuditSignedRangeSource("[-2, UNBOUNDED]", ", constraints: [{ in_range: [-3, UNBOUNDED] }]"))
		if r.Phase != "" {
			t.Fatalf("%s: %s", r.Phase, r.Problems)
		}
	})
	t.Run("signed-integer-limits", func(t *testing.T) {
		r := reAuditParse(t, reAuditSignedRangeSource("[-9223372036854775808, 9223372036854775807]", ""))
		if r.Phase != "" {
			t.Fatalf("%s: %s", r.Phase, r.Problems)
		}
		raw, _ := json.Marshal(r.Template.NodeTemplates["n"].Properties["p"])
		if !strings.Contains(string(raw), "-9223372036854775808") || !strings.Contains(string(raw), "9223372036854775807") {
			t.Fatalf("integer limits changed: %s", raw)
		}
	})
	t.Run("port-domain-stays-positive", func(t *testing.T) {
		source := strings.Replace(reAuditSignedRangeSource("{target: 1, target_range: [1, 2]}", ""), "type: range", "type: tosca.datatypes.network.PortSpec", 1)
		r := reAuditParse(t, source)
		if r.Phase != "" {
			t.Fatalf("%s: %s", r.Phase, r.Problems)
		}
		reAuditReject(t, strings.Replace(source, "target_range: [1, 2]", "target_range: [-2, -1]", 1), "rendering", "in_range constraint not satisfied")
	})
	t.Run("inherited-range", func(t *testing.T) {
		source := reAuditSignedRangeSource("[-2, -1]", "")
		source = strings.Replace(source, "node_types:", "data_types:\n  SignedRange:\n    derived_from: range\n    constraints: [{ in_range: [-3, 0] }]\nnode_types:", 1)
		source = strings.Replace(source, "type: range", "type: SignedRange", 1)
		r := reAuditParse(t, source)
		if r.Phase != "" {
			t.Fatalf("%s: %s", r.Phase, r.Problems)
		}
	})
}
func reAuditSignedRangeSource(value, constraints string) string {
	return fmt.Sprintf(`tosca_definitions_version: tosca_simple_yaml_1_3
node_types:
  N:
    derived_from: tosca.nodes.Root
    properties: { p: { type: range%s } }
topology_template:
  node_templates:
    n:
      type: N
      properties: { p: %s }
`, constraints, value)
}
