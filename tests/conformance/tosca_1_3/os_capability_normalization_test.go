package tosca_1_3_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tliron/exturl"
	"github.com/tliron/go-ard"
	"github.com/tliron/go-kutil/problems"
	"github.com/tliron/go-kutil/terminal"
	cloutjs "github.com/tliron/go-puccini/clout/js"
	cloututil "github.com/tliron/go-puccini/clout/util"
	"github.com/tliron/go-puccini/normal"
	"github.com/tliron/go-puccini/tests/conformance/internal/testsupport"
	"github.com/tliron/go-puccini/tosca/parser"
)

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 5.5.12.3, Additional Requirements
// Requirement: TOSCA13-5.5.12.3-001
// Expected: architecture, type, and distribution string values of the
// normative OperatingSystem capability normalize to lowercase.
// Category: positive, direct assignment, mixed case, boundary, normalization
func TestOperatingSystemCapabilityLowercaseNormalization(t *testing.T) {
	tests := []struct {
		name   string
		values map[string]string
		want   map[string]string
	}{
		{
			name: "uppercase",
			values: map[string]string{
				"architecture": "X86_64",
				"type":         "LINUX",
				"distribution": "UBUNTU",
			},
			want: map[string]string{
				"architecture": "x86_64",
				"type":         "linux",
				"distribution": "ubuntu",
			},
		},
		{
			name: "mixed-case",
			values: map[string]string{
				"architecture": "AaRcH64",
				"type":         "LiNuX",
				"distribution": "OpenSUSE",
			},
			want: map[string]string{
				"architecture": "aarch64",
				"type":         "linux",
				"distribution": "opensuse",
			},
		},
		{
			name: "already-lowercase",
			values: map[string]string{
				"architecture": "arm64",
				"type":         "linux",
				"distribution": "debian",
			},
			want: map[string]string{
				"architecture": "arm64",
				"type":         "linux",
				"distribution": "debian",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			serviceTemplate, parseProblems, err := testsupport.ParseSource(t, operatingSystemAssignmentTemplate(test.values))
			if err != nil {
				t.Fatalf("valid OperatingSystem assignment failed: %v\n%s", err, parseProblems)
			}
			properties := normalizedCapabilityProperties(t, serviceTemplate, "compute", "os")
			for name, want := range test.want {
				if got := normalizedString(t, properties[name]); got != want {
					t.Errorf("%s = %q, want %q", name, got, want)
				}
			}
		})
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 5.5.12.3, Additional Requirements
// Requirement: TOSCA13-5.5.12.3-001
// Expected: rendered default values, including defaults inherited through a
// derived capability type, are normalized by the same rule.
// Category: positive, default, inheritance, normalization
func TestOperatingSystemCapabilityDefaultsAndInheritance(t *testing.T) {
	tests := []struct {
		name           string
		capabilityType string
		want           string
	}{
		{"default", "DefaultOS", "freebsd"},
		{"inherited", "InheritedOS", "freebsd"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			source := `tosca_definitions_version: tosca_simple_yaml_1_3
capability_types:
  DefaultOS:
    derived_from: tosca.capabilities.OperatingSystem
    properties:
      distribution:
        type: string
        required: false
        default: FreeBSD
  InheritedOS:
    derived_from: DefaultOS
node_types:
  TestCompute:
    derived_from: tosca.nodes.Root
    capabilities:
      os:
        type: ` + test.capabilityType + `
topology_template:
  node_templates:
    compute:
      type: TestCompute
`
			serviceTemplate, parseProblems, err := testsupport.ParseSource(t, source)
			if err != nil {
				t.Fatalf("valid default/inheritance template failed: %v\n%s", err, parseProblems)
			}
			properties := normalizedCapabilityProperties(t, serviceTemplate, "compute", "os")
			if got := normalizedString(t, properties["distribution"]); got != test.want {
				t.Fatalf("distribution = %q, want %q", got, test.want)
			}
		})
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 5.5.12.3, Additional Requirements
// Requirement: TOSCA13-5.5.12.3-001
// Expected: normalization is restricted to the three named properties of
// OperatingSystem and types derived from it.
// Category: negative regression, scope, semantic preservation
func TestOperatingSystemCapabilityNormalizationScope(t *testing.T) {
	source := `tosca_definitions_version: tosca_simple_yaml_1_3
capability_types:
  ExtendedOS:
    derived_from: tosca.capabilities.OperatingSystem
    properties:
      label:
        type: string
        required: false
  Unrelated:
    derived_from: tosca.capabilities.Root
    properties:
      architecture:
        type: string
        required: false
node_types:
  TestNode:
    derived_from: tosca.nodes.Root
    capabilities:
      os:
        type: ExtendedOS
      unrelated:
        type: Unrelated
      endpoint:
        type: tosca.capabilities.Endpoint
topology_template:
  node_templates:
    node:
      type: TestNode
      capabilities:
        os:
          properties:
            type: LiNuX
            label: PreserveCase
        unrelated:
          properties:
            architecture: PreserveCase
        endpoint:
          properties:
            protocol: HTTPS
`
	serviceTemplate, parseProblems, err := testsupport.ParseSource(t, source)
	if err != nil {
		t.Fatalf("valid scope template failed: %v\n%s", err, parseProblems)
	}

	osProperties := normalizedCapabilityProperties(t, serviceTemplate, "node", "os")
	if got := normalizedString(t, osProperties["type"]); got != "linux" {
		t.Errorf("OperatingSystem type = %q, want %q", got, "linux")
	}
	if got := normalizedString(t, osProperties["label"]); got != "PreserveCase" {
		t.Errorf("unrelated OperatingSystem property changed to %q", got)
	}
	unrelatedProperties := normalizedCapabilityProperties(t, serviceTemplate, "node", "unrelated")
	if got := normalizedString(t, unrelatedProperties["architecture"]); got != "PreserveCase" {
		t.Errorf("unrelated capability property changed to %q", got)
	}
	endpointProperties := normalizedCapabilityProperties(t, serviceTemplate, "node", "endpoint")
	if got := normalizedString(t, endpointProperties["protocol"]); got != "HTTPS" {
		t.Errorf("unrelated normative capability property changed to %q", got)
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 3.6.10 Property definition; 5.5.12.1 Properties
// Requirement: TOSCA13-5.5.12.3-001
// Expected: a non-string assignment is rejected during rendering with a
// diagnostic identifying the OperatingSystem property and expected type.
// Category: negative, YAML type, parser phase
func TestOperatingSystemCapabilityLowercaseInvalidType(t *testing.T) {
	source := `tosca_definitions_version: tosca_simple_yaml_1_3
topology_template:
  node_templates:
    compute:
      type: tosca.nodes.Compute
      capabilities:
        os:
          properties:
            architecture: [x86_64]
`
	_, parseProblems, err := testsupport.ParseSource(t, source)
	if err == nil {
		t.Fatal("non-string OperatingSystem architecture was accepted")
	}
	if err.Error() != "parsing problems" {
		t.Fatalf("wrong parser phase: %v\n%s", err, parseProblems)
	}
	for _, expected := range []string{"architecture", "string"} {
		if !strings.Contains(parseProblems, expected) {
			t.Fatalf("diagnostic does not identify %q:\n%s", expected, parseProblems)
		}
	}
	if strings.Contains(parseProblems, "unsupported keyname") || strings.Contains(parseProblems, "unknown data type") {
		t.Fatalf("an unrelated earlier error masked the value type error:\n%s", parseProblems)
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 4.3.1 concat; 5.5.12.3 Additional Requirements
// Requirement: TOSCA13-5.5.12.3-001
// Expected: the normalized model preserves the intrinsic function object; its
// string result is lowercased only after function evaluation.
// Category: positive, function value, evaluation, normalization
func TestOperatingSystemCapabilityFunctionValue(t *testing.T) {
	source := operatingSystemFunctionTemplate()
	serviceTemplate, parseProblems, err := testsupport.ParseSource(t, source)
	if err != nil {
		t.Fatalf("valid function-valued assignment failed: %v\n%s", err, parseProblems)
	}

	properties := normalizedCapabilityProperties(t, serviceTemplate, "compute", "os")
	functionValue, ok := properties["architecture"].(*normal.FunctionCall)
	if !ok {
		t.Fatalf("architecture has type %T, want *normal.FunctionCall", properties["architecture"])
	}
	if functionValue.FunctionCall == nil || functionValue.FunctionCall.Name != "tosca.function.concat" {
		t.Fatalf("function representation was damaged: %#v", functionValue.FunctionCall)
	}

	clout, err := serviceTemplate.Compile()
	if err != nil {
		t.Fatalf("compile failed: %v", err)
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
		t.Fatalf("function evaluation failed:\n%s", evaluationProblems.ToString(false))
	}

	nodes := cloututil.GetToscaNodeTemplates(clout, "")
	if len(nodes) != 1 {
		t.Fatalf("compiled node count = %d, want 1", len(nodes))
	}
	for _, node := range nodes {
		capabilities := cloututil.GetToscaCapabilities(node, "")
		osCapability, ok := capabilities["os"].(ard.StringMap)
		if !ok {
			t.Fatalf("compiled OperatingSystem capability has type %T", capabilities["os"])
		}
		compiledProperties, ok := osCapability["properties"].(ard.StringMap)
		if !ok {
			t.Fatalf("compiled OperatingSystem properties have type %T", osCapability["properties"])
		}
		if got := compiledProperties["architecture"]; got != "x86_64" {
			t.Fatalf("evaluated architecture = %#v, want %q", got, "x86_64")
		}
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 5.5.12.3, Additional Requirements
// Requirement: TOSCA13-5.5.12.3-001
// Expected: repeated normalization is idempotent and deterministic.
// Category: boundary, idempotence, deterministic normalization
func TestOperatingSystemCapabilityNormalizationIdempotenceAndDeterminism(t *testing.T) {
	source := operatingSystemAssignmentTemplate(map[string]string{
		"architecture": "X86_64",
		"type":         "LiNuX",
		"distribution": "OpenSUSE",
	})
	parserContext, urlContext := parseContextForSource(t, source)
	defer urlContext.Release()

	first, err := parserContext.Parse(context.Background())
	if err != nil {
		t.Fatalf("initial parse failed: %v\n%s", err, parserContext.GetProblems().ToString(false))
	}
	second, ok := parserContext.Normalize()
	if !ok {
		t.Fatal("second normalization failed")
	}
	third, ok := parserContext.Normalize()
	if !ok {
		t.Fatal("third normalization failed")
	}

	firstJSON := normalizedCapabilityJSON(t, first, "compute", "os")
	secondJSON := normalizedCapabilityJSON(t, second, "compute", "os")
	thirdJSON := normalizedCapabilityJSON(t, third, "compute", "os")
	if firstJSON != secondJSON || secondJSON != thirdJSON {
		t.Fatalf("normalization is not idempotent/deterministic:\nfirst:  %s\nsecond: %s\nthird:  %s", firstJSON, secondJSON, thirdJSON)
	}
}

func operatingSystemAssignmentTemplate(values map[string]string) string {
	var properties strings.Builder
	for _, name := range []string{"architecture", "type", "distribution"} {
		if value, ok := values[name]; ok {
			properties.WriteString("            " + name + ": " + value + "\n")
		}
	}
	return `tosca_definitions_version: tosca_simple_yaml_1_3
topology_template:
  node_templates:
    compute:
      type: tosca.nodes.Compute
      capabilities:
        os:
          properties:
` + properties.String()
}

func operatingSystemFunctionTemplate() string {
	return `tosca_definitions_version: tosca_simple_yaml_1_3
topology_template:
  node_templates:
    compute:
      type: tosca.nodes.Compute
      capabilities:
        os:
          properties:
            architecture: { concat: [X86, _64] }
`
}

func normalizedCapabilityProperties(t *testing.T, serviceTemplate *normal.ServiceTemplate, nodeName string, capabilityName string) normal.Values {
	t.Helper()
	nodeTemplate := serviceTemplate.NodeTemplates[nodeName]
	if nodeTemplate == nil {
		t.Fatalf("normalized node template %q is absent", nodeName)
	}
	capability := nodeTemplate.Capabilities[capabilityName]
	if capability == nil {
		t.Fatalf("normalized capability %q is absent", capabilityName)
	}
	return capability.Properties
}

func normalizedString(t *testing.T, value normal.Value) string {
	t.Helper()
	primitive, ok := value.(*normal.Primitive)
	if !ok {
		t.Fatalf("normalized value has type %T, want *normal.Primitive", value)
	}
	stringValue, ok := primitive.Primitive.(string)
	if !ok {
		t.Fatalf("normalized primitive has type %T, want string", primitive.Primitive)
	}
	return stringValue
}

func normalizedCapabilityJSON(t *testing.T, serviceTemplate *normal.ServiceTemplate, nodeName string, capabilityName string) string {
	t.Helper()
	nodeTemplate := serviceTemplate.NodeTemplates[nodeName]
	if nodeTemplate == nil || nodeTemplate.Capabilities[capabilityName] == nil {
		t.Fatalf("normalized capability %s.%s is absent", nodeName, capabilityName)
	}
	data, err := json.Marshal(nodeTemplate.Capabilities[capabilityName])
	if err != nil {
		t.Fatalf("marshal normalized capability: %v", err)
	}
	return string(data)
}

func parseContextForSource(t *testing.T, source string) (*parser.Context, *exturl.Context) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "service-template.yaml")
	if err := os.WriteFile(path, []byte(source), 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	urlContext := exturl.NewContext()
	parserContext := parser.NewParser().NewContext()
	parserContext.URL = urlContext.NewFileURL(filepath.ToSlash(path))
	return parserContext, urlContext
}
