package tosca_1_3_corpus_test

import (
	"crypto/sha256"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

type caseExecution struct {
	ID         string   `yaml:"id"`
	Status     string   `yaml:"status"`
	Phase      string   `yaml:"phase"`
	Assertions []string `yaml:"executed_assertions"`
	SHA256     string   `yaml:"fixture_sha256"`
}
type evidenceExecution struct {
	Fingerprint string          `yaml:"source_fingerprint"`
	Cases       []caseExecution `yaml:"cases"`
}

func evidenceReasons(id string, reg predicateRegistry, cases []corpusCase, mutations mutationCoverage, execution evidenceExecution) []string {
	r, ok := registryIndex(reg)[id]
	if !ok {
		return []string{"no-reviewed-primary-predicate", "historical-labels-or-unreviewed-direct-tests-only"}
	}
	why := map[string]bool{}
	if !r.CompletePredicate {
		why["only-subset-of-historical-predicate-reviewed"] = true
	}
	linked, pos, neg := false, false, false
	es := map[string]caseExecution{}
	for _, e := range execution.Cases {
		es[e.ID] = e
	}
	for _, c := range cases {
		if !contains(c.Coverage.PrimaryRequirements, id) {
			continue
		}
		linked = true
		owned := []ownedAssertion{}
		for _, a := range c.Assertions {
			if contains(a.RequirementIDs, id) {
				owned = append(owned, a)
			}
		}
		if len(owned) == 0 {
			why["primary-without-owned-assertion"] = true
		}
		e := es[c.ID]
		if e.Status != "passed" {
			why["not-executed-owned-assertions"] = true
		}
		for _, a := range owned {
			if !contains(e.Assertions, a.ID) {
				why["not-executed-owned-assertions"] = true
			}
		}
		if c.Expected.Accepted {
			pos = true
			if strings.HasPrefix(r.EvidenceKind, "semantic-") || r.EvidenceKind == "normalization" {
				semantic := false
				for _, a := range owned {
					if a.Kind == "value" || a.Kind == "normalization" {
						semantic = true
					}
				}
				if !semantic {
					why["semantic-positive-lacks-value"] = true
				}
			}
		} else {
			neg = true
			phase, category, entity := false, false, false
			for _, a := range owned {
				if a.Kind == "phase" {
					phase = true
				}
				if a.Kind == "diagnostic" && a.Path == "category" {
					category = true
				}
				if a.Kind == "diagnostic" && a.Path == "entity" {
					entity = true
				}
			}
			if e.Phase != c.Expected.Phase || !phase {
				why["phase-not-proven"] = true
			}
			if !category || !entity {
				why["targeted-diagnostic-not-proven"] = true
			}
			sibling := false
			for _, s := range cases {
				if s.ID == c.Related.NearestValid && s.Expected.Accepted && es[s.ID].Status == "passed" {
					sibling = true
				}
			}
			if !sibling {
				why["nearest-valid-not-proven"] = true
			}
		}
	}
	if !linked {
		why["no-primary-fixture"] = true
	}
	if r.PositiveRequired && !pos {
		why["missing-positive"] = true
	}
	if r.NegativeRequired && !neg {
		why["missing-negative"] = true
	}
	if len(r.MutationPaths) == 0 {
		why["no-reviewed-mutation-path"] = true
	}
	for _, path := range r.MutationPaths {
		caught := false
		for _, m := range mutations.Paths {
			if m.ID == path && m.Status == "caught" && contains(m.Requirements, id) {
				caught = true
			}
		}
		if !caught {
			why["relevant-owned-mutation-not-caught"] = true
		}
	}
	return sortedKeys(why)
}
func executionFingerprint(t *testing.T) string {
	t.Helper()
	root := repositoryRoot(t)
	paths := []string{}
	for _, dir := range []string{"tests/corpus/tosca_1_3", "tosca", "normal", "clout", "assets/tosca"} {
		err := filepath.WalkDir(filepath.Join(root, dir), func(p string, e fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			ext := filepath.Ext(p)
			included := ext == ".go" || ext == ".yaml" || ext == ".js"
			if dir == "tests/corpus/tosca_1_3" {
				included = included || ext == ".yml" || ext == ".json" || ext == ".zip" || ext == ".csar"
			}
			if !e.IsDir() && included {
				paths = append(paths, p)
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	sort.Strings(paths)
	h := sha256.New()
	for _, p := range paths {
		b, err := os.ReadFile(p)
		if err != nil {
			t.Fatal(err)
		}
		rel, err := filepath.Rel(root, p)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Fprintf(h, "%s:%x\n", filepath.ToSlash(rel), sha256.Sum256(b))
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}
func TestEvidenceRecheck(t *testing.T) {
	root := repositoryRoot(t)
	base := filepath.Join(root, "docs/conformance/tosca-1.3")
	var report struct {
		Fingerprint string `yaml:"source_fingerprint"`
		Summary     struct {
			Historical   int `yaml:"historical_records"`
			Verified     int `yaml:"verified"`
			Insufficient int `yaml:"insufficient_evidence"`
		} `yaml:"summary"`
		Requirements []struct {
			ID      string   `yaml:"requirement_id"`
			Status  string   `yaml:"verification_status"`
			Reasons []string `yaml:"reasons"`
		} `yaml:"requirements"`
	}
	loadYAML(t, filepath.Join(base, "re-audit/evidence-recheck.yaml"), &report)
	var execution evidenceExecution
	loadYAML(t, filepath.Join(base, "re-audit/evidence-execution.yaml"), &execution)
	if report.Fingerprint != executionFingerprint(t) || report.Fingerprint != execution.Fingerprint {
		t.Fatal("stale execution fingerprint; regenerate evidence with --record-run")
	}
	var reg predicateRegistry
	loadYAML(t, filepath.Join(corpusRoot(t), "evidence-predicates.yaml"), &reg)
	var mutations mutationCoverage
	loadYAML(t, filepath.Join(base, "mutation-coverage.yaml"), &mutations)
	cases := loadManifest(t, corpusRoot(t)).Cases
	for _, c := range loadNonMUSTManifest(t, corpusRoot(t)).Cases {
		cases = append(cases, corpusCase{ID: c.ID, File: c.File, Expected: corpusExpected{Accepted: c.Expected.Accepted, Phase: c.Expected.Phase}, Coverage: c.Coverage, Assertions: c.Assertions, Related: c.Related})
	}
	es := map[string]caseExecution{}
	for _, e := range execution.Cases {
		if _, ok := es[e.ID]; ok {
			t.Fatal("duplicate execution")
		}
		es[e.ID] = e
	}
	for _, c := range cases {
		b, err := os.ReadFile(filepath.Join(corpusRoot(t), c.File))
		if err != nil {
			t.Fatal(err)
		}
		if es[c.ID].SHA256 != fmt.Sprintf("%x", sha256.Sum256(b)) {
			t.Fatalf("stale fixture %s", c.ID)
		}
	}
	var historical coverageCatalog
	loadYAML(t, filepath.Join(base, "coverage.yaml"), &historical)
	ids := map[string]bool{}
	for _, r := range historical.Requirements {
		if r.IncludedInMUSTDenomator {
			ids[r.RequirementID] = true
		}
	}
	if len(report.Requirements) != 223 || len(ids) != 223 {
		t.Fatal("historical membership changed")
	}
	verified := 0
	for _, r := range report.Requirements {
		if !ids[r.ID] {
			t.Fatalf("unexpected/duplicate ID %s", r.ID)
		}
		delete(ids, r.ID)
		reasons := evidenceReasons(r.ID, reg, cases, mutations, execution)
		want := "verified"
		if len(reasons) > 0 {
			want = "insufficient-evidence"
		} else {
			verified++
		}
		if r.Status != want || strings.Join(r.Reasons, "|") != strings.Join(reasons, "|") {
			t.Errorf("%s incorrect recheck %s: %v; want %s: %v", r.ID, r.Status, r.Reasons, want, reasons)
		}
	}
	if report.Summary.Historical != 223 || report.Summary.Verified != verified || report.Summary.Insufficient != 223-verified {
		t.Fatal("summary disagrees with independently recomputed evidence")
	}
	for file, want := range map[string]string{"requirements.yaml": "1861a9652dd659636702f29737818ac6f1c69355ee442e1bebe98739b21aabf0", "coverage.yaml": "1238b594dfeb73c07f517efe8315b8e38742bf24db5350c1a1a5bf4ded73e11a"} {
		b, err := os.ReadFile(filepath.Join(base, file))
		if err != nil {
			t.Fatal(err)
		}
		if fmt.Sprintf("%x", sha256.Sum256(b)) != want {
			t.Fatalf("historical %s changed outside catalog review", file)
		}
	}
}
func TestEvidenceVerificationBar(t *testing.T) {
	makeEvidence := func() (predicateRegistry, []corpusCase, mutationCoverage, evidenceExecution) {
		r := predicateRegistry{Requirements: []reviewedPredicate{{ID: "R", EvidenceKind: "semantic-value", CompletePredicate: true, PositiveRequired: true, MutationPaths: []string{"path"}}}}
		c := []corpusCase{{ID: "C", Expected: corpusExpected{Accepted: true}, Coverage: corpusCoverage{PrimaryRequirements: []string{"R"}}, Assertions: []ownedAssertion{{ID: "A", RequirementIDs: []string{"R"}, Kind: "value"}}}}
		m := mutationCoverage{Paths: []mutationPath{{ID: "path", Requirements: []string{"R"}, Status: "caught"}}}
		e := evidenceExecution{Cases: []caseExecution{{ID: "C", Status: "passed", Assertions: []string{"A"}}}}
		return r, c, m, e
	}
	r, c, m, e := makeEvidence()
	if why := evidenceReasons("R", r, c, m, e); len(why) != 0 {
		t.Fatal(why)
	}
	tests := map[string]func(*predicateRegistry, []corpusCase, *mutationCoverage, *evidenceExecution){
		"supporting-only": func(r *predicateRegistry, c []corpusCase, m *mutationCoverage, e *evidenceExecution) {
			c[0].Coverage.PrimaryRequirements = nil
			c[0].Coverage.SupportingRequirements = []string{"R"}
		},
		"accept-only": func(r *predicateRegistry, c []corpusCase, m *mutationCoverage, e *evidenceExecution) {
			c[0].Assertions[0].Kind = "accept"
		},
		"unexecuted": func(r *predicateRegistry, c []corpusCase, m *mutationCoverage, e *evidenceExecution) {
			e.Cases[0].Status = "unexecuted"
		},
		"partial-predicate": func(r *predicateRegistry, c []corpusCase, m *mutationCoverage, e *evidenceExecution) {
			r.Requirements[0].CompletePredicate = false
		},
		"uncaught-path": func(r *predicateRegistry, c []corpusCase, m *mutationCoverage, e *evidenceExecution) {
			m.Paths[0].Status = "missed"
		},
		"unrelated-path": func(r *predicateRegistry, c []corpusCase, m *mutationCoverage, e *evidenceExecution) {
			m.Paths[0].Requirements = []string{"other"}
		},
		"missing-negative": func(r *predicateRegistry, c []corpusCase, m *mutationCoverage, e *evidenceExecution) {
			r.Requirements[0].NegativeRequired = true
		},
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			r, c, m, e := makeEvidence()
			mutate(&r, c, &m, &e)
			if len(evidenceReasons("R", r, c, m, e)) == 0 {
				t.Fatal("insufficient evidence called verified")
			}
		})
	}
}
