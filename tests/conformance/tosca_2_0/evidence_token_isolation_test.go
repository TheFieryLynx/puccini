package tosca_2_0_test

import (
	"github.com/tliron/go-puccini/normal"
	"github.com/tliron/go-puccini/tests/conformance/internal/testsupport"
	"testing"
)

// TOSCA 2.0 observable characterization, not normative proof. Expected: the
// existing literal-map interpretation is unchanged by the 1.3 Credential fix.
// Category: positive, normalization, version isolation.
func TestEvidenceTokenDataPolicyIsolation(t *testing.T) {
	st, p, err := testsupport.ParseSource(t, `tosca_definitions_version: tosca_2_0
service_template:
  outputs:
    result: {value: {token: secret}}
`)
	if err != nil {
		t.Fatalf("%v\n%s", err, p)
	}
	if _, ok := st.Outputs["result"].(*normal.Map); !ok {
		t.Fatalf("2.0 token behavior changed: %#v", st.Outputs["result"])
	}
}
