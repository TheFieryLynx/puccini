package tosca_2_0_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/tliron/go-puccini/tests/conformance/internal/testsupport"
)

// TOSCA 2.0 §8.5 Requirement Assignment, §8.5.1 Keynames: exact count.
// Positive counts 0/1/2 retain their existing expansion; normalization must not
// acquire a TOSCA 1.3 occurrence interval. Negative 1.3 grammar is rejected.
func TestReAuditOccurrencePolicyIsolation(t *testing.T) {
	source := `tosca_definitions_version: tosca_2_0
capability_types:
  C: {}
node_types:
  Target:
    capabilities: { c: { type: C } }
  Source:
    requirements: [{ link: { capability: C, node: Target, count_range: [0, 5] } }]
service_template:
  node_templates:
    n:
      type: Source
      requirements: [{ link: { node: Target, count: %d } }]
`
	for _, count := range []int{0, 1, 2} {
		t.Run(fmt.Sprint(count), func(t *testing.T) {
			st, p, e := testsupport.ParseSource(t, fmt.Sprintf(source, count))
			if e != nil {
				t.Fatalf("%v %s", e, p)
			}
			reqs := st.NodeTemplates["n"].Requirements
			if len(reqs) != count {
				t.Fatalf("count=%d normalized=%d", count, len(reqs))
			}
			raw, _ := json.Marshal(reqs)
			if strings.Contains(string(raw), "occurrences") {
				t.Fatalf("1.3 representation leaked: %s", raw)
			}
		})
	}
	_, p, e := testsupport.ParseSource(t, strings.Replace(fmt.Sprintf(source, 1), "count: 1", "occurrences: [1, 2]", 1))
	if e == nil || !strings.Contains(p, "occurrences") {
		t.Fatalf("1.3 grammar leaked: %v %s", e, p)
	}
}
