package tosca_v1_3

import (
	"sort"

	"github.com/tliron/go-ard"
	"github.com/tliron/go-puccini/tosca/grammars/tosca_v2_0"
)

// validateArtifactDefinitionRequiredKeynames validates the effective long-form
// artifact definitions after ArtifactDefinitions.Inherit has supplied fields
// omitted by derived node types.
//
// TOSCA Simple Profile in YAML 1.3, sections 3.5.1 and 3.6.7.1.
func validateArtifactDefinitionRequiredKeynames(file *tosca_v2_0.File) {
	unique := make(map[*tosca_v2_0.ArtifactDefinition]struct{})
	for _, nodeType := range file.NodeTypes {
		for _, definition := range nodeType.ArtifactDefinitions {
			unique[definition] = struct{}{}
		}
	}

	definitions := make([]*tosca_v2_0.ArtifactDefinition, 0, len(unique))
	for definition := range unique {
		definitions = append(definitions, definition)
	}
	sort.Slice(definitions, func(i int, j int) bool {
		return definitions[i].Context.Path.String() < definitions[j].Context.Path.String()
	})

	for _, definition := range definitions {
		if !ard.IsMap(definition.Context.Data) {
			continue
		}
		if definition.ArtifactTypeName == nil {
			definition.Context.FieldChild("type", nil).ReportKeynameMissing()
		}
		if definition.File == nil {
			definition.Context.FieldChild("file", nil).ReportKeynameMissing()
		}
	}
}
