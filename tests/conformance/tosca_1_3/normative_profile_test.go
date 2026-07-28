package tosca_1_3_test

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/tliron/go-puccini/tests/conformance/internal/testsupport"
	"github.com/tliron/go-puccini/tosca/grammars/tosca_v1_3"
	"github.com/tliron/go-puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/go-puccini/tosca/parser"
)

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 8.5.2.3 Definition
// Requirement: TOSCA13-8.5.2.3-008
// Expected: the ordinary processor path loads the effective normative Port
// definition with order.required=true and preserves its related type, default,
// constraint, attribute, requirements, and parent.
// Category: positive, normative profile, effective definition
func TestNormativeNetworkPortEffectiveDefinition(t *testing.T) {
	parserContext, release := parseNormativeProfile(t)
	defer release()
	port := findProfileNodeType(t, parserContext, "tosca.nodes.network.Port")

	if port.ParentName == nil || *port.ParentName != "tosca.nodes.Root" {
		t.Fatalf("Port derived_from = %v", port.ParentName)
	}
	for _, name := range []string{"ip_address", "order", "is_default", "ip_range_start", "ip_range_end"} {
		if port.PropertyDefinitions[name] == nil {
			t.Fatalf("Port property %q is absent", name)
		}
	}

	order := port.PropertyDefinitions["order"]
	if !order.IsRequired() || order.Required == nil || !*order.Required {
		t.Fatalf("Port.order required = %v (effective %t), want explicit true", order.Required, order.IsRequired())
	}
	if order.DataTypeName == nil || *order.DataTypeName != "integer" {
		t.Fatalf("Port.order type = %v", order.DataTypeName)
	}
	if order.Default == nil || fmt.Sprint(order.Default.Context.Data) != "0" {
		t.Fatalf("Port.order default = %#v", order.Default)
	}
	if order.ValidationClause == nil {
		t.Fatal("Port.order greater_or_equal constraint is absent")
	}

	if port.AttributeDefinitions["ip_address"] == nil {
		t.Fatal("Port ip_address attribute is absent")
	}
	for name, capability := range map[string]string{
		"link":    "tosca.capabilities.network.Linkable",
		"binding": "tosca.capabilities.network.Bindable",
	} {
		requirement := port.RequirementDefinitions[name]
		if requirement == nil || requirement.TargetCapabilityTypeName == nil ||
			*requirement.TargetCapabilityTypeName != capability {
			t.Fatalf("Port requirement %q capability = %#v", name, requirement)
		}
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 8.5.2.1 Properties; 8.5.2.3 Definition
// Expected: the required property default supplies zero and explicit zero is
// accepted at the normative greater_or_equal boundary.
// Category: positive, boundary, rendering
func TestNormativeNetworkPortAssignments(t *testing.T) {
	for _, test := range []struct {
		name       string
		properties string
	}{
		{name: "default", properties: ""},
		{name: "explicit zero", properties: "\n      properties:\n        order: 0"},
	} {
		t.Run(test.name, func(t *testing.T) {
			serviceTemplate, parseProblems, err := testsupport.ParseSource(t, normativePortTemplate(test.properties))
			if err != nil {
				t.Fatalf("valid Port assignment failed: %v\n%s", err, parseProblems)
			}
			node := serviceTemplate.NodeTemplates["port"]
			if node == nil {
				t.Fatal("Port node did not normalize")
			}
			assertNormalizedPrimitive(t, node.Properties["order"], 0)
		})
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 3.6.10.6 Refining Property Definitions; 8.5.2.3 Definition
// Expected: a derived node type cannot refine the required Port.order property
// into an optional property.
// Category: negative, inheritance, refinement
func TestNormativeNetworkPortRequirednessRefinement(t *testing.T) {
	_, parseProblems, err := testsupport.ParseSource(t, `tosca_definitions_version: tosca_simple_yaml_1_3
node_types:
  InvalidOptionalPort:
    derived_from: tosca.nodes.network.Port
    properties:
      order:
        required: false
topology_template: {}
`)
	if err == nil {
		t.Fatal("derived type made the normative required Port.order property optional")
	}
	for _, expected := range []string{`properties["order"].required`, "cannot refine true to false"} {
		if !strings.Contains(parseProblems, expected) {
			t.Fatalf("requiredness diagnostic does not contain %q:\n%s", expected, parseProblems)
		}
	}
}

// Specification: TOSCA 1.3 section 8.5.2.3 only
// Expected: correcting the bundled 1.3 profile does not change TOSCA 2.0
// property requiredness or grammar behavior.
// Category: regression, version isolation
func TestNormativeProfileTosca20Isolation(t *testing.T) {
	_, parseProblems, err := testsupport.ParseSource(t, `tosca_definitions_version: tosca_2_0
imports:
  - profile: org.oasis-open.simple:2.0
    namespace: tosca
node_types:
  Version20Optional:
    properties:
      order:
        type: integer
        required: false
        default: 0
service_template:
  node_templates:
    node:
      type: Version20Optional
`)
	if err != nil {
		t.Fatalf("TOSCA 1.3 profile correction changed TOSCA 2.0 behavior: %v\n%s", err, parseProblems)
	}
}

func parseNormativeProfile(t *testing.T) (*parser.Context, func()) {
	t.Helper()
	parserContext, urlContext := reflectionParserContext(t, `tosca_definitions_version: tosca_simple_yaml_1_3
topology_template: {}
`)
	if _, err := parserContext.Parse(context.Background()); err != nil {
		problems := ""
		if parserContext.Root != nil {
			problems = parserContext.GetProblems().ToString(false)
		}
		urlContext.Release()
		t.Fatalf("ordinary TOSCA 1.3 processor path failed: %v\n%s", err, problems)
	}
	return parserContext, func() {
		_ = urlContext.Release()
	}
}

func findProfileNodeType(t *testing.T, parserContext *parser.Context, name string) *tosca_v2_0.NodeType {
	t.Helper()
	var files []string
	for _, file := range parserContext.Files {
		files = append(files, file.GetContext().URL.String())
		if profile, ok := file.EntityPtr.(*tosca_v1_3.File); ok {
			for _, nodeType := range profile.NodeTypes {
				if nodeType.Name == name {
					return nodeType
				}
			}
		}
	}
	sort.Strings(files)
	t.Fatalf("normative node type %q not found in processor files: %v", name, files)
	return nil
}

func normativePortTemplate(properties string) string {
	return `tosca_definitions_version: tosca_simple_yaml_1_3
topology_template:
  node_templates:
    port:
      type: tosca.nodes.network.Port` + properties + "\n"
}
