package grammars

import (
	"github.com/tliron/go-ard"
	"github.com/tliron/go-puccini/tosca/parsing"
)

func DetectGrammar(context *parsing.Context) bool {
	if context.Grammar == nil {
		versionContext, version := DetectGrammarVersion(context)
		switch {
		case versionContext == nil:
			context.FieldChild(definitionsVersionKey, nil).ReportKeynameMissing()
		case version == nil:
			// Type validation already reported the read-phase problem.
		default:
			if grammar, ok := grammars[*version]; ok {
				context.Grammar = grammar
			} else {
				versionContext.ReportToscaDefinitionsVersionUnsupported(SupportedDefinitionsVersions())
			}
		}
	}
	return context.Grammar != nil
}

func GetGrammar(context *parsing.Context) (*parsing.Grammar, *parsing.Context) {
	if versionContext, version := DetectGrammarVersion(context); versionContext != nil {
		if version != nil {
			if grammar, ok := grammars[*version]; ok {
				return grammar, nil
			}
		}
		return nil, versionContext
	}
	return nil, nil
}

func CompatibleGrammars(context1 *parsing.Context, context2 *parsing.Context) bool {
	grammar1, _ := GetGrammar(context1)
	grammar2, _ := GetGrammar(context2)
	return grammar1 == grammar2
}

func DetectGrammarVersion(context *parsing.Context) (*parsing.Context, *string) {
	if versionContext, ok := context.GetFieldChild(definitionsVersionKey); ok {
		if versionContext.ValidateType(ard.TypeString) {
			return versionContext, versionContext.ReadString()
		}
		return versionContext, nil
	}
	return nil, nil
}

func GetImplicitImportSpec(context *parsing.Context) (*parsing.ImportSpec, bool) {
	if _, version := DetectGrammarVersion(context); version != nil {
		if path, ok := implicitProfilePaths[*version]; ok {
			if url, err := context.URL.Context().NewValidInternalURL(path); err == nil {
				return &parsing.ImportSpec{URL: url, NameTransformer: nil, Implicit: true}, true
			} else {
				context.ReportError(err)
			}
		}
	}

	return nil, false
}
