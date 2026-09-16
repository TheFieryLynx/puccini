package tosca_v1_3

import (
	"math"
	"sort"

	"github.com/tliron/go-puccini/normal"
	"github.com/tliron/go-puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/go-puccini/tosca/parsing"
)

type requirementOccurrences struct{ Range *tosca_v2_0.Range }

func requirementDefinitionOccurrences(definition *tosca_v2_0.RequirementDefinition) *tosca_v2_0.Range {
	if definition.CountRange != nil {
		return definition.CountRange.Range
	}
	return &tosca_v2_0.Range{Lower: 1, Upper: 1}
}

// Section 3.8.2.2.2: an assignment interval must fall completely within
// the effective definition interval. Names and inherited definitions are known.
func validateRequirementOccurrences(service *ServiceTemplate) {
	if service == nil {
		return
	}
	nodes := append(tosca_v2_0.NodeTemplates(nil), service.NodeTemplates...)
	sort.Slice(nodes, func(i, j int) bool { return nodes[i].Name < nodes[j].Name })
	for _, node := range nodes {
		for _, assignment := range node.Requirements {
			state, ok := assignment.Context.GrammarData.(*requirementOccurrences)
			if !ok {
				continue
			}
			definition, ok := assignment.GetDefinition(node)
			if !ok {
				continue
			}
			bounds := requirementDefinitionOccurrences(definition)
			if state.Range == nil {
				lower := int64(bounds.Lower)
				assignment.Count = &lower
				continue
			}
			if state.Range.Lower < bounds.Lower || state.Range.Upper > bounds.Upper {
				assignment.Context.FieldChild("occurrences", nil).ReportValueMalformed("requirement occurrences", "must fall completely within the definition occurrences range")
			}
		}
	}
}

// Normalization is a Puccini representation choice: one declaration carries
// its whole interval, including [0,n]. It must not be expanded using only its
// lower bound, nor silently disappear when that bound is zero.
func normalizeRequirementOccurrences(entity, source parsing.EntityPtr, target any) bool {
	assignment := entity.(*tosca_v2_0.RequirementAssignment)
	node := source.(*tosca_v2_0.NodeTemplate)
	normalNode := target.(*normal.NodeTemplate)
	state, declared := assignment.Context.GrammarData.(*requirementOccurrences)
	if !declared {
		// Shared rendering can synthesize several unit assignments for the same
		// missing declaration. Its effective interval represents them together.
		for _, existing := range normalNode.Requirements {
			if existing.Name == assignment.Name {
				return true
			}
		}
	}
	var bounds *tosca_v2_0.Range
	if declared {
		bounds = state.Range
	}
	if bounds == nil {
		if definition, ok := assignment.GetDefinition(node); ok {
			bounds = requirementDefinitionOccurrences(definition)
		}
	}
	requirement := assignment.Normalize(node, normalNode)
	if bounds != nil {
		occurrence := &normal.OccurrenceRange{Lower: bounds.Lower}
		if bounds.Upper != math.MaxUint64 {
			upper := bounds.Upper
			occurrence.Upper = &upper
		}
		requirement.Occurrences = occurrence
	}
	return true
}
