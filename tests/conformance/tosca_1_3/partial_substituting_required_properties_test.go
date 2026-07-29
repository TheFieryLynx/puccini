package tosca_1_3_test

import (
	"strings"
	"testing"

	"github.com/tliron/go-puccini/normal"
	"github.com/tliron/go-puccini/tests/conformance/internal/testsupport"
)

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.8.8.3, Notes
// Requirement: TOSCA13-3.8.8.3-005
// Expected: every required property of every node template in a substituting
// topology has a valid assignment.
// Category: positive, assignment, normalization, direct
func TestPartialSubstitutingTemplateRequiredPropertiesAssigned(t *testing.T) {
	serviceTemplate, problems, err := testsupport.ParseSource(
		t,
		substitutingRequiredPropertyTemplate("required_setting: configured", false, false),
	)
	if err != nil {
		t.Fatalf("valid required property assignment was rejected: %v\n%s", err, problems)
	}
	node := serviceTemplate.NodeTemplates["implementation"]
	if node == nil {
		t.Fatal("internal node template is absent from normalized output")
	}
	value, ok := node.Properties["required_setting"].(*normal.Primitive)
	if !ok || value.Primitive != "configured" {
		t.Fatalf("normalized required property = %#v, want configured", node.Properties["required_setting"])
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.8.8.3, Notes
// Requirement: TOSCA13-3.8.8.3-005
// Expected: a missing required property makes the substituting template
// invalid during rendering.
// Category: negative, assignment, parser phase, direct
func TestPartialSubstitutingTemplateRejectsMissingRequiredProperty(t *testing.T) {
	problems := rejectSubstitutingRequiredProperty(
		t,
		substitutingRequiredPropertyTemplate("", false, false),
	)
	for _, fragment := range []string{
		`node_templates["implementation"].properties["required_setting"]`,
		"required property",
	} {
		if !strings.Contains(problems, fragment) {
			t.Fatalf("missing-property diagnostic does not contain %q:\n%s", fragment, problems)
		}
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.8.8.3, Notes
// Requirement: TOSCA13-3.8.8.3-005
// Expected: a property default is an already-defined valid assignment, while
// an absent optional property does not invalidate the template.
// Category: positive, default, optional boundary, direct
func TestPartialSubstitutingTemplateDefaultsAndOptionalProperties(t *testing.T) {
	serviceTemplate, problems, err := testsupport.ParseSource(
		t,
		substitutingRequiredPropertyTemplate("", true, false),
	)
	if err != nil {
		t.Fatalf("defaulted required property was rejected: %v\n%s", err, problems)
	}
	node := serviceTemplate.NodeTemplates["implementation"]
	value, ok := node.Properties["required_setting"].(*normal.Primitive)
	if !ok || value.Primitive != "from default" {
		t.Fatalf("defaulted required property = %#v", node.Properties["required_setting"])
	}
	if _, present := node.Properties["optional_setting"]; present {
		t.Fatal("absent optional property acquired a false assignment")
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.8.8.3, Notes
// Requirement: TOSCA13-3.8.8.3-005
// Expected: effective inherited property definitions participate in the
// all-node-template assignment check.
// Category: negative, inheritance, direct
func TestPartialSubstitutingTemplateInheritedRequiredProperty(t *testing.T) {
	problems := rejectSubstitutingRequiredProperty(
		t,
		substitutingRequiredPropertyTemplate("", false, true),
	)
	if !strings.Contains(problems, "required_setting") || !strings.Contains(problems, "required property") {
		t.Fatalf("inherited required property failed for the wrong reason:\n%s", problems)
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.8.8.3, Notes
// Requirement: TOSCA13-3.8.8.3-005
// Expected: every internal node is rendered, including a disconnected node.
// Category: negative, all-node coverage, direct
func TestPartialSubstitutingTemplateChecksEveryInternalNode(t *testing.T) {
	source := strings.Replace(
		substitutingRequiredPropertyTemplate("required_setting: configured", false, false),
		"  substitution_mappings:",
		`    disconnected:
      type: ImplementationNode
  substitution_mappings:`,
		1,
	)
	problems := rejectSubstitutingRequiredProperty(t, source)
	if !strings.Contains(problems, `node_templates["disconnected"].properties["required_setting"]`) {
		t.Fatalf("disconnected node was not checked:\n%s", problems)
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.8.8.3, Notes
// Requirement: TOSCA13-3.8.8.3-005
// Expected: a syntactically and semantically valid intrinsic value is a valid
// property assignment and remains a function through normalization.
// Category: positive, function assignment, normalization, direct
func TestPartialSubstitutingTemplateFunctionAssignmentIsValid(t *testing.T) {
	serviceTemplate, problems, err := testsupport.ParseSource(
		t,
		substitutingRequiredPropertyTemplate("required_setting: { get_input: setting }", false, false),
	)
	if err != nil {
		t.Fatalf("function-valued property assignment was rejected: %v\n%s", err, problems)
	}
	if _, ok := serviceTemplate.NodeTemplates["implementation"].Properties["required_setting"].(*normal.FunctionCall); !ok {
		t.Fatalf(
			"function-valued required property normalized as %T",
			serviceTemplate.NodeTemplates["implementation"].Properties["required_setting"],
		)
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.8.8.3, Notes
// Requirement: TOSCA13-3.8.8.3-005
// Expected: repeated processing reports the same missing assignment path.
// Category: negative, determinism, direct
func TestPartialSubstitutingRequiredPropertyDiagnosticIsDeterministic(t *testing.T) {
	source := substitutingRequiredPropertyTemplate("", false, false)
	first := substitutingRequiredPropertyDiagnosticBody(rejectSubstitutingRequiredProperty(t, source))
	second := substitutingRequiredPropertyDiagnosticBody(rejectSubstitutingRequiredProperty(t, source))
	if first != second {
		t.Fatalf("required-property diagnostic is nondeterministic:\nfirst:\n%s\nsecond:\n%s", first, second)
	}
}

func substitutingRequiredPropertyTemplate(assignment string, withDefault bool, inherited bool) string {
	baseType := ""
	implementationType := ""
	defaultValue := ""
	if withDefault {
		defaultValue = "\n        default: from default"
	}
	if inherited {
		baseType = `  ImplementationBase:
    derived_from: tosca.nodes.Root
    properties:
      required_setting:
        type: string
`
		implementationType = `  ImplementationNode:
    derived_from: ImplementationBase
    properties:
      optional_setting:
        type: string
        required: false
`
	} else {
		implementationType = `  ImplementationNode:
    derived_from: tosca.nodes.Root
    properties:
      required_setting:
        type: string` + defaultValue + `
      optional_setting:
        type: string
        required: false
`
	}
	propertyAssignment := ""
	if assignment != "" {
		propertyAssignment = "\n      properties:\n        " + assignment
	}

	return `tosca_definitions_version: tosca_simple_yaml_1_3
node_types:
  AbstractNode:
    derived_from: tosca.nodes.Root
` + baseType + implementationType + `
topology_template:
  inputs:
    setting:
      type: string
      default: from input
  node_templates:
    implementation:
      type: ImplementationNode` + propertyAssignment + `
  substitution_mappings:
    node_type: AbstractNode
`
}

func rejectSubstitutingRequiredProperty(t *testing.T, source string) string {
	t.Helper()
	_, problems, err := testsupport.ParseSource(t, source)
	if err == nil {
		t.Fatal("substituting template with a missing required property was accepted")
	}
	return problems
}

func substitutingRequiredPropertyDiagnosticBody(problems string) string {
	if index := strings.Index(problems, "topology_template"); index >= 0 {
		return problems[index:]
	}
	return problems
}
