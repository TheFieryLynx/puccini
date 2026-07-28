package tosca_v1_3

import (
	"github.com/tliron/go-ard"
	"github.com/tliron/go-puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/go-puccini/tosca/parsing"
)

// Instead of creating remapped scriptlets, use the original validation scriptlets directly
// This ensures the keys match what tosca_v2_0.ReadValidationClause expects

var ConstraintClauseScriptlets = cloneConstraintScriptlets()
var ConstraintClauseNativeArgumentIndexes = tosca_v2_0.ValidationClauseNativeArgumentIndexes

// Override the reader to ensure it has access to the correct operators
func ReadConstraintClause(context *parsing.Context) parsing.EntityPtr {
	if isConstraintScalar(context.Data) {
		context = context.Clone(ard.Map{"equal": context.Data})
	}
	validateConstraintOperand(context)
	clause := tosca_v2_0.ReadValidationClause(context).(*tosca_v2_0.ValidationClause)
	switch clause.Operator {
	case "length", "min_length", "max_length":
		clause.ValidateCollection = true
	}
	return clause
}

func normalizeConstraintList(constraints ard.List) ard.List {
	normalized := make(ard.List, len(constraints))
	for index, constraint := range constraints {
		if isConstraintScalar(constraint) {
			normalized[index] = ard.Map{"equal": constraint}
		} else {
			normalized[index] = constraint
		}
	}
	return normalized
}

func cloneConstraintScriptlets() map[string]string {
	scriptlets := make(map[string]string, len(tosca_v2_0.ValidationClauseScriptlets)+1)
	for name, scriptlet := range tosca_v2_0.ValidationClauseScriptlets {
		scriptlets[name] = scriptlet
	}
	scriptlets[parsing.MetadataValidationPrefix+"length"] = `
exports.validate = function(value, expected) {
    if (arguments.length !== 2 || !Number.isInteger(expected) || expected < 0)
        return false;
    if (typeof value === 'string' || Array.isArray(value))
        return value.length === expected;
    if (value !== null && typeof value === 'object')
        return Object.keys(value).length === expected;
    return false;
};`
	return scriptlets
}
