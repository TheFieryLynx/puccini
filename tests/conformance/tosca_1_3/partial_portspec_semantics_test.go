package tosca_1_3_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/tliron/go-puccini/normal"
	"github.com/tliron/go-puccini/tests/conformance/internal/testsupport"
)

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 5.3.11.3 Additional requirements
// Requirements: TOSCA13-5.3.11.3-001, TOSCA13-5.3.11.3-002,
// TOSCA13-5.3.11.3-003
// Expected: every PortSpec selects a source or target port and any declared
// source/target range contains its corresponding concrete port.
// Category: positive, boundary, rendering, normalization
func TestPartialPortSpecCrossPropertyRulesAccepted(t *testing.T) {
	tests := []struct {
		name       string
		spec       string
		wantField  string
		wantNumber int64
	}{
		{"target only", "{ target: 80 }", "target", 80},
		{"source only", "{ source: 53 }", "source", 53},
		{"target lower boundary", "{ target: 80, target_range: [ 80, 90 ] }", "target", 80},
		{"target upper boundary", "{ target: 90, target_range: [ 80, 90 ] }", "target", 90},
		{"source lower boundary", "{ source: 1000, source_range: [ 1000, 2000 ] }", "source", 1000},
		{"source upper boundary", "{ source: 2000, source_range: [ 1000, 2000 ] }", "source", 2000},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			serviceTemplate, parseProblems, err := testsupport.ParseSource(t, portSpecTemplate("tosca.datatypes.network.PortSpec", test.spec))
			if err != nil {
				t.Fatalf("valid PortSpec failed: %v\n%s", err, parseProblems)
			}
			spec, ok := serviceTemplate.NodeTemplates["node"].Properties["spec"].(*normal.Map)
			if !ok {
				t.Fatalf("normalized PortSpec has type %T", serviceTemplate.NodeTemplates["node"].Properties["spec"])
			}
			assertNormalizedMapPrimitive(t, spec, test.wantField, test.wantNumber)
		})
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 5.3.11.3 Additional requirements
// Requirement: TOSCA13-5.3.11.3-001
// Expected: a PortSpec with none of target, target_range, source, or
// source_range is rejected during value rendering.
// Category: negative, rendering
func TestPartialPortSpecRejectsNoPortFields(t *testing.T) {
	for _, spec := range []string{"{}", "{ protocol: tcp }"} {
		problems := parseRejectedPortSpec(t, "tosca.datatypes.network.PortSpec", spec)
		for _, expected := range []string{`properties["spec"]`, "PortSpec", "target"} {
			if !strings.Contains(problems, expected) {
				t.Fatalf("PortSpec diagnostic does not contain %q:\n%s", expected, problems)
			}
		}
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 5.3.11.3 Additional requirements
// Requirement: TOSCA13-5.3.11.3-002
// Expected: source_range requires source, and source is inclusively bounded by
// that range.
// Category: negative, boundary, rendering
func TestPartialPortSpecRejectsInvalidSourceRangePair(t *testing.T) {
	for _, spec := range []string{
		"{ source_range: [ 1000, 2000 ] }",
		"{ source: 999, source_range: [ 1000, 2000 ] }",
		"{ source: 2001, source_range: [ 1000, 2000 ] }",
	} {
		problems := parseRejectedPortSpec(t, "tosca.datatypes.network.PortSpec", spec)
		for _, expected := range []string{`properties["spec"]`, "source", "source_range"} {
			if !strings.Contains(problems, expected) {
				t.Fatalf("source-range diagnostic does not contain %q:\n%s", expected, problems)
			}
		}
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 5.3.11.3 Additional requirements
// Requirement: TOSCA13-5.3.11.3-003
// Expected: target_range requires target, and target is inclusively bounded by
// that range.
// Category: negative, boundary, rendering
func TestPartialPortSpecRejectsInvalidTargetRangePair(t *testing.T) {
	for _, spec := range []string{
		"{ target_range: [ 80, 90 ] }",
		"{ target: 79, target_range: [ 80, 90 ] }",
		"{ target: 91, target_range: [ 80, 90 ] }",
	} {
		problems := parseRejectedPortSpec(t, "tosca.datatypes.network.PortSpec", spec)
		for _, expected := range []string{`properties["spec"]`, "target", "target_range"} {
			if !strings.Contains(problems, expected) {
				t.Fatalf("target-range diagnostic does not contain %q:\n%s", expected, problems)
			}
		}
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 5.3.11.3 Additional requirements
// Expected: the PortSpec rules follow data-type inheritance but do not apply
// structurally to an unrelated type with the same property names.
// Category: positive, negative, inheritance, scope regression
func TestPartialPortSpecRulesFollowTypeIdentity(t *testing.T) {
	derived := `tosca_definitions_version: tosca_simple_yaml_1_3
data_types:
  DerivedPortSpec:
    derived_from: tosca.datatypes.network.PortSpec
node_types:
  Holder:
    derived_from: tosca.nodes.Root
    properties:
      spec:
        type: DerivedPortSpec
topology_template:
  node_templates:
    node:
      type: Holder
      properties:
        spec: { protocol: tcp }
`
	_, problems, err := testsupport.ParseSource(t, derived)
	if err == nil {
		t.Fatal("derived PortSpec bypassed the normative cross-property rule")
	}
	if !strings.Contains(problems, "PortSpec") {
		t.Fatalf("derived PortSpec diagnostic is not specific:\n%s", problems)
	}

	unrelated := `tosca_definitions_version: tosca_simple_yaml_1_3
data_types:
  Unrelated:
    properties:
      protocol:
        type: string
        required: false
      source_range:
        type: range
        required: false
node_types:
  Holder:
    derived_from: tosca.nodes.Root
    properties:
      spec:
        type: Unrelated
topology_template:
  node_templates:
    node:
      type: Holder
      properties:
        spec: { protocol: tcp }
`
	if _, unrelatedProblems, unrelatedErr := testsupport.ParseSource(t, unrelated); unrelatedErr != nil {
		t.Fatalf("PortSpec rule leaked to unrelated data type: %v\n%s", unrelatedErr, unrelatedProblems)
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 5.3.11.3 Additional requirements
// Requirements: TOSCA13-5.3.11.3-001, TOSCA13-5.3.11.3-002,
// TOSCA13-5.3.11.3-003
// Expected: PortSpec validation also runs for collection entry schemas, which
// is the normative Endpoint.ports representation.
// Category: positive, negative, nested rendering
func TestPartialPortSpecNestedMapEntryValidated(t *testing.T) {
	template := func(spec string) string {
		return fmt.Sprintf(`tosca_definitions_version: tosca_simple_yaml_1_3
node_types:
  Holder:
    derived_from: tosca.nodes.Root
    properties:
      specs:
        type: map
        entry_schema:
          type: tosca.datatypes.network.PortSpec
topology_template:
  node_templates:
    node:
      type: Holder
      properties:
        specs:
          first: %s
`, spec)
	}
	if _, problems, err := testsupport.ParseSource(t, template("{ target: 80 }")); err != nil {
		t.Fatalf("valid nested PortSpec failed: %v\n%s", err, problems)
	}
	_, problems, err := testsupport.ParseSource(t, template("{ protocol: tcp }"))
	if err == nil {
		t.Fatal("invalid nested PortSpec bypassed entry-schema validation")
	}
	for _, expected := range []string{`properties["specs"]["first"]`, "PortSpec"} {
		if !strings.Contains(problems, expected) {
			t.Fatalf("nested PortSpec diagnostic does not contain %q:\n%s", expected, problems)
		}
	}
}

// Specification: TOSCA 1.3 section 5.3.11.3 only
// Expected: the TOSCA 1.3 profile policy does not change the TOSCA 2.0
// rendering of a project-defined structurally similar data type.
// Category: cross-version regression
func TestPartialPortSpecRulesDoNotChangeTOSCA20(t *testing.T) {
	source := `tosca_definitions_version: tosca_2_0
data_types:
  PortSpec:
    properties:
      protocol:
        type: string
        required: false
node_types:
  Holder:
    properties:
      spec:
        type: PortSpec
service_template:
  node_templates:
    node:
      type: Holder
      properties:
        spec: { protocol: tcp }
`
	if _, problems, err := testsupport.ParseSource(t, source); err != nil {
		t.Fatalf("TOSCA 1.3 PortSpec policy changed TOSCA 2.0: %v\n%s", err, problems)
	}
}

func portSpecTemplate(dataType, spec string) string {
	return fmt.Sprintf(`tosca_definitions_version: tosca_simple_yaml_1_3
node_types:
  Holder:
    derived_from: tosca.nodes.Root
    properties:
      spec:
        type: %s
topology_template:
  node_templates:
    node:
      type: Holder
      properties:
        spec: %s
`, dataType, spec)
}

func parseRejectedPortSpec(t *testing.T, dataType, spec string) string {
	t.Helper()
	_, problems, err := testsupport.ParseSource(t, portSpecTemplate(dataType, spec))
	if err == nil {
		t.Fatalf("invalid PortSpec %s was accepted", spec)
	}
	return problems
}

func assertNormalizedMapPrimitive(t *testing.T, value *normal.Map, key string, want any) {
	t.Helper()
	for _, entry := range value.Entries {
		var entryKey any
		switch entry := entry.(type) {
		case *normal.Primitive:
			if keyValue, ok := entry.Key.(*normal.Primitive); ok {
				entryKey = keyValue.Primitive
			}
			if entryKey == key && fmt.Sprint(entry.Primitive) == fmt.Sprint(want) {
				return
			}
		}
	}
	t.Fatalf("normalized map does not contain %s=%v: %#v", key, want, value)
}
