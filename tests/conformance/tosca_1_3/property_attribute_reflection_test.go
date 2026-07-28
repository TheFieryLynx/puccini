package tosca_1_3_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/tliron/exturl"
	"github.com/tliron/go-ard"
	"github.com/tliron/go-kutil/problems"
	"github.com/tliron/go-kutil/terminal"
	cloutjs "github.com/tliron/go-puccini/clout/js"
	"github.com/tliron/go-puccini/normal"
	"github.com/tliron/go-puccini/tests/conformance/internal/testsupport"
	"github.com/tliron/go-puccini/tosca/grammars/tosca_v1_3"
	"github.com/tliron/go-puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/go-puccini/tosca/parser"
	"github.com/tliron/go-puccini/tosca/parsing"
)

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 3.6.10.1, 3.6.10.5, 3.6.12.1, 3.6.12.4
// Requirement: TOSCA13-3.6.10.5-001
// Expected: node, relationship, and capability properties create independent,
// same-named effective attribute definitions with the effective property type.
// Category: positive, definition, scope, type, mutation-safety
func TestPropertyAttributeReflectionEffectiveDefinitions(t *testing.T) {
	source := `tosca_definitions_version: tosca_simple_yaml_1_3
node_types:
  ReflectionNode:
    derived_from: tosca.nodes.Root
    properties:
      first:
        type: string
        description: desired first value
        required: false
        constraints:
          - min_length: 1
      second:
        type: integer
        required: false
    attributes:
      explicit_state:
        type: string
    capabilities:
      endpoint:
        type: ReflectionCapability
        properties:
          definition_property:
            type: string
            required: false
relationship_types:
  ReflectionRelationship:
    derived_from: tosca.relationships.Root
    properties:
      relationship_property:
        type: string
        required: false
capability_types:
  ReflectionCapability:
    derived_from: tosca.capabilities.Root
    properties:
      capability_property:
        type: string
        required: false
group_types:
  UnrelatedGroup:
    derived_from: tosca.groups.Root
    properties:
      project_data:
        type: string
        required: false
policy_types:
  UnrelatedPolicy:
    derived_from: tosca.policies.Root
    properties:
      project_data:
        type: string
        required: false
topology_template:
  node_templates:
    node:
      type: ReflectionNode
  relationship_templates:
    relationship:
      type: ReflectionRelationship
  groups:
    unrelated:
      type: UnrelatedGroup
      members: [node]
  policies:
    - unrelated:
        type: UnrelatedPolicy
`
	parserContext, serviceTemplate, parseProblems, release, err := parseReflectionSource(t, source)
	defer release()
	if err != nil {
		t.Fatalf("reflection template failed: %v\n%s", err, parseProblems)
	}

	root := reflectionRoot(t, parserContext)
	nodeType := findNodeType(t, root.NodeTypes, "ReflectionNode")
	firstProperty := nodeType.PropertyDefinitions["first"]
	firstAttribute := nodeType.AttributeDefinitions["first"]
	if firstAttribute == nil {
		t.Fatal("property first did not create an effective reflected attribute")
	}
	if firstAttribute == firstProperty.AttributeDefinition {
		t.Fatal("reflected attribute aliases the mutable property definition")
	}
	if firstAttribute.DataType != firstProperty.DataType || firstAttribute.DataTypeName == nil ||
		firstProperty.DataTypeName == nil || *firstAttribute.DataTypeName != *firstProperty.DataTypeName {
		t.Fatalf("reflected type = %v/%v, property type = %v/%v",
			firstAttribute.DataTypeName, firstAttribute.DataType,
			firstProperty.DataTypeName, firstProperty.DataType)
	}
	if firstAttribute.Description == nil || *firstAttribute.Description != "desired first value" {
		t.Fatalf("reflected description = %v", firstAttribute.Description)
	}
	if firstAttribute.Default != nil {
		t.Fatal("property default was copied into the reflected attribute definition")
	}
	if firstAttribute.ValidationClause != nil {
		t.Fatal("property constraints were copied into the reflected attribute definition")
	}
	for _, name := range []string{"first", "second", "explicit_state"} {
		if nodeType.AttributeDefinitions[name] == nil {
			t.Fatalf("effective node attribute %q is absent", name)
		}
	}

	relationshipType := findRelationshipType(t, root.RelationshipTypes, "ReflectionRelationship")
	if relationshipType.AttributeDefinitions["relationship_property"] == nil {
		t.Fatal("relationship property was not reflected")
	}
	capabilityType := findCapabilityType(t, root.CapabilityTypes, "ReflectionCapability")
	if capabilityType.AttributeDefinitions["capability_property"] == nil {
		t.Fatal("capability type property was not reflected")
	}
	capabilityDefinition := nodeType.CapabilityDefinitions["endpoint"]
	if capabilityDefinition == nil || capabilityDefinition.AttributeDefinitions["definition_property"] == nil {
		t.Fatal("node capability definition property was not reflected")
	}

	group := serviceTemplate.Groups["unrelated"]
	if group == nil {
		t.Fatal("unrelated group did not normalize")
	}
	groupJSON, marshalErr := json.Marshal(group)
	if marshalErr != nil {
		t.Fatalf("marshal unrelated group: %v", marshalErr)
	}
	if strings.Contains(string(groupJSON), `"attributes"`) {
		t.Fatalf("project-specific group property was reflected outside normative scope: %s", groupJSON)
	}
	policy := serviceTemplate.Policies["unrelated"]
	if policy == nil {
		t.Fatal("unrelated policy did not normalize")
	}
	policyJSON, marshalErr := json.Marshal(policy)
	if marshalErr != nil {
		t.Fatalf("marshal unrelated policy: %v", marshalErr)
	}
	if strings.Contains(string(policyJSON), `"attributes"`) {
		t.Fatalf("project-specific policy property was reflected outside normative scope: %s", policyJSON)
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 3.6.10.5, 3.6.11.2, 3.6.12.1, 3.6.13.2
// Requirement: TOSCA13-3.6.10.5-001
// Expected: assigned, defaulted, null, absent optional, and function-valued
// properties have deterministic reflected current-value behavior.
// Category: positive, negative, assignment, default, null, function, evaluation
func TestPropertyAttributeReflectionValues(t *testing.T) {
	source := `tosca_definitions_version: tosca_simple_yaml_1_3
node_types:
  ValueNode:
    derived_from: tosca.nodes.Root
    properties:
      assigned:
        type: string
      defaulted:
        type: string
        default: from-default
      optional:
        type: string
        required: false
        constraints:
          - min_length: 1
      computed:
        type: string
topology_template:
  node_templates:
    node:
      type: ValueNode
      properties:
        assigned: explicit
        computed: { concat: [function-, value] }
  outputs:
    reflected:
      value: { get_attribute: [node, computed] }
`
	serviceTemplate, parseProblems, err := testsupport.ParseSource(t, source)
	if err != nil {
		t.Fatalf("value reflection failed: %v\n%s", err, parseProblems)
	}
	node := serviceTemplate.NodeTemplates["node"]
	if node == nil {
		t.Fatal("normalized node is absent")
	}
	assertNormalizedPrimitive(t, node.Attributes["assigned"], "explicit")
	assertNormalizedPrimitive(t, node.Attributes["defaulted"], "from-default")
	assertNormalizedPrimitive(t, node.Attributes["optional"], nil)
	if _, ok := node.Attributes["computed"].(*normal.FunctionCall); !ok {
		t.Fatalf("function-valued reflected attribute has type %T", node.Attributes["computed"])
	}
	if node.Attributes["computed"] == node.Properties["computed"] {
		t.Fatal("normalized reflected function value aliases the normalized property value")
	}

	clout, compileErr := serviceTemplate.Compile()
	if compileErr != nil {
		t.Fatalf("compile reflected function: %v", compileErr)
	}
	urlContext := exturl.NewContext()
	defer urlContext.Release()
	evaluationProblems := problems.NewProblems(terminal.NewStylist(false))
	execContext := cloutjs.ExecContext{
		Clout:      clout,
		Problems:   evaluationProblems,
		URLContext: urlContext,
		Format:     "yaml",
	}
	execContext.Coerce()
	if !evaluationProblems.Empty() {
		t.Fatalf("reflected function evaluation failed:\n%s", evaluationProblems.ToString(false))
	}
	tosca, ok := clout.Properties["tosca"].(ard.StringMap)
	if !ok {
		t.Fatalf("compiled TOSCA properties have type %T", clout.Properties["tosca"])
	}
	outputs, ok := tosca["outputs"].(ard.StringMap)
	if !ok {
		t.Fatalf("compiled outputs have type %T", tosca["outputs"])
	}
	if got := outputs["reflected"]; got != "function-value" {
		t.Fatalf("evaluated reflected function = %#v, want %q", got, "function-value")
	}

	nullSource := `tosca_definitions_version: tosca_simple_yaml_1_3
node_types:
  NullNode:
    derived_from: tosca.nodes.Root
    properties:
      value:
        type: string
topology_template:
  node_templates:
    node:
      type: NullNode
      properties:
        value: null
`
	nullTemplate, nullProblems, nullErr := testsupport.ParseSource(t, nullSource)
	if nullErr != nil {
		t.Fatalf("null property did not follow the accepted property semantics: %v\n%s", nullErr, nullProblems)
	}
	nullProperty, ok := nullTemplate.NodeTemplates["node"].Properties["value"].(*normal.Primitive)
	if !ok {
		t.Fatalf("normalized null property has type %T", nullTemplate.NodeTemplates["node"].Properties["value"])
	}
	assertNormalizedPrimitive(t, nullTemplate.NodeTemplates["node"].Attributes["value"], nullProperty.Primitive)
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 3.5.1, 3.6.10.5, 3.6.10.6, 3.6.12.1
// Requirement: TOSCA13-3.6.10.5-001
// Expected: inherited, refined, and newly declared effective properties each
// produce exactly one reflected attribute with the final effective type.
// Category: positive, inheritance, refinement, type
func TestPropertyAttributeReflectionInheritanceAndRefinement(t *testing.T) {
	source := `tosca_definitions_version: tosca_simple_yaml_1_3
data_types:
  BaseRecord:
    derived_from: tosca.datatypes.Root
    properties:
      base:
        type: string
  DerivedRecord:
    derived_from: BaseRecord
    properties:
      extra:
        type: string
node_types:
  ParentNode:
    derived_from: tosca.nodes.Root
    properties:
      inherited:
        type: string
        required: false
      refined:
        type: BaseRecord
        required: false
  ChildNode:
    derived_from: ParentNode
    properties:
      refined:
        type: DerivedRecord
      added:
        type: integer
        required: false
topology_template:
  node_templates:
    node:
      type: ChildNode
`
	parserContext, _, parseProblems, release, err := parseReflectionSource(t, source)
	defer release()
	if err != nil {
		t.Fatalf("inheritance/refinement reflection failed: %v\n%s", err, parseProblems)
	}
	root := reflectionRoot(t, parserContext)
	child := findNodeType(t, root.NodeTypes, "ChildNode")
	for _, name := range []string{"inherited", "refined", "added"} {
		if child.PropertyDefinitions[name] == nil || child.AttributeDefinitions[name] == nil {
			t.Fatalf("effective property/attribute pair %q is incomplete", name)
		}
	}
	refinedProperty := child.PropertyDefinitions["refined"]
	refinedAttribute := child.AttributeDefinitions["refined"]
	if refinedAttribute.DataType != refinedProperty.DataType ||
		refinedAttribute.DataTypeName == nil || *refinedAttribute.DataTypeName != "DerivedRecord" {
		t.Fatalf("refined attribute did not use final effective type: %#v", refinedAttribute.DataTypeName)
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 3.6.10.5, 3.6.12.1, 3.6.12.4
// Requirement: TOSCA13-3.6.10.5-001
// Expected: a compatible explicit same-name attribute is the single effective
// declaration; an incompatible declaration is rejected deterministically.
// Category: positive, negative, explicit attribute, conflict, diagnostic
func TestPropertyAttributeReflectionExplicitAttributeConflict(t *testing.T) {
	compatible := `tosca_definitions_version: tosca_simple_yaml_1_3
node_types:
  CompatibleNode:
    derived_from: tosca.nodes.Root
    properties:
      reflected:
        type: string
        required: false
    attributes:
      reflected:
        type: string
        description: explicit attribute wins
      unrelated:
        type: integer
topology_template:
  node_templates:
    node:
      type: CompatibleNode
      properties:
        reflected: value
`
	parserContext, serviceTemplate, parseProblems, release, err := parseReflectionSource(t, compatible)
	defer release()
	if err != nil {
		t.Fatalf("compatible explicit attribute failed: %v\n%s", err, parseProblems)
	}
	root := reflectionRoot(t, parserContext)
	nodeType := findNodeType(t, root.NodeTypes, "CompatibleNode")
	reflected := nodeType.AttributeDefinitions["reflected"]
	if reflected.Description == nil || *reflected.Description != "explicit attribute wins" {
		t.Fatalf("explicit attribute definition was not preserved: %v", reflected.Description)
	}
	assertNormalizedPrimitive(t, serviceTemplate.NodeTemplates["node"].Attributes["reflected"], "value")
	if nodeType.AttributeDefinitions["unrelated"] == nil {
		t.Fatal("different-name explicit attribute was affected")
	}

	incompatible := `tosca_definitions_version: tosca_simple_yaml_1_3
node_types:
  ConflictNode:
    derived_from: tosca.nodes.Root
    properties:
      reflected:
        type: string
        required: false
    attributes:
      reflected:
        type: integer
topology_template:
  node_templates:
    node:
      type: ConflictNode
`
	firstProblems := parseRejectedReflectionSource(t, incompatible)
	secondProblems := parseRejectedReflectionSource(t, incompatible)
	firstBody := reflectionDiagnosticBody(firstProblems)
	secondBody := reflectionDiagnosticBody(secondProblems)
	if firstBody != secondBody {
		t.Fatalf("conflict diagnostic is nondeterministic:\nfirst:\n%s\nsecond:\n%s", firstProblems, secondProblems)
	}
	for _, expected := range []string{"reflected", "property", "attribute", "string", "integer"} {
		if !strings.Contains(firstProblems, expected) {
			t.Fatalf("conflict diagnostic does not contain %q:\n%s", expected, firstProblems)
		}
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 3.6.10.5, 3.6.13.2, 4.5.1
// Requirement: TOSCA13-3.6.10.5-001
// Expected: get_attribute resolves reflected attributes and rejects unknown
// names using effective attribute definitions.
// Category: positive, negative, function resolution, SELF
func TestPropertyAttributeReflectionGetAttributeResolution(t *testing.T) {
	source := `tosca_definitions_version: tosca_simple_yaml_1_3
node_types:
  FunctionNode:
    derived_from: tosca.nodes.Root
    properties:
      reflected:
        type: string
relationship_types:
  FunctionRelationship:
    derived_from: tosca.relationships.Root
    attributes:
      from_source:
        type: string
      from_target:
        type: string
topology_template:
  node_templates:
    node:
      type: FunctionNode
      properties:
        reflected: value
      attributes:
        state: { get_attribute: [SELF, reflected] }
  relationship_templates:
    relationship:
      type: FunctionRelationship
      attributes:
        from_source: { get_attribute: [SOURCE, reflected] }
        from_target: { get_attribute: [TARGET, reflected] }
  outputs:
    result:
      value: { get_attribute: [node, reflected] }
`
	parserContext, serviceTemplate, parseProblems, release, err := parseReflectionSource(t, source)
	defer release()
	if err != nil {
		t.Fatalf("get_attribute did not resolve reflected definition: %v\n%s", err, parseProblems)
	}
	call := normalizedOutputFunction(t, serviceTemplate)
	if call.Name != "tosca.function.get_attribute" {
		t.Fatalf("normalized function name = %q", call.Name)
	}
	root := reflectionRoot(t, parserContext)
	var relationship *tosca_v2_0.RelationshipTemplate
	for _, candidate := range root.ServiceTemplate.RelationshipTemplates {
		if candidate.Name == "relationship" {
			relationship = candidate
			break
		}
	}
	if relationship == nil {
		t.Fatal("SOURCE/TARGET relationship context is absent")
	}
	for _, name := range []string{"from_source", "from_target"} {
		value := relationship.Attributes[name]
		if value == nil {
			t.Fatalf("relationship attribute %q is absent", name)
		}
		if _, ok := value.Context.Data.(*parsing.FunctionCall); !ok {
			t.Fatalf("relationship attribute %q has parser value type %T", name, value.Context.Data)
		}
	}

	unknown := strings.ReplaceAll(source, "[node, reflected]", "[node, absent]")
	if _, unknownProblems, unknownErr := testsupport.ParseSource(t, unknown); unknownErr == nil {
		t.Fatal("unknown reflected attribute was accepted")
	} else if !strings.Contains(unknownProblems, `attribute_name "absent" not found`) {
		t.Fatalf("unknown reflected attribute failed for the wrong reason:\n%s", unknownProblems)
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 3.6.10.5, 3.6.12.1
// Requirement: TOSCA13-3.6.10.5-001
// Expected: node, relationship, and capability normalized representations
// contain reflected values and repeated normalization is idempotent.
// Category: positive, normalization, relationship, capability, idempotence,
// determinism
func TestPropertyAttributeReflectionNormalization(t *testing.T) {
	source := `tosca_definitions_version: tosca_simple_yaml_1_3
capability_types:
  ValueCapability:
    derived_from: tosca.capabilities.Root
    properties:
      value:
        type: string
        default: capability-default
relationship_types:
  ValueRelationship:
    derived_from: tosca.relationships.Root
    properties:
      value:
        type: string
        default: relationship-default
node_types:
  ValueNode:
    derived_from: tosca.nodes.Root
    properties:
      zeta:
        type: string
        default: z
      alpha:
        type: string
        default: a
    capabilities:
      endpoint:
        type: ValueCapability
  TargetNode:
    derived_from: tosca.nodes.Root
    capabilities:
      endpoint:
        type: tosca.capabilities.Root
  SourceNode:
    derived_from: tosca.nodes.Root
    requirements:
      - link:
          capability: tosca.capabilities.Root
          relationship: ValueRelationship
          occurrences: [0, 1]
topology_template:
  node_templates:
    node:
      type: ValueNode
    target:
      type: TargetNode
    source:
      type: SourceNode
      requirements:
        - link:
            node: target
            capability: endpoint
            relationship:
              type: ValueRelationship
  relationship_templates:
    relationship:
      type: ValueRelationship
`
	parserContext, urlContext := reflectionParserContext(t, source)
	defer urlContext.Release()
	first, err := parserContext.Parse(context.Background())
	if err != nil {
		t.Fatalf("initial reflection parse failed: %v\n%s", err, parserContext.GetProblems().ToString(false))
	}
	second, ok := parserContext.Normalize()
	if !ok {
		t.Fatal("second reflection normalization failed")
	}
	third, ok := parserContext.Normalize()
	if !ok {
		t.Fatal("third reflection normalization failed")
	}
	firstJSON := reflectionNormalizedJSON(t, first)
	secondJSON := reflectionNormalizedJSON(t, second)
	thirdJSON := reflectionNormalizedJSON(t, third)
	if firstJSON != secondJSON || secondJSON != thirdJSON {
		t.Fatalf("reflection normalization is not idempotent/deterministic:\nfirst: %s\nsecond: %s\nthird: %s",
			firstJSON, secondJSON, thirdJSON)
	}
	node := first.NodeTemplates["node"]
	assertNormalizedPrimitive(t, node.Attributes["alpha"], "a")
	assertNormalizedPrimitive(t, node.Attributes["zeta"], "z")
	assertNormalizedPrimitive(t, node.Capabilities["endpoint"].Attributes["value"], "capability-default")
	sourceNode := first.NodeTemplates["source"]
	if sourceNode == nil || len(sourceNode.Requirements) != 1 || sourceNode.Requirements[0].Relationship == nil {
		t.Fatalf("normalized relationship assignment is absent: %#v", sourceNode)
	}
	assertNormalizedPrimitive(t, sourceNode.Requirements[0].Relationship.Attributes["value"], "relationship-default")
}

// Specification: TOSCA Version 2.0 is independently selected and is not a
// normative source for TOSCA13-3.6.10.5-001.
// Expected: the TOSCA 1.3 policy hook does not reflect TOSCA 2.0 properties.
// Category: negative, cross-version, regression, version isolation
func TestPropertyAttributeReflectionDoesNotChangeTosca20(t *testing.T) {
	source := `tosca_definitions_version: tosca_2_0
imports:
  - profile: org.oasis-open.simple:2.0
    namespace: tosca
node_types:
  Version20Node:
    properties:
      version_specific:
        type: string
        default: property-only
service_template:
  node_templates:
    node:
      type: Version20Node
`
	serviceTemplate, parseProblems, err := testsupport.ParseSource(t, source)
	if err != nil {
		t.Fatalf("TOSCA 2.0 regression template failed: %v\n%s", err, parseProblems)
	}
	node := serviceTemplate.NodeTemplates["node"]
	if node == nil {
		t.Fatal("TOSCA 2.0 node did not normalize")
	}
	if _, exists := node.Attributes["version_specific"]; exists {
		t.Fatal("TOSCA 1.3 property reflection leaked into TOSCA 2.0")
	}
}

func parseReflectionSource(t *testing.T, source string) (*parser.Context, *normal.ServiceTemplate, string, func(), error) {
	t.Helper()
	parserContext, urlContext := reflectionParserContext(t, source)
	serviceTemplate, err := parserContext.Parse(context.Background())
	parseProblems := ""
	if parserContext.Root != nil {
		parseProblems = parserContext.GetProblems().ToString(false)
	}
	return parserContext, serviceTemplate, parseProblems, func() {
		_ = urlContext.Release()
	}, err
}

func reflectionParserContext(t *testing.T, source string) (*parser.Context, *exturl.Context) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "service-template.yaml")
	if err := os.WriteFile(path, []byte(source), 0o600); err != nil {
		t.Fatalf("write reflection fixture: %v", err)
	}
	urlContext := exturl.NewContext()
	parserContext := parser.NewParser().NewContext()
	parserContext.URL = urlContext.NewFileURL(filepath.ToSlash(path))
	return parserContext, urlContext
}

func parseRejectedReflectionSource(t *testing.T, source string) string {
	t.Helper()
	_, parseProblems, err := testsupport.ParseSource(t, source)
	if err == nil {
		t.Fatal("incompatible explicit reflected attribute was accepted")
	}
	return parseProblems
}

func reflectionRoot(t *testing.T, parserContext *parser.Context) *tosca_v1_3.ServiceFile {
	t.Helper()
	if parserContext.Root == nil {
		t.Fatal("parser root is absent")
	}
	root, ok := parserContext.Root.EntityPtr.(*tosca_v1_3.ServiceFile)
	if !ok {
		t.Fatalf("parser root has type %T", parserContext.Root.EntityPtr)
	}
	return root
}

func findNodeType(t *testing.T, types tosca_v2_0.NodeTypes, name string) *tosca_v2_0.NodeType {
	t.Helper()
	for _, type_ := range types {
		if type_.Name == name {
			return type_
		}
	}
	t.Fatalf("node type %q not found", name)
	return nil
}

func findRelationshipType(t *testing.T, types tosca_v2_0.RelationshipTypes, name string) *tosca_v2_0.RelationshipType {
	t.Helper()
	for _, type_ := range types {
		if type_.Name == name {
			return type_
		}
	}
	t.Fatalf("relationship type %q not found", name)
	return nil
}

func findCapabilityType(t *testing.T, types tosca_v2_0.CapabilityTypes, name string) *tosca_v2_0.CapabilityType {
	t.Helper()
	for _, type_ := range types {
		if type_.Name == name {
			return type_
		}
	}
	t.Fatalf("capability type %q not found", name)
	return nil
}

func assertNormalizedPrimitive(t *testing.T, value normal.Value, want any) {
	t.Helper()
	primitive, ok := value.(*normal.Primitive)
	if !ok {
		t.Fatalf("normalized value has type %T, want *normal.Primitive", value)
	}
	if !reflect.DeepEqual(primitive.Primitive, want) {
		t.Fatalf("normalized primitive = %#v, want %#v", primitive.Primitive, want)
	}
}

func reflectionNormalizedJSON(t *testing.T, serviceTemplate *normal.ServiceTemplate) string {
	t.Helper()
	data, err := json.Marshal(serviceTemplate)
	if err != nil {
		t.Fatalf("marshal normalized reflection template: %v", err)
	}
	return string(data)
}

func reflectionDiagnosticBody(diagnostic string) string {
	if index := strings.Index(diagnostic, `node_types[`); index >= 0 {
		return diagnostic[index:]
	}
	return diagnostic
}
