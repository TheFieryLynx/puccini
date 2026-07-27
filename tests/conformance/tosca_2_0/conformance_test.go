package tosca_2_0_test

import (
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
