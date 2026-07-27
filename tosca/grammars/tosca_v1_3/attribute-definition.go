package tosca_v1_3

import (
	"github.com/tliron/go-puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/go-puccini/tosca/parsing"
)

//
// AttributeDefinition
//
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.12
//

// ([parsing.Reader] signature)
func ReadAttributeDefinition(context *parsing.Context) parsing.EntityPtr {
	context.SetReadTag("Metadata", "")

	return tosca_v2_0.ReadAttributeDefinition(context).(*tosca_v2_0.AttributeDefinition)
}
