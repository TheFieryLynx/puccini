package tosca_1_3_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tliron/exturl"
	"github.com/tliron/go-puccini/normal"
	"github.com/tliron/go-puccini/tosca/parser"
)

// Direct remediation tests observe the existing phase API at each boundary.
// This helper does not change parser behavior or the general corpus runner.
type reAuditResult struct {
	Context  *parser.Context
	Template *normal.ServiceTemplate
	Phase    string
	Problems string
}

func reAuditParse(t *testing.T, source string) reAuditResult {
	t.Helper()
	path := filepath.Join(t.TempDir(), "template.yaml")
	if err := os.WriteFile(path, []byte(source), 0o600); err != nil {
		t.Fatal(err)
	}
	u := exturl.NewContext()
	t.Cleanup(func() { _ = u.Release() })
	p := parser.NewParser().NewContext()
	p.URL = u.NewFileURL(filepath.ToSlash(path))
	r := reAuditResult{Context: p}
	check := func(phase string) bool {
		if p.Root == nil {
			r.Phase = phase
			return true
		}
		p.MergeProblems()
		if !p.GetProblems().Empty() {
			r.Phase, r.Problems = phase, p.GetProblems().ToString(false)
			return true
		}
		return false
	}
	ok := p.ReadRoot(context.Background(), p.URL, p.Bases, "")
	if check("read") || !ok {
		r.Phase = "read"
		return r
	}
	p.AddNamespaces()
	p.LookupNames()
	if check("namespace") {
		return r
	}
	p.AddHierarchies()
	if check("hierarchy") {
		return r
	}
	p.Inherit(nil)
	if check("inheritance") {
		return r
	}
	p.SetInputs(p.Inputs)
	p.Render()
	if check("rendering") {
		return r
	}
	r.Template, ok = p.Normalize()
	if check("normalization") || !ok {
		r.Phase = "normalization"
	}
	return r
}

func reAuditReject(t *testing.T, source, phase string, fragments ...string) reAuditResult {
	t.Helper()
	r := reAuditParse(t, source)
	if r.Phase != phase {
		t.Fatalf("expected %s rejection, got phase %q:\n%s", phase, r.Phase, r.Problems)
	}
	for _, fragment := range fragments {
		if !strings.Contains(r.Problems, fragment) {
			t.Fatalf("expected target diagnostic %q:\n%s", fragment, r.Problems)
		}
	}
	return r
}
