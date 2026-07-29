package tosca_1_3_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/tliron/go-puccini/tests/conformance/internal/testsupport"
)

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 8.5.1.1 Properties
// Requirement: TOSCA13-8.5.1.1-036
// Expected: flat and vlan networks with physical_network are accepted and the
// concrete values survive normalization.
// Category: positive, rendering, normalization
func TestPartialNetworkPhysicalNetworkConditionalRequirement(t *testing.T) {
	for _, networkType := range []string{"flat", "vlan"} {
		serviceTemplate, problems, err := testsupport.ParseSource(t, networkProfileTemplate(
			"tosca.nodes.network.Network",
			fmt.Sprintf("{ network_type: %s, physical_network: physnet1 }", networkType),
		))
		if err != nil {
			t.Fatalf("valid %s network failed: %v\n%s", networkType, err, problems)
		}
		node := serviceTemplate.NodeTemplates["network"]
		assertNormalizedPrimitive(t, node.Properties["network_type"], networkType)
		assertNormalizedPrimitive(t, node.Properties["physical_network"], "physnet1")
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 8.5.1.1 Properties
// Requirement: TOSCA13-8.5.1.1-036
// Expected: flat or vlan without physical_network is rejected after effective
// node property rendering.
// Category: negative, rendering
func TestPartialFlatAndVlanNetworkRequirePhysicalNetwork(t *testing.T) {
	for _, networkType := range []string{"flat", "vlan"} {
		problems := parseRejectedNetworkProfile(t, "tosca.nodes.network.Network", fmt.Sprintf("{ network_type: %s }", networkType))
		for _, expected := range []string{`properties["physical_network"]`, "network_type", networkType} {
			if !strings.Contains(problems, expected) {
				t.Fatalf("%s network diagnostic does not contain %q:\n%s", networkType, expected, problems)
			}
		}
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 8.5.1.1 Properties
// Requirement: TOSCA13-8.5.1.1-036
// Expected: other or absent network_type values do not make physical_network
// conditionally required.
// Category: positive, boundary
func TestPartialOtherNetworkTypesDoNotRequirePhysicalNetwork(t *testing.T) {
	for _, properties := range []string{"{}", "{ network_type: gre }", "{ network_type: vxlan }", "{ network_type: custom }"} {
		if _, problems, err := testsupport.ParseSource(t, networkProfileTemplate("tosca.nodes.network.Network", properties)); err != nil {
			t.Fatalf("non-flat/vlan network %s failed: %v\n%s", properties, err, problems)
		}
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 8.5.1.1 Properties
// Requirement: TOSCA13-8.5.1.1-036
// Expected: the conditional follows node type inheritance and effective
// property defaults.
// Category: inheritance, defaults, negative, positive
func TestPartialDerivedNetworkRetainsConditionalRule(t *testing.T) {
	valid := `tosca_definitions_version: tosca_simple_yaml_1_3
node_types:
  ValidFlat:
    derived_from: tosca.nodes.network.Network
    properties:
      network_type:
        default: flat
      physical_network:
        default: physnet1
topology_template:
  node_templates:
    network:
      type: ValidFlat
`
	if _, problems, err := testsupport.ParseSource(t, valid); err != nil {
		t.Fatalf("valid derived network defaults failed: %v\n%s", err, problems)
	}

	invalid := `tosca_definitions_version: tosca_simple_yaml_1_3
node_types:
  InvalidFlat:
    derived_from: tosca.nodes.network.Network
    properties:
      network_type:
        default: flat
topology_template:
  node_templates:
    network:
      type: InvalidFlat
`
	_, problems, err := testsupport.ParseSource(t, invalid)
	if err == nil {
		t.Fatal("derived flat network default bypassed physical_network requirement")
	}
	if !strings.Contains(problems, "physical_network") {
		t.Fatalf("derived network diagnostic is not specific:\n%s", problems)
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 8.5.1.1 Properties
// Requirement: TOSCA13-8.5.1.1-036
// Expected: structurally similar unrelated node types are not subject to the
// normative Network policy.
// Category: scope regression
func TestPartialNetworkRuleDoesNotApplyToUnrelatedNode(t *testing.T) {
	source := `tosca_definitions_version: tosca_simple_yaml_1_3
node_types:
  Unrelated:
    derived_from: tosca.nodes.Root
    properties:
      network_type:
        type: string
        required: false
      physical_network:
        type: string
        required: false
topology_template:
  node_templates:
    network:
      type: Unrelated
      properties:
        network_type: flat
`
	if _, problems, err := testsupport.ParseSource(t, source); err != nil {
		t.Fatalf("Network policy leaked to unrelated node: %v\n%s", err, problems)
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 8.5.1.1 Properties
// Requirement: TOSCA13-8.5.1.1-036
// Expected: an unresolved intrinsic network_type remains deferred instead of
// being treated as a concrete flat/vlan value.
// Category: function value, rendering
func TestPartialNetworkConditionalDefersFunctionValue(t *testing.T) {
	source := `tosca_definitions_version: tosca_simple_yaml_1_3
topology_template:
  inputs:
    kind:
      type: string
      default: gre
  node_templates:
    network:
      type: tosca.nodes.network.Network
      properties:
        network_type: { get_input: kind }
`
	if _, problems, err := testsupport.ParseSource(t, source); err != nil {
		t.Fatalf("function-valued network_type was validated prematurely: %v\n%s", err, problems)
	}
}

// Specification: TOSCA 1.3 section 8.5.1.1 only
// Expected: a structurally similar TOSCA 2.0 project node retains existing
// behavior because the version-specific hook is not installed.
// Category: cross-version regression
func TestPartialNetworkRuleDoesNotChangeTOSCA20(t *testing.T) {
	source := `tosca_definitions_version: tosca_2_0
node_types:
  Network:
    properties:
      network_type:
        type: string
        required: false
      physical_network:
        type: string
        required: false
service_template:
  node_templates:
    network:
      type: Network
      properties:
        network_type: flat
`
	if _, problems, err := testsupport.ParseSource(t, source); err != nil {
		t.Fatalf("TOSCA 1.3 Network policy changed TOSCA 2.0: %v\n%s", err, problems)
	}
}

func networkProfileTemplate(nodeType, properties string) string {
	return fmt.Sprintf(`tosca_definitions_version: tosca_simple_yaml_1_3
topology_template:
  node_templates:
    network:
      type: %s
      properties: %s
`, nodeType, properties)
}

func parseRejectedNetworkProfile(t *testing.T, nodeType, properties string) string {
	t.Helper()
	_, problems, err := testsupport.ParseSource(t, networkProfileTemplate(nodeType, properties))
	if err == nil {
		t.Fatalf("invalid %s properties %s were accepted", nodeType, properties)
	}
	return problems
}
