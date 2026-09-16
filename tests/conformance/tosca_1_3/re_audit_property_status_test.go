package tosca_1_3_test

import (
	"fmt"
	"testing"

	v13 "github.com/tliron/go-puccini/tosca/grammars/tosca_v1_3"
)

// TOSCA 1.3 §3.6.10.2 Keynames and §3.6.10.3 Status values, F05.
// Positive/negative domain boundaries, read-phase diagnostics, effective defaults
// and inherited status. Status does not alter the assigned property value.
func TestReAuditPropertyStatus(t *testing.T) {
	for _, status := range []string{"supported", "unsupported", "experimental", "deprecated", ""} {
		t.Run("valid-"+status, func(t *testing.T) {
			field := ""
			if status != "" {
				field = ", status: " + status
			}
			r := reAuditParse(t, reAuditPropertySource("{ type: integer, default: 7"+field+" }", "{}", "{}", ""))
			reAuditPropertyValue(t, r, "7")
			expected := status
			if expected == "" {
				expected = "supported"
			}
			file := r.Context.Root.EntityPtr.(*v13.ServiceFile)
			for _, typ := range file.NodeTypes {
				p := typ.PropertyDefinitions["p"]
				if p == nil || p.Status == nil || *p.Status != expected {
					t.Fatalf("%s effective status = %#v, want %s", typ.Name, p, expected)
				}
			}
		})
	}
	for _, status := range []string{"unknown", "Supported", "SUPPORTED", "''", "null", "7", "[]", "{}"} {
		t.Run("invalid-"+status, func(t *testing.T) {
			source := reAuditPropertySource(fmt.Sprintf("{ type: integer, default: 7, status: %s }", status), "{}", "", "")
			reAuditReject(t, source, "read", "status")
		})
	}
	t.Run("invalid-refinement", func(t *testing.T) {
		reAuditReject(t, reAuditPropertySource("{ type: integer, default: 7 }", "{ status: unknown }", "", ""), "read", "status", "unknown")
	})
}
