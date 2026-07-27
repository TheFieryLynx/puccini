package tosca_v2_0

import (
	"github.com/tliron/go-puccini/tosca/parsing"
)

//
// EventFilter
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.21
//

type EventFilter struct {
	*Entity `name:"event filter" json:"-" yaml:"-"`

	NodeTemplateNameOrTypeName *string `read:"node"`
	RequirementName            *string `read:"requirement"`
	CapabilityName             *string `read:"capability"`

	NodeTemplate *NodeTemplate `lookup:"node,NodeTemplateNameOrTypeName" traverse:"ignore" json:"-" yaml:"-"`
	NodeType     *NodeType     `lookup:"node,NodeTemplateNameOrTypeName" traverse:"ignore" json:"-" yaml:"-"`
}

func NewEventFilter(context *parsing.Context) *EventFilter {
	return &EventFilter{Entity: NewEntity(context)}
}

// ([parsing.Reader] signature)
func ReadEventFilter(context *parsing.Context) parsing.EntityPtr {
	self := NewEventFilter(context)
	context.ReadFields(self)
	return self
}
