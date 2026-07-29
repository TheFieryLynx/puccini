package tosca_v1_3

import (
	"github.com/tliron/go-puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/go-puccini/tosca/parsing"
)

func validateWorkflowStepOperationHost(entityPtr parsing.EntityPtr) {
	step := entityPtr.(*tosca_v2_0.WorkflowStepDefinition)

	if step.TargetNodeRequirementName != nil {
		if step.TargetNodeTemplate == nil {
			// Target lookup reports the primary error.
			return
		}
		if step.OperationHost == nil {
			step.Context.FieldChild("operation_host", nil).ReportPathf(
				0,
				"operation_host is required when target_relationship selects a relationship",
			)
			return
		}

		host := *step.OperationHost
		if host != "SOURCE" && host != "TARGET" {
			step.Context.FieldChild("operation_host", host).ReportValueMalformed(
				"operation_host",
				"must be SOURCE or TARGET for a relationship target",
			)
		}
		return
	}

	if step.TargetGroup != nil {
		// A group host is optional. Its node type/template value is governed by
		// the separate group-host rule.
		return
	}

	if step.TargetNodeTemplate != nil && step.OperationHost != nil {
		step.Context.FieldChild("operation_host", *step.OperationHost).ReportValueMalformed(
			"operation_host",
			"is only valid for relationship or group targets",
		)
	}
}
