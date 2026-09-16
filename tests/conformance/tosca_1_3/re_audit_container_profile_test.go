package tosca_1_3_test

import (
	"strings"
	"testing"

	v13 "github.com/tliron/go-puccini/tosca/grammars/tosca_v1_3"
)

// TOSCA 1.3 §5.9.12.1 Definition and §5.9.13.1 Definition, F07/F08.
// Positive effective inherited definitions and assignments; negative compatibility
// and case-sensitive lookup. Runs the ordinary implicit-profile parser path.
func TestReAuditContainerProfile(t *testing.T) {
	source := `tosca_definitions_version: tosca_simple_yaml_1_3
node_types:
  RuntimeChild: { derived_from: tosca.nodes.Container.Runtime }
  ApplicationChild: { derived_from: tosca.nodes.Container.Application }
`
	r := reAuditParse(t, source)
	if r.Phase != "" {
		t.Fatalf("%s: %s", r.Phase, r.Problems)
	}
	file := r.Context.Root.EntityPtr.(*v13.ServiceFile)
	for _, typ := range file.NodeTypes {
		if typ.Name == "RuntimeChild" {
			cap := typ.CapabilityDefinitions["host"]
			if cap == nil || cap.CapabilityType == nil || cap.CapabilityType.Name != "tosca.capabilities.Compute" {
				t.Errorf("effective runtime host: %#v", cap)
			}
			if len(cap.ValidSourceNodeTypes) != 1 || cap.ValidSourceNodeTypes[0].Name != "tosca.nodes.Container.Application" {
				t.Fatalf("host valid source types: %#v", cap.ValidSourceNodeTypes)
			}
		}
		if typ.Name == "ApplicationChild" {
			req := typ.RequirementDefinitions["network"]
			if req == nil || req.TargetCapabilityType == nil || req.TargetCapabilityType.Name != "tosca.capabilities.Endpoint" {
				t.Errorf("effective application network: %#v", req)
			}
		}
	}

	for _, prefix := range []string{"", "tosca:"} {
		t.Run("lookup-"+prefix, func(t *testing.T) {
			named := strings.ReplaceAll(source, "tosca.nodes.Container.", prefix+"Container.")
			result := reAuditParse(t, named)
			if result.Phase != "" {
				t.Fatalf("%s: %s", result.Phase, result.Problems)
			}
			types := result.Context.Root.EntityPtr.(*v13.ServiceFile).NodeTypes
			for _, typ := range types {
				if typ.Name == "RuntimeChild" && typ.CapabilityDefinitions["host"].CapabilityType.Name != "tosca.capabilities.Compute" {
					t.Fatal("effective shorthand runtime host mismatch")
				}
			}
		})
	}
	reAuditReject(t, strings.Replace(source, "tosca.nodes.Container.Runtime", "tosca.nodes.Container.runtime", 1), "namespace", "Container.runtime")
}
func TestReAuditContainerAssignments(t *testing.T) {
	source := `tosca_definitions_version: tosca_simple_yaml_1_3
node_types:
  App: { derived_from: tosca.nodes.Container.Application }
  Runtime: { derived_from: tosca.nodes.Container.Runtime }
  Endpoint:
    derived_from: tosca.nodes.Root
    capabilities:
      network: { type: tosca.capabilities.Endpoint }
topology_template:
  node_templates:
    vm: { type: tosca.nodes.Compute }
    runtime:
      type: Runtime
      requirements: [{ host: vm }]
    storage:
      type: tosca.nodes.Storage.BlockStorage
      properties: { name: disk, size: 1 GiB }
    endpoint: { type: Endpoint }
    app:
      type: App
      requirements: [{ host: runtime }, { storage: storage }, { network: endpoint }]
`
	r := reAuditParse(t, source)
	if r.Phase != "" {
		t.Fatalf("%s: %s", r.Phase, r.Problems)
	}
	app := r.Template.NodeTemplates["app"]
	if len(app.Requirements) != 3 {
		t.Fatalf("normalized requirements: %#v", app.Requirements)
	}
	found := false
	for _, req := range app.Requirements {
		if req.Name == "network" {
			found = true
			if req.NodeTemplate == nil || req.NodeTemplate.Name != "endpoint" {
				t.Fatalf("network target: %#v", req)
			}
		}
	}
	if !found {
		t.Fatal("network requirement lost")
	}
	reAuditReject(t, strings.Replace(source, "App: { derived_from: tosca.nodes.Container.Application }", "App: { derived_from: tosca.nodes.Container.Application, requirements: [{network: {capability: tosca.capabilities.Compute}}] }", 1), "inheritance", "must be derived from", "Endpoint")
}
