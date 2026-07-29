package tosca_v1_3

import (
	"sort"

	"github.com/tliron/go-puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/go-puccini/tosca/parsing"
)

// TOSCA Simple Profile in YAML 1.3 §3.8.13.4 requires mappings for
// every effective property, capability, and requirement definition.
func validateSubstitutionMappingCoverage(entityPtr parsing.EntityPtr) {
	substitution, ok := entityPtr.(*tosca_v2_0.SubstitutionMappings)
	if !ok || substitution.NodeType == nil {
		return
	}

	for _, name := range sortedPropertyDefinitionNames(substitution.NodeType.PropertyDefinitions) {
		if _, mapped := substitution.PropertyMappings[name]; !mapped {
			substitution.Context.
				FieldChild("properties", nil).
				MapChild(name, nil).
				ReportValueRequired("property mapping")
		}
	}

	for _, name := range sortedCapabilityDefinitionNames(substitution.NodeType.CapabilityDefinitions) {
		if _, mapped := substitution.CapabilityMappings[name]; !mapped {
			substitution.Context.
				FieldChild("capabilities", nil).
				MapChild(name, nil).
				ReportValueRequired("capability mapping")
		}
	}

	mappedRequirements := make(map[string]struct{}, len(substitution.RequirementMappings))
	for _, mapping := range substitution.RequirementMappings {
		if mapping != nil {
			mappedRequirements[mapping.Name] = struct{}{}
		}
	}
	for _, name := range sortedRequirementDefinitionNames(substitution.NodeType.RequirementDefinitions) {
		if _, mapped := mappedRequirements[name]; !mapped {
			substitution.Context.
				FieldChild("requirements", nil).
				MapChild(name, nil).
				ReportValueRequired("requirement mapping")
		}
	}
}

func sortedPropertyDefinitionNames(definitions tosca_v2_0.PropertyDefinitions) []string {
	names := make([]string, 0, len(definitions))
	for name := range definitions {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func sortedCapabilityDefinitionNames(definitions tosca_v2_0.CapabilityDefinitions) []string {
	names := make([]string, 0, len(definitions))
	for name := range definitions {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func sortedRequirementDefinitionNames(definitions tosca_v2_0.RequirementDefinitions) []string {
	names := make([]string, 0, len(definitions))
	for name := range definitions {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
