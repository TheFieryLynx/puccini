package parsing

import "github.com/tliron/go-ard"

type DataDefinitionValidator func(EntityPtr)
type DataValueValidator func(*Context, EntityPtr, EntityPtr)
type CapabilityDefinitionRefinementValidator func(EntityPtr, EntityPtr)
type CapabilityAssignmentValidator func(EntityPtr, EntityPtr)
type NodeTemplateValidator func(EntityPtr)
type GroupTypeValidator func(EntityPtr)
type RequirementAssignmentValidator func(EntityPtr)

// Return true when the versioned policy handled target compatibility.
type RequirementTargetValidator func(EntityPtr, EntityPtr, EntityPtr) bool
type AttributeDefinitionValidator func(EntityPtr)
type WorkflowStepDefinitionValidator func(EntityPtr)
type TemplateCopyValidator func(*Context, string, ard.Value)
type SubstitutionMappingsValidator func(EntityPtr)
type PropertyDefinitionInheritor func(EntityPtr, EntityPtr)
type PropertyAssignmentValidator func(*Context, EntityPtr)
type AdditionalNormativeNames func(string, bool) []string
type RequirementAssignmentNormalizer func(EntityPtr, EntityPtr, any) bool

//
// Grammar
//

type Grammar struct {
	Readers                                 Readers
	InvalidNamespaceCharacters              string
	DataDefinitionValidator                 DataDefinitionValidator
	DataValueValidator                      DataValueValidator
	CapabilityDefinitionRefinementValidator CapabilityDefinitionRefinementValidator
	CapabilityAssignmentValidator           CapabilityAssignmentValidator
	NodeTemplateValidator                   NodeTemplateValidator
	GroupTypeValidator                      GroupTypeValidator
	RequirementAssignmentValidator          RequirementAssignmentValidator
	RequirementTargetValidator              RequirementTargetValidator
	AttributeDefinitionValidator            AttributeDefinitionValidator
	WorkflowStepDefinitionValidator         WorkflowStepDefinitionValidator
	TemplateCopyValidator                   TemplateCopyValidator
	SubstitutionMappingsValidator           SubstitutionMappingsValidator
	PropertyDefinitionInheritor             PropertyDefinitionInheritor
	PropertyAssignmentValidator             PropertyAssignmentValidator
	AdditionalNormativeNames                AdditionalNormativeNames
	RequirementAssignmentNormalizer         RequirementAssignmentNormalizer
}

func NewGrammar() Grammar {
	return Grammar{
		Readers: make(Readers),
	}
}

func (self *Grammar) RegisterReader(name string, reader Reader) {
	self.Readers[name] = reader
}
