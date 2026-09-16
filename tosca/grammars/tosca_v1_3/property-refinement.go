package tosca_v1_3

import (
	"sort"

	"github.com/tliron/go-ard"
	"github.com/tliron/go-puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/go-puccini/tosca/parsing"
)

// TOSCA 1.3 sections 3.6.10.6 and 3.6.14.2. State is owned by the
// declaring property context, never by a global map or another version.
type propertyRefinement struct {
	Fixed           *tosca_v2_0.Value
	DeclaredDefault *tosca_v2_0.Value
	Refined         bool
}

func isShortPropertyRefinement(data ard.Value) bool {
	fields, ok := data.(ard.Map)
	if !ok {
		return true
	}
	if len(fields) == 0 {
		return false
	}
	for _, key := range []string{"type", "description", "required", "default", "value", "status", "constraints", "validation", "key_schema", "entry_schema", "external-schema", "metadata"} {
		if _, exists := fields[key]; exists {
			return false
		}
	}
	return true
}

func validateRefinementDefault(definition *tosca_v2_0.PropertyDefinition) {
	state, ok := definition.Context.GrammarData.(*propertyRefinement)
	if !ok || state.DeclaredDefault == nil || state.DeclaredDefault == definition.Default || definition.DataType == nil {
		return
	}
	// Even a fallback hidden by a fixed value must be type/constraint compatible.
	state.DeclaredDefault.RenderProperty(definition.DataType, definition)
}

func inheritPropertyDefinition(childEntity, parentEntity parsing.EntityPtr) {
	child := childEntity.(*tosca_v2_0.PropertyDefinition)
	parent := parentEntity.(*tosca_v2_0.PropertyDefinition)
	state, ok := child.Context.GrammarData.(*propertyRefinement)
	if !ok {
		return // Parameter definitions have their own inheritance contract.
	}
	state.Refined = true
	if parentState, ok := parent.Context.GrammarData.(*propertyRefinement); ok && parentState.Fixed != nil {
		if state.Fixed != nil && !constraintValuesEqual(state.Fixed.Context.Data, parentState.Fixed.Context.Data) {
			state.Fixed.Context.ReportValueMalformed("fixed property refinement", "cannot change a fixed property value")
		}
		state.Fixed = parentState.Fixed
	}
	if parent.ValidationClause != nil && child.ValidationClause != nil && child.ValidationClause != parent.ValidationClause {
		combined := tosca_v2_0.NewValidationClause(child.Context)
		combined.Operator = "and"
		combined.Arguments = ard.List{parent.ValidationClause, child.ValidationClause}
		child.ValidationClause = combined
	}
	if state.Fixed != nil {
		// The existing default application machinery supplies the initial
		// value; finality remains an explicit 1.3 assignment policy below.
		child.Default = state.Fixed
	}
}

func validateFilePropertyRefinements(file *tosca_v2_0.File) {
	unique := make(map[*tosca_v2_0.PropertyDefinition]bool)
	add := func(properties tosca_v2_0.PropertyDefinitions) {
		for _, property := range properties {
			unique[property] = true
		}
	}
	for _, typ := range file.DataTypes {
		add(typ.PropertyDefinitions)
	}
	for _, typ := range file.ArtifactTypes {
		add(typ.PropertyDefinitions)
	}
	for _, typ := range file.CapabilityTypes {
		add(typ.PropertyDefinitions)
	}
	for _, typ := range file.RelationshipTypes {
		add(typ.PropertyDefinitions)
	}
	for _, typ := range file.GroupTypes {
		add(typ.PropertyDefinitions)
	}
	for _, typ := range file.PolicyTypes {
		add(typ.PropertyDefinitions)
	}
	for _, typ := range file.NodeTypes {
		add(typ.PropertyDefinitions)
		for _, capability := range typ.CapabilityDefinitions {
			add(capability.PropertyDefinitions)
		}
	}
	properties := make([]*tosca_v2_0.PropertyDefinition, 0, len(unique))
	for property := range unique {
		properties = append(properties, property)
	}
	sort.Slice(properties, func(i, j int) bool { return properties[i].Context.Path.String() < properties[j].Context.Path.String() })
	for _, property := range properties {
		// Section 3.6.10.3: apply the default only after explicit status
		// values have had the opportunity to propagate through inheritance.
		if property.Status == nil {
			status := "supported"
			property.Status = &status
		}
		if state, ok := property.Context.GrammarData.(*propertyRefinement); ok && state.Fixed != nil && !state.Refined {
			state.Fixed.Context.ReportValueMalformed("property value", "value is only permitted in a property refinement")
		}
	}
}

func validateFixedPropertyAssignment(context *parsing.Context, entity parsing.EntityPtr) {
	definition := entity.(*tosca_v2_0.PropertyDefinition)
	state, ok := definition.Context.GrammarData.(*propertyRefinement)
	if !ok || state.Fixed == nil {
		return
	}
	if !constraintValuesEqual(propertyValueData(context.Data), propertyValueData(state.Fixed.Context.Data)) {
		context.ReportValueMalformed("fixed property assignment", "cannot change a fixed property value")
	}
}

func propertyValueData(data any) any {
	switch value := data.(type) {
	case *tosca_v2_0.Value:
		return propertyValueData(value.Context.Data)
	case *tosca_v2_0.ValueList:
		return propertyValueData(value.Slice)
	case *tosca_v2_0.ValueMap:
		return propertyValueData(value.Map)
	case ard.Map:
		result := make(ard.Map, len(value))
		for key, item := range value {
			result[propertyValueData(key)] = propertyValueData(item)
		}
		return result
	case ard.List:
		result := make(ard.List, len(value))
		for index, item := range value {
			result[index] = propertyValueData(item)
		}
		return result
	default:
		return data
	}
}
