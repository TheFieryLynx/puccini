package tosca_1_3_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/tliron/go-puccini/normal"
	"github.com/tliron/go-puccini/tests/conformance/internal/testsupport"
)

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 3.6.3.3, 3.6.10.5
// Requirements: TOSCA13-3.6.3.3-001, TOSCA13-3.6.10.5-004
// Expected: a bare scalar is equal and accepts the same value.
// Category: positive, property definition, assignment, direct
func TestBareConstraintMeansEqual(t *testing.T) {
	assertConstraintAccepted(t, scalarConstraintSource("integer", "7", "7"))
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.6.3.3, Additional Requirements
// Requirement: TOSCA13-3.6.3.3-001
// Expected: a different assignment violates bare equal.
// Category: negative, assignment, direct
func TestBareConstraintRejectsDifferentValue(t *testing.T) {
	assertConstraintRejected(t, scalarConstraintSource("integer", "7", "8"), "equal", "constraint")
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.6.3.3, Additional Requirements
// Requirement: TOSCA13-3.6.3.3-001
// Expected: bare and explicit equal have identical acceptance and normalization.
// Category: boundary, normalization, direct
func TestBareAndExplicitEqualAreEquivalent(t *testing.T) {
	bare := parseConstraintTemplate(t, scalarConstraintSource("integer", "7", "7"))
	explicit := parseConstraintTemplate(t, scalarConstraintSource("integer", "equal: 7", "7"))
	bareSignature := normalizedPropertyValidatorSignature(t, bare, "value")
	explicitSignature := normalizedPropertyValidatorSignature(t, explicit, "value")
	if bareSignature != explicitSignature {
		t.Fatalf("bare and explicit equal normalize differently:\n%s\n%s", bareSignature, explicitSignature)
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.6.3.3, Additional Requirements
// Requirement: TOSCA13-3.6.3.3-001
// Expected: an inherited bare-equal constraint remains effective.
// Category: inheritance, negative, direct
func TestInheritedBareEqualConstraint(t *testing.T) {
	source := `tosca_definitions_version: tosca_simple_yaml_1_3
node_types:
  ParentNode:
    derived_from: tosca.nodes.Root
    properties:
      value:
        type: integer
        constraints: [7]
  ChildNode:
    derived_from: ParentNode
topology_template:
  node_templates:
    node:
      type: ChildNode
      properties:
        value: 8
`
	assertConstraintRejected(t, source, "equal", "constraint")
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.6.3.3, Additional Requirements
// Requirement: TOSCA13-3.6.3.3-002
// Expected: length measures list and map entry counts.
// Category: positive, collection size, direct
func TestLengthConstraintUsesListAndMapSize(t *testing.T) {
	assertConstraintAccepted(t, collectionConstraintSource("list", "string", "length: 2", "[one, two]"))
	assertConstraintAccepted(t, collectionConstraintSource("map", "integer", "length: 2", "{one: 1, two: 2}"))
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.6.3.3, Additional Requirements
// Requirement: TOSCA13-3.6.3.3-002
// Expected: a collection with the wrong number of entries is rejected.
// Category: negative, collection size, direct
func TestLengthConstraintRejectsWrongCollectionSize(t *testing.T) {
	assertConstraintRejected(t, collectionConstraintSource("list", "string", "length: 2", "[one]"), "length", "constraint")
	assertConstraintRejected(t, collectionConstraintSource("map", "integer", "length: 2", "{one: 1}"), "length", "constraint")
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.6.3.3, Additional Requirements
// Requirement: TOSCA13-3.6.3.3-002
// Expected: zero is a valid exact size for an empty collection.
// Category: boundary, collection size, direct
func TestLengthConstraintEmptyCollection(t *testing.T) {
	assertConstraintAccepted(t, collectionConstraintSource("list", "string", "length: 0", "[]"))
	assertConstraintAccepted(t, collectionConstraintSource("map", "integer", "length: 0", "{}"))
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.6.3.3, Additional Requirements
// Requirement: TOSCA13-3.6.3.3-002
// Expected: the normalized validator applies to the collection, not each entry.
// Category: normalization, collection size, direct
func TestLengthConstraintNormalizesOnCollection(t *testing.T) {
	template := parseConstraintTemplate(t, collectionConstraintSource("list", "string", "length: 2", "[one, two]"))
	property, ok := template.NodeTemplates["node"].Properties["value"].(*normal.List)
	if !ok || property.ValueMeta == nil || len(property.ValueMeta.Validators) == 0 {
		t.Fatal("length validator is absent from collection metadata")
	}
	if property.ValueMeta.Element != nil && len(property.ValueMeta.Element.Validators) != 0 {
		t.Fatal("length validator was incorrectly attached to collection elements")
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.6.3.1, Operator keynames
// Requirement: TOSCA13-3.6.3.3-003
// Expected: compatible operators and operands are accepted.
// Category: positive, operand type compatibility, direct
func TestConstraintOperandTypeCompatibility(t *testing.T) {
	assertConstraintAccepted(t, scalarConstraintSource("integer", "greater_or_equal: 1", "2"))
	assertConstraintAccepted(t, scalarConstraintSource("string", "pattern: '[a-z]+'", "abc"))
	assertConstraintAccepted(t, collectionConstraintSource("list", "string", "min_length: 1", "[one]"))
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 3.6.3.1, 3.6.3.3
// Requirement: TOSCA13-3.6.3.3-003
// Expected: operands have the operator-mandated type even without an assignment.
// Category: negative, operand type compatibility, direct
func TestConstraintRejectsIncompatibleOperand(t *testing.T) {
	assertConstraintRejected(t, definitionOnlyConstraintSource("integer", "greater_than: wrong"), "greater_than", "integer")
	assertConstraintRejected(t, definitionOnlyConstraintSource("string", "length: wrong"), "length", "integer")
	assertConstraintRejected(t, definitionOnlyConstraintSource("string", "pattern: '['"), "pattern operand")
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 3.6.3.1, 3.6.3.3
// Requirement: TOSCA13-3.6.3.3-003
// Expected: length cannot constrain a non-string, non-list, non-map datatype.
// Category: negative, operator applicability, direct
func TestLengthConstraintRejectsNonCollectionType(t *testing.T) {
	assertConstraintRejected(t, definitionOnlyConstraintSource("integer", "length: 1"), "length", "integer")
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.6.10.5, Additional Requirements
// Requirement: TOSCA13-3.6.10.5-004
// Expected: property operators compatible with the effective type are accepted.
// Category: positive, property definition, direct
func TestPropertyConstraintCompatibleWithType(t *testing.T) {
	assertConstraintAccepted(t, definitionOnlyConstraintSource("integer", "in_range: [1, 10]"))
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.6.10.5, Additional Requirements
// Requirement: TOSCA13-3.6.10.5-004
// Expected: a property operator incompatible with its type is rejected.
// Category: negative, property definition, direct
func TestPropertyConstraintRejectsIncompatibleType(t *testing.T) {
	assertConstraintRejected(t, definitionOnlyConstraintSource("integer", "pattern: '[0-9]+'"), "pattern", "integer")
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 3.6.3, 3.6.10.5
// Requirement: TOSCA13-3.6.10.5-004
// Expected: a property default is checked against its constraints.
// Category: negative, default, direct
func TestPropertyDefaultConstraintEvaluated(t *testing.T) {
	source := `tosca_definitions_version: tosca_simple_yaml_1_3
node_types:
  TestNode:
    derived_from: tosca.nodes.Root
    properties:
      value:
        type: integer
        default: 0
        constraints:
          - greater_than: 0
`
	assertConstraintRejected(t, source, "greater_than", "constraint")
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 3.6.10.5, 3.6.10.6
// Requirement: TOSCA13-3.6.10.5-004
// Expected: a refinement is checked against its inherited effective type.
// Category: inheritance, negative, direct
func TestRefinedPropertyConstraintUsesEffectiveType(t *testing.T) {
	source := `tosca_definitions_version: tosca_simple_yaml_1_3
node_types:
  ParentNode:
    derived_from: tosca.nodes.Root
    properties:
      value:
        type: integer
        required: false
  ChildNode:
    derived_from: ParentNode
    properties:
      value:
        constraints:
          - pattern: '[0-9]+'
`
	assertConstraintRejected(t, source, "pattern", "integer")
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.6.14.3, Additional Requirements
// Requirement: TOSCA13-3.6.14.3-003
// Expected: compatible input and output parameter constraints are accepted.
// Category: positive, parameter definition, direct
func TestParameterConstraintCompatibleWithType(t *testing.T) {
	source := parameterConstraintSource("greater_or_equal: 1", "max_length: 4")
	assertConstraintAccepted(t, source)
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.6.14.3, Additional Requirements
// Requirement: TOSCA13-3.6.14.3-003
// Expected: incompatible input and output parameter constraints are rejected.
// Category: negative, parameter definition, direct
func TestParameterConstraintRejectsIncompatibleType(t *testing.T) {
	assertConstraintRejected(t, parameterConstraintSource("pattern: '[0-9]+'", "greater_than: 1"), "constraint", "incompatible")
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.6.14.3, Additional Requirements
// Requirement: TOSCA13-3.6.14.3-003
// Expected: a compatible parameter constraint remains in normalized metadata.
// Category: normalization, direct
func TestParameterConstraintMetadataPreserved(t *testing.T) {
	template := parseConstraintTemplate(t, parameterConstraintSource("greater_or_equal: 1", "max_length: 4"))
	input, ok := template.Inputs["count"].(*normal.Primitive)
	if !ok || input.ValueMeta == nil || len(input.ValueMeta.Validators) == 0 {
		t.Fatal("input constraint validator is absent from normalized metadata")
	}
	output, ok := template.Outputs["name"].(*normal.Primitive)
	if !ok || output.ValueMeta == nil || len(output.ValueMeta.Validators) == 0 {
		t.Fatal("output constraint validator is absent from normalized metadata")
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.7.6.3, Additional Requirements
// Requirement: TOSCA13-3.7.6.3-002
// Expected: a datatype constraint compatible with derived_from is accepted.
// Category: positive, hierarchy, direct
func TestDatatypeConstraintCompatibleWithParent(t *testing.T) {
	assertConstraintAccepted(t, datatypeConstraintSource("integer", "greater_or_equal: 1"))
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.7.6.3, Additional Requirements
// Requirement: TOSCA13-3.7.6.3-002
// Expected: an incompatible datatype constraint is rejected.
// Category: negative, hierarchy, direct
func TestDatatypeConstraintRejectsIncompatibleParent(t *testing.T) {
	assertConstraintRejected(t, datatypeConstraintSource("integer", "pattern: '[0-9]+'"), "pattern", "integer")
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.7.6.3, Additional Requirements
// Requirement: TOSCA13-3.7.6.3-002
// Expected: an inherited compatible datatype constraint remains effective.
// Category: inheritance, assignment, direct
func TestInheritedDatatypeConstraint(t *testing.T) {
	source := `tosca_definitions_version: tosca_simple_yaml_1_3
data_types:
  Positive:
    derived_from: integer
    constraints:
      - greater_than: 0
  MorePositive:
    derived_from: Positive
node_types:
  TestNode:
    derived_from: tosca.nodes.Root
    properties:
      value:
        type: MorePositive
topology_template:
  node_templates:
    node:
      type: TestNode
      properties:
        value: 0
`
	assertConstraintRejected(t, source, "greater_than", "constraint")
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.3.6.2, Additional requirements
// Requirement: TOSCA13-3.3.6.2-003
// Expected: values and bounds compare in canonical units.
// Category: positive, scalar-unit, direct
func TestScalarUnitConstraintConvertsUnits(t *testing.T) {
	assertConstraintAccepted(t, scalarUnitConstraintSource("1536 MiB", "1 GiB", "2 GiB"))
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.3.6.2, Additional requirements
// Requirement: TOSCA13-3.3.6.2-003
// Expected: the inclusive boundary is equal after unit conversion.
// Category: boundary, scalar-unit, direct
func TestScalarUnitConstraintInclusiveBoundary(t *testing.T) {
	assertConstraintAccepted(t, scalarUnitConstraintSource("1024 MiB", "1 GiB", "2 GiB"))
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.3.6.2, Additional requirements
// Requirement: TOSCA13-3.3.6.2-003
// Expected: canonical comparison rejects a value below the lower bound.
// Category: negative, scalar-unit, direct
func TestScalarUnitConstraintRejectsOutOfRangeValue(t *testing.T) {
	assertConstraintRejected(t, scalarUnitConstraintSource("1023 MiB", "1 GiB", "2 GiB"), "in_range", "constraint")
	assertConstraintRejected(t, scalarUnitConstraintSource("2049 MiB", "1 GiB", "2 GiB"), "in_range", "constraint")
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.3.6.2, Additional requirements
// Requirement: TOSCA13-3.3.6.2-003
// Expected: an inherited scalar-unit constraint evaluates an assignment.
// Category: inheritance, negative, direct
func TestInheritedScalarUnitConstraint(t *testing.T) {
	source := `tosca_definitions_version: tosca_simple_yaml_1_3
node_types:
  ParentNode:
    derived_from: tosca.nodes.Root
    properties:
      size:
        type: scalar-unit.size
        constraints:
          - greater_or_equal: 1 GiB
  ChildNode:
    derived_from: ParentNode
topology_template:
  node_templates:
    node:
      type: ChildNode
      properties:
        size: 1023 MiB
`
	assertConstraintRejected(t, source, "greater_or_equal", "constraint")
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 3.3.6.2, 3.6.3.3
// Requirement: TOSCA13-3.3.6.2-003
// Expected: normalized scalar-unit validator metadata is deterministic.
// Category: normalization, determinism, direct
func TestScalarUnitConstraintMetadataDeterministic(t *testing.T) {
	source := scalarUnitConstraintSource("1536 MiB", "1 GiB", "2 GiB")
	firstTemplate := parseConstraintTemplate(t, source)
	secondTemplate := parseConstraintTemplate(t, source)
	first := normalizedPropertyValidatorSignature(t, firstTemplate, "value")
	second := normalizedPropertyValidatorSignature(t, secondTemplate, "value")
	if first != second {
		t.Fatal("normalized scalar-unit constraint output is nondeterministic")
	}
}

// Specification: TOSCA Version 2.0
// Category: cross-version regression
// Expected: TOSCA 1.3 constraint policy does not activate for TOSCA 2.0.
func TestTosca20ConstraintBehaviorUnchanged(t *testing.T) {
	source := `tosca_definitions_version: tosca_2_0
node_types:
  TestNode:
    properties:
      count:
        type: integer
        validation:
          $greater_or_equal: [$value, 1]
service_template:
  node_templates:
    node:
      type: TestNode
      properties:
        count: 2
`
	assertConstraintAccepted(t, source)
}

func TestTosca20ValidationGrammarUnchanged(t *testing.T) {
	source := `tosca_definitions_version: tosca_2_0
data_types:
  Count:
    data_type: integer
    validation: 1
`
	assertConstraintRejected(t, source, "validation", "map")
}

func TestTosca20DatatypeConstraintBehaviorUnchanged(t *testing.T) {
	TestTosca20ConstraintBehaviorUnchanged(t)
}

func TestPropertiesWithoutConstraintsUnchanged(t *testing.T) {
	assertConstraintAccepted(t, definitionOnlyConstraintSource("integer", ""))
}

func TestUnconstrainedParametersUnchanged(t *testing.T) {
	source := `tosca_definitions_version: tosca_simple_yaml_1_3
topology_template:
  inputs:
    count:
      type: integer
      required: false
`
	assertConstraintAccepted(t, source)
}

func TestBareEqualConstraintNormalizesDeterministically(t *testing.T) {
	TestBareAndExplicitEqualAreEquivalent(t)
}

func scalarConstraintSource(typeName string, constraint string, assignment string) string {
	return `tosca_definitions_version: tosca_simple_yaml_1_3
node_types:
  TestNode:
    derived_from: tosca.nodes.Root
    properties:
      value:
        type: ` + typeName + `
        constraints:
          - ` + constraint + `
topology_template:
  node_templates:
    node:
      type: TestNode
      properties:
        value: ` + assignment + `
`
}

func definitionOnlyConstraintSource(typeName string, constraint string) string {
	constraintBlock := ""
	if constraint != "" {
		constraintBlock = "\n        constraints:\n          - " + constraint
	}
	return `tosca_definitions_version: tosca_simple_yaml_1_3
node_types:
  TestNode:
    derived_from: tosca.nodes.Root
    properties:
      value:
        type: ` + typeName + `
        required: false` + constraintBlock + `
`
}

func collectionConstraintSource(typeName string, entryType string, constraint string, assignment string) string {
	return `tosca_definitions_version: tosca_simple_yaml_1_3
node_types:
  TestNode:
    derived_from: tosca.nodes.Root
    properties:
      value:
        type: ` + typeName + `
        entry_schema:
          type: ` + entryType + `
        constraints:
          - ` + constraint + `
topology_template:
  node_templates:
    node:
      type: TestNode
      properties:
        value: ` + assignment + `
`
}

func parameterConstraintSource(inputConstraint string, outputConstraint string) string {
	return `tosca_definitions_version: tosca_simple_yaml_1_3
topology_template:
  inputs:
    count:
      type: integer
      required: false
      default: 2
      constraints:
        - ` + inputConstraint + `
  outputs:
    name:
      type: string
      value: test
      constraints:
        - ` + outputConstraint + `
`
}

func datatypeConstraintSource(parent string, constraint string) string {
	return `tosca_definitions_version: tosca_simple_yaml_1_3
data_types:
  Constrained:
    derived_from: ` + parent + `
    constraints:
      - ` + constraint + `
`
}

func scalarUnitConstraintSource(value string, lower string, upper string) string {
	return scalarConstraintSource("scalar-unit.size", "in_range: ["+lower+", "+upper+"]", value)
}

func parseConstraintTemplate(t *testing.T, source string) *normal.ServiceTemplate {
	t.Helper()
	template, problems, err := testsupport.ParseSource(t, source)
	if err != nil {
		t.Fatalf("constraint template was rejected: %v\n%s", err, problems)
	}
	return template
}

func normalizedPropertyValidatorSignature(t *testing.T, template *normal.ServiceTemplate, name string) string {
	t.Helper()
	property, ok := template.NodeTemplates["node"].Properties[name].(*normal.Primitive)
	if !ok || property.ValueMeta == nil || len(property.ValueMeta.Validators) == 0 {
		t.Fatalf("normalized property %q has no constraint validator", name)
	}
	validator := property.ValueMeta.Validators[0].FunctionCall
	signature, err := json.Marshal(struct {
		Name      string
		Arguments []any
	}{validator.Name, validator.Arguments})
	if err != nil {
		t.Fatal(err)
	}
	return string(signature)
}

func assertConstraintAccepted(t *testing.T, source string) {
	t.Helper()
	if _, problems, err := testsupport.ParseSource(t, source); err != nil {
		t.Fatalf("constraint template was rejected: %v\n%s", err, problems)
	}
}

func assertConstraintRejected(t *testing.T, source string, fragments ...string) {
	t.Helper()
	_, problems, err := testsupport.ParseSource(t, source)
	if err == nil {
		t.Fatal("invalid constraint template was accepted")
	}
	for _, fragment := range fragments {
		if !strings.Contains(problems, fragment) {
			t.Fatalf("constraint diagnostic does not contain %q:\n%s", fragment, problems)
		}
	}
}
