package tosca_v1_3

import (
	"github.com/tliron/go-puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/go-puccini/tosca/parsing"
)

//
// NodeTemplate
//
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.8.3
//

// ([parsing.Reader] signature)
func ReadNodeTemplate(context *parsing.Context) parsing.EntityPtr {
	self := tosca_v2_0.NewNodeTemplate(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	switch self.Name {
	case "SELF", "SOURCE", "TARGET", "HOST":
		context.Clone(self.Name).ReportValueInvalid("node template name", "reserved")
	}
	return self
}
