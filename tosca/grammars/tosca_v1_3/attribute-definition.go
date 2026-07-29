package tosca_v1_3

import (
	"sort"
	"strings"

	"github.com/tliron/go-ard"
	"github.com/tliron/go-puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/go-puccini/tosca/parsing"
	"github.com/tliron/yamlkeys"
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

type attributeDefaultProvenanceIssue struct {
	context *parsing.Context
	reason  string
}

func validateAttributeDefaultProvenance(entityPtr parsing.EntityPtr) {
	definition := entityPtr.(*tosca_v2_0.AttributeDefinition)
	if definition.Default == nil {
		return
	}
	if strings.HasPrefix(definition.Context.URL.String(), "internal:/profiles/simple/1.3/") {
		// The pinned normative profile explicitly requires literal "initial"
		// defaults for the Root node and relationship state attributes.
		return
	}

	data, ok := definition.Context.Data.(ard.Map)
	if !ok {
		return
	}
	if _, declared := data["default"]; !declared {
		// An inherited default is validated at its declaring definition.
		return
	}

	tosca_v2_0.ParseFunctionCalls(definition.Default.Context)
	hasActualState, issue := analyzeAttributeDefaultProvenance(
		definition.Default.Context,
		definition.Default.Context.Data,
		false,
	)
	if issue == nil && !hasActualState {
		issue = &attributeDefaultProvenanceIssue{
			context: definition.Default.Context,
			reason:  "hard-coded value is invalid",
		}
	}
	if issue != nil {
		issue.context.ReportPathf(
			0,
			"attribute default must be derived from an actual-state attribute or operation output; %s",
			issue.reason,
		)
	}
}

func analyzeAttributeDefaultProvenance(
	context *parsing.Context,
	data ard.Value,
	calculationArgument bool,
) (bool, *attributeDefaultProvenanceIssue) {
	switch value := data.(type) {
	case *parsing.FunctionCall:
		name := strings.TrimPrefix(value.Name, parsing.MetadataFunctionPrefix)
		switch name {
		case "get_attribute", "get_operation_output":
			return true, nil
		case "get_input", "get_property":
			return false, &attributeDefaultProvenanceIssue{
				context: context,
				reason:  name + " derives desired state and is invalid",
			}
		}

		hasActualState := false
		for index, argument := range value.Arguments {
			actualState, issue := analyzeAttributeDefaultProvenance(
				context.ListChild(index, argument),
				argument,
				true,
			)
			if issue != nil {
				return false, issue
			}
			hasActualState = hasActualState || actualState
		}
		if !hasActualState {
			return false, &attributeDefaultProvenanceIssue{
				context: context,
				reason:  name + " is not calculated from actual state",
			}
		}
		return true, nil

	case ard.List:
		if len(value) == 0 {
			return false, &attributeDefaultProvenanceIssue{context: context, reason: "hard-coded value is invalid"}
		}
		hasActualState := false
		for index, item := range value {
			actualState, issue := analyzeAttributeDefaultProvenance(context.ListChild(index, item), item, calculationArgument)
			if issue != nil {
				return false, issue
			}
			hasActualState = hasActualState || actualState
		}
		return hasActualState, nil

	case ard.Map:
		if len(value) == 0 {
			return false, &attributeDefaultProvenanceIssue{context: context, reason: "hard-coded value is invalid"}
		}
		type mapEntry struct {
			key  ard.Value
			name string
		}
		entries := make([]mapEntry, 0, len(value))
		for key := range value {
			entries = append(entries, mapEntry{key: key, name: yamlkeys.KeyString(key)})
		}
		sort.Slice(entries, func(first int, second int) bool {
			return entries[first].name < entries[second].name
		})

		hasActualState := false
		for _, entry := range entries {
			item := value[entry.key]
			actualState, issue := analyzeAttributeDefaultProvenance(
				context.MapChild(entry.key, item),
				item,
				calculationArgument,
			)
			if issue != nil {
				return false, issue
			}
			hasActualState = hasActualState || actualState
		}
		return hasActualState, nil

	default:
		if calculationArgument {
			return false, nil
		}
		return false, &attributeDefaultProvenanceIssue{
			context: context,
			reason:  "hard-coded value is invalid",
		}
	}
}
