package tosca_1_3_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/tliron/go-puccini/normal"
	v13 "github.com/tliron/go-puccini/tosca/grammars/tosca_v1_3"
)

// Specification: TOSCA 1.3, 3.6.10.6 Refining Property Definitions;
// 3.6.14.2 Grammar. F02. Expected: fixed values survive inheritance and
// normalization; changes to a final value are rejected at the earliest phase.
// Category: positive, negative, notation, inheritance, boundary, normalization.
func TestReAuditFixedPropertyRefinement(t *testing.T) {
	for _, notation := range []string{"{ value: 7 }", "7", "{ value: 7, default: 3 }"} {
		t.Run(notation, func(t *testing.T) {
			source := reAuditPropertySource("{ type: integer, default: 1 }", notation, "", "")
			r := reAuditParse(t, source)
			reAuditPropertyValue(t, r, "7")
			attribute, ok := r.Template.NodeTemplates["n"].Attributes["p"].(*normal.Primitive)
			if !ok || fmt.Sprint(attribute.Primitive) != "7" {
				t.Fatal("fixed value was not reflected as attribute p")
			}
		})
	}
	t.Run("multi-level-final-value", func(t *testing.T) {
		r := reAuditParse(t, reAuditPropertySource("{ type: integer, default: 1 }", "{ value: 7 }", "{ default: 8 }", ""))
		reAuditPropertyValue(t, r, "7")
	})
	t.Run("same-final-value", func(t *testing.T) {
		reAuditPropertyValue(t, reAuditParse(t, reAuditPropertySource("{ type: integer }", "{ value: 7 }", "{ value: 7 }", "p: 7")), "7")
	})
	t.Run("changed-refinement", func(t *testing.T) {
		reAuditReject(t, reAuditPropertySource("{ type: integer }", "{ value: 7 }", "{ value: 8 }", ""), "inheritance", "fixed", "p")
	})
	t.Run("changed-assignment", func(t *testing.T) {
		reAuditReject(t, reAuditPropertySource("{ type: integer }", "{ value: 7 }", "", "p: 8"), "rendering", "fixed", "p")
	})
	t.Run("wrong-fixed-type", func(t *testing.T) {
		reAuditReject(t, reAuditPropertySource("{ type: integer }", "{ value: wrong }", "", ""), "rendering", "integer")
	})
	t.Run("value-on-new-property", func(t *testing.T) {
		reAuditReject(t, reAuditPropertySource("{ type: integer, value: 7 }", "{}", "", ""), "inheritance", "value", "refinement")
	})
	t.Run("default-remains-overridable", func(t *testing.T) {
		reAuditPropertyValue(t, reAuditParse(t, reAuditPropertySource("{ type: integer, default: 1 }", "{ default: 3 }", "", "p: 9")), "9")
	})
	t.Run("invalid-unused-default", func(t *testing.T) {
		reAuditReject(t, reAuditPropertySource("{ type: integer }", "{ value: 7, default: wrong }", "", ""), "rendering", "default", "integer")
	})
	for _, notation := range []string{"{ value: [1, 2] }", "[1, 2]"} {
		t.Run("list-"+notation, func(t *testing.T) {
			r := reAuditParse(t, reAuditPropertySource("{ type: list, entry_schema: integer }", notation, "", "p: [1, 2]"))
			if r.Phase != "" {
				t.Fatalf("list fixed value failed in %s: %s", r.Phase, r.Problems)
			}
			list, ok := r.Template.NodeTemplates["n"].Properties["p"].(*normal.List)
			if !ok || len(list.Entries) != 2 {
				t.Fatalf("normalized list: %#v", list)
			}
			assertNormalizedPrimitive(t, list.Entries[0], 1)
			assertNormalizedPrimitive(t, list.Entries[1], 2)
		})
	}
	for _, notation := range []string{"{ value: { item: 7 } }", "{ item: 7 }"} {
		t.Run("map-"+notation, func(t *testing.T) {
			r := reAuditParse(t, reAuditPropertySource("{ type: map, entry_schema: integer }", notation, "", "p: { item: 7 }"))
			if r.Phase != "" {
				t.Fatalf("map fixed value failed in %s: %s", r.Phase, r.Problems)
			}
			m, ok := r.Template.NodeTemplates["n"].Properties["p"].(*normal.Map)
			if !ok || len(m.Entries) != 1 {
				t.Fatalf("normalized map: %#v", m)
			}
			assertNormalizedPrimitive(t, m.Entries[0], 7)
		})
	}
	t.Run("changed-list", func(t *testing.T) {
		reAuditReject(t, reAuditPropertySource("{ type: list, entry_schema: integer }", "{ value: [1, 2] }", "", "p: [1, 3]"), "rendering", "fixed")
	})
	t.Run("fixed-expression", func(t *testing.T) {
		source := reAuditPropertySource("{ type: integer }", "{ get_input: fixed }", "", "")
		source = strings.Replace(source, "topology_template:\n", "topology_template:\n  inputs:\n    fixed: { type: integer, default: 7 }\n", 1)
		r := reAuditParse(t, source)
		if r.Phase != "" {
			t.Fatalf("expression failed in %s: %s", r.Phase, r.Problems)
		}
		fn, ok := r.Template.NodeTemplates["n"].Properties["p"].(*normal.FunctionCall)
		if !ok || !strings.HasSuffix(fn.FunctionCall.Name, "get_input") {
			t.Fatalf("normalized expression: %#v", fn)
		}
	})
}

// Specification: TOSCA 1.3, 3.6.10.6 Refining Property Definitions. F03.
// Expected: the intersection of all ancestor and child constraints applies;
// only compatible type narrowing and optional-to-required refinement are valid.
// Category: positive, negative, nearest sibling, multi-level, boundary.
func TestReAuditInheritedPropertyConstraints(t *testing.T) {
	for _, value := range []int{5, 10, 11, 12, 15, 20, 21, 25} {
		t.Run(fmt.Sprint(value), func(t *testing.T) {
			source := reAuditPropertySource("{ type: integer, constraints: [ { greater_or_equal: 10 } ] }", "{ constraints: [ { less_or_equal: 20 } ] }", "{ constraints: [ { greater_or_equal: 12 } ] }", fmt.Sprintf("p: %d", value))
			if value >= 12 && value <= 20 {
				r := reAuditParse(t, source)
				reAuditPropertyValue(t, r, fmt.Sprint(value))
				encoded, err := json.Marshal(r.Template.NodeTemplates["n"].Properties["p"])
				if err != nil {
					t.Fatal(err)
				}
				for _, operator := range []string{"greater_or_equal", "less_or_equal"} {
					if !strings.Contains(string(encoded), operator) {
						t.Fatalf("normalized inherited validator %s lost: %s", operator, encoded)
					}
				}
			} else {
				reAuditReject(t, source, "rendering", "constraint not satisfied", "p")
			}
		})
	}
	t.Run("no-grandchild-constraint", func(t *testing.T) {
		reAuditReject(t, reAuditPropertySource("{ type: integer, constraints: [ { greater_or_equal: 10 } ] }", "{ constraints: [ { less_or_equal: 20 } ] }", "", "p: 5"), "rendering", "greater_or_equal constraint not satisfied")
	})
	t.Run("required-cannot-become-optional", func(t *testing.T) {
		reAuditReject(t, reAuditPropertySource("{ type: integer, required: true }", "{ required: false }", "", "p: 7"), "inheritance", "cannot refine")
	})
	t.Run("incompatible-type", func(t *testing.T) {
		reAuditReject(t, reAuditPropertySource("{ type: integer }", "{ type: string }", "", "p: wrong"), "inheritance", "incompatible")
	})
	t.Run("type-narrowing", func(t *testing.T) {
		source := reAuditPropertySource("{ type: integer }", "{ type: Positive }", "", "p: 7")
		source = strings.Replace(source, "node_types:\n", "data_types:\n  Positive:\n    derived_from: integer\n    constraints: [{ greater_than: 0 }]\nnode_types:\n", 1)
		r := reAuditParse(t, source)
		reAuditPropertyValue(t, r, "7")
		file := r.Context.Root.EntityPtr.(*v13.ServiceFile)
		for _, typ := range file.NodeTypes {
			if typ.Name == "Child" && typ.PropertyDefinitions["p"].DataType.Name != "Positive" {
				t.Fatal("effective narrowed datatype lost")
			}
		}
		widen := reAuditPropertySource("{ type: Positive }", "{ type: integer }", "", "p: 7")
		widen = strings.Replace(widen, "node_types:\n", "data_types:\n  Positive:\n    derived_from: integer\n    constraints: [{ greater_than: 0 }]\nnode_types:\n", 1)
		reAuditReject(t, widen, "inheritance", "must be derived from", "Positive")
	})
}

func reAuditPropertySource(parent, child, grandchild, assignment string) string {
	last := "Child"
	grand := ""
	if grandchild != "" {
		last = "Grandchild"
		grand = "  Grandchild:\n    derived_from: Child\n    properties:\n      p: " + grandchild + "\n"
	}
	return fmt.Sprintf(`tosca_definitions_version: tosca_simple_yaml_1_3
node_types:
  Parent:
    derived_from: tosca.nodes.Root
    properties:
      p: %s
  Child:
    derived_from: Parent
    properties:
      p: %s
%stopology_template:
  node_templates:
    n:
      type: %s
      properties: { %s }
`, parent, child, grand, last, assignment)
}

func reAuditPropertyValue(t *testing.T, r reAuditResult, expected string) {
	t.Helper()
	if r.Phase != "" || r.Template == nil {
		t.Fatalf("expected accepted template, failed in %s:\n%s", r.Phase, r.Problems)
	}
	value, ok := r.Template.NodeTemplates["n"].Properties["p"].(*normal.Primitive)
	if !ok || fmt.Sprint(value.Primitive) != expected {
		t.Fatalf("normalized p = %#v, want %s", value, expected)
	}
}
