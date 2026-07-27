package tosca_v1_3

import (
	"sync"

	"github.com/tliron/go-ard"
	"github.com/tliron/go-puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/go-puccini/tosca/parsing"
)

//
// TriggerDefinitionCondition
//
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.22.2
//

type TriggerDefinitionCondition struct {
	*tosca_v2_0.Entity `name:"trigger definition condition" json:"-" yaml:"-"`

	Constraint  *tosca_v2_0.ConditionClause `read:"constraint,ConditionClause"`
	Period      *tosca_v2_0.Value           `read:"period,Value"`
	Evaluations *int                        `read:"evaluations"`
	Method      *string                     `read:"method"`

	renderOnce sync.Once
}

func NewTriggerDefinitionCondition(context *parsing.Context) *TriggerDefinitionCondition {
	return &TriggerDefinitionCondition{Entity: tosca_v2_0.NewEntity(context)}
}

// ([parsing.Reader] signature)
func ReadTriggerDefinitionCondition(context *parsing.Context) parsing.EntityPtr {
	self := NewTriggerDefinitionCondition(context)

	if context.Is(ard.TypeMap) {
		map_ := context.Data.(ard.Map)
		if _, ok := map_["constraint"]; ok {
			context.ValidateUnsupportedFields(context.ReadFields(self))
		} else {
			self.Constraint = tosca_v2_0.ReadConditionClause(context).(*tosca_v2_0.ConditionClause)
		}
	} else if context.ValidateType(ard.TypeMap, ard.TypeList) {
		context = context.Clone(ard.Map{"and": context.Data})
		self.Constraint = tosca_v2_0.ReadConditionClause(context).(*tosca_v2_0.ConditionClause)
	}

	return self
}

// ([parsing.Renderable] interface)
func (self *TriggerDefinitionCondition) Render() {
	self.renderOnce.Do(self.render)
}

func (self *TriggerDefinitionCondition) render() {
	if self.Period != nil {
		self.Period.RenderDataType("scalar-unit.time")
	}
}
