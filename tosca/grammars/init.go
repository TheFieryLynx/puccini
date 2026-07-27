package grammars

import (
	"github.com/tliron/go-puccini/tosca/grammars/tosca_v1_3"
	"github.com/tliron/go-puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/go-puccini/tosca/parsing"

	_ "github.com/tliron/go-puccini/assets/tosca/profiles"
)

func init() {
	grammars = map[string]*parsing.Grammar{
		toscaSimpleYAML13: &tosca_v1_3.Grammar,
		tosca20:           &tosca_v2_0.Grammar,
	}
	implicitProfilePaths = map[string]string{
		toscaSimpleYAML13: "/profiles/simple/1.3/profile.yaml",
		tosca20:           "/profiles/implicit/2.0/profile.yaml",
	}
}
