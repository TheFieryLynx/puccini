package tosca_1_3_corpus_test

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/tliron/exturl"
	"github.com/tliron/go-puccini/normal"
	"github.com/tliron/go-puccini/tosca/parser"
)

// Identities are test-side parser API boundaries, never diagnostic text.
// CSAR opening/metadata validation occurs inside ReadRoot. Normalization is
// an internal representation phase, not an invented normative error category.
type processingPhase string

const (
	phaseRead          processingPhase = "read"
	phaseNamespaces    processingPhase = "namespaces"
	phaseHierarchy     processingPhase = "hierarchy"
	phaseInheritance   processingPhase = "inheritance"
	phaseRendering     processingPhase = "rendering"
	phaseNormalization processingPhase = "normalization"
)

type corpusResult struct {
	Template *normal.ServiceTemplate
	Parser   *parser.Context
	Phase    processingPhase
	Problems string
	Err      error
	Reached  []processingPhase
}

func parseCorpusPhases(t *testing.T, path string) corpusResult {
	t.Helper()
	absolute, err := filepath.Abs(path)
	if err != nil {
		t.Fatal(err)
	}
	u := exturl.NewContext()
	t.Cleanup(func() { _ = u.Release() })
	p := parser.NewParser().NewContext()
	p.URL = u.NewFileURL(filepath.ToSlash(absolute))
	r := corpusResult{Parser: p}
	check := func(phase processingPhase, ok bool) bool {
		r.Phase = phase
		r.Reached = append(r.Reached, phase)
		p.MergeProblems()
		r.Problems = p.GetProblems().ToString(false)
		if !ok || !p.GetProblems().Empty() {
			r.Err = fmt.Errorf("processor failed at %s", phase)
			return true
		}
		return false
	}
	if check(phaseRead, p.ReadRoot(context.Background(), p.URL, p.Bases, "")) {
		return r
	}
	p.AddNamespaces()
	p.LookupNames()
	if check(phaseNamespaces, true) {
		return r
	}
	p.AddHierarchies()
	if check(phaseHierarchy, true) {
		return r
	}
	p.Inherit(nil)
	if check(phaseInheritance, true) {
		return r
	}
	p.SetInputs(p.Inputs)
	p.Render()
	if check(phaseRendering, true) {
		return r
	}
	var ok bool
	r.Template, ok = p.Normalize()
	check(phaseNormalization, ok)
	return r
}
func checkExpectedPhase(expected string, r corpusResult) error {
	if r.Err == nil {
		return fmt.Errorf("expected rejection at %s, got acceptance", expected)
	}
	if processingPhase(expected) != r.Phase {
		return fmt.Errorf("wrong phase: expected %s, actual %s", expected, r.Phase)
	}
	return nil
}
