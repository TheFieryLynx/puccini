package tosca_v1_3

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/tliron/go-ard"
	"github.com/tliron/go-puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/go-puccini/tosca/parsing"
	"github.com/tliron/yamlkeys"
)

type constraintSpec struct {
	operator string
	operands ard.List
	context  *parsing.Context
}

type constraintComparable interface {
	Compare(any) (int, error)
}

// validateConstraintDefinition is installed only on the TOSCA 1.3 grammar.
// The shared rendering hooks provide lifecycle placement without selecting a
// language version.
func validateConstraintDefinition(entity parsing.EntityPtr) {
	switch definition := entity.(type) {
	case *tosca_v2_0.PropertyDefinition:
		validateRefinementDefault(definition)
		validateConstraintCompatibility(definition.ValidationClause, definition.DataType, definition)
	case *tosca_v2_0.DataType:
		if definition.Parent != nil {
			validateConstraintCompatibility(definition.ValidationClause, definition.Parent, nil)
		}
	}
}

func validateConstraintValue(context *parsing.Context, dataTypeEntity parsing.EntityPtr, definitionEntity parsing.EntityPtr) {
	if _, function := context.Data.(*parsing.FunctionCall); function {
		return
	}
	dataType, ok := dataTypeEntity.(*tosca_v2_0.DataType)
	if !ok || dataType == nil {
		return
	}
	definition, _ := definitionEntity.(*tosca_v2_0.PropertyDefinition)

	if dataType.ValidationClause != nil {
		validateConstraintClauseValue(context, context.Data, dataType.ValidationClause, dataType, definition)
	}
	if definition != nil && definition.ValidationClause != nil && definition.ValidationClause != dataType.ValidationClause {
		validateConstraintClauseValue(context, context.Data, definition.ValidationClause, dataType, definition)
	}
	validatePortSpec(context, dataType)
}

func validateConstraintCompatibility(clause *tosca_v2_0.ValidationClause, dataType *tosca_v2_0.DataType, definition tosca_v2_0.DataDefinition) {
	if clause == nil || dataType == nil {
		return
	}
	if isRangeDataType(dataType) && clause.Operator == "in_range" {
		// The two operands constrain the bounds of the range value; neither
		// operand is itself a range-valued native argument.
		clause.NativeArgumentIndexes = nil
	}
	if isCollectionDataType(dataType) {
		clause.ValidateCollection = true
	}
	for _, constraint := range collectConstraintSpecs(clause) {
		if !constraintOperatorCompatible(constraint.operator, dataType) {
			reportConstraintIncompatible(constraint.context, constraint.operator, effectiveConstraintTypeName(dataType))
			continue
		}
		validateConstraintOperands(constraint, dataType, definition)
	}
}

func validateConstraintOperands(constraint constraintSpec, dataType *tosca_v2_0.DataType, definition tosca_v2_0.DataDefinition) {
	switch constraint.operator {
	case "length", "min_length", "max_length":
		if len(constraint.operands) != 1 {
			constraint.context.ReportValueMalformed("constraint", constraint.operator+" requires one integer operand")
			return
		}
		if value, ok := integerOperand(constraint.operands[0]); !ok || value < 0 {
			constraint.context.FieldChild(constraint.operator, constraint.operands[0]).
				ReportValueMalformed("integer operand", "non-negative integer required")
		}
	case "pattern":
		if len(constraint.operands) != 1 {
			constraint.context.ReportValueMalformed("constraint", "pattern requires one string operand")
			return
		}
		pattern, ok := constraint.operands[0].(string)
		if !ok {
			constraint.context.FieldChild(constraint.operator, constraint.operands[0]).
				ReportValueWrongType(ard.TypeString)
			return
		}
		if _, err := regexp.Compile(pattern); err != nil {
			constraint.context.FieldChild(constraint.operator, pattern).
				ReportValueMalformed("pattern operand", err.Error())
		}
	case "schema":
		// External-schema compilation has its own TOSCA 1.3 policy and error
		// categories. Its declaration operand was structurally checked at read.
	case "equal", "greater_than", "greater_or_equal", "less_than", "less_or_equal":
		if len(constraint.operands) != 1 {
			constraint.context.ReportValueMalformed("constraint", constraint.operator+" requires one operand")
			return
		}
		renderConstraintOperand(constraint, 0, dataType, definition)
	case "in_range":
		if len(constraint.operands) != 2 {
			constraint.context.ReportValueMalformed("constraint", "in_range requires two operands")
			return
		}
		if isRangeDataType(dataType) {
			lower, lowerOK := integerOperand(constraint.operands[0])
			upper, upperOK := integerOperand(constraint.operands[1])
			if !lowerOK || !upperOK || lower < 0 || upper < lower {
				constraint.context.ReportValueMalformed("constraint", "in_range on range requires two ordered non-negative integer operands")
			}
			return
		}
		renderConstraintOperand(constraint, 0, dataType, definition)
		renderConstraintOperand(constraint, 1, dataType, definition)
	case "valid_values":
		if len(constraint.operands) == 0 {
			constraint.context.ReportValueMalformed("constraint", "valid_values requires at least one operand")
			return
		}
		for index := range constraint.operands {
			renderConstraintOperand(constraint, index, dataType, definition)
		}
	}
}

func renderConstraintOperand(constraint constraintSpec, index int, dataType *tosca_v2_0.DataType, definition tosca_v2_0.DataDefinition) any {
	operand := constraint.operands[index]
	context := constraint.context.FieldChild(constraint.operator, operand)
	if len(constraint.operands) > 1 {
		context = context.ListChild(index, operand)
	}
	value := tosca_v2_0.ReadValue(context).(*tosca_v2_0.Value)
	value.Render(dataType, definition, true, false)
	return value.Context.Data
}

func validateConstraintClauseValue(
	context *parsing.Context,
	value any,
	clause *tosca_v2_0.ValidationClause,
	dataType *tosca_v2_0.DataType,
	definition tosca_v2_0.DataDefinition,
) {
	for _, constraint := range collectConstraintSpecs(clause) {
		if !constraintOperatorCompatible(constraint.operator, dataType) {
			continue
		}
		if !evaluateConstraint(value, constraint, dataType, definition) {
			context.ReportValueMalformed(
				"constraint",
				fmt.Sprintf("%s constraint not satisfied", constraint.operator),
			)
		}
	}
}

func evaluateConstraint(value any, constraint constraintSpec, dataType *tosca_v2_0.DataType, definition tosca_v2_0.DataDefinition) bool {
	switch constraint.operator {
	case "equal":
		if len(constraint.operands) != 1 {
			return false
		}
		return constraintValuesEqual(value, renderConstraintOperand(constraint, 0, dataType, definition))
	case "greater_than", "greater_or_equal", "less_than", "less_or_equal":
		if len(constraint.operands) != 1 {
			return false
		}
		comparison, ok := compareConstraintValues(value, renderConstraintOperand(constraint, 0, dataType, definition))
		if !ok {
			return false
		}
		switch constraint.operator {
		case "greater_than":
			return comparison > 0
		case "greater_or_equal":
			return comparison >= 0
		case "less_than":
			return comparison < 0
		default:
			return comparison <= 0
		}
	case "in_range":
		if len(constraint.operands) != 2 {
			return false
		}
		if rangeValue, ok := value.(*tosca_v2_0.Range); ok {
			lower, lowerOK := integerOperand(constraint.operands[0])
			upper, upperOK := integerOperand(constraint.operands[1])
			return lowerOK && upperOK && lower >= 0 && upper >= lower &&
				rangeValue.Lower >= uint64(lower) && rangeValue.Upper <= uint64(upper)
		}
		lower, lowerOK := compareConstraintValues(value, renderConstraintOperand(constraint, 0, dataType, definition))
		upper, upperOK := compareConstraintValues(value, renderConstraintOperand(constraint, 1, dataType, definition))
		return lowerOK && upperOK && lower >= 0 && upper <= 0
	case "valid_values":
		for index := range constraint.operands {
			if constraintValuesEqual(value, renderConstraintOperand(constraint, index, dataType, definition)) {
				return true
			}
		}
		return false
	case "length", "min_length", "max_length":
		if len(constraint.operands) != 1 {
			return false
		}
		expected, ok := integerOperand(constraint.operands[0])
		if !ok || expected < 0 {
			return false
		}
		actual, ok := collectionLength(value)
		if !ok {
			return false
		}
		switch constraint.operator {
		case "length":
			return actual == expected
		case "min_length":
			return actual >= expected
		default:
			return actual <= expected
		}
	case "pattern":
		if len(constraint.operands) != 1 {
			return false
		}
		text, textOK := value.(string)
		pattern, patternOK := constraint.operands[0].(string)
		if !textOK || !patternOK {
			return false
		}
		expression, err := regexp.Compile("^(?:" + pattern + ")$")
		return err == nil && expression.MatchString(text)
	case "schema":
		return true
	default:
		return true
	}
}

func collectConstraintSpecs(clause *tosca_v2_0.ValidationClause) []constraintSpec {
	if clause == nil {
		return nil
	}
	if clause.Operator == "and" {
		var constraints []constraintSpec
		for _, argument := range clause.Arguments {
			constraints = append(constraints, collectConstraintValue(argument, clause.Context)...)
		}
		return constraints
	}
	return []constraintSpec{{
		operator: strings.TrimPrefix(clause.Operator, "$"),
		operands: append(ard.List(nil), clause.Arguments...),
		context:  clause.Context,
	}}
}

func collectConstraintValue(value any, context *parsing.Context) []constraintSpec {
	switch value := value.(type) {
	case *tosca_v2_0.ValidationClause:
		return collectConstraintSpecs(value)
	case ard.Map:
		var constraints []constraintSpec
		for key, operand := range value {
			operator := strings.TrimPrefix(yamlkeys.KeyString(key), "$")
			if operator == "and" {
				if list, ok := operand.(ard.List); ok {
					for _, nested := range list {
						constraints = append(constraints, collectConstraintValue(nested, context)...)
					}
				}
				continue
			}
			operands, ok := operand.(ard.List)
			if !ok {
				operands = ard.List{operand}
			}
			constraints = append(constraints, constraintSpec{operator, operands, context})
		}
		return constraints
	default:
		if isConstraintScalar(value) {
			return []constraintSpec{{"equal", ard.List{value}, context}}
		}
		return nil
	}
}

func constraintOperatorCompatible(operator string, dataType *tosca_v2_0.DataType) bool {
	switch operator {
	case "equal", "valid_values", "schema":
		return true
	case "greater_than", "greater_or_equal", "less_than", "less_or_equal", "in_range":
		if operator == "in_range" && isRangeDataType(dataType) {
			return true
		}
		return isComparableDataType(dataType)
	case "length", "min_length", "max_length":
		return isSizedDataType(dataType)
	case "pattern":
		return effectiveConstraintTypeName(dataType) == string(ard.TypeString)
	default:
		return true
	}
}

func isRangeDataType(dataType *tosca_v2_0.DataType) bool {
	for current := dataType; current != nil; current = current.Parent {
		if current.Name == "range" || parsing.GetCanonicalName(current) == "range" {
			return true
		}
	}
	return false
}

func validatePortSpec(context *parsing.Context, dataType *tosca_v2_0.DataType) {
	if !isDataTypeOrDerivedFrom(dataType, "tosca.datatypes.network.PortSpec") {
		return
	}
	fields, ok := context.Data.(ard.Map)
	if !ok {
		return
	}

	hasPortField := false
	for _, name := range []string{"target", "target_range", "source", "source_range"} {
		if _, present := fields[name]; present {
			hasPortField = true
			break
		}
	}
	if !hasPortField {
		context.ReportValueMalformed(
			"PortSpec",
			"at least one of target, target_range, source, or source_range is required",
		)
		return
	}

	validatePortRangePair(context, fields, "source", "source_range")
	validatePortRangePair(context, fields, "target", "target_range")
}

func validatePortRangePair(context *parsing.Context, fields ard.Map, portName, rangeName string) {
	rangeData, rangePresent := fields[rangeName]
	if !rangePresent {
		return
	}
	rangeValue, ok := renderedValueData(rangeData).(*tosca_v2_0.Range)
	if !ok {
		return
	}

	portData, portPresent := fields[portName]
	if !portPresent {
		context.MapChild(portName, nil).ReportValueMalformed(
			"PortSpec",
			fmt.Sprintf("%s is required when %s is specified", portName, rangeName),
		)
		return
	}
	port, ok := integerOperand(renderedValueData(portData))
	if !ok || port < 0 {
		return
	}
	if !rangeValue.InRange(uint64(port)) {
		context.MapChild(portName, port).ReportValueMalformed(
			"PortSpec",
			fmt.Sprintf("%s must be within %s", portName, rangeName),
		)
	}
}

func renderedValueData(value any) any {
	if value, ok := value.(*tosca_v2_0.Value); ok {
		return value.Context.Data
	}
	return value
}

func isDataTypeOrDerivedFrom(dataType *tosca_v2_0.DataType, canonicalName string) bool {
	for current := dataType; current != nil; current = current.Parent {
		if current.Name == canonicalName || parsing.GetCanonicalName(current) == canonicalName {
			return true
		}
	}
	return false
}

func isComparableDataType(dataType *tosca_v2_0.DataType) bool {
	switch effectiveConstraintTypeName(dataType) {
	case string(ard.TypeInteger), string(ard.TypeFloat), string(ard.TypeTimestamp), string(ard.TypeString), "version",
		"scalar-unit.size", "scalar-unit.time", "scalar-unit.frequency", "scalar-unit.bitrate":
		return true
	default:
		return false
	}
}

func isSizedDataType(dataType *tosca_v2_0.DataType) bool {
	switch effectiveConstraintTypeName(dataType) {
	case string(ard.TypeString), string(ard.TypeList), string(ard.TypeMap):
		return true
	default:
		return false
	}
}

func isCollectionDataType(dataType *tosca_v2_0.DataType) bool {
	switch effectiveConstraintTypeName(dataType) {
	case string(ard.TypeList), string(ard.TypeMap):
		return true
	default:
		return false
	}
}

func effectiveConstraintTypeName(dataType *tosca_v2_0.DataType) string {
	if dataType == nil {
		return "<unresolved>"
	}
	if internal, ok := dataType.GetInternalTypeName(); ok {
		return string(internal)
	}
	return dataType.Name
}

func reportConstraintIncompatible(context *parsing.Context, operator string, dataType string) {
	context.ReportValueMalformed(
		"constraint",
		fmt.Sprintf("%s is incompatible with data type %s", operator, dataType),
	)
}

func integerOperand(value any) (int, bool) {
	switch value := value.(type) {
	case int:
		return value, true
	case int32:
		return int(value), true
	case int64:
		return int(value), true
	default:
		return 0, false
	}
}

func collectionLength(value any) (int, bool) {
	switch value := value.(type) {
	case string:
		return utf8.RuneCountInString(value), true
	case ard.List:
		return len(value), true
	case ard.Map:
		return len(value), true
	case *tosca_v2_0.ValueList:
		return len(value.Slice), true
	case *tosca_v2_0.ValueMap:
		return len(value.Map), true
	default:
		return 0, false
	}
}

func constraintValuesEqual(left any, right any) bool {
	if comparison, ok := compareConstraintValues(left, right); ok {
		return comparison == 0
	}
	return reflect.DeepEqual(left, right)
}

func compareConstraintValues(left any, right any) (int, bool) {
	if comparable, ok := left.(constraintComparable); ok {
		comparison, err := comparable.Compare(right)
		return comparison, err == nil
	}
	switch left := left.(type) {
	case int:
		right, ok := right.(int)
		return CompareInt64(int64(left), int64(right)), ok
	case int32:
		right, ok := right.(int32)
		return CompareInt64(int64(left), int64(right)), ok
	case int64:
		right, ok := right.(int64)
		return CompareInt64(left, right), ok
	case float32:
		right, ok := right.(float32)
		return CompareFloat64(float64(left), float64(right)), ok
	case float64:
		right, ok := right.(float64)
		return CompareFloat64(left, right), ok
	case string:
		right, ok := right.(string)
		if !ok {
			return 0, false
		}
		return strings.Compare(left, right), true
	default:
		return 0, false
	}
}
