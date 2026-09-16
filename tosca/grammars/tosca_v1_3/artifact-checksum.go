package tosca_v1_3

import (
	"sort"

	"github.com/tliron/go-puccini/tosca/grammars/tosca_v2_0"
)

// Sections 3.5.1 and 3.6.7.1: checksum_algorithm is conditionally required.
// Resolve the effective pair after type inheritance. Template artifact Copy
// runs during rendering; inspect its same fallback without advancing that phase.
func validateArtifactChecksums(file *tosca_v2_0.File, service *ServiceTemplate) {
	effective := make(map[*tosca_v2_0.ArtifactDefinition][2]*string)
	for _, typ := range file.NodeTypes {
		for _, artifact := range typ.ArtifactDefinitions {
			effective[artifact] = [2]*string{artifact.Checksum, artifact.ChecksumAlgorithm}
		}
	}
	if service != nil {
		for _, node := range service.NodeTemplates {
			for name, artifact := range node.Artifacts {
				checksum, algorithm := artifact.Checksum, artifact.ChecksumAlgorithm
				if node.NodeType != nil {
					if parent := node.NodeType.ArtifactDefinitions[name]; parent != nil {
						if checksum == nil {
							checksum = parent.Checksum
						}
						if algorithm == nil {
							algorithm = parent.ChecksumAlgorithm
						}
					}
				}
				effective[artifact.ArtifactDefinition] = [2]*string{checksum, algorithm}
			}
		}
	}
	definitions := make([]*tosca_v2_0.ArtifactDefinition, 0, len(effective))
	for definition := range effective {
		definitions = append(definitions, definition)
	}
	sort.Slice(definitions, func(i, j int) bool {
		return definitions[i].Context.Path.String() < definitions[j].Context.Path.String()
	})
	for _, definition := range definitions {
		pair := effective[definition]
		if pair[0] != nil && pair[1] == nil {
			definition.Context.FieldChild("checksum_algorithm", nil).ReportKeynameMissing()
		}
	}
}
