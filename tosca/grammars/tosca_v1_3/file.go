package tosca_v1_3

import (
	"github.com/tliron/go-puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/go-puccini/tosca/parsing"
)

//
// File
//
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.10
//

// ([parsing.Reader] signature)
func ReadFile(context *parsing.Context) parsing.EntityPtr {
	context.SetReadTag("Profile", "namespace")

	self := tosca_v2_0.NewFile(context)
	context.ScriptletNamespace.Merge(DefaultScriptletNamespace)
	ignore := []string{"dsl_definitions"}
	if context.HasQuirk(parsing.QuirkImportsTopologyTemplateIgnore) {
		ignore = append(ignore, "topology_template")
	}
	if context.HasQuirk(parsing.QuirkAnnotationsIgnore) {
		ignore = append(ignore, "annotation_types")
	}
	context.ValidateUnsupportedFields(append(context.ReadFields(self), ignore...))
	validateNamespaceDeclarations(context, self.Profile, self.Imports)
	if self.Profile != nil {
		context.CanonicalNamespace = self.Profile
	}
	return self
}
