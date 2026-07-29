package parsing

import "github.com/tliron/go-ard"

type DataDefinitionValidator func(EntityPtr)
type DataValueValidator func(*Context, EntityPtr, EntityPtr)
type CapabilityDefinitionRefinementValidator func(EntityPtr, EntityPtr)
type CapabilityAssignmentValidator func(EntityPtr, EntityPtr)
type GroupTypeValidator func(EntityPtr)
type RequirementAssignmentValidator func(EntityPtr)
type AttributeDefinitionValidator func(EntityPtr)
type WorkflowStepDefinitionValidator func(EntityPtr)
type TemplateCopyValidator func(*Context, string, ard.Value)
type SubstitutionMappingsValidator func(EntityPtr)

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
	GroupTypeValidator                      GroupTypeValidator
	RequirementAssignmentValidator          RequirementAssignmentValidator
	AttributeDefinitionValidator            AttributeDefinitionValidator
	WorkflowStepDefinitionValidator         WorkflowStepDefinitionValidator
	TemplateCopyValidator                   TemplateCopyValidator
	SubstitutionMappingsValidator           SubstitutionMappingsValidator
}

func NewGrammar() Grammar {
	return Grammar{
		Readers: make(Readers),
	}
}

func (self *Grammar) RegisterReader(name string, reader Reader) {
	self.Readers[name] = reader
}
