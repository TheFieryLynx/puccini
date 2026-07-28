package tosca_v1_3

import (
	"github.com/tliron/go-ard"
	"github.com/tliron/go-puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/go-puccini/tosca/parsing"
)

// ([parsing.Reader] signature)
func ReadValue(context *parsing.Context) parsing.EntityPtr {
	ParseFunctionCall(context)
	return tosca_v2_0.NewValue(context)
}

// ([parsing.Reader] signature)
func ReadAttributeValue(context *parsing.Context) parsing.EntityPtr {
	self := tosca_v2_0.NewValue(context)

	// Unpack long notation (only for attributes)
	// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.13.2.2
	if context.Is(ard.TypeMap) {
		map_ := context.Data.(ard.Map)
		if len(map_) == 2 {
			if description, ok := map_["description"]; ok {
				if value, ok := map_["value"]; ok {
					self.Description = context.FieldChild("description", description).ReadString()
					context.Data = value
				}
			}
		}
	}

	ParseFunctionCall(context)

	return self
}
