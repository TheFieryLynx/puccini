package tosca_1_3_test

import (
	"strings"
	"testing"

	"github.com/tliron/go-puccini/tests/conformance/internal/testsupport"
)

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 8.2 Connectivity semantics
// Requirement: TOSCA13-8.2-002
// Expected: the effective ConnectsTo relationship targets Endpoint
// capabilities, and an explicit source requirement using ConnectsTo resolves
// to the target Endpoint and preserves the relationship type in normalization.
// Category: positive, normative profile, resolution, normalization, direct
func TestVerificationConnectsToEndpointSemantics(t *testing.T) {
	parserContext, release := parseNormativeProfile(t)
	defer release()

	connectsTo := findProfileRelationshipType(t, parserContext, "tosca.relationships.ConnectsTo")
	if connectsTo.ValidCapabilityTypeNames == nil ||
		len(*connectsTo.ValidCapabilityTypeNames) != 1 ||
		(*connectsTo.ValidCapabilityTypeNames)[0] != "tosca.capabilities.Endpoint" {
		t.Fatalf("ConnectsTo valid_target_types = %v, want [tosca.capabilities.Endpoint]",
			connectsTo.ValidCapabilityTypeNames)
	}
	if len(connectsTo.ValidCapabilityTypes) != 1 ||
		connectsTo.ValidCapabilityTypes[0].Name != "tosca.capabilities.Endpoint" {
		t.Fatalf("ConnectsTo effective target capability types = %#v, want Endpoint",
			connectsTo.ValidCapabilityTypes)
	}

	serviceTemplate, problems, err := testsupport.ParseSource(t, connectivityTemplate(
		"tosca.capabilities.Endpoint",
		"endpoint",
	))
	if err != nil {
		t.Fatalf("valid ConnectsTo requirement failed: %v\n%s", err, problems)
	}
	source := serviceTemplate.NodeTemplates["source"]
	if source == nil || len(source.Requirements) != 1 {
		t.Fatalf("normalized source requirements = %#v, want one", source)
	}
	requirement := source.Requirements[0]
	if requirement.NodeTemplate == nil || requirement.NodeTemplate.Name != "target" {
		t.Fatalf("ConnectsTo target node = %#v, want target", requirement.NodeTemplate)
	}
	if requirement.CapabilityTypeName == nil ||
		!strings.HasSuffix(*requirement.CapabilityTypeName, "::Endpoint") {
		t.Fatalf("ConnectsTo target capability type = %v, want Endpoint",
			requirement.CapabilityTypeName)
	}
	if requirement.Relationship == nil {
		t.Fatal("ConnectsTo relationship did not normalize")
	}
	hasConnectsTo := false
	for name := range requirement.Relationship.Types {
		if strings.HasSuffix(name, "::ConnectsTo") {
			hasConnectsTo = true
		}
	}
	if !hasConnectsTo {
		t.Fatalf("normalized relationship types = %#v, want ConnectsTo",
			requirement.Relationship.Types)
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 8.2 Connectivity semantics
// Requirement: TOSCA13-8.2-002
// Expected: ConnectsTo rejects a target capability that is not Endpoint or
// derived from Endpoint, at requirement-assignment rendering.
// Category: negative, resolution, diagnostic, direct
func TestVerificationConnectsToRejectsNonEndpoint(t *testing.T) {
	_, problems, err := testsupport.ParseSource(t, connectivityTemplate(
		"tosca.capabilities.Container",
		"container",
	))
	if err == nil {
		t.Fatal("ConnectsTo accepted a non-Endpoint target capability")
	}
	for _, expected := range []string{
		`requirements{0}.capability`,
		"tosca::Container",
		"tosca::ConnectsTo",
		"tosca.capabilities.Endpoint",
	} {
		if !strings.Contains(problems, expected) {
			t.Fatalf("ConnectsTo diagnostic does not contain %q:\n%s", expected, problems)
		}
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 8.3.2 Specifying layer 4 ports
// Requirement: TOSCA13-8.3.2-001
// Expected: Endpoint exposes a single PortDef and a map of PortSpec values;
// assigned layer-4 port values survive rendering and normalization.
// Category: positive, normative profile, assignment, normalization, direct
func TestVerificationEndpointLayer4PortSpecifications(t *testing.T) {
	parserContext, release := parseNormativeProfile(t)
	defer release()

	endpoint := findProfileCapabilityType(t, parserContext, "tosca.capabilities.Endpoint")
	port := endpoint.PropertyDefinitions["port"]
	if port == nil || port.DataTypeName == nil ||
		*port.DataTypeName != "tosca.datatypes.network.PortDef" {
		t.Fatalf("Endpoint.port = %#v, want PortDef", port)
	}
	ports := endpoint.PropertyDefinitions["ports"]
	if ports == nil || ports.DataTypeName == nil || *ports.DataTypeName != "map" ||
		ports.EntrySchema == nil || ports.EntrySchema.DataTypeName == nil ||
		*ports.EntrySchema.DataTypeName != "tosca.datatypes.network.PortSpec" {
		t.Fatalf("Endpoint.ports = %#v, want map of PortSpec", ports)
	}

	serviceTemplate, problems, err := testsupport.ParseSource(t, `tosca_definitions_version: tosca_simple_yaml_1_3
node_types:
  example.Holder:
    derived_from: tosca.nodes.Root
    capabilities:
      endpoint:
        type: tosca.capabilities.Endpoint
topology_template:
  node_templates:
    holder:
      type: example.Holder
      capabilities:
        endpoint:
          properties:
            ports:
              http:
                protocol: tcp
                source: 1024
                target: 8080
`)
	if err != nil {
		t.Fatalf("valid Endpoint port specification failed: %v\n%s", err, problems)
	}
	value := serviceTemplate.NodeTemplates["holder"].Capabilities["endpoint"].Properties["ports"]
	if value == nil {
		t.Fatal("Endpoint.ports did not normalize")
	}
}

func connectivityTemplate(capabilityType string, capabilityName string) string {
	return `tosca_definitions_version: tosca_simple_yaml_1_3
node_types:
  example.Target:
    derived_from: tosca.nodes.Root
    capabilities:
      ` + capabilityName + `:
        type: ` + capabilityType + `
  example.Source:
    derived_from: tosca.nodes.Root
    requirements:
      - connect:
          capability: ` + capabilityType + `
          node: example.Target
          relationship: tosca.relationships.ConnectsTo
topology_template:
  node_templates:
    target:
      type: example.Target
    source:
      type: example.Source
      requirements:
        - connect:
            node: target
            capability: ` + capabilityName + `
            relationship: tosca.relationships.ConnectsTo
`
}
