package tosca_v1_3

import (
	"strings"

	"github.com/tliron/go-ard"
	"github.com/tliron/go-puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/go-puccini/tosca/parsing"
	"github.com/tliron/yamlkeys"
)

// ReadArtifactDefinition retains the TOSCA 1.3 grammar shape while deferring
// effective required-key validation until after inheritance.
//
// TOSCA Simple Profile in YAML 1.3, sections 3.5.1, 3.6.7.1, and 3.6.7.2.2.
func ReadArtifactDefinition(context *parsing.Context) parsing.EntityPtr {
	return tosca_v2_0.ReadArtifactDefinition(context)
}

// ReadEventFilter supplies the TOSCA 1.3 requiredness that is not shared by
// the historical structural entity.
//
// TOSCA Simple Profile in YAML 1.3, section 3.6.21.1.
func ReadEventFilter(context *parsing.Context) parsing.EntityPtr {
	self := tosca_v2_0.ReadEventFilter(context).(*tosca_v2_0.EventFilter)
	context.ValidateUnsupportedFields([]string{"node", "requirement", "capability"})
	if self.NodeTemplateNameOrTypeName == nil {
		context.FieldChild("node", nil).ReportKeynameMissing()
	}
	return self
}

// ReadInterfaceDefinition validates the cardinality of collection keynames
// only when they are explicitly present. The keynames themselves remain
// optional, including the deprecated pre-1.3 operation shorthand.
//
// TOSCA Simple Profile in YAML 1.3, section 3.6.20.2.2.
func ReadInterfaceDefinition(context *parsing.Context) parsing.EntityPtr {
	self := tosca_v2_0.ReadInterfaceDefinition(context).(*tosca_v2_0.InterfaceDefinition)
	validateNonEmptyMapField(context, "operations")
	validateNonEmptyMapField(context, "notifications")
	return self
}

// ReadInterfaceType keeps the TOSCA 1.3 context restriction out of the shared
// TOSCA 2.0 reader while retaining the shared structural representation.
//
// TOSCA Simple Profile in YAML 1.3, sections 3.7.5.2 and 3.7.5.4.
func ReadInterfaceType(context *parsing.Context) parsing.EntityPtr {
	validateInterfaceTypeOperationNames(context)
	validateInterfaceTypeImplementations(context, "operations")
	validateInterfaceTypeImplementations(context, "notifications")
	self := tosca_v2_0.ReadInterfaceType(context).(*tosca_v2_0.InterfaceType)
	validateNonEmptyMapField(context, "operations")
	validateNonEmptyMapField(context, "notifications")
	return self
}

// ReadGroup validates the cardinality of the optional members keyname.
//
// TOSCA Simple Profile in YAML 1.3, section 3.8.5.2.
func ReadGroup(context *parsing.Context) parsing.EntityPtr {
	self := tosca_v2_0.ReadGroup(context).(*tosca_v2_0.Group)
	if data, present := mapField(context, "members"); present {
		if members, ok := data.(ard.List); ok && len(members) == 0 {
			context.FieldChild("members", data).ReportValueMalformed("members", "must contain one or more entries")
		}
	}
	return self
}

// ReadWorkflowActivityCallOperation enforces the operation field in extended
// notation without changing scalar short notation.
//
// TOSCA Simple Profile in YAML 1.3, section 3.6.23.3.1.
func ReadWorkflowActivityCallOperation(context *parsing.Context) parsing.EntityPtr {
	self := tosca_v2_0.ReadWorkflowActivityCallOperation(context).(*tosca_v2_0.WorkflowActivityCallOperation)
	if context.Is(ard.TypeMap) && self.InterfaceAndOperation == nil {
		context.FieldChild("operation", nil).ReportKeynameMissing()
	}
	return self
}

// ReadWorkflowActivityDefinition adds the TOSCA 1.3 extended inline notation.
// Other activity forms retain the shared structural representation.
//
// TOSCA Simple Profile in YAML 1.3, sections 3.6.23.1-3.6.23.4.
func ReadWorkflowActivityDefinition(context *parsing.Context) parsing.EntityPtr {
	self := tosca_v2_0.NewWorkflowActivityDefinition(context)

	if !context.ValidateType(ard.TypeMap) {
		return self
	}
	map_ := context.Data.(ard.Map)
	if len(map_) != 1 {
		context.ReportValueMalformed("workflow activity definition", "map length not 1")
		return self
	}

	for key, value := range map_ {
		operator := yamlkeys.KeyString(key)
		childContext := context.FieldChild(operator, value)
		switch operator {
		case "delegate":
			self.DelegateWorkflowDefinitionName = childContext.ReadString()
		case "inline":
			if childContext.Is(ard.TypeMap) {
				inline := childContext.Data.(ard.Map)
				childContext.ValidateUnsupportedFields([]string{"workflow", "inputs"})
				if workflow, present := inline["workflow"]; present {
					self.InlineWorkflowDefinitionName = childContext.FieldChild("workflow", workflow).ReadString()
				} else {
					childContext.FieldChild("workflow", nil).ReportKeynameMissing()
				}
				if inputs, present := inline["inputs"]; present {
					inputContext := childContext.FieldChild("inputs", inputs)
					inputContext.ReadMapItems(Grammar.Readers["Value"], func(ard.Value) {})
				}
			} else {
				self.InlineWorkflowDefinitionName = childContext.ReadString()
			}
		case "set_state":
			self.SetNodeState = childContext.ReadString()
		case "call_operation":
			self.CallOperation = ReadWorkflowActivityCallOperation(childContext).(*tosca_v2_0.WorkflowActivityCallOperation)
		default:
			childContext.ReportValueMalformed("workflow activity definition", "unsupported operator")
		}
	}
	return self
}

func validateNonEmptyMapField(context *parsing.Context, key string) {
	if data, present := mapField(context, key); present {
		if map_, ok := data.(ard.Map); ok && len(map_) == 0 {
			context.FieldChild(key, data).ReportValueMalformed(key, "must contain one or more entries")
		}
	}
}

func validateInterfaceTypeOperationNames(context *parsing.Context) {
	data, present := mapField(context, "operations")
	if !present {
		return
	}
	operations, ok := data.(ard.Map)
	if !ok {
		return
	}
	for name, definition := range operations {
		if yamlkeys.KeyString(name) == "inputs" {
			context.FieldChild("operations", data).
				MapChild(name, definition).
				ReportValueMalformed("operation name", "inputs is reserved")
		}
	}
}

func validateInterfaceTypeImplementations(context *parsing.Context, key string) {
	data, present := mapField(context, key)
	if !present {
		return
	}
	definitions, ok := data.(ard.Map)
	if !ok {
		return
	}
	for name, definition := range definitions {
		definitionMap, ok := definition.(ard.Map)
		if !ok {
			continue
		}
		implementation, present := definitionMap["implementation"]
		if present {
			context.FieldChild(key, data).
				MapChild(name, definition).
				FieldChild("implementation", implementation).
				ReportValueMalformed("interface type", "implementation keyname is invalid")
		}
	}
}

func mapField(context *parsing.Context, key string) (ard.Value, bool) {
	map_, ok := context.Data.(ard.Map)
	if !ok {
		return nil, false
	}
	value, present := map_[key]
	return value, present
}

func isConstraintScalar(value ard.Value) bool {
	if value == nil {
		return false
	}
	return !ard.IsMap(value) && !ard.IsList(value)
}

func validateConstraintOperand(context *parsing.Context) {
	map_, ok := context.Data.(ard.Map)
	if !ok || len(map_) != 1 {
		return
	}
	for key, value := range map_ {
		operator := strings.TrimPrefix(yamlkeys.KeyString(key), "$")
		operandContext := context.FieldChild(yamlkeys.KeyString(key), value)
		switch operator {
		case "and":
			list, ok := value.(ard.List)
			if !ok {
				operandContext.ReportValueWrongType(ard.TypeList)
				break
			}
			for index, clause := range list {
				validateConstraintOperand(operandContext.ListChild(index, clause))
			}
		case "in_range":
			list, ok := value.(ard.List)
			if !ok {
				operandContext.ReportValueWrongType(ard.TypeList)
			} else if len(list) != 2 {
				operandContext.ReportValueWrongLength("in_range operand", 2)
			} else {
				for index, item := range list {
					if !isConstraintScalar(item) {
						operandContext.ListChild(index, item).ReportValueMalformed("scalar operand", "required scalar value")
					}
				}
			}
		case "valid_values":
			list, ok := value.(ard.List)
			if !ok {
				operandContext.ReportValueWrongType(ard.TypeList)
			} else if len(list) == 0 {
				operandContext.ReportValueMalformed("valid_values operand", "must contain one or more entries")
			}
		case "pattern", "schema":
			if value == nil {
				operandContext.ReportValueMalformed(operator+" operand", "required value")
			} else {
				operandContext.ValidateType(ard.TypeString)
			}
		case "equal", "greater_than", "greater_or_equal", "less_than", "less_or_equal",
			"length", "min_length", "max_length":
			if !isConstraintScalar(value) {
				operandContext.ReportValueMalformed(operator+" operand", "required scalar value")
			}
		}
	}
}
