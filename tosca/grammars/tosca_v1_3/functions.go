package tosca_v1_3

import (
	"fmt"
	"strings"

	"github.com/tliron/go-ard"
	"github.com/tliron/go-puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/go-puccini/tosca/parsing"
	"github.com/tliron/yamlkeys"
)

// TOSCA Simple Profile in YAML 1.3 intrinsic-function grammar.
//
// Sections 4.3.1.1-4.6.1.2 and 4.8.1.1-4.8.1.2 define function-specific
// argument shapes. Validation belongs to the read phase because it depends
// only on the YAML shape and scalar types. Keeping it in this package prevents
// TOSCA 1.3 grammar from changing the independently selected TOSCA 2.0 path.
var intrinsicFunctionNames = map[string]struct{}{
	"concat":               {},
	"join":                 {},
	"token":                {},
	"get_input":            {},
	"get_property":         {},
	"get_attribute":        {},
	"get_operation_output": {},
	"get_nodes_of_type":    {},
	"get_artifact":         {},
}

func ParseFunctionCall(context *parsing.Context) bool {
	if functionCall, ok := context.Data.(*parsing.FunctionCall); ok {
		validateIntrinsicFunctionCall(context, functionCall, nil)
		return true
	}

	rawArgument, _, _ := intrinsicFunctionEntry(context.Data)

	if !tosca_v2_0.ParseFunctionCall(context) {
		return false
	}

	if functionCall, ok := context.Data.(*parsing.FunctionCall); ok {
		validateIntrinsicFunctionCall(context, functionCall, rawArgument)
	}
	return true
}

func intrinsicFunctionEntry(data ard.Value) (ard.Value, string, int) {
	map_, ok := data.(ard.Map)
	if !ok {
		return nil, "", 0
	}

	var argument ard.Value
	var functionName string
	count := 0
	for key, value := range map_ {
		name := yamlkeys.KeyString(key)
		if _, ok := intrinsicFunctionNames[name]; ok {
			argument = value
			functionName = name
			count++
		}
	}
	return argument, functionName, count
}

func validateIntrinsicFunctionCall(context *parsing.Context, functionCall *parsing.FunctionCall, rawArgument ard.Value) {
	functionName := functionCall.Name
	const prefix = parsing.MetadataFunctionPrefix
	if len(functionName) >= len(prefix) && functionName[:len(prefix)] == prefix {
		functionName = functionName[len(prefix):]
	}
	if _, ok := intrinsicFunctionNames[functionName]; !ok {
		return
	}

	arguments := functionCall.Arguments
	malformed := func(reason string) {
		context.ReportValueMalformed(functionName+" function", reason)
	}

	switch functionName {
	case "concat":
		if _, ok := rawArgument.(ard.List); rawArgument != nil && !ok {
			malformed("arguments must be a YAML sequence")
		} else if len(arguments) < 1 {
			malformed("requires at least one string value expression")
		} else {
			for index, argument := range arguments {
				if !isStringExpression(argument) {
					malformed(fmt.Sprintf("argument %d must be a string or string value expression", index+1))
					break
				}
			}
		}

	case "join":
		if _, ok := rawArgument.(ard.List); rawArgument != nil && !ok {
			malformed("arguments must be a YAML sequence")
		} else if len(arguments) < 1 {
			malformed("requires list of string value expressions")
		} else if len(arguments) > 2 {
			malformed("accepts at most 2 arguments")
		} else {
			switch list := arguments[0].(type) {
			case *parsing.FunctionCall:
				// A string value expression may itself produce the required
				// list at evaluation time.
			case ard.List:
				if len(list) == 0 {
					malformed("first argument list requires one or more string value expressions")
				} else {
					for index, item := range list {
						if !isStringExpression(item) {
							malformed(fmt.Sprintf("list element %d must be a string or string value expression", index+1))
							break
						}
					}
				}
			default:
				malformed("first argument must be a list or a string value expression that produces a list")
			}
			if len(arguments) == 2 && !isString(arguments[1]) {
				malformed("delimiter must be a string")
			}
		}

	case "token":
		if _, ok := rawArgument.(ard.List); rawArgument != nil && !ok {
			malformed("arguments must be a YAML sequence")
		} else if len(arguments) != 3 {
			switch len(arguments) {
			case 0:
				malformed("requires string_with_tokens, string_of_token_chars, and substring_index")
			case 1:
				malformed("requires string_of_token_chars and substring_index")
			case 2:
				malformed("requires substring_index; token requires exactly 3 arguments")
			default:
				malformed("requires exactly 3 arguments")
			}
		} else if !isStringExpression(arguments[0]) {
			malformed("string_with_tokens must be a string")
		} else if !isStringExpression(arguments[1]) {
			malformed("string_of_token_chars must be a string")
		} else if !isIntegerExpression(arguments[2]) {
			malformed("substring_index must be an integer")
		}

	case "get_input":
		if len(arguments) < 1 {
			malformed("requires input_property_name")
		} else if !isString(arguments[0]) {
			malformed("input_property_name must be a string")
		} else {
			for index, argument := range arguments[1:] {
				if !isStringOrInteger(argument) {
					malformed(fmt.Sprintf("nested input argument %d must be a string or integer", index+2))
					break
				}
			}
		}

	case "get_nodes_of_type":
		if len(arguments) < 1 {
			malformed("requires node_type_name")
		} else if len(arguments) > 1 {
			malformed("requires exactly 1 argument")
		} else if !isString(arguments[0]) {
			malformed("node_type_name must be a string")
		}

	case "get_property", "get_attribute":
		valueName := "property"
		if functionName == "get_attribute" {
			valueName = "attribute"
		}
		if len(arguments) < 1 {
			malformed("requires modelable_entity_name")
		} else if !isString(arguments[0]) {
			malformed("modelable_entity_name must be a string")
		} else if len(arguments) < 2 {
			malformed("requires " + valueName + "_name")
		} else if !isString(arguments[1]) {
			malformed(valueName + "_name must be a string")
		} else {
			for index, argument := range arguments[2:] {
				if !isStringOrInteger(argument) {
					malformed(fmt.Sprintf("%s path argument %d must be a string or integer", valueName, index+3))
					break
				}
			}
		}

	case "get_operation_output":
		names := []string{"modelable_entity_name", "interface_name", "operation_name", "output_variable_name"}
		if len(arguments) != len(names) {
			if len(arguments) < len(names) {
				malformed("requires " + names[len(arguments)])
			} else {
				malformed("requires exactly 4 arguments")
			}
		} else {
			for index, argument := range arguments {
				if !isString(argument) {
					malformed(names[index] + " must be a string")
					break
				}
			}
		}

	case "get_artifact":
		if len(arguments) < 1 {
			malformed("requires modelable_entity_name")
		} else if !isString(arguments[0]) {
			malformed("modelable_entity_name must be a string")
		} else if len(arguments) < 2 {
			malformed("requires artifact_name")
		} else if !isString(arguments[1]) {
			malformed("artifact_name must be a string")
		} else if len(arguments) > 4 {
			malformed("accepts at most 4 arguments")
		} else if len(arguments) >= 3 && !isString(arguments[2]) {
			malformed("location must be a string")
		} else if len(arguments) == 4 && !isBoolean(arguments[3]) {
			malformed("remove must be a boolean")
		}
	}

	validateNestedIntrinsicFunctions(context, arguments)
}

func validateNestedIntrinsicFunctions(context *parsing.Context, arguments []any) {
	for index, argument := range arguments {
		switch value := argument.(type) {
		case *parsing.FunctionCall:
			validateIntrinsicFunctionCall(context.ListChild(index, value), value, nil)
		case ard.List:
			for nestedIndex, nested := range value {
				if functionCall, ok := nested.(*parsing.FunctionCall); ok {
					validateIntrinsicFunctionCall(context.ListChild(index, value).ListChild(nestedIndex, functionCall), functionCall, nil)
				}
			}
		}
	}
}

func isString(value any) bool {
	_, ok := value.(string)
	return ok
}

func isBoolean(value any) bool {
	_, ok := value.(bool)
	return ok
}

func isInteger(value any) bool {
	switch value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return true
	default:
		return false
	}
}

func isStringExpression(value any) bool {
	if isString(value) {
		return true
	}
	_, ok := value.(*parsing.FunctionCall)
	return ok
}

func isIntegerExpression(value any) bool {
	if isInteger(value) {
		return true
	}
	_, ok := value.(*parsing.FunctionCall)
	return ok
}

func isStringOrInteger(value any) bool {
	return isString(value) || isInteger(value)
}

// ValidateServiceTemplateFunctions performs checks that require completed
// namespace lookup, hierarchy construction, and inheritance.
func ValidateServiceTemplateFunctions(serviceTemplate *ServiceTemplate) {
	for _, definition := range serviceTemplate.OutputDefinitions {
		validateParameterFunction(serviceTemplate, definition)
	}

	for _, nodeTemplate := range serviceTemplate.NodeTemplates {
		validateValues(serviceTemplate, nodeTemplate.Properties)
		validateValues(serviceTemplate, nodeTemplate.Attributes)
		validateInterfaceFunctions(serviceTemplate, nodeTemplate.Interfaces)
	}

	for _, relationshipTemplate := range serviceTemplate.RelationshipTemplates {
		validateValues(serviceTemplate, relationshipTemplate.Properties)
		validateValues(serviceTemplate, relationshipTemplate.Attributes)
		validateInterfaceFunctions(serviceTemplate, relationshipTemplate.Interfaces)
	}

	for _, workflow := range serviceTemplate.WorkflowDefinitions {
		for _, step := range workflow.StepDefinitions {
			for _, activity := range step.ActivityDefinitions {
				if activity.CallOperation != nil {
					validateValues(serviceTemplate, activity.CallOperation.Inputs)
				}
			}
		}
	}
}

func validateParameterFunction(serviceTemplate *ServiceTemplate, definition *tosca_v2_0.ParameterDefinition) {
	if definition != nil && definition.Value != nil {
		validateValueFunction(serviceTemplate, definition.Value)
	}
}

func validateValues(serviceTemplate *ServiceTemplate, values tosca_v2_0.Values) {
	for _, value := range values {
		validateValueFunction(serviceTemplate, value)
	}
}

func validateInterfaceFunctions(serviceTemplate *ServiceTemplate, interfaces tosca_v2_0.InterfaceAssignments) {
	for _, interfaceAssignment := range interfaces {
		validateValues(serviceTemplate, interfaceAssignment.Inputs)
		for _, operation := range interfaceAssignment.Operations {
			validateValues(serviceTemplate, operation.Inputs)
		}
	}
}

func validateValueFunction(serviceTemplate *ServiceTemplate, value *tosca_v2_0.Value) {
	if value == nil {
		return
	}
	if functionCall, ok := value.Context.Data.(*parsing.FunctionCall); ok {
		validateFunctionResolution(serviceTemplate, value.Context, functionCall)
	}
}

func validateFunctionResolution(serviceTemplate *ServiceTemplate, context *parsing.Context, functionCall *parsing.FunctionCall) {
	functionName := strings.TrimPrefix(functionCall.Name, parsing.MetadataFunctionPrefix)
	arguments := functionCall.Arguments
	malformed := func(reason string) {
		context.ReportValueMalformed(functionName+" function", reason)
	}

	switch functionName {
	case "get_input":
		if inputName, ok := arguments[0].(string); ok {
			if _, found := serviceTemplate.InputDefinitions[inputName]; !found {
				malformed(fmt.Sprintf("input_property_name %q not found", inputName))
			}
		}

	case "get_nodes_of_type":
		if len(arguments) == 1 {
			if nodeTypeName, ok := arguments[0].(string); ok {
				entity, found := context.Namespace.Lookup(nodeTypeName)
				if !found {
					malformed(fmt.Sprintf("node_type_name %q not found", nodeTypeName))
				} else if _, ok := entity.(*tosca_v2_0.NodeType); !ok {
					malformed(fmt.Sprintf("node_type_name %q is not a Node Type", nodeTypeName))
				}
			}
		}

	case "get_property", "get_attribute", "get_operation_output", "get_artifact":
		entity, ok := resolveModelableEntity(serviceTemplate, context, arguments[0].(string))
		if !ok {
			malformed(fmt.Sprintf("modelable_entity_name %q not found or is used in an invalid context", arguments[0]))
			break
		}
		if entity.node == nil && entity.relationship == nil {
			// SOURCE and TARGET in an inline relationship assignment are
			// resolved against runtime endpoints, not static templates.
			break
		}

		switch functionName {
		case "get_property":
			name := arguments[1].(string)
			if !entity.hasProperty(name) && !entity.hasRequirementOrCapability(name) {
				malformed(fmt.Sprintf("property_name %q not found in %q", name, entity.name()))
			}
		case "get_attribute":
			name := arguments[1].(string)
			if !entity.hasAttribute(name) && !entity.hasRequirementOrCapability(name) {
				malformed(fmt.Sprintf("attribute_name %q not found in %q", name, entity.name()))
			}
		case "get_operation_output":
			interfaceName := arguments[1].(string)
			operationName := arguments[2].(string)
			outputName := arguments[3].(string)
			interfaceDefinition, found := entity.interfaceDefinition(interfaceName)
			if !found {
				malformed(fmt.Sprintf("interface_name %q not found in %q", interfaceName, entity.name()))
			} else if operationDefinition, found := interfaceDefinition.OperationDefinitions[operationName]; !found {
				malformed(fmt.Sprintf("operation_name %q not found in interface %q", operationName, interfaceName))
			} else if _, found := operationDefinition.OutputDefinitions[outputName]; !found {
				malformed(fmt.Sprintf("output_variable_name %q not found in operation %q", outputName, operationName))
			}
		case "get_artifact":
			name := arguments[1].(string)
			if !entity.hasArtifact(name) {
				malformed(fmt.Sprintf("artifact_name %q not found in %q", name, entity.name()))
			}
		}
	}

	for index, argument := range arguments {
		if nested, ok := argument.(*parsing.FunctionCall); ok {
			validateFunctionResolution(serviceTemplate, context.ListChild(index, nested), nested)
		}
	}
}

type modelableEntity struct {
	node         *tosca_v2_0.NodeTemplate
	relationship *tosca_v2_0.RelationshipTemplate
}

func resolveModelableEntity(serviceTemplate *ServiceTemplate, context *parsing.Context, name string) (modelableEntity, bool) {
	switch name {
	case "SELF":
		if node := contextNodeTemplate(serviceTemplate, context); node != nil {
			return modelableEntity{node: node}, true
		}
		if relationship := contextRelationshipTemplate(serviceTemplate, context); relationship != nil {
			return modelableEntity{relationship: relationship}, true
		}
		return modelableEntity{}, false
	case "SOURCE", "TARGET":
		if contextRelationshipTemplate(serviceTemplate, context) != nil || strings.Contains(context.Path.String(), ".requirements") {
			// The endpoint is selected at runtime, but the keyword is valid in
			// relationship context.
			return modelableEntity{}, true
		}
		return modelableEntity{}, false
	case "HOST":
		if node := contextNodeTemplate(serviceTemplate, context); node != nil {
			// HOST traversal itself is orchestrator runtime behavior. It does
			// not necessarily select the node containing the function.
			return modelableEntity{}, true
		}
		return modelableEntity{}, false
	default:
		for _, node := range serviceTemplate.NodeTemplates {
			if node.Name == name {
				return modelableEntity{node: node}, true
			}
		}
		for _, relationship := range serviceTemplate.RelationshipTemplates {
			if relationship.Name == name {
				return modelableEntity{relationship: relationship}, true
			}
		}
		return modelableEntity{}, false
	}
}

func contextNodeTemplate(serviceTemplate *ServiceTemplate, context *parsing.Context) *tosca_v2_0.NodeTemplate {
	path := context.Path.String()
	for _, node := range serviceTemplate.NodeTemplates {
		if strings.Contains(path, fmt.Sprintf("node_templates[%q]", node.Name)) {
			return node
		}
	}
	return nil
}

func contextRelationshipTemplate(serviceTemplate *ServiceTemplate, context *parsing.Context) *tosca_v2_0.RelationshipTemplate {
	path := context.Path.String()
	for _, relationship := range serviceTemplate.RelationshipTemplates {
		if strings.Contains(path, fmt.Sprintf("relationship_templates[%q]", relationship.Name)) {
			return relationship
		}
	}
	return nil
}

func (entity modelableEntity) name() string {
	if entity.node != nil {
		return entity.node.Name
	}
	if entity.relationship != nil {
		return entity.relationship.Name
	}
	return "relationship endpoint"
}

func (entity modelableEntity) hasProperty(name string) bool {
	if entity.node != nil && entity.node.NodeType != nil {
		_, ok := entity.node.NodeType.PropertyDefinitions[name]
		return ok
	}
	if entity.relationship != nil && entity.relationship.RelationshipType != nil {
		_, ok := entity.relationship.RelationshipType.PropertyDefinitions[name]
		return ok
	}
	return false
}

func (entity modelableEntity) hasAttribute(name string) bool {
	if entity.node != nil && entity.node.NodeType != nil {
		_, ok := entity.node.NodeType.AttributeDefinitions[name]
		return ok
	}
	if entity.relationship != nil && entity.relationship.RelationshipType != nil {
		_, ok := entity.relationship.RelationshipType.AttributeDefinitions[name]
		return ok
	}
	return false
}

func (entity modelableEntity) hasRequirementOrCapability(name string) bool {
	if entity.node == nil || entity.node.NodeType == nil {
		return false
	}
	if _, ok := entity.node.NodeType.CapabilityDefinitions[name]; ok {
		return true
	}
	for _, requirement := range entity.node.NodeType.RequirementDefinitions {
		if requirement.Name == name {
			return true
		}
	}
	return false
}

func (entity modelableEntity) interfaceDefinition(name string) (*tosca_v2_0.InterfaceDefinition, bool) {
	if entity.node != nil && entity.node.NodeType != nil {
		definition, ok := entity.node.NodeType.InterfaceDefinitions[name]
		return definition, ok
	}
	if entity.relationship != nil && entity.relationship.RelationshipType != nil {
		definition, ok := entity.relationship.RelationshipType.InterfaceDefinitions[name]
		return definition, ok
	}
	return nil, false
}

func (entity modelableEntity) hasArtifact(name string) bool {
	if entity.node == nil {
		return false
	}
	if _, ok := entity.node.Artifacts[name]; ok {
		return true
	}
	if entity.node.NodeType != nil {
		_, ok := entity.node.NodeType.ArtifactDefinitions[name]
		return ok
	}
	return false
}
