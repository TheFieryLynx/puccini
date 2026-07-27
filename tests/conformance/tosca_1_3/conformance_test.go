package tosca_1_3_test

import (
	"strings"
	"testing"

	"github.com/tliron/go-puccini/tests/conformance/internal/testsupport"
)

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.10, Service Template definition
// Expected: accepted and normalized
// Category: positive, minimal, normalization
func TestMinimalServiceTemplate(t *testing.T) {
	serviceTemplate, problems, err := testsupport.ParseFile(t, "minimal.yaml")
	if err != nil {
		t.Fatalf("minimal TOSCA 1.3 template failed: %v\n%s", err, problems)
	}
	if serviceTemplate == nil {
		t.Fatal("minimal TOSCA 1.3 template was not normalized")
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3, TOSCA Simple Profile definitions in YAML
// Expected: accepted and normalized
// Category: positive, complete, regression
func TestFullExistingExample(t *testing.T) {
	serviceTemplate, problems, err := testsupport.ParseFile(t, "../../../examples/1.3/requirements-and-capabilities.yaml")
	if err != nil {
		t.Fatalf("full TOSCA 1.3 example failed: %v\n%s", err, problems)
	}
	if serviceTemplate == nil || len(serviceTemplate.NodeTemplates) == 0 {
		t.Fatal("full TOSCA 1.3 example did not normalize its node templates")
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.6.22.2, Additional keynames for the extended condition notation
// Expected: the 1.3 extended trigger-condition example is accepted
// Category: positive, version isolation, regression
func TestExtendedTriggerConditionRemainsVersionSpecific(t *testing.T) {
	_, problems, err := testsupport.ParseFile(t, "../../../examples/1.3/policies-and-groups.yaml")
	if err != nil {
		t.Fatalf("TOSCA 1.3 extended trigger condition failed: %v\n%s", err, problems)
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.10.1, Keynames; 3.10.2, Grammar
// Expected: rejected in the read phase because service_template is not a 1.3 top-level keyname
// Category: negative, cross-version
func TestRejectsTosca20ServiceTemplateShape(t *testing.T) {
	_, problems, err := testsupport.ParseSource(t, `tosca_definitions_version: tosca_simple_yaml_1_3
service_template:
  node_templates: {}
`)
	if err == nil {
		t.Fatal("TOSCA 1.3 accepted the TOSCA 2.0 service_template keyname")
	}
	if !strings.Contains(problems, "service_template") || !strings.Contains(problems, "unsupported keyname") {
		t.Fatalf("wrong rejection reason:\n%s", problems)
	}
}
