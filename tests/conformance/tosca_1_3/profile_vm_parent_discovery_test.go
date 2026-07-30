package tosca_1_3_test

import (
	"testing"

	"github.com/tliron/go-puccini/tosca/grammars/tosca_v1_3"
	"github.com/tliron/go-puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/go-puccini/tosca/parser"
)

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 5.4.3.4.1 Definition
// Expected: the ordinary processor path loads Deployment.Image.VM with
// derived_from set to tosca.artifacts.Deployment.Image.
// Category: positive, direct, normative profile, hierarchy
func TestProfileDeploymentImageVMParent(t *testing.T) {
	parserContext, release := parseNormativeProfile(t)
	defer release()
	artifact := findProfileArtifactType(
		t, parserContext, "tosca.artifacts.Deployment.Image.VM")
	if artifact.ParentName == nil ||
		*artifact.ParentName != "tosca.artifacts.Deployment.Image" {
		t.Fatalf("Deployment.Image.VM derived_from = %v, want tosca.artifacts.Deployment.Image",
			artifact.ParentName)
	}
}

func findProfileArtifactType(
	t *testing.T,
	parserContext *parser.Context,
	name string,
) *tosca_v2_0.ArtifactType {
	t.Helper()
	for _, file := range parserContext.Files {
		if profile, ok := file.EntityPtr.(*tosca_v1_3.File); ok {
			for _, type_ := range profile.ArtifactTypes {
				if type_.Name == name {
					return type_
				}
			}
		}
	}
	t.Fatalf("normative artifact type %q is absent", name)
	return nil
}
