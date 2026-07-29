package tosca_1_3_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/tliron/go-puccini/tests/conformance/internal/testsupport"
)

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 5.5.7.4 Additional requirements
// Requirement: TOSCA13-5.5.7.4-001
// Expected: an explicitly assigned Endpoint supplies port or ports.
// Category: positive, rendering, normalization
func TestPartialEndpointPortOrPortsAccepted(t *testing.T) {
	for _, properties := range []string{
		"{ port: 8080 }",
		"{ ports: { http: { target: 8080 } } }",
		"{ port: 8080, ports: { http: { target: 8080 } } }",
	} {
		if _, problems, err := testsupport.ParseSource(t, capabilityProfileTemplate(
			"tosca.capabilities.Endpoint",
			"properties: "+properties,
		)); err != nil {
			t.Fatalf("valid Endpoint failed: %v\n%s", err, problems)
		}
	}

	serviceTemplate, problems, err := testsupport.ParseSource(t, capabilityProfileTemplate(
		"tosca.capabilities.Endpoint",
		"properties: { port: 8080 }",
	))
	if err != nil {
		t.Fatalf("normalized Endpoint failed: %v\n%s", err, problems)
	}
	assertNormalizedPrimitive(t, serviceTemplate.NodeTemplates["node"].Capabilities["profile"].Properties["port"], 8080)
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 5.5.7.4 Additional requirements
// Requirement: TOSCA13-5.5.7.4-001
// Expected: an explicit Endpoint assignment with neither port nor ports is
// rejected after effective property rendering.
// Category: negative, rendering
func TestPartialEndpointWithoutPortOrPortsRejected(t *testing.T) {
	for _, assignment := range []string{"{}", "properties: { protocol: tcp }"} {
		problems := parseRejectedCapabilityProfile(t, "tosca.capabilities.Endpoint", assignment)
		for _, expected := range []string{`capabilities["profile"]`, "Endpoint", "port", "ports"} {
			if !strings.Contains(problems, expected) {
				t.Fatalf("Endpoint diagnostic does not contain %q:\n%s", expected, problems)
			}
		}
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 5.5.7.4 Additional requirements
// Requirement: TOSCA13-5.5.7.4-001
// Expected: the rule follows capability type inheritance and does not apply
// to an unrelated capability with structurally similar properties.
// Category: inheritance, negative, scope regression
func TestPartialEndpointRuleFollowsCapabilityTypeIdentity(t *testing.T) {
	derived := `tosca_definitions_version: tosca_simple_yaml_1_3
capability_types:
  DerivedEndpoint:
    derived_from: tosca.capabilities.Endpoint
node_types:
  Holder:
    derived_from: tosca.nodes.Root
    capabilities:
      profile:
        type: DerivedEndpoint
topology_template:
  node_templates:
    node:
      type: Holder
      capabilities:
        profile: {}
`
	_, problems, err := testsupport.ParseSource(t, derived)
	if err == nil {
		t.Fatal("derived Endpoint bypassed the port-or-ports rule")
	}
	if !strings.Contains(problems, "Endpoint") {
		t.Fatalf("derived Endpoint diagnostic is not specific:\n%s", problems)
	}

	unrelated := `tosca_definitions_version: tosca_simple_yaml_1_3
capability_types:
  Unrelated:
    derived_from: tosca.capabilities.Root
    properties:
      port:
        type: integer
        required: false
      ports:
        type: map
        entry_schema:
          type: string
        required: false
node_types:
  Holder:
    derived_from: tosca.nodes.Root
    capabilities:
      profile:
        type: Unrelated
topology_template:
  node_templates:
    node:
      type: Holder
      capabilities:
        profile: {}
`
	if _, unrelatedProblems, unrelatedErr := testsupport.ParseSource(t, unrelated); unrelatedErr != nil {
		t.Fatalf("Endpoint rule leaked to unrelated capability: %v\n%s", unrelatedErr, unrelatedProblems)
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 5.5.7.4 Additional requirements
// Requirement: TOSCA13-5.5.7.4-001
// Expected: absence of a source capability assignment is not reinterpreted as
// an explicit empty Endpoint merely because normalization materializes all
// capability definitions.
// Category: boundary, source-document fidelity
func TestPartialEndpointOmittedAssignmentIsNotExplicitEmptyValue(t *testing.T) {
	source := `tosca_definitions_version: tosca_simple_yaml_1_3
node_types:
  Holder:
    derived_from: tosca.nodes.Root
    capabilities:
      profile:
        type: tosca.capabilities.Endpoint
topology_template:
  node_templates:
    node:
      type: Holder
`
	if _, problems, err := testsupport.ParseSource(t, source); err != nil {
		t.Fatalf("omitted capability assignment became an explicit invalid value: %v\n%s", err, problems)
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 5.5.13.1 Properties
// Requirement: TOSCA13-5.5.13.1-012
// Expected: default_instances lies inclusively between effective min_instances
// and max_instances.
// Category: positive, boundary, defaults, rendering
func TestPartialScalableDefaultInstancesInsideRange(t *testing.T) {
	for _, properties := range []string{
		"{ default_instances: 1 }",
		"{ min_instances: 2, max_instances: 4, default_instances: 2 }",
		"{ min_instances: 2, max_instances: 4, default_instances: 4 }",
		"{ min_instances: 2, max_instances: 4, default_instances: 3 }",
	} {
		if _, problems, err := testsupport.ParseSource(t, capabilityProfileTemplate(
			"tosca.capabilities.Scalable",
			"properties: "+properties,
		)); err != nil {
			t.Fatalf("valid Scalable assignment failed: %v\n%s", err, problems)
		}
	}

	serviceTemplate, problems, err := testsupport.ParseSource(t, capabilityProfileTemplate(
		"tosca.capabilities.Scalable",
		"properties: { min_instances: 2, max_instances: 4, default_instances: 3 }",
	))
	if err != nil {
		t.Fatalf("normalized Scalable failed: %v\n%s", err, problems)
	}
	capability := serviceTemplate.NodeTemplates["node"].Capabilities["profile"]
	assertNormalizedPrimitive(t, capability.Properties["min_instances"], 2)
	assertNormalizedPrimitive(t, capability.Properties["max_instances"], 4)
	assertNormalizedPrimitive(t, capability.Properties["default_instances"], 3)
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 5.5.13.1 Properties
// Requirement: TOSCA13-5.5.13.1-012
// Expected: values below min_instances and above max_instances are rejected
// with a property-specific rendering diagnostic.
// Category: negative, boundary, rendering
func TestPartialScalableDefaultInstancesOutsideRangeRejected(t *testing.T) {
	for _, properties := range []string{
		"{ min_instances: 2, max_instances: 4, default_instances: 1 }",
		"{ min_instances: 2, max_instances: 4, default_instances: 5 }",
		"{ default_instances: 2 }",
	} {
		problems := parseRejectedCapabilityProfile(
			t,
			"tosca.capabilities.Scalable",
			"properties: "+properties,
		)
		for _, expected := range []string{`capabilities["profile"].properties["default_instances"]`, "min_instances", "max_instances"} {
			if !strings.Contains(problems, expected) {
				t.Fatalf("Scalable diagnostic does not contain %q:\n%s", expected, problems)
			}
		}
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 5.5.13.1 Properties
// Requirement: TOSCA13-5.5.13.1-012
// Expected: inherited/refined min and max defaults participate in the same
// effective assignment check.
// Category: inheritance, refinement, defaults
func TestPartialScalableUsesInheritedRefinedDefaults(t *testing.T) {
	source := `tosca_definitions_version: tosca_simple_yaml_1_3
capability_types:
  RefinedScalable:
    derived_from: tosca.capabilities.Scalable
    properties:
      min_instances:
        default: 2
      max_instances:
        default: 4
node_types:
  Holder:
    derived_from: tosca.nodes.Root
    capabilities:
      profile:
        type: RefinedScalable
topology_template:
  node_templates:
    node:
      type: Holder
      capabilities:
        profile:
          properties:
            default_instances: 3
`
	if _, problems, err := testsupport.ParseSource(t, source); err != nil {
		t.Fatalf("valid refined Scalable assignment failed: %v\n%s", err, problems)
	}

	invalid := strings.Replace(source, "default_instances: 3", "default_instances: 5", 1)
	if _, invalidProblems, invalidErr := testsupport.ParseSource(t, invalid); invalidErr == nil {
		t.Fatal("refined Scalable bounds were ignored")
	} else if !strings.Contains(invalidProblems, "default_instances") {
		t.Fatalf("refined Scalable diagnostic is not specific:\n%s", invalidProblems)
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 5.5.7.4 Additional requirements; 5.5.13.1 Properties
// Expected: presence checks accept an intrinsic-valued Endpoint port, while
// cross-property numeric comparison is deferred until Scalable values are
// concrete rather than treating function syntax as a final integer.
// Category: function value, rendering
func TestPartialCapabilityProfileDefersFunctionValues(t *testing.T) {
	source := `tosca_definitions_version: tosca_simple_yaml_1_3
node_types:
  Holder:
    derived_from: tosca.nodes.Root
    capabilities:
      endpoint:
        type: tosca.capabilities.Endpoint
      scalable:
        type: tosca.capabilities.Scalable
topology_template:
  inputs:
    port:
      type: integer
      default: 8080
    count:
      type: integer
      default: 2
  node_templates:
    node:
      type: Holder
      capabilities:
        endpoint:
          properties:
            port: { get_input: port }
        scalable:
          properties:
            min_instances: 1
            max_instances: 3
            default_instances: { get_input: count }
`
	if _, problems, err := testsupport.ParseSource(t, source); err != nil {
		t.Fatalf("function-valued capability properties were validated prematurely: %v\n%s", err, problems)
	}
}

// Specification: TOSCA 1.3 sections 5.5.7.4 and 5.5.13.1 only
// Expected: project-defined TOSCA 2.0 capabilities with the same property
// names retain their existing behavior.
// Category: cross-version regression
func TestPartialCapabilityProfileRulesDoNotChangeTOSCA20(t *testing.T) {
	source := `tosca_definitions_version: tosca_2_0
capability_types:
  Endpoint:
    properties:
      port:
        type: integer
        required: false
  Scalable:
    properties:
      min_instances:
        type: integer
        default: 2
      max_instances:
        type: integer
        default: 4
      default_instances:
        type: integer
        required: false
node_types:
  Holder:
    capabilities:
      endpoint:
        type: Endpoint
      scalable:
        type: Scalable
service_template:
  node_templates:
    node:
      type: Holder
      capabilities:
        endpoint: {}
        scalable:
          properties:
            default_instances: 5
`
	if _, problems, err := testsupport.ParseSource(t, source); err != nil {
		t.Fatalf("TOSCA 1.3 capability profile policy changed TOSCA 2.0: %v\n%s", err, problems)
	}
}

func capabilityProfileTemplate(capabilityType, assignment string) string {
	return fmt.Sprintf(`tosca_definitions_version: tosca_simple_yaml_1_3
node_types:
  Holder:
    derived_from: tosca.nodes.Root
    capabilities:
      profile:
        type: %s
topology_template:
  node_templates:
    node:
      type: Holder
      capabilities:
        profile:
          %s
`, capabilityType, assignment)
}

func parseRejectedCapabilityProfile(t *testing.T, capabilityType, assignment string) string {
	t.Helper()
	_, problems, err := testsupport.ParseSource(t, capabilityProfileTemplate(capabilityType, assignment))
	if err == nil {
		t.Fatalf("invalid %s assignment was accepted", capabilityType)
	}
	return problems
}
