package tosca_v1_3

import (
	"fmt"
	"sort"
	"strings"

	"github.com/tliron/go-ard"
	"github.com/tliron/go-puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/go-puccini/tosca/parsing"
)

// validateRequirementTarget handles concrete TOSCA 1.3 assignments (§3.8.2)
// after tuple defaults and relationship rendering. All type references and
// inherited definitions are available; target rendering order is irrelevant.
// Abstract provider selection is intentionally left to the existing path.
func validateRequirementTarget(sourceEntity, assignmentEntity, definitionEntity parsing.EntityPtr) bool {
	source := sourceEntity.(*tosca_v2_0.NodeTemplate)
	assignment := assignmentEntity.(*tosca_v2_0.RequirementAssignment)
	definition := definitionEntity.(*tosca_v2_0.RequirementDefinition)
	target := assignment.TargetNodeTemplate
	if target == nil {
		return false
	}
	if target.NodeType == nil || definition.TargetCapabilityType == nil {
		return true
	}
	state, ok := assignment.Context.GrammarData.(*requirementOccurrences)
	if !ok {
		return true
	}
	state.SelectedCapability = nil
	state.CandidateCapabilities = nil

	requested := "<implicit>"
	explicit := false
	if data, ok := assignment.Context.Data.(ard.Map); ok {
		if value, ok := data["capability"].(string); ok {
			requested, explicit = value, true
		}
	}
	names := make([]string, 0, len(target.NodeType.CapabilityDefinitions))
	for name := range target.NodeType.CapabilityDefinitions {
		names = append(names, name)
	}
	sort.Strings(names)
	actual := make([]string, 0, len(names))
	for _, name := range names {
		cap := target.NodeType.CapabilityDefinitions[name]
		if cap.CapabilityType != nil {
			actual = append(actual, name+":"+cap.CapabilityType.Name)
		}
	}
	report := func(reason string) {
		assignment.Context.ReportPathf(0,
			"[TOSCA13-TARGET-MATCH] source=%q requirement=%q target=%q requested=%q expected=%q actual=[%s]: %s",
			source.Name, assignment.Name, target.Name, requested, definition.TargetCapabilityType.Name, strings.Join(actual, ", "), reason)
	}
	if assignment.TargetNodeFilter != nil {
		report("node_filter requires a Node Type; a Node Template is invalid")
		return true
	}
	hierarchy := assignment.Context.Hierarchy
	if required := definition.TargetNodeType; required != nil && !hierarchy.IsCompatible(required, target.NodeType) {
		report(fmt.Sprintf("node type %q must derive from %q", target.NodeType.Name, required.Name))
		return true
	}
	// Local symbolic names take precedence. A missing symbol may instead be a
	// resolved type selector; a missing symbol must never fall back to the default.
	_, symbolic := target.NodeType.CapabilityDefinitions[requested]
	state.CapabilityIsSymbol = explicit && symbolic
	if explicit && !symbolic && !state.CapabilityIsType {
		report("explicit capability does not exist")
		return true
	}
	var relationship *tosca_v2_0.RelationshipType
	if assignment.Relationship != nil {
		relationship = assignment.Relationship.GetType(definition.RelationshipDefinition)
	}
	var candidates []string
	for _, name := range names {
		capability := target.NodeType.CapabilityDefinitions[name]
		actualType := capability.CapabilityType
		if actualType == nil {
			continue
		}
		if explicit && symbolic && name != requested {
			continue
		}
		if !hierarchy.IsCompatible(definition.TargetCapabilityType, actualType) {
			continue
		}
		if explicit && !symbolic && !hierarchy.IsCompatible(assignment.TargetCapabilityType, actualType) {
			continue
		}
		if relationship != nil && relationship.ValidCapabilityTypeNames != nil {
			valid := false
			for _, allowed := range relationship.ValidCapabilityTypes {
				if hierarchy.IsCompatible(allowed, actualType) {
					valid = true
					break
				}
			}
			if !valid {
				continue
			}
		}
		// valid_source_types belongs to the capability definition/type in 1.3.
		if len(capability.ValidSourceNodeTypes) > 0 {
			valid := false
			for _, allowed := range capability.ValidSourceNodeTypes {
				if hierarchy.IsCompatible(allowed, source.NodeType) {
					valid = true
					break
				}
			}
			if !valid {
				continue
			}
		}
		candidates = append(candidates, name)
	}
	state.CandidateCapabilities = candidates
	if len(candidates) == 0 {
		reason := "no compatible capability"
		if relationship != nil && relationship.ValidCapabilityTypeNames != nil {
			reason += fmt.Sprintf(" (relationship %q valid_target_types=%v)", relationship.Name, *relationship.ValidCapabilityTypeNames)
		}
		report(reason)
	} else if len(candidates) == 1 {
		state.SelectedCapability = &candidates[0]
	}
	// Several candidates are valid but unbound. No map-order tie-break (§3.8.2;
	// normalization policy documented in decision 0015).
	return true
}
