package tosca_1_3_corpus_test

import (
	"fmt"
	"sort"

	"github.com/tliron/exturl"
	"github.com/tliron/go-kutil/problems"
	"github.com/tliron/go-kutil/terminal"
	cloutjs "github.com/tliron/go-puccini/clout/js"
	v13 "github.com/tliron/go-puccini/tosca/grammars/tosca_v1_3"
)

// Project only the effective fields explicitly used by owned assertions.
// This reads resolved entities after inheritance/rendering, never raw YAML.
func effectiveEvidence(r corpusResult) (any, error) {
	if r.Parser == nil || r.Parser.Root == nil {
		return nil, fmt.Errorf("no effective root")
	}
	root, ok := r.Parser.Root.EntityPtr.(*v13.ServiceFile)
	if !ok {
		return nil, fmt.Errorf("not a 1.3 root")
	}
	nodes := map[string]any{}
	for _, n := range root.NodeTypes {
		caps := map[string]any{}
		props := map[string]any{}
		for name, c := range n.CapabilityDefinitions {
			sources := []string{}
			for _, s := range c.ValidSourceNodeTypes {
				sources = append(sources, s.Name)
			}
			sort.Strings(sources)
			caps[name] = map[string]any{"validSourceTypes": sources}
		}
		for name, p := range n.PropertyDefinitions {
			props[name] = map[string]any{"required": p.IsRequired()}
		}
		nodes[n.Name] = map[string]any{"capabilities": caps, "properties": props}
	}
	workflows := map[string]any{}
	if root.ServiceTemplate != nil {
		for name, w := range root.ServiceTemplate.WorkflowDefinitions {
			steps := map[string]any{}
			for sn, s := range w.StepDefinitions {
				if s.OperationHost != nil {
					steps[sn] = map[string]any{"operationHost": *s.OperationHost}
				}
			}
			workflows[name] = map[string]any{"steps": steps}
		}
	}
	return map[string]any{"effective": map[string]any{"nodeTypes": nodes, "workflows": workflows}}, nil
}
func evaluatedEvidence(r corpusResult) (any, error) {
	c, err := r.Template.Compile()
	if err != nil {
		return nil, err
	}
	u := exturl.NewContext()
	defer u.Release()
	ps := problems.NewProblems(terminal.NewStylist(false))
	e := cloutjs.ExecContext{Clout: c, Problems: ps, URLContext: u, Format: "yaml"}
	e.Coerce()
	if !ps.Empty() {
		return nil, fmt.Errorf("function evaluation: %s", ps.ToString(false))
	}
	return map[string]any{"evaluated": c.Properties["tosca"]}, nil
}
