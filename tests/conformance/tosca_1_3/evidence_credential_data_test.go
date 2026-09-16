package tosca_1_3_test

import (
	"github.com/tliron/go-puccini/normal"
	"testing"
)

// TOSCA 1.3 §5.3.6.1–2 Credential and §4.3.3.1 token grammar.
// Expected: string token maps are data, including collection entries; sequence
// token calls remain functions. Categories: positive, negative, rendering,
// nearest-valid, collections, inheritance, normalization. No runtime token validation.
func TestEvidenceCredentialData(t *testing.T) {
	for _, tc := range []struct{ name, definition, value string }{
		{"direct", "{type: tosca.datatypes.Credential}", "{token: secret}"},
		{"list", "{type: list, entry_schema: tosca.datatypes.Credential}", "[{token: secret}]"},
		{"map", "{type: map, entry_schema: tosca.datatypes.Credential}", "{item: {token: secret}}"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := reAuditParse(t, credentialDataSource(tc.definition, tc.value))
			if r.Phase != "" {
				t.Fatalf("%s: %s", r.Phase, r.Problems)
			}
			v := r.Template.NodeTemplates["n"].Properties["p"]
			if tc.name == "list" {
				v = v.(*normal.List).Entries[0]
			}
			if tc.name == "map" {
				v = v.(*normal.Map).Entries[0]
			}
			fields := map[string]any{}
			for _, entry := range v.(*normal.Map).Entries {
				p := entry.(*normal.Primitive)
				fields[p.Key.(*normal.Primitive).Primitive.(string)] = p.Primitive
			}
			if fields["token"] != "secret" || fields["token_type"] != "password" {
				t.Fatalf("effective Credential: %#v", fields)
			}
		})
	}
	valid := credentialDataSource("{type: string}", "secret")
	r := reAuditParse(t, valid)
	if r.Phase != "" {
		t.Fatal(r.Problems)
	}
	if p := r.Template.NodeTemplates["n"].Properties["p"].(*normal.Primitive); p.Primitive != "secret" {
		t.Fatal(p)
	}
	reAuditReject(t, credentialDataSource("{type: string}", "{token: secret}"), "rendering", "instead of", "p")
	r = reAuditParse(t, credentialDataSource("{type: string}", `{token: ['a:b', ':', 1]}`))
	if r.Phase != "" {
		t.Fatal(r.Problems)
	}
	if _, ok := r.Template.NodeTemplates["n"].Properties["p"].(*normal.FunctionCall); !ok {
		t.Fatal("sequence token ceased to be a function")
	}
	reAuditReject(t, credentialDataSource("{type: string}", `{token: ['a:b', ':']}`), "read", "substring_index")
}
func credentialDataSource(definition, value string) string {
	return `tosca_definitions_version: tosca_simple_yaml_1_3
node_types:
  Holder:
    derived_from: tosca.nodes.Root
    properties:
      p: ` + definition + `
topology_template:
  node_templates:
    n:
      type: Holder
      properties:
        p: ` + value + "\n"
}
