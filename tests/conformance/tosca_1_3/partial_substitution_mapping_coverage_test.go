package tosca_1_3_test

import (
	"strings"
	"testing"

	"github.com/tliron/go-puccini/tests/conformance/internal/testsupport"
)

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.8.13.4, Additional requirements
// Requirement: TOSCA13-3.8.13.4-001
// Expected: every effective property, capability, and requirement of the
// substituted node type has a mapping.
// Category: positive, inheritance, resolution, normalization, direct
func TestPartialSubstitutionMappingCoversEffectiveNodeType(t *testing.T) {
	serviceTemplate, problems, err := testsupport.ParseSource(
		t,
		substitutionCoverageTemplate(""),
	)
	if err != nil {
		t.Fatalf("complete substitution mapping was rejected: %v\n%s", err, problems)
	}
	substitution := serviceTemplate.Substitution
	if substitution == nil {
		t.Fatal("normalized substitution mapping is absent")
	}
	for _, name := range []string{"inherited_property", "local_property"} {
		if _, ok := substitution.InputPointers[name]; !ok {
			t.Fatalf("normalized property mapping %q is absent", name)
		}
	}
	for _, name := range []string{"inherited_capability", "local_capability"} {
		if _, ok := substitution.CapabilityPointers[name]; !ok {
			t.Fatalf("normalized capability mapping %q is absent", name)
		}
	}
	for _, name := range []string{"inherited_requirement", "local_requirement"} {
		if _, ok := substitution.RequirementPointer[name]; !ok {
			t.Fatalf("normalized requirement mapping %q is absent", name)
		}
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.8.13.4, Additional requirements
// Requirement: TOSCA13-3.8.13.4-001
// Expected: omitting any property definition from substitution properties is
// rejected, including an optional inherited property.
// Category: negative, property, inheritance, boundary, direct
func TestPartialSubstitutionMappingRejectsMissingProperties(t *testing.T) {
	for _, name := range []string{"inherited_property", "local_property"} {
		t.Run(name, func(t *testing.T) {
			problems := rejectSubstitutionCoverage(t, substitutionCoverageTemplate(name))
			assertMissingSubstitutionMapping(t, problems, "properties", name, "property")
		})
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.8.13.4, Additional requirements
// Requirement: TOSCA13-3.8.13.4-001
// Expected: omitting any capability definition from substitution capabilities
// is rejected, including an inherited capability.
// Category: negative, capability, inheritance, direct
func TestPartialSubstitutionMappingRejectsMissingCapabilities(t *testing.T) {
	for _, name := range []string{"inherited_capability", "local_capability"} {
		t.Run(name, func(t *testing.T) {
			problems := rejectSubstitutionCoverage(t, substitutionCoverageTemplate(name))
			assertMissingSubstitutionMapping(t, problems, "capabilities", name, "capability")
		})
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.8.13.4, Additional requirements
// Requirement: TOSCA13-3.8.13.4-001
// Expected: omitting any requirement definition from substitution requirements
// is rejected even when its lower occurrence bound is zero.
// Category: negative, requirement, inheritance, optional boundary, direct
func TestPartialSubstitutionMappingRejectsMissingRequirements(t *testing.T) {
	for _, name := range []string{"inherited_requirement", "local_requirement"} {
		t.Run(name, func(t *testing.T) {
			problems := rejectSubstitutionCoverage(t, substitutionCoverageTemplate(name))
			assertMissingSubstitutionMapping(t, problems, "requirements", name, "requirement")
		})
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.8.13.4, Additional requirements
// Requirement: TOSCA13-3.8.13.4-001
// Expected: the completeness rule does not require mappings for attributes or
// interfaces, which are not named by the MUST sentence.
// Category: positive, scope boundary, direct
func TestPartialSubstitutionMappingCoverageScopeIsExact(t *testing.T) {
	if _, problems, err := testsupport.ParseSource(t, substitutionCoverageTemplate("")); err != nil {
		t.Fatalf("unmapped attribute changed mapping completeness scope: %v\n%s", err, problems)
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.8.13.4, Additional requirements
// Requirement: TOSCA13-3.8.13.4-001
// Expected: supplied mappings continue to be reference-checked independently
// of the missing-mapping check.
// Category: negative, resolution regression, direct
func TestPartialSubstitutionMappingUnknownDefinitionStillRejected(t *testing.T) {
	source := strings.Replace(
		substitutionCoverageTemplate(""),
		"    node_type: AbstractNode\n    properties:\n",
		"    node_type: AbstractNode\n    properties:\n      absent: [ inherited_input ]\n",
		1,
	)
	_, problems, err := testsupport.ParseSource(t, source)
	if err == nil {
		t.Fatal("unknown supplied substitution property mapping was accepted")
	}
	if !strings.Contains(problems, "absent") || !strings.Contains(problems, "property") {
		t.Fatalf("unknown supplied mapping failed for the wrong reason:\n%s", problems)
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.8.13.4, Additional requirements
// Requirement: TOSCA13-3.8.13.4-001
// Expected: multiple missing definitions are reported in stable category/name
// order across processing operations.
// Category: negative, determinism, diagnostic order, direct
func TestPartialSubstitutionMappingCoverageDiagnosticIsDeterministic(t *testing.T) {
	source := substitutionCoverageTemplate("all")
	first := substitutionCoverageDiagnosticBody(rejectSubstitutionCoverage(t, source))
	second := substitutionCoverageDiagnosticBody(rejectSubstitutionCoverage(t, source))
	if first != second {
		t.Fatalf("mapping-coverage diagnostic is nondeterministic:\nfirst:\n%s\nsecond:\n%s", first, second)
	}

	ordered := []string{
		`capabilities["inherited_capability"]`,
		`capabilities["local_capability"]`,
		`properties["inherited_property"]`,
		`properties["local_property"]`,
		`requirements["inherited_requirement"]`,
		`requirements["local_requirement"]`,
	}
	position := -1
	for _, fragment := range ordered {
		next := strings.Index(first, fragment)
		if next <= position {
			t.Fatalf("diagnostic %q is absent or out of order:\n%s", fragment, first)
		}
		position = next
	}
}

// Specification: TOSCA Version 2.0
// Cross-version isolation: the TOSCA 1.3 completeness policy must not silently
// change the shared TOSCA 2.0 substitution renderer.
// Category: positive, cross-version regression
func TestPartialSubstitutionMappingCoverageDoesNotChangeTOSCA20(t *testing.T) {
	source := `tosca_definitions_version: tosca_2_0
node_types:
  AbstractNode:
    properties:
      setting:
        type: string
  InternalNode: {}
service_template:
  node_templates:
    internal:
      type: InternalNode
  substitution_mappings:
    node_type: AbstractNode
`
	if _, problems, err := testsupport.ParseSource(t, source); err != nil {
		t.Fatalf("TOSCA 2.0 substitution behavior changed: %v\n%s", err, problems)
	}
}

func substitutionCoverageTemplate(omit string) string {
	mapping := func(name string, line string) string {
		if omit == name || omit == "all" {
			return ""
		}
		return line
	}

	propertyMappings :=
		mapping("inherited_property", "      inherited_property: [ inherited_input ]\n") +
			mapping("local_property", "      local_property: [ local_input ]\n")
	capabilityMappings :=
		mapping("inherited_capability", "      inherited_capability: [ internal, inherited_endpoint ]\n") +
			mapping("local_capability", "      local_capability: [ internal, local_endpoint ]\n")
	requirementMappings :=
		mapping("inherited_requirement", "      inherited_requirement: [ internal, external_inherited ]\n") +
			mapping("local_requirement", "      local_requirement: [ internal, external_local ]\n")
	if propertyMappings == "" {
		propertyMappings = "      {}\n"
	}
	if capabilityMappings == "" {
		capabilityMappings = "      {}\n"
	}
	if requirementMappings == "" {
		requirementMappings = "      {}\n"
	}

	return `tosca_definitions_version: tosca_simple_yaml_1_3
capability_types:
  ExampleEndpoint: {}
relationship_types:
  ExampleLink:
    valid_target_types: [ ExampleEndpoint ]
node_types:
  AbstractBase:
    properties:
      inherited_property:
        type: string
        required: false
    attributes:
      observed:
        type: string
    capabilities:
      inherited_capability: ExampleEndpoint
    requirements:
      - inherited_requirement:
          capability: ExampleEndpoint
          relationship: ExampleLink
          occurrences: [ 0, 1 ]
  AbstractNode:
    derived_from: AbstractBase
    properties:
      local_property:
        type: string
    capabilities:
      local_capability: ExampleEndpoint
    requirements:
      - local_requirement:
          capability: ExampleEndpoint
          relationship: ExampleLink
          occurrences: [ 0, 1 ]
  InternalNode:
    capabilities:
      inherited_endpoint: ExampleEndpoint
      local_endpoint: ExampleEndpoint
    requirements:
      - external_inherited:
          capability: ExampleEndpoint
          relationship: ExampleLink
          occurrences: [ 0, 1 ]
      - external_local:
          capability: ExampleEndpoint
          relationship: ExampleLink
          occurrences: [ 0, 1 ]
topology_template:
  inputs:
    inherited_input:
      type: string
      default: inherited
    local_input:
      type: string
      default: local
  node_templates:
    internal:
      type: InternalNode
  substitution_mappings:
    node_type: AbstractNode
    properties:
` + propertyMappings + `    capabilities:
` + capabilityMappings + `    requirements:
` + requirementMappings
}

func rejectSubstitutionCoverage(t *testing.T, source string) string {
	t.Helper()
	_, problems, err := testsupport.ParseSource(t, source)
	if err == nil {
		t.Fatal("incomplete substitution mapping was accepted")
	}
	return problems
}

func assertMissingSubstitutionMapping(t *testing.T, problems string, category string, name string, kind string) {
	t.Helper()
	for _, fragment := range []string{
		`substitution_mappings.` + category + `["` + name + `"]`,
		"required " + kind + " mapping",
	} {
		if !strings.Contains(problems, fragment) {
			t.Fatalf("missing %s mapping diagnostic does not contain %q:\n%s", kind, fragment, problems)
		}
	}
}

func substitutionCoverageDiagnosticBody(problems string) string {
	if index := strings.Index(problems, "topology_template"); index >= 0 {
		return problems[index:]
	}
	return problems
}
