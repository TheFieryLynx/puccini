package tosca_2_0_test

import (
	"testing"

	"github.com/tliron/go-puccini/normal"
	"github.com/tliron/go-puccini/tests/conformance/internal/testsupport"
)

// Specification: TOSCA 2.0, 6.4.2 Type Derivation; 9.4 Property Definition.
// Expected: ordinary default refinement and explicit assignment keep their
// existing precedence; TOSCA 1.3 policy hooks do not run for this grammar.
// Category: positive, inheritance, normalization, cross-version regression.
func TestReAuditPropertyPolicyIsolation(t *testing.T) {
	source := `tosca_definitions_version: tosca_2_0
node_types:
  Parent:
    properties:
      p: { type: integer, default: 1 }
  Child:
    derived_from: Parent
    properties:
      p: { default: 3 }
service_template:
  node_templates:
    inherited:
      type: Child
    assigned:
      type: Child
      properties: { p: 9 }
`
	st, problems, err := testsupport.ParseSource(t, source)
	if err != nil {
		t.Fatalf("2.0 property regression: %v\n%s", err, problems)
	}
	for name, expected := range map[string]int{"inherited": 3, "assigned": 9} {
		p, ok := st.NodeTemplates[name].Properties["p"].(*normal.Primitive)
		if !ok || p.Primitive != expected {
			t.Fatalf("%s.p = %#v, want %d", name, p, expected)
		}
	}
}
