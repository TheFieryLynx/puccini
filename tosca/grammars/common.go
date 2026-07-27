package grammars

import (
	"github.com/tliron/go-puccini/tosca/parsing"
)

const (
	definitionsVersionKey = "tosca_definitions_version"
	toscaSimpleYAML13     = "tosca_simple_yaml_1_3"
	tosca20               = "tosca_2_0"
)

var grammars map[string]*parsing.Grammar
var implicitProfilePaths map[string]string

// SupportedDefinitionsVersions returns the complete set of accepted values
// for the mandatory tosca_definitions_version keyname.
func SupportedDefinitionsVersions() []string {
	return []string{toscaSimpleYAML13, tosca20}
}
