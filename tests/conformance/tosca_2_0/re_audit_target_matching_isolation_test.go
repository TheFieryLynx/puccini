package tosca_2_0_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/tliron/exturl"
	problemspkg "github.com/tliron/go-kutil/problems"
	"github.com/tliron/go-kutil/terminal"
	cloutjs "github.com/tliron/go-puccini/clout/js"
	"github.com/tliron/go-puccini/tests/conformance/internal/testsupport"
)

// Observable isolation regression for the shared rendering dispatch. This
// locks existing 2.0 behavior; it is not evidence of 2.0 normative conformance.
// No TOSCA 1.3 target-binding policy may be enabled by the shared hook.
func TestReAuditTargetMatchingPolicyIsolation(t *testing.T) {
	source := `tosca_definitions_version: tosca_2_0
capability_types:
  C: {}
node_types:
  Target:
    capabilities: {c: {type: C}}
  Source:
    requirements: [{link: {capability: C, node: Target}}]
service_template:
  node_templates:
    source:
      type: Source
      requirements: [{link: {node: target}}]
    target: {type: Target}
`
	st, p, e := testsupport.ParseSource(t, source)
	if e != nil {
		t.Fatalf("%v %s", e, p)
	}
	req := st.NodeTemplates["source"].Requirements[0]
	if req.NodeTemplate == nil || req.NodeTemplate.Name != "target" || req.CapabilityName != nil || req.Occurrences != nil || req.CapabilityNames != nil {
		t.Fatalf("1.3 policy leaked: %#v", req)
	}
	graph, err := st.Compile()
	if err != nil {
		t.Fatal(err)
	}
	u := exturl.NewContext()
	defer u.Release()
	problems := problemspkg.NewProblems(terminal.NewStylist(false))
	exec := cloutjs.ExecContext{Clout: graph, Problems: problems, URLContext: u, Format: "yaml"}
	exec.Resolve()
	if !problems.Empty() {
		t.Fatalf("2.0 resolution: %s", problems.ToString(false))
	}
	found := false
	for _, v := range graph.Vertexes {
		if v.Properties["name"] != "source" {
			continue
		}
		for _, e := range v.EdgesOut {
			if e.Properties["name"] == "link" {
				found = true
				if e.Properties["capability"] != "c" {
					t.Fatalf("2.0 binding changed: %#v", e.Properties)
				}
			}
		}
	}
	if !found {
		t.Fatal("2.0 relationship missing")
	}
	raw, err := json.Marshal(req)
	if err != nil || strings.Contains(string(raw), "capabilityNames") {
		t.Fatalf("1.3 candidate representation leaked: %s %v", raw, err)
	}
}
