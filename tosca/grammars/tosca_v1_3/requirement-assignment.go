package tosca_v1_3

import (
	"github.com/tliron/go-ard"
	"github.com/tliron/go-puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/go-puccini/tosca/parsing"
)

//
// RequirementAssignment
//
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.8.2
//

// ([parsing.Reader] signature)
func ReadRequirementAssignment(context *parsing.Context) parsing.EntityPtr {
	context.SetReadTag("Count", "")
	context.SetReadTag("Directives", "")
	context.SetReadTag("Optional", "")
	//context.SetReadTag("Allocation", "")

	self := tosca_v2_0.NewRequirementAssignment(context)

	if context.Is(ard.TypeMap) {
		// Long notation
		context.ValidateUnsupportedFields(append(context.ReadFields(self), "occurrences"))

		if occurrences := ard.With(self.Context.Data).Get("occurrences"); occurrences != ard.NoNode {
			occurrences_ := tosca_v2_0.ReadRange(context.FieldChild("occurrences", occurrences.Value)).(*tosca_v2_0.Range)
			lower := int64(occurrences_.Lower)
			self.Count = &lower
			// TODO: have no idea what to do with max bound in "occurrences" keyname
		}
	} else if context.ValidateType(ard.TypeMap, ard.TypeString) {
		// Short notation
		self.TargetNodeTemplateNameOrTypeName = context.FieldChild("node", context.Data).ReadString()
	}

	return self
}

func validateRequirementAssignmentNodeFilter(entityPtr parsing.EntityPtr) {
	assignment := entityPtr.(*tosca_v2_0.RequirementAssignment)
	if assignment.TargetNodeFilter == nil {
		return
	}

	data, ok := assignment.Context.Data.(ard.Map)
	if !ok {
		return
	}
	nodeData, declared := data["node"]
	if !declared {
		assignment.Context.FieldChild("node_filter", data["node_filter"]).ReportPathf(
			0,
			"node_filter requires an explicit node keyname whose value is a Node Type",
		)
		return
	}

	if assignment.TargetNodeType != nil {
		return
	}
	if assignment.TargetNodeTemplate != nil {
		assignment.Context.FieldChild("node", nodeData).ReportValueMalformed(
			"node",
			"must name a Node Type when node_filter is present; a Node Template is invalid",
		)
	}
}
