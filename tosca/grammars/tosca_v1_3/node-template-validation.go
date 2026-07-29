package tosca_v1_3

import (
	"github.com/tliron/go-puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/go-puccini/tosca/parsing"
)

func validateNetworkNodeTemplate(entity parsing.EntityPtr) {
	nodeTemplate, ok := entity.(*tosca_v2_0.NodeTemplate)
	if !ok || nodeTemplate == nil || nodeTemplate.NodeType == nil {
		return
	}
	if !isNodeTypeOrDerivedFrom(nodeTemplate.NodeType, "tosca.nodes.network.Network") {
		return
	}

	networkTypeValue, present := nodeTemplate.Properties["network_type"]
	if !present || networkTypeValue == nil {
		return
	}
	networkType, concrete := networkTypeValue.Context.Data.(string)
	if !concrete || (networkType != "flat" && networkType != "vlan") {
		return
	}
	if _, hasPhysicalNetwork := nodeTemplate.Properties["physical_network"]; hasPhysicalNetwork {
		return
	}

	nodeTemplate.Context.FieldChild("properties", nil).MapChild("physical_network", nil).ReportValueMalformed(
		"physical_network",
		"required when network_type is "+networkType,
	)
}

func isNodeTypeOrDerivedFrom(nodeType *tosca_v2_0.NodeType, canonicalName string) bool {
	for current := nodeType; current != nil; current = current.Parent {
		if current.Name == canonicalName || parsing.GetCanonicalName(current) == canonicalName {
			return true
		}
	}
	return false
}
