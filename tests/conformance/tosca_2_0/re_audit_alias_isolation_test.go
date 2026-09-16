package tosca_2_0_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tliron/go-puccini/normal"
	"github.com/tliron/go-puccini/tests/conformance/internal/testsupport"
)

// TOSCA 2.0 import/name-resolution regression (clause 6.8 Imports and Namespaces).
// This deliberately uses Puccini's existing normative metadata extension as a
// sentinel: 1.3 alias additions must not change the 2.0 transformer's behavior.
// Positive imported custom type/property semantics; negative 1.3-only alias.
func TestReAuditAliasPolicyIsolation(t *testing.T) {
	path := filepath.Join(t.TempDir(), "types.yaml")
	err := os.WriteFile(path, []byte(`tosca_definitions_version: tosca_2_0
node_types:
  tosca.nodes.Storage.BlockStorage:
    metadata: { tosca.normative: 'true' }
    properties:
      p: { type: integer, default: 7 }
`), 0600)
	if err != nil {
		t.Fatal(err)
	}
	source := fmt.Sprintf(`tosca_definitions_version: tosca_2_0
imports: [{ url: '%s' }]
service_template:
  node_templates:
    n: { type: Storage.BlockStorage }
`, filepath.ToSlash(path))
	st, problems, err := testsupport.ParseSource(t, source)
	if err != nil {
		t.Fatalf("2.0 import: %v %s", err, problems)
	}
	p := st.NodeTemplates["n"].Properties["p"].(*normal.Primitive)
	if p.Primitive != 7 {
		t.Fatalf("2.0 default changed: %#v", p)
	}
	_, problems, err = testsupport.ParseSource(t, strings.Replace(source, "type: Storage.BlockStorage", "type: BlockStorage", 1))
	if err == nil || !strings.Contains(problems, "unknown node type: BlockStorage") {
		t.Fatalf("1.3 alias leaked into 2.0: %v %s", err, problems)
	}
}
