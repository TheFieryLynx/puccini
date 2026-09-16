// tosca_v1_3/property_definition_adapter.go
package tosca_v1_3

import (
	"github.com/tliron/go-ard"
	"github.com/tliron/go-puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/go-puccini/tosca/parsing"
)

// Embeds all fields/methods from tosca_v2_0.PropertyDefinition
type PropertyDefinition struct {
	*tosca_v2_0.PropertyDefinition
}

// Bridge method: allows 1.3 code to call this
func (p *PropertyDefinition) GetConstraintClauses() ConstraintClauses {
	if p.ValidationClause == nil {
		return nil
	}
	return ConstraintClauses{p.ValidationClause}
}

// [parsing.Reader] signature
func ReadPropertyDefinition(ctx *parsing.Context) parsing.EntityPtr {
	state := &propertyRefinement{}
	ctx.GrammarData = state
	if isShortPropertyRefinement(ctx.Data) {
		// Section 3.6.10.6 uses parameter grammar, including its short
		// fixed-value notation. Inheritance determines whether it is a refinement.
		self := tosca_v2_0.NewPropertyDefinition(ctx)
		state.Fixed = tosca_v2_0.ReadValue(ctx.FieldChild("value", ctx.Data)).(*tosca_v2_0.Value)
		return self
	}

	// Convert "constraints" list (1.x) to "validation" (2.0)
	if m, ok := ctx.Data.(ard.Map); ok {
		// "validation" is TOSCA 2.0 grammar. Do not let the shared
		// representation make it an undocumented TOSCA 1.3 spelling.
		if validation, ok := m["validation"]; ok {
			ctx.FieldChild("validation", validation).ReportKeynameUnsupported()
			delete(m, "validation")
		}
		if c, ok := m["constraints"].(ard.List); ok && len(c) > 0 {
			c = normalizeConstraintList(c)
			if len(c) == 1 {
				m["validation"] = c[0]
			} else {
				m["validation"] = ard.Map{"$and": c}
			}
			delete(m, "constraints")
		}
	}

	// Metadata supported in TOSCA 1.3
	// ctx.SetReadTag("Metadata", "") // Removed: metadata is supported in 1.3

	// Extend only the 1.3 reader with the parameter value field. The shared
	// property model and the 2.0 reader do not recognize this keyname.
	reader := struct {
		*tosca_v2_0.PropertyDefinition
		Value *tosca_v2_0.Value `read:"value,Value"`
	}{PropertyDefinition: tosca_v2_0.NewPropertyDefinition(ctx)}
	ignore := ctx.ReadFields(&reader)
	if ctx.HasQuirk(parsing.QuirkAnnotationsIgnore) {
		ignore = append(ignore, "annotations")
	}
	ctx.ValidateUnsupportedFields(ignore)
	v2prop := reader.PropertyDefinition
	if v2prop.Status != nil {
		switch *v2prop.Status {
		case "supported", "unsupported", "experimental", "deprecated":
		default:
			ctx.FieldChild("status", *v2prop.Status).ReportValueMalformed("property status", "expected supported, unsupported, experimental, or deprecated")
		}
	}
	state.Fixed = reader.Value
	state.DeclaredDefault = v2prop.Default

	// Return the v2.0 entity to match NodeType.Properties type
	return v2prop
}
