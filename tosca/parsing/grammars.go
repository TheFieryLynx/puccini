package parsing

type DataDefinitionValidator func(EntityPtr)
type DataValueValidator func(*Context, EntityPtr, EntityPtr)
type CapabilityDefinitionRefinementValidator func(EntityPtr, EntityPtr)
type GroupTypeValidator func(EntityPtr)
type RequirementAssignmentValidator func(EntityPtr)

//
// Grammar
//

type Grammar struct {
	Readers                                 Readers
	InvalidNamespaceCharacters              string
	DataDefinitionValidator                 DataDefinitionValidator
	DataValueValidator                      DataValueValidator
	CapabilityDefinitionRefinementValidator CapabilityDefinitionRefinementValidator
	GroupTypeValidator                      GroupTypeValidator
	RequirementAssignmentValidator          RequirementAssignmentValidator
}

func NewGrammar() Grammar {
	return Grammar{
		Readers: make(Readers),
	}
}

func (self *Grammar) RegisterReader(name string, reader Reader) {
	self.Readers[name] = reader
}
