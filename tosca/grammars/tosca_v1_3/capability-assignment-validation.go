package tosca_v1_3

import (
	"github.com/tliron/go-puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/go-puccini/tosca/parsing"
)

func validateCapabilityAssignmentProfile(assignmentEntity parsing.EntityPtr, definitionEntity parsing.EntityPtr) {
	assignment, assignmentOK := assignmentEntity.(*tosca_v2_0.CapabilityAssignment)
	definition, definitionOK := definitionEntity.(*tosca_v2_0.CapabilityDefinition)
	if !assignmentOK || !definitionOK || assignment == nil || definition == nil || definition.CapabilityType == nil {
		return
	}

	if isCapabilityTypeOrDerivedFrom(definition.CapabilityType, "tosca.capabilities.Endpoint") {
		validateEndpointAssignment(assignment)
	}
	if isCapabilityTypeOrDerivedFrom(definition.CapabilityType, "tosca.capabilities.Scalable") {
		validateScalableAssignment(assignment)
	}
}

func validateEndpointAssignment(assignment *tosca_v2_0.CapabilityAssignment) {
	// CapabilityAssignments.Render materializes every definition for
	// normalization. A nil source context distinguishes that internal object
	// from an explicit YAML assignment, including an explicit empty map.
	if assignment.Context.Data == nil {
		return
	}
	if _, hasPort := assignment.Properties["port"]; hasPort {
		return
	}
	if _, hasPorts := assignment.Properties["ports"]; hasPorts {
		return
	}
	assignment.Context.FieldChild("properties", nil).ReportValueMalformed(
		"Endpoint",
		"one of port or ports must be provided",
	)
}

func validateScalableAssignment(assignment *tosca_v2_0.CapabilityAssignment) {
	defaultValue, hasDefault := capabilityIntegerValue(assignment.Properties, "default_instances")
	if !hasDefault {
		return
	}
	minimum, hasMinimum := capabilityIntegerValue(assignment.Properties, "min_instances")
	maximum, hasMaximum := capabilityIntegerValue(assignment.Properties, "max_instances")
	if !hasMinimum || !hasMaximum {
		// Intrinsic function values remain deferred until they are concrete.
		return
	}
	if defaultValue < minimum || defaultValue > maximum {
		assignment.Context.FieldChild("properties", nil).MapChild("default_instances", defaultValue).ReportValueMalformed(
			"Scalable default_instances",
			"must be between min_instances and max_instances",
		)
	}
}

func capabilityIntegerValue(values tosca_v2_0.Values, name string) (int64, bool) {
	value, ok := values[name]
	if !ok || value == nil {
		return 0, false
	}
	switch data := value.Context.Data.(type) {
	case int:
		return int64(data), true
	case int32:
		return int64(data), true
	case int64:
		return data, true
	default:
		return 0, false
	}
}

func isCapabilityTypeOrDerivedFrom(capabilityType *tosca_v2_0.CapabilityType, canonicalName string) bool {
	for current := capabilityType; current != nil; current = current.Parent {
		if current.Name == canonicalName || parsing.GetCanonicalName(current) == canonicalName {
			return true
		}
	}
	return false
}
