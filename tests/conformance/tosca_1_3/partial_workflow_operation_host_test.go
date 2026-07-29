package tosca_1_3_test

import (
	"strings"
	"testing"

	"github.com/tliron/go-puccini/tests/conformance/internal/testsupport"
)

// TOSCA 1.3 §3.6.27.1, "Keynames".
// Category: positive, negative, boundary, relationship/group resolution,
// normalization, and cross-version regression.
// Expected: relationship steps require SOURCE or TARGET operation_host;
// group steps may omit it; ordinary node steps do not use it.
func TestPartialWorkflowOperationHostAcceptsRelationshipEndpoints(t *testing.T) {
	for _, host := range []string{"SOURCE", "TARGET"} {
		t.Run(host, func(t *testing.T) {
			serviceTemplate, parseProblems, err := testsupport.ParseSource(
				t,
				workflowOperationHostTemplate("app", true, host),
			)
			if err != nil {
				t.Fatalf("relationship operation_host %s failed: %v\n%s", host, err, parseProblems)
			}
			if got := serviceTemplate.Workflows["workflow"].Steps["step"].Host; got != host {
				t.Fatalf("normalized operation_host = %q, want %q", got, host)
			}
		})
	}
}

func TestPartialWorkflowOperationHostRejectsMissingRelationshipHost(t *testing.T) {
	problems := rejectWorkflowOperationHost(
		t,
		workflowOperationHostTemplate("app", true, ""),
	)
	for _, fragment := range []string{
		`workflows["workflow"].steps["step"].operation_host`,
		"required",
		"target_relationship",
	} {
		if !strings.Contains(problems, fragment) {
			t.Fatalf("missing relationship host diagnostic does not contain %q:\n%s", fragment, problems)
		}
	}
}

func TestPartialWorkflowOperationHostRejectsInvalidRelationshipHosts(t *testing.T) {
	for _, host := range []string{"SELF", "source", "server"} {
		t.Run(host, func(t *testing.T) {
			problems := rejectWorkflowOperationHost(
				t,
				workflowOperationHostTemplate("app", true, host),
			)
			for _, fragment := range []string{
				`workflows["workflow"].steps["step"].operation_host`,
				host,
				"SOURCE",
				"TARGET",
			} {
				if !strings.Contains(problems, fragment) {
					t.Fatalf("invalid relationship host diagnostic does not contain %q:\n%s", fragment, problems)
				}
			}
		})
	}
}

func TestPartialWorkflowOperationHostNodeTargetApplicability(t *testing.T) {
	if _, parseProblems, err := testsupport.ParseSource(
		t,
		workflowOperationHostTemplate("app", false, ""),
	); err != nil {
		t.Fatalf("ordinary node step without operation_host failed: %v\n%s", err, parseProblems)
	}

	problems := rejectWorkflowOperationHost(
		t,
		workflowOperationHostTemplate("app", false, "SELF"),
	)
	if !strings.Contains(problems, "operation_host") ||
		!strings.Contains(problems, "only valid for relationship or group targets") {
		t.Fatalf("ordinary node operation_host failed for the wrong reason:\n%s", problems)
	}
}

func TestPartialWorkflowOperationHostGroupTargetIsOptional(t *testing.T) {
	for _, host := range []string{"", "server"} {
		t.Run("host-"+host, func(t *testing.T) {
			if _, parseProblems, err := testsupport.ParseSource(
				t,
				workflowOperationHostTemplate("targets", false, host),
			); err != nil {
				t.Fatalf("group operation_host %q failed: %v\n%s", host, err, parseProblems)
			}
		})
	}
}

func TestPartialWorkflowOperationHostUnknownTargetStillFailsLookup(t *testing.T) {
	problems := rejectWorkflowOperationHost(
		t,
		workflowOperationHostTemplate("missing", false, ""),
	)
	if !strings.Contains(problems, "target") || !strings.Contains(problems, "unknown") {
		t.Fatalf("unknown workflow target failed for the wrong reason:\n%s", problems)
	}
	if strings.Contains(problems, "only valid for relationship or group targets") {
		t.Fatalf("unknown target also produced a misleading operation_host diagnostic:\n%s", problems)
	}
}

func TestPartialWorkflowOperationHostDiagnosticIsDeterministic(t *testing.T) {
	source := workflowOperationHostTemplate("app", true, "SELF")
	first := rejectWorkflowOperationHost(t, source)
	second := rejectWorkflowOperationHost(t, source)
	if workflowOperationHostDiagnosticBody(first) != workflowOperationHostDiagnosticBody(second) {
		t.Fatalf("operation_host diagnostic is nondeterministic:\nfirst:\n%s\nsecond:\n%s", first, second)
	}
}

func TestPartialWorkflowOperationHostDoesNotChangeTOSCA20(t *testing.T) {
	source := `tosca_definitions_version: tosca_2_0
description: TOSCA 2.0 workflow operation_host regression
imports:
  - profile: org.oasis-open.simple:2.0
    namespace: tosca
node_types:
  example.Node:
    derived_from: tosca:Root
service_template:
  node_templates:
    node:
      type: example.Node
  workflows:
    workflow:
      steps:
        step:
          target: node
          operation_host: SELF
          activities: []
`
	if _, parseProblems, err := testsupport.ParseSource(t, source); err != nil {
		t.Fatalf("TOSCA 2.0 workflow operation_host behavior changed: %v\n%s", err, parseProblems)
	}
}

func workflowOperationHostTemplate(target string, relationship bool, host string) string {
	targetRelationship := ""
	if relationship {
		targetRelationship = `
          target_relationship: dependency`
	}
	operationHost := ""
	if host != "" {
		operationHost = `
          operation_host: ` + host
	}

	return `tosca_definitions_version: tosca_simple_yaml_1_3
description: TOSCA 1.3 workflow operation_host conformance
node_types:
  example.Node:
    derived_from: tosca.nodes.Root
group_types:
  example.Group:
    members: [tosca.nodes.Root]
topology_template:
  node_templates:
    server:
      type: tosca.nodes.Root
    app:
      type: example.Node
      requirements:
        - dependency: server
  groups:
    targets:
      type: example.Group
      members: [app, server]
  workflows:
    workflow:
      steps:
        step:
          target: ` + target + targetRelationship + operationHost + `
          activities: []
`
}

func rejectWorkflowOperationHost(t *testing.T, source string) string {
	t.Helper()
	_, parseProblems, err := testsupport.ParseSource(t, source)
	if err == nil {
		t.Fatalf("invalid workflow operation_host was accepted")
	}
	return parseProblems
}

func workflowOperationHostDiagnosticBody(problems string) string {
	if index := strings.Index(problems, "@"); index >= 0 {
		return problems[index:]
	}
	return problems
}
