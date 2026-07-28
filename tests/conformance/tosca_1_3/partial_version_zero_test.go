package tosca_1_3_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/tliron/go-puccini/normal"
	"github.com/tliron/go-puccini/tests/conformance/internal/testsupport"
)

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.3.2.5, Additional Requirements
// Requirements: TOSCA13-3.3.2.5-001
// Expected: 0, 0.0, and 0.0.0 are accepted as the same unspecified-version sentinel.
// Category: positive, boundary, normalization, direct
func TestVersionZeroMeansUnspecified(t *testing.T) {
	for _, spelling := range []string{"0", "0.0", "0.0.0"} {
		t.Run(spelling, func(t *testing.T) {
			serviceTemplate, parseProblems, err := testsupport.ParseSource(t, version13Template(spelling))
			if err != nil {
				t.Fatalf("zero version %q was rejected: %v\n%s", spelling, err, parseProblems)
			}

			value := normalizedVersionValue(t, serviceTemplate)
			stringer, ok := value.(fmt.Stringer)
			if !ok {
				t.Fatalf("normalized zero version has type %T, want fmt.Stringer", value)
			}
			if got := stringer.String(); got != "" {
				t.Fatalf("zero version %q normalized as supplied version %q, want unspecified sentinel", spelling, got)
			}
		})
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.3.2.5, Additional Requirements
// Requirements: TOSCA13-3.3.2.5-001
// Expected: all three normative zero spellings have equivalent effective values.
// Category: positive, boundary, equivalence, direct
func TestVersionZeroSpellingsAreEquivalent(t *testing.T) {
	var canonical []string
	for _, spelling := range []string{"0", "0.0", "0.0.0"} {
		serviceTemplate, parseProblems, err := testsupport.ParseSource(t, version13Template(spelling))
		if err != nil {
			t.Fatalf("zero version %q was rejected: %v\n%s", spelling, err, parseProblems)
		}
		canonical = append(canonical, normalizedVersionValue(t, serviceTemplate).(fmt.Stringer).String())
	}
	if canonical[0] != canonical[1] || canonical[1] != canonical[2] {
		t.Fatalf("zero spellings are not equivalent: %#v", canonical)
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.3.2.5, Additional Requirements
// Requirements: TOSCA13-3.3.2.5-002
// Expected: zero with a qualifier is rejected during value rendering.
// Category: negative, qualifier, parser phase, direct
func TestQualifiedZeroVersionRejected(t *testing.T) {
	problems := parseRejectedVersion(t, "0.0.0.alpha")
	for _, fragment := range []string{"release", "version", "qualifier"} {
		if !strings.Contains(problems, fragment) {
			t.Fatalf("qualified-zero diagnostic does not contain %q:\n%s", fragment, problems)
		}
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.3.2.5, Additional Requirements
// Requirements: TOSCA13-3.3.2.5-002
// Expected: repeated qualified-zero rejection has a stable diagnostic body.
// Category: negative, determinism, direct
func TestQualifiedZeroVersionDiagnosticDeterministic(t *testing.T) {
	first := versionDiagnosticBody(parseRejectedVersion(t, "0.0.0.alpha"))
	second := versionDiagnosticBody(parseRejectedVersion(t, "0.0.0.alpha"))
	if first != second {
		t.Fatalf("qualified-zero diagnostic is nondeterministic:\nfirst:\n%s\nsecond:\n%s", first, second)
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.3.2.5, Additional Requirements
// Requirements: TOSCA13-3.3.2.5-002
// Expected: a qualifier remains valid for a non-zero version.
// Category: positive, qualifier, regression, direct
func TestNonZeroQualifiedVersionAccepted(t *testing.T) {
	serviceTemplate, parseProblems, err := testsupport.ParseSource(t, version13Template("1.2.3.alpha-4"))
	if err != nil {
		t.Fatalf("non-zero qualified version was rejected: %v\n%s", err, parseProblems)
	}
	if got := normalizedVersionValue(t, serviceTemplate).(fmt.Stringer).String(); got != "1.2.3.alpha-4" {
		t.Fatalf("qualified version normalized as %q", got)
	}
}

// Specification: TOSCA Version 2.0
// Section: version data type (cross-version isolation)
// Expected: the TOSCA 1.3 unspecified-version policy does not alter TOSCA 2.0.
// Category: positive, cross-version regression
func TestTosca20VersionZeroBehaviorUnchanged(t *testing.T) {
	serviceTemplate, parseProblems, err := testsupport.ParseSource(t, version20Template("0.0"))
	if err != nil {
		t.Fatalf("TOSCA 2.0 zero version regressed: %v\n%s", err, parseProblems)
	}
	if got := normalizedVersionValue(t, serviceTemplate).(fmt.Stringer).String(); got != "0.0" {
		t.Fatalf("TOSCA 2.0 zero version changed to %q", got)
	}
}

func parseRejectedVersion(t *testing.T, value string) string {
	t.Helper()
	_, parseProblems, err := testsupport.ParseSource(t, version13Template(value))
	if err == nil {
		t.Fatalf("invalid qualified zero version %q was accepted", value)
	}
	return parseProblems
}

func versionDiagnosticBody(problems string) string {
	if index := strings.Index(problems, "topology_template"); index >= 0 {
		return problems[index:]
	}
	return problems
}

func normalizedVersionValue(t *testing.T, serviceTemplate *normal.ServiceTemplate) any {
	t.Helper()
	node := serviceTemplate.NodeTemplates["node"]
	if node == nil {
		t.Fatal("normalized node template is missing")
	}
	value, ok := node.Properties["release"].(*normal.Primitive)
	if !ok {
		t.Fatalf("normalized release has type %T, want *normal.Primitive", node.Properties["release"])
	}
	return value.Primitive
}

func version13Template(version string) string {
	return `tosca_definitions_version: tosca_simple_yaml_1_3
node_types:
  VersionNode:
    derived_from: tosca.nodes.Root
    properties:
      release:
        type: version
topology_template:
  node_templates:
    node:
      type: VersionNode
      properties:
        release: "` + version + `"
`
}

func version20Template(version string) string {
	return `tosca_definitions_version: tosca_2_0
node_types:
  VersionNode:
    properties:
      release:
        type: version
service_template:
  node_templates:
    node:
      type: VersionNode
      properties:
        release: "` + version + `"
`
}
