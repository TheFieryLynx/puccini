package tosca_1_3_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

// TOSCA 1.3 §3.7.3.1 Keynames; §3.8.2.1 Keynames and §3.8.2.2.2
// Extended Notation (containment sentence adjoining §3.8.2.3), F09.
// Positive, negative, boundaries, definition inheritance/refinement and
// normalized lower/upper assertions. Range is not a TOSCA 2.0 count.
func TestReAuditRequirementOccurrences(t *testing.T) {
	for _, tc := range []struct {
		name, base, refinement, assignment string
		lower                              int
		upper                              any
	}{
		{"finite", "[0, 5]", "", "[1, 2]", 1, float64(2)},
		{"zero", "[0, 5]", "", "[0, 0]", 0, float64(0)},
		{"optional", "[0, 5]", "", "[0, 2]", 0, float64(2)},
		{"multiple", "[0, 5]", "", "[2, 4]", 2, float64(4)},
		{"unbounded", "[0, UNBOUNDED]", "", "[1, UNBOUNDED]", 1, nil},
		{"zero-fallback", "[0, 0]", "", "", 0, float64(0)},
		{"multiple-fallback", "[2, 4]", "", "", 2, float64(4)},
		{"definition-fallback", "[0, 5]", "", "", 0, float64(5)},
		{"inherited-refinement", "[0, 5]", "[1, 4]", "[2, 3]", 2, float64(3)},
		{"refinement-fallback", "[0, 5]", "[1, 4]", "", 1, float64(4)},
		{"default", "", "", "", 1, float64(1)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := reAuditParse(t, reAuditOccurrencesSource(tc.base, tc.refinement, tc.assignment))
			if r.Phase != "" {
				t.Fatalf("%s: %s", r.Phase, r.Problems)
			}
			reqs := r.Template.NodeTemplates["n"].Requirements
			if len(reqs) != 1 {
				t.Fatalf("want one declaration retaining multiplicity, got %d", len(reqs))
			}
			raw, err := json.Marshal(reqs[0])
			if err != nil {
				t.Fatal(err)
			}
			var m map[string]any
			if err = json.Unmarshal(raw, &m); err != nil {
				t.Fatal(err)
			}
			bounds, ok := m["occurrences"].(map[string]any)
			if !ok || bounds["lower"] != float64(tc.lower) || bounds["upper"] != tc.upper {
				t.Fatalf("normalized occurrence bounds: %s", raw)
			}
		})
	}

	t.Run("implicit-declaration", func(t *testing.T) {
		source := reAuditOccurrencesSource("[2, 4]", "", "")
		source = strings.Replace(source, "      requirements: [{ link: { node: Target } }]\n", "", 1)
		r := reAuditParse(t, source)
		if r.Phase != "" {
			t.Fatalf("%s: %s", r.Phase, r.Problems)
		}
		reqs := r.Template.NodeTemplates["n"].Requirements
		if len(reqs) != 1 || reqs[0].Occurrences.Lower != 2 || *reqs[0].Occurrences.Upper != 4 {
			t.Fatalf("implicit effective range: %#v", reqs)
		}
	})
	t.Run("separate-explicit-bindings", func(t *testing.T) {
		source := reAuditOccurrencesSource("[0, 5]", "", "[1, 2]")
		source = strings.Replace(source, "requirements: [{ link: { node: Target, occurrences: [1, 2] } }]", "requirements: [{ link: { node: Target, occurrences: [1, 2] } }, { link: { node: Target, occurrences: [1, 2] } }]", 1)
		r := reAuditParse(t, source)
		if r.Phase != "" {
			t.Fatalf("%s: %s", r.Phase, r.Problems)
		}
		reqs := r.Template.NodeTemplates["n"].Requirements
		if len(reqs) != 2 {
			t.Fatalf("explicit bindings merged: %d", len(reqs))
		}
		for _, req := range reqs {
			if req.Occurrences == nil || req.Occurrences.Lower != 1 || *req.Occurrences.Upper != 2 {
				t.Fatalf("%#v", req)
			}
		}
	})
	for _, r := range []string{"[0, 2]", "[1, 4]", "[2, UNBOUNDED]"} {
		t.Run("outside-"+r, func(t *testing.T) {
			reAuditReject(t, reAuditOccurrencesSource("[1, 3]", "", r), "inheritance", "occurrences", "within")
		})
	}
	for _, r := range []string{"[-1, 2]", "[2, 1]", "[1]", "[1, 2, 3]", "[1, invalid]"} {
		t.Run("malformed-"+r, func(t *testing.T) { reAuditReject(t, reAuditOccurrencesSource("[0, 5]", "", r), "read", "occurrences") })
	}
}
func reAuditOccurrencesSource(base, refinement, assignment string) string {
	b, r, a := "", "", ""
	if base != "" {
		b = ", occurrences: " + base
	}
	if refinement != "" {
		r = "\n    requirements: [{ link: { occurrences: " + refinement + " } }]"
	}
	if assignment != "" {
		a = ", occurrences: " + assignment
	}
	return fmt.Sprintf(`tosca_definitions_version: tosca_simple_yaml_1_3
node_types:
  Target:
    derived_from: tosca.nodes.Root
    capabilities: { endpoint: { type: tosca.capabilities.Node } }
  Parent:
    derived_from: tosca.nodes.Root
    requirements: [{ link: { capability: tosca.capabilities.Node, node: Target%s } }]
  Child:
    derived_from: Parent%s
topology_template:
  node_templates:
    n:
      type: Child
      requirements: [{ link: { node: Target%s } }]
`, b, r, a)
}
