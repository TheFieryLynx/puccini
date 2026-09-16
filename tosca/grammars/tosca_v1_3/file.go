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

type File struct {
	*tosca_v2_0.File `name:"file"`
}

func NewFile(context *parsing.Context) *File {
	return &File{File: tosca_v2_0.NewFile(context)}
}

// ([parsing.Reader] signature)
func ReadFile(context *parsing.Context) parsing.EntityPtr {
	context.SetReadTag("Profile", "namespace")

	self := NewFile(context)
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

// ([parsing.Inheritable] interface)
func (self *File) Inherit() {
	validateFilePropertyRefinements(self.File)
	validateArtifactChecksums(self.File, nil)
}

// ([parsing.Renderable] interface)
func (self *File) Render() {
	validateArtifactDefinitionRequiredKeynames(self.File)
	validateExternalPropertySchemas(self.File)
	reflectFilePropertyDefinitions(self.File)
}
