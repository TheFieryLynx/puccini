package tosca_v1_3

import (
	"github.com/tliron/go-puccini/normal"
	"github.com/tliron/go-puccini/tosca/csar"
	"github.com/tliron/go-puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/go-puccini/tosca/parsing"
)

//
// ServiceFile
//
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.10
//

type ServiceFile struct {
	*tosca_v2_0.File `name:"service file"`

	ServiceTemplate *ServiceTemplate `read:"topology_template,ServiceTemplate"`
}

func NewServiceFile(context *parsing.Context) *ServiceFile {
	return &ServiceFile{File: tosca_v2_0.NewFile(context)}
}

// ([parsing.Reader] signature)
func ReadServiceFile(context *parsing.Context) parsing.EntityPtr {
	context.SetReadTag("Profile", "namespace")

	self := NewServiceFile(context)
	context.ScriptletNamespace.Merge(DefaultScriptletNamespace)
	ignore := []string{"dsl_definitions"}
	if context.HasQuirk(parsing.QuirkAnnotationsIgnore) {
		ignore = append(ignore, "annotation_types")
	}
	context.ValidateUnsupportedFields(append(context.ReadFields(self), ignore...))
	validateCSAR(context, self)
	validateNamespaceDeclarations(context, self.Profile, self.Imports)
	if self.Profile != nil {
		context.CanonicalNamespace = self.Profile
	}
	return self
}

func validateCSAR(context *parsing.Context, serviceFile *ServiceFile) {
	csarContext := context.CSAR
	if csarContext == nil {
		return
	}

	if csarContext.MetaPresent {
		if err := csar.ValidateTOSCA13Meta(csarContext.MetaFileVersion, csarContext.EntryDefinitions); err != nil {
			context.ReportError(err)
		}
		return
	}

	if !csarContext.RootFallback {
		return
	}
	if serviceFile.Metadata == nil {
		context.FieldChild("metadata", nil).ReportKeynameMissing()
		return
	}
	if _, ok := serviceFile.Metadata["template_name"]; !ok {
		context.FieldChild("metadata", nil).FieldChild("template_name", nil).ReportKeynameMissing()
	}
	if _, ok := serviceFile.Metadata["template_version"]; !ok {
		context.FieldChild("metadata", nil).FieldChild("template_version", nil).ReportKeynameMissing()
	}
}

// ([parsing.NamespaceValidator] interface)
func (self *ServiceFile) ValidateNamespace() {
	validateImportedDefinitionIdentities(self.Context)
}

// ([parsing.Inheritable] interface)
func (self *ServiceFile) Inherit() {
	validateFilePropertyRefinements(self.File)
	validateArtifactChecksums(self.File, self.ServiceTemplate)
}

// ([parsing.Renderable] interface)
func (self *ServiceFile) Render() {
	validateArtifactDefinitionRequiredKeynames(self.File)
	validateExternalPropertySchemas(self.File)
	reflectFilePropertyDefinitions(self.File)
}

// normal.Normalizable interface
func (self *ServiceFile) NormalizeServiceTemplate() *normal.ServiceTemplate {
	normalServiceTemplate := normal.NewServiceTemplate()

	if self.Description != nil {
		normalServiceTemplate.Description = *self.Description
	}

	normalServiceTemplate.ScriptletNamespace = self.Context.ScriptletNamespace

	self.File.Normalize(normalServiceTemplate)
	if self.ServiceTemplate != nil {
		self.ServiceTemplate.Normalize(normalServiceTemplate)
	}
	normalizeOperatingSystemCapabilities(normalServiceTemplate)

	return normalServiceTemplate
}
