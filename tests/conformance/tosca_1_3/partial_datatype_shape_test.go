package tosca_1_3_test

import (
	"strings"
	"testing"

	"github.com/tliron/go-puccini/tests/conformance/internal/testsupport"
)

// TOSCA 1.3 §3.7.6.3, "Additional Requirements".
// Category: positive, negative, boundary, hierarchy, and cross-version regression.
// Expected: a data type has a valid parent or a property, and an explicitly
// supplied properties map is never empty.
func TestPartialDataTypeShapeAcceptsEachRequiredAlternative(t *testing.T) {
	tests := []struct {
		name       string
		definition string
	}{
		{
			name: "derived-from",
			definition: `
    derived_from: string`,
		},
		{
			name: "property-definition",
			definition: `
    properties:
      value:
        type: string`,
		},
		{
			name: "derived-and-property",
			definition: `
    derived_from: tosca.datatypes.Root
    properties:
      value:
        type: string`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, parseProblems, err := testsupport.ParseSource(t, dataTypeShapeTemplate("Example", test.definition))
			if err != nil {
				t.Fatalf("valid data type definition failed: %v\n%s", err, parseProblems)
			}
		})
	}
}

func TestPartialDataTypeShapeRejectsMissingParentAndProperties(t *testing.T) {
	problems := rejectDataTypeShape(t, dataTypeShapeTemplate("Empty", `
    description: no parent and no property definitions`))
	for _, fragment := range []string{`data_types["Empty"]`, "derived_from", "property"} {
		if !strings.Contains(problems, fragment) {
			t.Fatalf("missing-alternative diagnostic does not contain %q:\n%s", fragment, problems)
		}
	}
}

func TestPartialDataTypeShapeRejectsExplicitEmptyProperties(t *testing.T) {
	tests := []struct {
		name       string
		definition string
	}{
		{
			name: "without-parent",
			definition: `
    properties: {}`,
		},
		{
			name: "with-parent",
			definition: `
    derived_from: tosca.datatypes.Root
    properties: {}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			problems := rejectDataTypeShape(t, dataTypeShapeTemplate("EmptyProperties", test.definition))
			for _, fragment := range []string{"properties", "one or more"} {
				if !strings.Contains(problems, fragment) {
					t.Fatalf("empty-properties diagnostic does not contain %q:\n%s", fragment, problems)
				}
			}
		})
	}
}

func TestPartialDataTypeShapeKeepsStructuralAndHierarchyValidation(t *testing.T) {
	t.Run("properties-wrong-type", func(t *testing.T) {
		problems := rejectDataTypeShape(t, dataTypeShapeTemplate("WrongProperties", `
    derived_from: tosca.datatypes.Root
    properties: []`))
		if !strings.Contains(problems, "properties") || !strings.Contains(problems, "map") {
			t.Fatalf("wrong properties type failed for the wrong reason:\n%s", problems)
		}
	})

	t.Run("unknown-parent", func(t *testing.T) {
		problems := rejectDataTypeShape(t, dataTypeShapeTemplate("UnknownParent", `
    derived_from: example.DoesNotExist`))
		if !strings.Contains(problems, "derived_from") || !strings.Contains(problems, "unknown data type") {
			t.Fatalf("unknown parent failed for the wrong reason:\n%s", problems)
		}
	})

	t.Run("invalid-property-definition", func(t *testing.T) {
		problems := rejectDataTypeShape(t, dataTypeShapeTemplate("InvalidProperty", `
    properties:
      value: {}`))
		if !strings.Contains(problems, "type") || !strings.Contains(problems, "missing") {
			t.Fatalf("invalid property definition failed for the wrong reason:\n%s", problems)
		}
	})
}

func TestPartialDataTypeShapeDiagnosticIsDeterministic(t *testing.T) {
	source := dataTypeShapeTemplate("Empty", `
    properties: {}`)
	first := rejectDataTypeShape(t, source)
	second := rejectDataTypeShape(t, source)
	if dataTypeShapeDiagnosticBody(first) != dataTypeShapeDiagnosticBody(second) {
		t.Fatalf("data type shape diagnostic is nondeterministic:\nfirst:\n%s\nsecond:\n%s", first, second)
	}
}

func TestPartialDataTypeShapeDoesNotChangeTOSCA20(t *testing.T) {
	source := `tosca_definitions_version: tosca_2_0
description: TOSCA 2.0 data type shape regression
data_types:
  Example: {}
`
	if _, parseProblems, err := testsupport.ParseSource(t, source); err != nil {
		t.Fatalf("TOSCA 2.0 data type behavior changed: %v\n%s", err, parseProblems)
	}
}

func dataTypeShapeTemplate(name, definition string) string {
	return `tosca_definitions_version: tosca_simple_yaml_1_3
description: TOSCA 1.3 data type shape conformance
data_types:
  ` + name + `:` + definition + `
`
}

func rejectDataTypeShape(t *testing.T, source string) string {
	t.Helper()
	_, parseProblems, err := testsupport.ParseSource(t, source)
	if err == nil {
		t.Fatalf("invalid data type definition was accepted")
	}
	return parseProblems
}

func dataTypeShapeDiagnosticBody(problems string) string {
	if index := strings.Index(problems, "@"); index >= 0 {
		return problems[index:]
	}
	return problems
}
