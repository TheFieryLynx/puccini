package tosca_v1_3

import (
	"github.com/tliron/go-ard"
	"github.com/tliron/go-puccini/tosca/parsing"
)

// TOSCA Simple Profile in YAML 1.3 §§3.8.3.3 and 3.8.4.3 require
// the source of a node- or relationship-template copy to be complete.
func validateTemplateCopySource(context *parsing.Context, sourceName string, source ard.Value) {
	sourceMap, ok := source.(ard.Map)
	if !ok {
		return
	}
	if _, copied := sourceMap["copy"]; copied {
		context.FieldChild("copy", sourceName).ReportValueInvalid(
			"source template",
			"must not itself use copy",
		)
	}
}
