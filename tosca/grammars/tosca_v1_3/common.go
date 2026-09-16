package tosca_v1_3

import (
	"github.com/tliron/go-puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/go-puccini/tosca/parsing"
)

type (
	ConstraintClause  = tosca_v2_0.ValidationClause
	ConstraintClauses = tosca_v2_0.ValidationClauses
)

var Grammar = parsing.NewGrammar()

var DefaultScriptletNamespace = parsing.NewScriptletNamespace()

func init() {
	Grammar.RegisterReader("$Root", ReadServiceFile) // override
	Grammar.RegisterReader("$File", ReadFile)        // override

	Grammar.RegisterReader("Artifact", tosca_v2_0.ReadArtifact)
	Grammar.RegisterReader("ArtifactDefinition", ReadArtifactDefinition) // override: 1.3 extended notation required fields
	Grammar.RegisterReader("ArtifactType", tosca_v2_0.ReadArtifactType)
	Grammar.RegisterReader("AttributeDefinition", ReadAttributeDefinition)      // override
	Grammar.RegisterReader("AttributeMapping", tosca_v2_0.ReadAttributeMapping) // introduced in TOSCA 1.3
	Grammar.RegisterReader("AttributeValue", ReadAttributeValue)                // override
	Grammar.RegisterReader("CapabilityAssignment", tosca_v2_0.ReadCapabilityAssignment)
	Grammar.RegisterReader("CapabilityDefinition", ReadCapabilityDefinition) // override
	Grammar.RegisterReader("CapabilityFilter", tosca_v2_0.ReadCapabilityFilter)
	Grammar.RegisterReader("CapabilityMapping", tosca_v2_0.ReadCapabilityMapping)
	Grammar.RegisterReader("CapabilityType", tosca_v2_0.ReadCapabilityType)
	Grammar.RegisterReader("ConditionClause", tosca_v2_0.ReadConditionClause)
	Grammar.RegisterReader("ConditionClauseAnd", tosca_v2_0.ReadConditionClauseAnd)
	Grammar.RegisterReader("ConstraintClause", ReadConstraintClause) //override
	Grammar.RegisterReader("ValidationClause", ReadConstraintClause) // override: 1.3 constraint operand grammar
	Grammar.RegisterReader("DataType", ReadDataType)
	Grammar.RegisterReader("EventFilter", ReadEventFilter) // override: 1.3 required node
	Grammar.RegisterReader("Group", ReadGroup)             // override: 1.3 members cardinality
	Grammar.RegisterReader("GroupType", tosca_v2_0.ReadGroupType)
	Grammar.RegisterReader("Import", ReadImport) // override
	Grammar.RegisterReader("InterfaceAssignment", tosca_v2_0.ReadInterfaceAssignment)
	Grammar.RegisterReader("InterfaceDefinition", ReadInterfaceDefinition) // override: 1.3 collection cardinality
	Grammar.RegisterReader("InterfaceMapping", ReadInterfaceMapping)       // override - TOSCA 1.3 format
	Grammar.RegisterReader("InterfaceType", ReadInterfaceType)             // override: 1.3 implementation restriction
	Grammar.RegisterReader("Metadata", tosca_v2_0.ReadMetadata)
	Grammar.RegisterReader("NodeFilter", tosca_v2_0.ReadNodeFilter)
	Grammar.RegisterReader("NodeTemplate", ReadNodeTemplate) // override
	Grammar.RegisterReader("NodeType", tosca_v2_0.ReadNodeType)
	Grammar.RegisterReader("NotificationAssignment", tosca_v2_0.ReadNotificationAssignment) // introduced in TOSCA 1.3
	Grammar.RegisterReader("NotificationDefinition", tosca_v2_0.ReadNotificationDefinition) // introduced in TOSCA 1.3
	Grammar.RegisterReader("OperationAssignment", tosca_v2_0.ReadOperationAssignment)
	Grammar.RegisterReader("OperationDefinition", tosca_v2_0.ReadOperationDefinition)
	Grammar.RegisterReader("OutputMapping", tosca_v2_0.ReadOutputMapping) // introduced in TOSCA 1.3
	Grammar.RegisterReader("InterfaceImplementation", tosca_v2_0.ReadInterfaceImplementation)
	Grammar.RegisterReader("ParameterDefinition", ReadParameterDefinition) // override
	Grammar.RegisterReader("Policy", tosca_v2_0.ReadPolicy)
	Grammar.RegisterReader("PolicyType", tosca_v2_0.ReadPolicyType)
	Grammar.RegisterReader("PropertyDefinition", ReadPropertyDefinition)
	Grammar.RegisterReader("PropertyFilter", tosca_v2_0.ReadPropertyFilter)
	Grammar.RegisterReader("PropertyMapping", tosca_v2_0.ReadPropertyMapping)
	Grammar.RegisterReader("range", ReadRange)
	Grammar.RegisterReader("RangeEntity", tosca_v2_0.ReadRangeEntity)
	Grammar.RegisterReader("RelationshipAssignment", tosca_v2_0.ReadRelationshipAssignment)
	Grammar.RegisterReader("RelationshipDefinition", tosca_v2_0.ReadRelationshipDefinition)
	Grammar.RegisterReader("RelationshipTemplate", tosca_v2_0.ReadRelationshipTemplate)
	Grammar.RegisterReader("RelationshipType", ReadRelationshipType) // override
	Grammar.RegisterReader("Repository", tosca_v2_0.ReadRepository)
	Grammar.RegisterReader("RequirementAssignment", ReadRequirementAssignment) // override
	Grammar.RegisterReader("RequirementDefinition", ReadRequirementDefinition) // override
	Grammar.RegisterReader("RequirementMapping", tosca_v2_0.ReadRequirementMapping)
	Grammar.RegisterReader("ServiceTemplate", ReadServiceTemplate)       // override: TOSCA 1.3 function resolution
	Grammar.RegisterReader("scalar-unit.bitrate", ReadScalarUnitBitrate) // introduced in TOSCA 1.3
	Grammar.RegisterReader("scalar-unit.frequency", ReadScalarUnitFrequency)
	Grammar.RegisterReader("scalar-unit.size", ReadScalarUnitSize)
	Grammar.RegisterReader("scalar-unit.time", ReadScalarUnitTime)
	Grammar.RegisterReader("Schema", ReadSchema)
	Grammar.RegisterReader("SubstitutionMappings", ReadSubstitutionMappings) // override
	Grammar.RegisterReader("timestamp", tosca_v2_0.ReadTimestamp)
	Grammar.RegisterReader("TriggerDefinition", ReadTriggerDefinition) // override
	Grammar.RegisterReader("TriggerDefinitionCondition", ReadTriggerDefinitionCondition)
	Grammar.RegisterReader("Value", ReadValue)                                                 // override: TOSCA 1.3 intrinsic-function grammar
	Grammar.RegisterReader("version", ReadVersion)                                             // override: TOSCA 1.3 zero-version semantics
	Grammar.RegisterReader("WorkflowActivityCallOperation", ReadWorkflowActivityCallOperation) // override: 1.3 extended notation
	Grammar.RegisterReader("WorkflowActivityDefinition", ReadWorkflowActivityDefinition)       // override: 1.3 inline notation
	Grammar.RegisterReader("WorkflowDefinition", tosca_v2_0.ReadWorkflowDefinition)
	Grammar.RegisterReader("WorkflowPreconditionDefinition", tosca_v2_0.ReadWorkflowPreconditionDefinition)
	Grammar.RegisterReader("WorkflowStepDefinition", tosca_v2_0.ReadWorkflowStepDefinition)

	DefaultScriptletNamespace.RegisterScriptlets(tosca_v2_0.FunctionScriptlets, nil)
	DefaultScriptletNamespace.RegisterScriptlets(ConstraintClauseScriptlets, ConstraintClauseNativeArgumentIndexes)
	DefaultScriptletNamespace.RegisterScriptlet(
		operatingSystemLowercaseConverterName,
		operatingSystemLowercaseConverterScript,
		nil,
	)

	Grammar.InvalidNamespaceCharacters = ":"
	Grammar.DataDefinitionValidator = validateConstraintDefinition
	Grammar.PropertyDefinitionInheritor = inheritPropertyDefinition
	Grammar.PropertyAssignmentValidator = validateFixedPropertyAssignment
	Grammar.AdditionalNormativeNames = additionalNormativeNames
	Grammar.RequirementAssignmentNormalizer = normalizeRequirementOccurrences
	Grammar.DataValueValidator = validateConstraintValue
	Grammar.CapabilityDefinitionRefinementValidator = validateCapabilityDefinitionRefinement
	Grammar.CapabilityAssignmentValidator = validateCapabilityAssignmentProfile
	Grammar.NodeTemplateValidator = validateNetworkNodeTemplate
	Grammar.GroupTypeValidator = validateGroupTypeMembers
	Grammar.RequirementAssignmentValidator = validateRequirementAssignmentNodeFilter
	Grammar.RequirementTargetValidator = validateRequirementTarget
	Grammar.AttributeDefinitionValidator = validateAttributeDefaultProvenance
	Grammar.WorkflowStepDefinitionValidator = validateWorkflowStepOperationHost
	Grammar.TemplateCopyValidator = validateTemplateCopySource
	Grammar.SubstitutionMappingsValidator = validateSubstitutionMappingCoverage
}

func CompareUint32(v1 uint32, v2 uint32) int {
	if v1 < v2 {
		return -1
	} else if v1 > v2 {
		return 1
	}
	return 0
}

func CompareUint64(v1 uint64, v2 uint64) int {
	if v1 < v2 {
		return -1
	} else if v1 > v2 {
		return 1
	}
	return 0
}

func CompareInt64(v1 int64, v2 int64) int {
	if v1 < v2 {
		return -1
	} else if v1 > v2 {
		return 1
	}
	return 0
}

func CompareFloat64(v1 float64, v2 float64) int {
	if v1 < v2 {
		return -1
	} else if v1 > v2 {
		return 1
	}
	return 0
}
