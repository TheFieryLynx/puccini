package tosca_1_3_test

import (
	"strings"
	"testing"

	"github.com/tliron/go-puccini/tests/conformance/internal/testsupport"
)

// TOSCA 1.3 §3.7.11.4, "Additional Requirements".
// Category: positive, negative, boundary, inheritance, hierarchy, and
// cross-version regression.
// Expected: every explicitly declared members list contains node types from
// one hierarchy, including ancestors, descendants, and sibling node types.
func TestPartialGroupMemberHomogeneityAcceptsOneHierarchy(t *testing.T) {
	tests := []struct {
		name    string
		members string
	}{
		{name: "siblings", members: "example.Left, example.Right"},
		{name: "ancestor-and-descendant", members: "example.Branch, example.Leaf"},
		{name: "same-type", members: "example.Left, example.Left"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			source := groupMemberHomogeneityTemplate(test.members)
			if _, parseProblems, err := testsupport.ParseSource(t, source); err != nil {
				t.Fatalf("homogeneous group members failed: %v\n%s", err, parseProblems)
			}
		})
	}
}

func TestPartialGroupMemberHomogeneityRejectsDifferentHierarchies(t *testing.T) {
	problems := rejectGroupMemberHomogeneity(t, groupMemberHomogeneityTemplate(
		"example.Left, example.OtherRoot",
	))
	for _, fragment := range []string{
		`group_types["example.Group"].members[1]`,
		"example.OtherRoot",
		"same type hierarchy",
		"example.Left",
	} {
		if !strings.Contains(problems, fragment) {
			t.Fatalf("heterogeneous members diagnostic does not contain %q:\n%s", fragment, problems)
		}
	}
}

func TestPartialGroupMemberHomogeneityAcceptsSingleMember(t *testing.T) {
	if _, parseProblems, err := testsupport.ParseSource(t, groupMemberHomogeneityTemplate("example.Left")); err != nil {
		t.Fatalf("single group member failed: %v\n%s", err, parseProblems)
	}
}

func TestPartialGroupMemberHomogeneityInheritsMembers(t *testing.T) {
	source := groupMemberHomogeneityTemplate("example.Left, example.Right") + `
  example.DerivedGroup:
    derived_from: example.Group
`
	if _, parseProblems, err := testsupport.ParseSource(t, source); err != nil {
		t.Fatalf("inherited homogeneous members failed: %v\n%s", err, parseProblems)
	}
}

func TestPartialGroupMemberHomogeneityUnknownTypeStillFailsLookup(t *testing.T) {
	problems := rejectGroupMemberHomogeneity(t, groupMemberHomogeneityTemplate("example.Left, example.Missing"))
	if !strings.Contains(problems, "members") || !strings.Contains(problems, "unknown node type") {
		t.Fatalf("unknown member type failed for the wrong reason:\n%s", problems)
	}
}

func TestPartialGroupMemberHomogeneityDiagnosticIsDeterministic(t *testing.T) {
	source := groupMemberHomogeneityTemplate("example.Left, example.OtherRoot")
	first := rejectGroupMemberHomogeneity(t, source)
	second := rejectGroupMemberHomogeneity(t, source)
	if groupMemberDiagnosticBody(first) != groupMemberDiagnosticBody(second) {
		t.Fatalf("group member diagnostic is nondeterministic:\nfirst:\n%s\nsecond:\n%s", first, second)
	}
}

func TestPartialGroupMemberHomogeneityDoesNotChangeTOSCA20(t *testing.T) {
	source := `tosca_definitions_version: tosca_2_0
description: TOSCA 2.0 group member homogeneity regression
node_types:
  example.First: {}
  example.Second: {}
group_types:
  example.Group:
    members: [example.First, example.Second]
`
	if _, parseProblems, err := testsupport.ParseSource(t, source); err != nil {
		t.Fatalf("TOSCA 2.0 group member behavior changed: %v\n%s", err, parseProblems)
	}
}

func groupMemberHomogeneityTemplate(members string) string {
	return `tosca_definitions_version: tosca_simple_yaml_1_3
description: TOSCA 1.3 group member homogeneity conformance
node_types:
  example.Branch: {}
  example.Left:
    derived_from: example.Branch
  example.Right:
    derived_from: example.Branch
  example.Leaf:
    derived_from: example.Left
  example.OtherRoot: {}
group_types:
  example.Group:
    members: [` + members + `]
`
}

func rejectGroupMemberHomogeneity(t *testing.T, source string) string {
	t.Helper()
	_, parseProblems, err := testsupport.ParseSource(t, source)
	if err == nil {
		t.Fatalf("heterogeneous group member types were accepted")
	}
	return parseProblems
}

func groupMemberDiagnosticBody(problems string) string {
	if index := strings.Index(problems, "@"); index >= 0 {
		return problems[index:]
	}
	return problems
}
