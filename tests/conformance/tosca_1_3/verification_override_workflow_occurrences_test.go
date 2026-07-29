package tosca_1_3_test

import (
	"testing"

	"github.com/tliron/go-puccini/normal"
	"github.com/tliron/go-puccini/tosca/parser"
)

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 3.6.17.3 Additional requirements; 3.6.25.2.1 And clause;
// 3.7.3.3 Additional Requirements; 5.8.1 Additional Requirements
// Requirements: TOSCA13-3.6.17.3-002, TOSCA13-3.6.25.2.1-002,
// TOSCA13-3.7.3.3-002, TOSCA13-5.8.1-004
// Expected: subtype and template operation implementations override inherited
// implementations; an and clause retains every conjunct; omitted requirement
// occurrences render as the closed range [1,1].
// Category: positive, inheritance, workflow, rendering, normalization, direct
func TestVerificationOverrideWorkflowOccurrences(t *testing.T) {
	t.Run("TOSCA13-3.6.17.3-002", func(t *testing.T) {
		serviceTemplate, problems, err := parseOverrideWorkflowOccurrences(t)
		if err != nil {
			t.Fatalf("operation subtype override failed: %v\n%s", err, problems)
		}
		operation := serviceTemplate.NodeTemplates["subtype"].Interfaces["Test"].Operations["run"]
		if operation.Implementation != "scripts/subtype.sh" {
			t.Fatalf("subtype implementation = %q, want subtype script", operation.Implementation)
		}
	})

	t.Run("TOSCA13-3.6.25.2.1-002", func(t *testing.T) {
		parserContext, _, problems, release, err := parseOverrideWorkflowOccurrencesContext(t)
		defer release()
		if err != nil {
			t.Fatalf("workflow and clause failed: %v\n%s", err, problems)
		}
		root := reflectionRoot(t, parserContext)
		preconditions := root.ServiceTemplate.WorkflowDefinitions["deploy"].PreconditionDefinitions
		if len(preconditions) != 1 || preconditions[0].ConditionClause == nil {
			t.Fatalf("workflow preconditions = %#v, want one condition", preconditions)
		}
		outer := preconditions[0].ConditionClause
		if outer.Operator == nil || *outer.Operator != "and" || len(outer.ConditionClauses) != 1 {
			t.Fatalf("outer condition = operator %v with %d clauses, want and/1",
				outer.Operator, len(outer.ConditionClauses))
		}
		conjunction := outer.ConditionClauses[0]
		if conjunction.Operator == nil || *conjunction.Operator != "and" ||
			len(conjunction.ConditionClauses) != 2 {
			t.Fatalf("explicit conjunction = operator %v with %d clauses, want and/2",
				conjunction.Operator, len(conjunction.ConditionClauses))
		}
		for index, clause := range conjunction.ConditionClauses {
			if clause.AttributeName == nil || *clause.AttributeName != "state" ||
				len(clause.ValidationClauses) != 1 {
				t.Fatalf("conjunct %d does not retain its state assertion: %#v", index, clause)
			}
		}
	})

	t.Run("TOSCA13-3.7.3.3-002", func(t *testing.T) {
		parserContext, _, problems, release, err := parseOverrideWorkflowOccurrencesContext(t)
		defer release()
		if err != nil {
			t.Fatalf("requirement occurrence default failed: %v\n%s", err, problems)
		}
		root := reflectionRoot(t, parserContext)
		nodeType := findNodeType(t, root.NodeTypes, "example.Parent")
		requirement := nodeType.RequirementDefinitions["host"]
		if requirement == nil || requirement.CountRange == nil || requirement.CountRange.Range == nil {
			t.Fatalf("effective host occurrence range is absent: %#v", requirement)
		}
		if lower, upper := requirement.CountRange.Range.Lower, requirement.CountRange.Range.Upper; lower != 1 || upper != 1 {
			t.Fatalf("default occurrences = [%d,%d], want [1,1]", lower, upper)
		}
	})

	t.Run("TOSCA13-5.8.1-004", func(t *testing.T) {
		serviceTemplate, problems, err := parseOverrideWorkflowOccurrences(t)
		if err != nil {
			t.Fatalf("operation template override failed: %v\n%s", err, problems)
		}
		operation := serviceTemplate.NodeTemplates["template"].Interfaces["Test"].Operations["run"]
		if operation.Implementation != "scripts/template.sh" {
			t.Fatalf("template implementation = %q, want template script", operation.Implementation)
		}
	})
}

func parseOverrideWorkflowOccurrences(t *testing.T) (*normal.ServiceTemplate, string, error) {
	t.Helper()
	_, serviceTemplate, problems, release, err := parseOverrideWorkflowOccurrencesContext(t)
	release()
	return serviceTemplate, problems, err
}

func parseOverrideWorkflowOccurrencesContext(t *testing.T) (*parser.Context, *normal.ServiceTemplate, string, func(), error) {
	t.Helper()
	return parseReflectionSource(t, overrideWorkflowOccurrencesSource)
}

const overrideWorkflowOccurrencesSource = `tosca_definitions_version: tosca_simple_yaml_1_3
description: TOSCA 1.3 override, workflow, and occurrences verification
interface_types:
  example.Interface:
    operations:
      run: {}
node_types:
  example.Parent:
    derived_from: tosca.nodes.Root
    requirements:
      - host:
          capability: tosca.capabilities.Container
    interfaces:
      Test:
        type: example.Interface
        operations:
          run: scripts/parent.sh
  example.Child:
    derived_from: example.Parent
    interfaces:
      Test:
        type: example.Interface
        operations:
          run: scripts/subtype.sh
topology_template:
  node_templates:
    subtype:
      type: example.Child
    template:
      type: example.Child
      interfaces:
        Test:
          operations:
            run: scripts/template.sh
  workflows:
    deploy:
      preconditions:
        ready:
          target: subtype
          condition:
            - and:
                - state:
                    - equal: started
                - state:
                    - equal: available
      steps: {}
`
