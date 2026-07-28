package tosca_v1_3

import (
	"github.com/tliron/go-ard"
	"github.com/tliron/go-puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/go-puccini/tosca/parsing"
)

//
// Version
//
// TOSCA Simple Profile in YAML 1.3, section 3.3.2.5:
// zero denotes an unspecified version and cannot have a qualifier.
//

// ReadVersion applies the TOSCA 1.3-only effective-value policy after using
// the shared lexical reader. The shared TOSCA 2.0 reader remains unchanged.
func ReadVersion(context *parsing.Context) parsing.EntityPtr {
	// The normative zero spellings include "0", while the ordinary lexical
	// grammar starts with major.minor. Canonicalize only for lexical parsing
	// and restore the source value immediately so the YAML document is not
	// mutated.
	var original any
	if context.Is(ard.TypeString) && context.Data.(string) == "0" {
		original = context.Data
		context.Data = "0.0"
		defer func() {
			context.Data = original
		}()
	}

	version := tosca_v2_0.ReadVersion(context).(*tosca_v2_0.Version)
	if original != nil {
		version.OriginalString = original.(string)
	}

	if version.Major == 0 && version.Minor == 0 && version.Fix == 0 {
		if version.Qualifier != "" {
			context.ReportValueMalformed("version", "zero version must not have a qualifier")
			return version
		}

		// An empty canonical string is the effective unspecified-version
		// sentinel. OriginalString preserves the user's lexical spelling.
		version.CanonicalString = ""
	}

	return version
}
