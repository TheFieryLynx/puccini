package tosca_v1_3

import (
	"github.com/tliron/go-ard"
	"github.com/tliron/go-puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/go-puccini/tosca/parsing"
)

//
// RequirementDefinition
//
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.7.3
//

// ([parsing.Reader] signature)
func ReadRequirementDefinition(context *parsing.Context) parsing.EntityPtr {
	context.SetReadTag("CountRange", "occurrences,RangeEntity")

	self := tosca_v2_0.ReadRequirementDefinition(context).(*tosca_v2_0.RequirementDefinition)
	self.DefaultCountRange = ard.List{1, 1}
	return self
}
