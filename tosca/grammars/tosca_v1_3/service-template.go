package tosca_v1_3

import (
	"sync"

	"github.com/tliron/go-puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/go-puccini/tosca/parsing"
)

//
// ServiceTemplate
//
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.9
//

// ServiceTemplate keeps the TOSCA 1.3 rendering lifecycle in the versioned
// grammar while reusing the structural entity model shared by the existing
// grammar implementations.
type ServiceTemplate struct {
	*tosca_v2_0.ServiceTemplate

	validateFunctionsOnce sync.Once
}

func NewServiceTemplate(context *parsing.Context) *ServiceTemplate {
	return &ServiceTemplate{ServiceTemplate: tosca_v2_0.NewServiceTemplate(context)}
}

// ([parsing.Reader] signature)
func ReadServiceTemplate(context *parsing.Context) parsing.EntityPtr {
	self := NewServiceTemplate(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

// ([parsing.Renderable] interface)
func (self *ServiceTemplate) Render() {
	self.ServiceTemplate.Render()
	self.validateFunctionsOnce.Do(func() {
		ValidateServiceTemplateFunctions(self)
	})
}
