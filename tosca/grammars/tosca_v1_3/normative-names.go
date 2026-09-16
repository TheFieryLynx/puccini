package tosca_v1_3

// Section 5.2 and the named-type tables contain aliases that cannot be
// inferred by dropping lowercase URI segments. Add only those exact spellings;
// preserve the existing canonical names and the namespace-shortcuts quirk.
func additionalNormativeNames(name string, shortcuts bool) []string {
	var short, qualified string
	switch name {
	case "tosca.artifacts.Implementation.Bash":
		short, qualified = "Bash", "tosca:Bash"
	case "tosca.artifacts.Implementation.Python":
		short, qualified = "Python", "tosca:Python"
	case "tosca.capabilities.network.Bindable":
		short, qualified = "network.Bindable", "tosca:network.Bindable"
	case "tosca.nodes.Abstract.Storage":
		short = "AbstractStorage" // tosca:Abstract.Storage is already registered.
	case "tosca.nodes.Storage.ObjectStorage":
		short, qualified = "ObjectStorage", "tosca:ObjectStorage"
	case "tosca.nodes.Storage.BlockStorage":
		short, qualified = "BlockStorage", "tosca:BlockStorage"
	case "tosca.relationships.network.BindsTo":
		short = "network.BindsTo" // tosca:BindsTo is already registered.
	}
	var names []string
	if qualified != "" {
		names = append(names, qualified)
	}
	if shortcuts && short != "" {
		names = append(names, short)
	}
	return names
}
