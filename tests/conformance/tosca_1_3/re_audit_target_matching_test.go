package tosca_1_3_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tliron/exturl"
	problemspkg "github.com/tliron/go-kutil/problems"
	"github.com/tliron/go-kutil/terminal"
	cloutjs "github.com/tliron/go-puccini/clout/js"
	"github.com/tliron/go-puccini/normal"
)

// TOSCA 1.3 §§3.7.2/3.7.3/3.7.7/3.7.9/3.7.10 and 3.8.2: positive,
// negative, inheritance and normalization of concrete requirement targets.
// Titles: Capability definition; Requirement definition; Capability Type; Node
// Type; Relationship Type; Requirement assignment. Source: pinned 1.3 OS HTML.
// Negative matching cases must reach rendering and identify the entire edge.
const targetMatchingSource = `tosca_definitions_version: tosca_simple_yaml_1_3
capability_types:
  Base: { derived_from: tosca.capabilities.Root }
  Derived: { derived_from: Base }
  Sibling: { derived_from: tosca.capabilities.Root }
relationship_types:
  Rel: { derived_from: tosca.relationships.Root, valid_target_types: [Base] }
  DerivedRel: { derived_from: Rel }
node_types:
  TargetBase: { derived_from: tosca.nodes.Root }
  Target:
    derived_from: TargetBase
    properties: { size: { type: integer, default: 4 } }
    capabilities:
      good: { type: Base }
  TargetChild: { derived_from: Target }
  TargetGrandchild: { derived_from: TargetChild }
  Source:
    derived_from: tosca.nodes.Root
    requirements: [{ link: { capability: Base, node: TargetBase, relationship: Rel, occurrences: [1, 2] } }]
  SourceChild: { derived_from: Source }
topology_template:
  node_templates:
    source:
      type: Source
      requirements: [{ link: { node: target } }]
    target: { type: Target }
`

func targetRequirement(t *testing.T, source, capability string) (*normal.Requirement, reAuditResult) {
	t.Helper()
	r := reAuditParse(t, source)
	if r.Phase != "" {
		t.Fatalf("%s: %s", r.Phase, r.Problems)
	}
	reqs := r.Template.NodeTemplates["source"].Requirements
	if len(reqs) != 1 {
		t.Fatalf("requirements: %#v", reqs)
	}
	req := reqs[0]
	if req.NodeTemplate == nil || req.NodeTemplate.Name != "target" {
		t.Fatalf("wrong target: %#v", req)
	}
	if capability == "" {
		if req.CapabilityName != nil {
			t.Fatalf("ambiguous selection must remain deferred, got %q", *req.CapabilityName)
		}
	} else if req.CapabilityName == nil || *req.CapabilityName != capability {
		t.Fatalf("expected capability %q, got %#v", capability, req.CapabilityName)
	}
	if req.Occurrences == nil || req.Occurrences.Lower != 1 || req.Occurrences.Upper == nil || *req.Occurrences.Upper != 2 {
		t.Fatalf("occurrence bounds lost: %#v", req.Occurrences)
	}
	return req, r
}

func TestReAuditTargetMatching(t *testing.T) {
	cases := []struct{ name, from, to, selected string }{
		{"exact", "", "", "good"},
		{"derived", "good: { type: Base }", "good: { type: Derived }", "good"},
		{"multiple-one-match", "good: { type: Base }", "bad: { type: Sibling }\n      good: { type: Base }", "good"},
		{"explicit-name", "node: target }", "node: target, capability: good }", "good"},
		{"explicit-type", "node: target }", "node: target, capability: Base }", "good"},
		{"short", "{ node: target }", "target", "good"},
		{"inherited-capability", "target: { type: Target }", "target: { type: TargetChild }", "good"},
		{"multilevel", "target: { type: Target }", "target: { type: TargetGrandchild }", "good"},
		{"multiple-deferred", "good: { type: Base }", "a: { type: Derived }\n      z: { type: Base }", ""},
		{"relationship-inherited", "relationship: Rel", "relationship: DerivedRel", "good"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			source := targetMatchingSource
			if c.from != "" {
				source = strings.Replace(source, c.from, c.to, 1)
			}
			targetRequirement(t, source, c.selected)
		})
	}
	negatives := []struct{ name, from, to, reason, requested string }{
		{"no-match", "good: { type: Base }", "good: { type: Sibling }", "no compatible capability", "<implicit>"},
		{"missing-name", "node: target }", "node: target, capability: missing }", "explicit capability does not exist", "missing"},
		{"wrong-name-type", "node: target }", "node: target, capability: feature }", "no compatible capability", "feature"},
		{"wrong-node-type", "derived_from: TargetBase", "derived_from: tosca.nodes.Root", "node type", "<implicit>"},
		{"relationship-targets", "valid_target_types: [Base]", "valid_target_types: [Sibling]", "valid_target_types", "<implicit>"},
	}
	for _, c := range negatives {
		t.Run(c.name, func(t *testing.T) {
			source := strings.Replace(targetMatchingSource, c.from, c.to, 1)
			reAuditReject(t, source, "rendering", "[TOSCA13-TARGET-MATCH]", `source="source"`, `requirement="link"`, `target="target"`, `requested="`+c.requested+`"`, "expected=\"Base\"", "actual=", c.reason)
		})
	}
	t.Run("reverse-subtype", func(t *testing.T) {
		source := strings.Replace(targetMatchingSource, "capability: Base, node:", "capability: Derived, node:", 1)
		reAuditReject(t, source, "rendering", "[TOSCA13-TARGET-MATCH]", "expected=\"Derived\"", "good:Base")
	})
	t.Run("sibling-not-derived", func(t *testing.T) {
		source := strings.Replace(targetMatchingSource, "good: { type: Base }", "good: { type: Sibling }", 1)
		reAuditReject(t, source, "rendering", "[TOSCA13-TARGET-MATCH]", "good:Sibling")
	})
	t.Run("relationship-actual-derived", func(t *testing.T) {
		source := strings.Replace(targetMatchingSource, "good: { type: Base }", "good: { type: Derived }", 1)
		source = strings.Replace(source, "valid_target_types: [Base]", "valid_target_types: [Derived]", 1)
		targetRequirement(t, source, "good")
	})
	t.Run("explicit-incompatible-no-fallback", func(t *testing.T) {
		source := strings.Replace(targetMatchingSource, "good: { type: Base }", "bad: { type: Sibling }\n      good: { type: Base }", 1)
		source = strings.Replace(source, "node: target }", "node: target, capability: bad }", 1)
		reAuditReject(t, source, "rendering", "[TOSCA13-TARGET-MATCH]", `requested="bad"`, "bad:Sibling", "good:Base")
	})
	t.Run("explicit-type-narrows", func(t *testing.T) {
		source := strings.Replace(targetMatchingSource, "node: target }", "node: target, capability: Derived }", 1)
		reAuditReject(t, source, "rendering", "[TOSCA13-TARGET-MATCH]", `requested="Derived"`, "good:Base")
	})
	t.Run("relationship-template", func(t *testing.T) {
		source := strings.Replace(targetMatchingSource, "  node_templates:", "  relationship_templates:\n    edge: { type: DerivedRel }\n  node_templates:", 1)
		source = strings.Replace(source, "node: target }", "node: target, relationship: edge }", 1)
		targetRequirement(t, source, "good")
	})
}

// §3.8.2.2.3 forbids concrete node_filter, independently of the predicate.
// Abstract filter evaluation is orchestrator behavior (§3.6.5), not a parser
// rejection when a particular local candidate fails. Preserve its conjunction.
func TestReAuditTargetMatchingFilters(t *testing.T) {
	for _, filter := range []string{
		"{ properties: [{size: {equal: 4}}] }",
		"{ properties: [{size: {equal: 9}}] }",
		"{ capabilities: [{missing: {properties: [{size: {equal: 4}}]}}] }",
	} {
		t.Run(filter, func(t *testing.T) {
			source := strings.Replace(targetMatchingSource, "node: target }", "node: target, node_filter: "+filter+" }", 1)
			reAuditReject(t, source, "rendering", "node_filter", "Node Template is invalid", `node_templates["source"]`, "link")
		})
	}
	t.Run("inherited-value-still-invalid", func(t *testing.T) {
		source := strings.Replace(targetMatchingSource, "target: { type: Target }", "target: { type: TargetGrandchild }", 1)
		source = strings.Replace(source, "node: target }", "node: target, node_filter: {properties: [{size: {equal: 4}}]} }", 1)
		reAuditReject(t, source, "rendering", "Node Template is invalid")
	})
	t.Run("abstract-preserves-conjunction", func(t *testing.T) {
		source := strings.Replace(targetMatchingSource, "node: target }", "node: Target, node_filter: {properties: [{size: {equal: 4}}]} }", 1)
		r := reAuditParse(t, source)
		if r.Phase != "" {
			t.Fatalf("%s %s", r.Phase, r.Problems)
		}
		req := r.Template.NodeTemplates["source"].Requirements[0]
		if req.NodeTemplate != nil || req.CapabilityTypeName == nil || len(req.NodeTemplatePropertyValidation) == 0 {
			t.Fatalf("abstract restrictions lost: %#v", req)
		}
	})
}

func TestReAuditTargetMatchingInheritance(t *testing.T) {
	for _, actual := range []string{"Derived", "Sibling"} {
		t.Run(actual, func(t *testing.T) {
			source := strings.Replace(targetMatchingSource, "TargetChild: { derived_from: Target }", "TargetChild: { derived_from: Target, capabilities: {good: {type: "+actual+"}} }", 1)
			source = strings.Replace(source, "target: { type: Target }", "target: { type: TargetGrandchild }", 1)
			if actual == "Derived" {
				targetRequirement(t, source, "good")
			} else {
				reAuditReject(t, source, "inheritance", "must be derived from", "Base")
			}
		})
	}
	t.Run("requirement-narrowing", func(t *testing.T) {
		source := strings.Replace(targetMatchingSource, "SourceChild: { derived_from: Source }", "SourceChild: { derived_from: Source, requirements: [{link: {capability: Derived}}] }", 1)
		source = strings.Replace(source, "type: Source\n", "type: SourceChild\n", 1)
		reAuditReject(t, source, "rendering", "[TOSCA13-TARGET-MATCH]", `expected="Derived"`)
		targetRequirement(t, strings.Replace(source, "good: { type: Base }", "good: { type: Derived }", 1), "good")
	})
	t.Run("requirement-widening", func(t *testing.T) {
		source := strings.Replace(targetMatchingSource, "SourceChild: { derived_from: Source }", "SourceChild: { derived_from: Source, requirements: [{link: {capability: tosca.capabilities.Root}}] }", 1)
		reAuditReject(t, source, "inheritance", "must be derived from", "Base")
	})
}

func TestReAuditTargetMatchingOccurrences(t *testing.T) {
	targetRequirement(t, strings.Replace(targetMatchingSource, "node: target }", "node: target, occurrences: [1, 2] }", 1), "good")
	reAuditReject(t, strings.Replace(targetMatchingSource, "node: target }", "node: target, occurrences: [1, 3] }", 1), "inheritance", "occurrences", "fall completely within")
	reAuditReject(t, strings.Replace(targetMatchingSource, "requirements: [{ link: { node: target } }]", "requirements: [{link: target}, {link: target}, {link: target}]", 1), "rendering", "number of requirement", "link", "assignments")
	source := strings.Replace(targetMatchingSource, "      requirements: [{ link: { node: target } }]\n", "", 1)
	r := reAuditParse(t, source)
	if r.Phase != "" {
		t.Fatalf("%s %s", r.Phase, r.Problems)
	}
	req := r.Template.NodeTemplates["source"].Requirements[0]
	if req.NodeTemplate != nil || req.Occurrences == nil || req.Occurrences.Lower != 1 {
		t.Fatalf("missing abstract obligation: %#v", req)
	}
}

func TestReAuditTargetMatchingImports(t *testing.T) {
	for _, prefix := range []string{"", "remote:"} {
		t.Run(prefix, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "types.yaml")
			imported := `tosca_definitions_version: tosca_simple_yaml_1_3
capability_types:
  ImportedBase: { derived_from: tosca.capabilities.Root }
  ImportedDerived: { derived_from: ImportedBase }
node_types:
  ImportedTarget:
    derived_from: tosca.nodes.Root
    capabilities: {good: {type: ImportedDerived}}
`
			if err := os.WriteFile(path, []byte(imported), 0600); err != nil {
				t.Fatal(err)
			}
			imp := fmt.Sprintf("imports: [{file: %q", path)
			if prefix != "" {
				imp += ", namespace_prefix: remote"
			}
			imp += "}]\n"
			source := strings.Replace(targetMatchingSource, "capability_types:", imp+"capability_types:", 1)
			source = strings.Replace(source, "capability: Base, node: TargetBase", "capability: "+prefix+"ImportedBase, node: "+prefix+"ImportedTarget", 1)
			source = strings.Replace(source, "relationship: Rel, ", "", 1)
			source = strings.Replace(source, "target: { type: Target }", "target: { type: "+prefix+"ImportedTarget }", 1)
			targetRequirement(t, source, "good")
		})
	}
}

func TestReAuditTargetMatchingDeterminism(t *testing.T) {
	for _, caps := range []string{"bad: {type: Sibling}\n      good: {type: Base}", "a: {type: Base}\n      z: {type: Derived}"} {
		selected := "good"
		if strings.HasPrefix(caps, "a:") {
			selected = ""
		}
		source := strings.Replace(targetMatchingSource, "good: { type: Base }", caps, 1)
		for i := 0; i < 16; i++ {
			if i%2 == 1 {
				source = strings.Replace(source, "  node_templates:\n", "  node_templates:\n    target: { type: Target }\n", 1)
				source = strings.TrimSuffix(source, "    target: { type: Target }\n")
			}
			_, r := targetRequirement(t, source, selected)
			again, ok := r.Context.Normalize()
			if !ok {
				t.Fatal("repeat normalization failed")
			}
			req := again.NodeTemplates["source"].Requirements[0]
			if req.NodeTemplate.Name != "target" || (selected == "" && req.CapabilityName != nil) || (selected != "" && (req.CapabilityName == nil || *req.CapabilityName != selected)) {
				t.Fatalf("non-deterministic result: %#v", req)
			}
			source = strings.Replace(targetMatchingSource, "good: { type: Base }", caps, 1)
		}
	}
}

func TestReAuditTargetMatchingSelectorBoundaries(t *testing.T) {
	t.Run("case-sensitive", func(t *testing.T) {
		source := strings.Replace(targetMatchingSource, "node: target }", "node: target, capability: Good }", 1)
		reAuditReject(t, source, "rendering", "[TOSCA13-TARGET-MATCH]", `requested="Good"`, "explicit capability does not exist")
	})
	t.Run("name-type-collision", func(t *testing.T) {
		source := strings.Replace(targetMatchingSource, "good: { type: Base }", "Sibling: { type: Base }", 1)
		source = strings.Replace(source, "node: target }", "node: target, capability: Sibling }", 1)
		req, _ := targetRequirement(t, source, "Sibling")
		if req.CapabilityTypeName == nil || *req.CapabilityTypeName != "Base" {
			t.Fatalf("symbol was confused with type: %#v", req.CapabilityTypeName)
		}
	})
	t.Run("valid-source-type", func(t *testing.T) {
		source := strings.Replace(targetMatchingSource, "good: { type: Base }", "good: { type: Base, valid_source_types: [Source] }", 1)
		targetRequirement(t, source, "good")
		source = strings.Replace(source, "valid_source_types: [Source]", "valid_source_types: [Target]", 1)
		reAuditReject(t, source, "rendering", "[TOSCA13-TARGET-MATCH]", "no compatible capability")
	})
	t.Run("no-relationship-definition", func(t *testing.T) {
		source := strings.Replace(targetMatchingSource, "relationship: Rel, ", "", 1)
		targetRequirement(t, source, "good")
		reAuditReject(t, strings.Replace(source, "good: { type: Base }", "good: { type: Sibling }", 1), "rendering", "[TOSCA13-TARGET-MATCH]", "good:Sibling")
	})
	t.Run("incompatible-without-relationship", func(t *testing.T) {
		source := strings.Replace(targetMatchingSource, "relationship: Rel, ", "", 1)
		source = strings.Replace(source, "good: { type: Base }", "good: { type: Sibling }", 1)
		reAuditReject(t, source, "rendering", "[TOSCA13-TARGET-MATCH]", "good:Sibling")
	})
	t.Run("filter-and-type-invalid", func(t *testing.T) {
		source := strings.Replace(targetMatchingSource, "good: { type: Base }", "good: { type: Sibling }", 1)
		source = strings.Replace(source, "node: target }", "node: target, node_filter: {properties: [{size: {equal: 4}}]} }", 1)
		reAuditReject(t, source, "rendering", "[TOSCA13-TARGET-MATCH]", `source="source"`, `requirement="link"`, `target="target"`, "good:Sibling", "node_filter")
	})
}

// Preserve the admissible set when the declared type alone is less restrictive
// than relationship or source constraints. No arbitrary concrete binding.
func TestReAuditTargetMatchingCandidateNormalization(t *testing.T) {
	source := strings.Replace(targetMatchingSource, "good: { type: Base }", "a: {type: Derived}\n      z: {type: Derived}\n      excluded: {type: Base}", 1)
	source = strings.Replace(source, "valid_target_types: [Base]", "valid_target_types: [Derived]", 1)
	req, _ := targetRequirement(t, source, "")
	raw, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), `"capabilityNames":["a","z"]`) {
		t.Fatalf("admissible capability set lost: %v", req.CapabilityNames)
	}
}

func TestReAuditTargetMatchingCandidateResolution(t *testing.T) {
	source := strings.Replace(targetMatchingSource, "good: { type: Base }", "y: {type: Derived}\n      z: {type: Derived}\n      bad: {type: Base}", 1)
	source = strings.Replace(source, "valid_target_types: [Base]", "valid_target_types: [Derived]", 1)
	_, r := targetRequirement(t, source, "")
	graph, err := r.Template.Compile()
	if err != nil {
		t.Fatal(err)
	}
	u := exturl.NewContext()
	defer u.Release()
	problems := problemspkg.NewProblems(terminal.NewStylist(false))
	exec := cloutjs.ExecContext{Clout: graph, Problems: problems, URLContext: u, Format: "yaml"}
	exec.Resolve()
	if !problems.Empty() {
		t.Fatalf("resolution: %s", problems.ToString(false))
	}
	found := false
	for _, vertex := range graph.Vertexes {
		if vertex.Properties["name"] != "source" {
			continue
		}
		for _, edge := range vertex.EdgesOut {
			if edge.Properties["name"] != "link" {
				continue
			}
			found = true
			selected := edge.Properties["capability"]
			if selected != "y" && selected != "z" {
				t.Fatalf("resolver selected excluded capability: %v", selected)
			}
		}
	}
	if !found {
		t.Fatal("resolved relationship missing")
	}
}
