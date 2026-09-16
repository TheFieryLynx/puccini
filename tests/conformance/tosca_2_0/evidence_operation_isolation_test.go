package tosca_2_0_test

import (
	"github.com/tliron/go-puccini/tests/conformance/internal/testsupport"
	"testing"
)

// TOSCA 2.0 observable-regression characterization, not new conformance evidence.
// The 1.3 interpretation of §3.6.17.3 must not silently change the 2.0 grammar.
// Expected: existing parent and explicitly overridden implementation values;
// the existing local-description-only behavior is unchanged. Category: isolation.
func TestEvidenceOperationPolicyIsolation(t *testing.T) {
	st, p, err := testsupport.ParseSource(t, `tosca_definitions_version: tosca_2_0
interface_types:
  Test:
    operations:
      run: {}
node_types:
  Parent:
    interfaces:
      Test:
        type: Test
        operations:
          run: parent.sh
  Child:
    derived_from: Parent
    interfaces:
      Test:
        operations:
          run: {description: local}
service_template:
  node_templates:
    parent: {type: Parent}
    child: {type: Child}
`)
	if err != nil {
		t.Fatalf("%v\n%s", err, p)
	}
	for n, want := range map[string]string{"parent": "parent.sh", "child": ""} {
		if got := st.NodeTemplates[n].Interfaces["Test"].Operations["run"].Implementation; got != want {
			t.Fatalf("%s: %q != %q", n, got, want)
		}
	}
}
