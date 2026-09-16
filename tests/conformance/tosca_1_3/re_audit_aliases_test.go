package tosca_1_3_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	v13 "github.com/tliron/go-puccini/tosca/grammars/tosca_v1_3"
	"github.com/tliron/go-puccini/tosca/parsing"
)

var reAuditAliases = []struct{ category, field, uri, short, qualified string }{
	{"artifact_types", "ArtifactTypes", "tosca.artifacts.Implementation.Bash", "Bash", "tosca:Bash"},
	{"artifact_types", "ArtifactTypes", "tosca.artifacts.Implementation.Python", "Python", "tosca:Python"},
	{"capability_types", "CapabilityTypes", "tosca.capabilities.network.Bindable", "network.Bindable", "tosca:network.Bindable"},
	{"node_types", "NodeTypes", "tosca.nodes.Abstract.Storage", "AbstractStorage", "tosca:Abstract.Storage"},
	{"node_types", "NodeTypes", "tosca.nodes.Storage.ObjectStorage", "ObjectStorage", "tosca:ObjectStorage"},
	{"node_types", "NodeTypes", "tosca.nodes.Storage.BlockStorage", "BlockStorage", "tosca:BlockStorage"},
	{"relationship_types", "RelationshipTypes", "tosca.relationships.network.BindsTo", "network.BindsTo", "tosca:BindsTo"},
}

// TOSCA 1.3 §5.2 Naming Conventions; tables §5.4.4.3/4, §5.5.14,
// §5.9.9/10/11 and §8.5.5; §5.2.1 Case Sensitivity, F16.
// Positive full URI, shorthand, qualified names and custom derivation; negative
// case siblings must fail at namespace resolution. All aliases resolve to the
// same effective type object as the full URI, not merely to a similarly named type.
func TestReAuditNormativeAliases(t *testing.T) {
	for _, tc := range reAuditAliases {
		for _, name := range []string{tc.uri, tc.short, tc.qualified} {
			t.Run(name, func(t *testing.T) {
				source := fmt.Sprintf("tosca_definitions_version: tosca_simple_yaml_1_3\n%s:\n  Custom: { derived_from: '%s' }\n", tc.category, name)
				r := reAuditParse(t, source)
				if r.Phase != "" {
					t.Fatalf("%s: %s", r.Phase, r.Problems)
				}
				f := r.Context.Root.EntityPtr.(*v13.ServiceFile)
				types := reflect.ValueOf(f.File).Elem().FieldByName(tc.field)
				if types.Len() != 1 {
					t.Fatalf("custom types: %d", types.Len())
				}
				parent := types.Index(0).Elem().FieldByName("Parent").Interface()
				if parsing.GetContext(parent).Name != tc.uri {
					t.Fatalf("resolved parent %s, want %s", parsing.GetContext(parent).Name, tc.uri)
				}
				ns := f.Context.Namespace
				full, ok := ns.LookupForType(tc.uri, reflect.TypeOf(parent))
				if !ok || full != parent {
					t.Fatal("URI and alias identify different types")
				}
				reAuditReject(t, strings.Replace(source, "'"+name+"'", "'"+strings.ToLower(name)+"'", 1), "namespace", strings.ToLower(name))
			})
		}
	}
}
