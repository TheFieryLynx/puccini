package tosca_2_0_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tliron/go-puccini/tests/conformance/internal/testsupport"
)

// Specification: TOSCA 2.0
// Section: 6.9, Service Template Definition; 6.9.1, Service Template Grammar
// Expected: accepted and normalized
// Category: positive, minimal, normalization
func TestMinimalServiceTemplate(t *testing.T) {
	serviceTemplate, problems, err := testsupport.ParseFile(t, "minimal.yaml")
	if err != nil {
		t.Fatalf("minimal TOSCA 2.0 template failed: %v\n%s", err, problems)
	}
	if serviceTemplate == nil {
		t.Fatal("minimal TOSCA 2.0 template was not normalized")
	}
}

// Specification: TOSCA 2.0
// Section: 6, TOSCA File Definition; 6.9, Service Template Definition
// Expected: accepted and normalized
// Category: positive, complete, regression
func TestFullExistingExample(t *testing.T) {
	serviceTemplate, problems, err := testsupport.ParseFile(t, "../../../examples/2.0/copy.yaml")
	if err != nil {
		t.Fatalf("full TOSCA 2.0 example failed: %v\n%s", err, problems)
	}
	if serviceTemplate == nil || len(serviceTemplate.NodeTemplates) == 0 {
		t.Fatal("full TOSCA 2.0 example did not normalize its node templates")
	}
}

// Specification: TOSCA 2.0
// Section: 16.5, Trigger Definition
// Expected: the 1.3 extended trigger-condition keyname is rejected in the read phase
// Category: negative, version isolation, regression
func TestRejectsTosca13ExtendedTriggerCondition(t *testing.T) {
	_, problems, err := testsupport.ParseSource(t, `tosca_definitions_version: tosca_2_0
policy_types:
  Monitor:
    triggers:
      threshold:
        event: changed
        condition:
          constraint:
            state:
              - equal: active
        action:
          - set_state: running
`)
	if err == nil {
		t.Fatal("TOSCA 2.0 accepted the TOSCA 1.3 extended trigger condition")
	}
	if !strings.Contains(problems, "constraint") || !strings.Contains(problems, "unsupported operator") {
		t.Fatalf("wrong rejection reason:\n%s", problems)
	}
}

// Specification: TOSCA 2.0
// Section: 6.1, Keynames; 6.9.1, Service Template Grammar
// Expected: rejected in the read phase because topology_template is not a 2.0 top-level keyname
// Category: negative, cross-version
func TestRejectsTosca13TopologyTemplateShape(t *testing.T) {
	_, problems, err := testsupport.ParseSource(t, `tosca_definitions_version: tosca_2_0
topology_template:
  node_templates: {}
`)
	if err == nil {
		t.Fatal("TOSCA 2.0 accepted the TOSCA 1.3 topology_template keyname")
	}
	if !strings.Contains(problems, "topology_template") || !strings.Contains(problems, "unsupported keyname") {
		t.Fatalf("wrong rejection reason:\n%s", problems)
	}
}

// Specification: TOSCA 2.0
// Section: 10.1, Function definition and invocation
// Expected: the TOSCA 2.0 $-prefixed function path remains selected and is
// unaffected by TOSCA 1.3 argument and reference validators
// Category: positive, function, normalization, cross-version regression
func TestTosca20FunctionPathRemainsIsolated(t *testing.T) {
	serviceTemplate, problems, err := testsupport.ParseSource(t, `tosca_definitions_version: tosca_2_0
service_template:
  outputs:
    result:
      value: { $concat: [TOSCA, " ", "2.0"] }
`)
	if err != nil {
		t.Fatalf("TOSCA 2.0 function path regressed: %v\n%s", err, problems)
	}
	if serviceTemplate == nil || serviceTemplate.Outputs["result"] == nil {
		t.Fatal("TOSCA 2.0 function was not normalized")
	}
}

// Specification: TOSCA 2.0
// Sections: 6.4, Namespace; 6.5, Import
// Expected: sorting in the version-neutral namespace merge infrastructure
// does not change TOSCA 2.0 import or qualified-name resolution
// Category: positive, import, namespace, cross-version regression
func TestTosca20NamespaceMergeRemainsIsolated(t *testing.T) {
	directory := t.TempDir()
	writeFixture := func(name string, source string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(directory, name), []byte(source), 0o600); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}

	writeFixture("one.yaml", `tosca_definitions_version: tosca_2_0
node_types:
  One: {}
`)
	writeFixture("two.yaml", `tosca_definitions_version: tosca_2_0
node_types:
  Two: {}
`)
	writeFixture("main.yaml", `tosca_definitions_version: tosca_2_0
imports:
  - url: two.yaml
    namespace: second
  - url: one.yaml
    namespace: first
node_types:
  FirstChild:
    derived_from: first:One
  SecondChild:
    derived_from: second:Two
`)

	if _, problems, err := testsupport.ParseFile(t, filepath.Join(directory, "main.yaml")); err != nil {
		t.Fatalf("TOSCA 2.0 namespace merge regressed: %v\n%s", err, problems)
	}
}

// Specification: TOSCA 2.0
// Section: 16.5, Trigger Definition
// Expected: the version-neutral required-field reader distinguishes an absent
// required sequence key from an explicitly present empty sequence.
// Category: positive, negative, required field, shared-reader regression
func TestTosca20RequiredSequencePresenceRegression(t *testing.T) {
	valid := `tosca_definitions_version: tosca_2_0
policy_types:
  Monitor:
    triggers:
      threshold:
        event: changed
        action: []
`
	if _, problems, err := testsupport.ParseSource(t, valid); err != nil {
		t.Fatalf("present empty required sequence regressed: %v\n%s", err, problems)
	}

	invalid := `tosca_definitions_version: tosca_2_0
policy_types:
  Monitor:
    triggers:
      threshold:
        event: changed
`
	_, problems, err := testsupport.ParseSource(t, invalid)
	if err == nil {
		t.Fatal("absent required sequence key was accepted")
	}
	if !strings.Contains(problems, "action") || !strings.Contains(problems, "missing required keyname") {
		t.Fatalf("wrong required-sequence diagnostic:\n%s", problems)
	}
}
