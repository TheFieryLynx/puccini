package tosca_v2_0

import (
	"github.com/tliron/go-ard"
	"github.com/tliron/go-puccini/normal"
	"github.com/tliron/go-puccini/tosca/parsing"
)

//
// OperationDefinition
//
// [TOSCA-v2.0] @ 11.4
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.17
//

type OperationDefinition struct {
	*Entity `name:"operation definition"`
	Name    string

	Description       *string                  `read:"description"`
	Implementation    *InterfaceImplementation `read:"implementation,InterfaceImplementation"`
	InputDefinitions  ParameterDefinitions     `read:"inputs,ParameterDefinition"`
	OutputDefinitions ParameterDefinitions     `read:"outputs,ParameterDefinition"` // changed from OutputMappings to ParameterDefinitions
}

func NewOperationDefinition(context *parsing.Context) *OperationDefinition {
	return &OperationDefinition{
		Entity:            NewEntity(context),
		Name:              context.Name,
		InputDefinitions:  make(ParameterDefinitions),
		OutputDefinitions: make(ParameterDefinitions),
	}
}

// ([parsing.Reader] signature)
func ReadOperationDefinition(context *parsing.Context) parsing.EntityPtr {
	self := NewOperationDefinition(context)

	if context.Is(ard.TypeMap) {
		// Long notation
		context.ValidateUnsupportedFields(context.ReadFields(self))
	} else if context.ValidateType(ard.TypeMap, ard.TypeString) {
		// Short notation
		self.Implementation = ReadInterfaceImplementation(context.FieldChild("implementation", context.Data)).(*InterfaceImplementation)
	}

	return self
}

// ([parsing.Mappable] interface)
func (self *OperationDefinition) GetKey() string {
	return self.Name
}

func (self *OperationDefinition) Inherit(parentDefinition *OperationDefinition) {
	logInherit.Debugf("operation definition: %s", self.Name)

	if (self.Description == nil) && (parentDefinition.Description != nil) {
		self.Description = parentDefinition.Description
	}

	self.InputDefinitions.Inherit(parentDefinition.InputDefinitions)
	self.OutputDefinitions.Inherit(parentDefinition.OutputDefinitions)
	if inherit := self.Context.Grammar.OperationDefinitionInheritor; inherit != nil {
		inherit(self, parentDefinition)
	}
}

func (self *OperationDefinition) Normalize(normalOperation *normal.Operation) {
	if self.Description != nil {
		normalOperation.Description = *self.Description
	}

	if self.Implementation != nil {
		self.Implementation.NormalizeOperation(normalOperation)
	}

	// TODO: input definitions
	//self.InputDefinitions.Normalize(o.Inputs)

	// TODO: output definitions
	//self.OutputDefinitions.Normalize(o.Outputs)
}

//
// OperationDefinitions
//

type OperationDefinitions map[string]*OperationDefinition

func (self OperationDefinitions) Inherit(parentDefinitions OperationDefinitions) {
	for name, definition := range parentDefinitions {
		if _, ok := self[name]; !ok {
			self[name] = definition
		}
	}

	for name, definition := range self {
		if parentDefinition, ok := parentDefinitions[name]; ok {
			if definition != parentDefinition {
				definition.Inherit(parentDefinition)
			}
		}
	}
}
