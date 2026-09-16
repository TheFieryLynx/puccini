package tosca_1_3_corpus_test

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type mutationCoverage struct {
	Paths []mutationPath `yaml:"production_paths"`
}
type mutationPath struct {
	ID            string            `yaml:"id"`
	Category      string            `yaml:"category"`
	Symbols       []string          `yaml:"symbols"`
	Requirements  []string          `yaml:"requirements"`
	Cases         []string          `yaml:"detecting_cases"`
	Assertions    []string          `yaml:"detecting_assertions"`
	Detected      []string          `yaml:"detected_assertions"`
	Status        string            `yaml:"status"`
	SourceSHA256  string            `yaml:"source_sha256"`
	FixtureSHA256 map[string]string `yaml:"fixture_sha256"`
	Patch         struct {
		File   string `yaml:"file"`
		Before string `yaml:"before"`
		After  string `yaml:"after"`
	} `yaml:"patch"`
}

var requiredMutationCategories = []string{"structural-grammar", "namespaces-imports", "hierarchy", "inheritance", "refinement", "constraints", "assignments", "intrinsic-functions", "requirements-capabilities", "interfaces-operations", "workflows", "substitution-mappings", "normative-profile", "csar", "normalization"}

func validateMutationCoverage(m mutationCoverage, cases []corpusCase, read func(string) ([]byte, error)) error {
	ids := map[string]corpusCase{}
	for _, c := range cases {
		ids[c.ID] = c
	}
	categories := map[string]bool{}
	paths := map[string]bool{}
	for _, p := range m.Paths {
		if p.ID == "" || paths[p.ID] {
			return fmt.Errorf("duplicate/empty mutation path")
		}
		paths[p.ID] = true
		if p.Status != "caught" || len(p.Symbols) == 0 || len(p.Requirements) == 0 || len(p.Cases) == 0 || len(p.Assertions) == 0 || len(p.Detected) == 0 {
			return fmt.Errorf("%s lacks caught evidence", p.ID)
		}
		source, err := read(p.Patch.File)
		if err != nil {
			return err
		}
		if fmt.Sprintf("%x", sha256.Sum256(source)) != p.SourceSHA256 {
			return fmt.Errorf("stale mutation source %s", p.ID)
		}
		for _, cid := range p.Cases {
			c, ok := ids[cid]
			if !ok {
				return fmt.Errorf("missing detecting fixture %s", cid)
			}
			b, err := read("tests/corpus/tosca_1_3/" + c.File)
			if err != nil {
				return err
			}
			if fmt.Sprintf("%x", sha256.Sum256(b)) != p.FixtureSHA256[cid] {
				return fmt.Errorf("stale detecting fixture %s", cid)
			}
		}
		detectedRequirements := map[string]bool{}
		for _, ref := range p.Assertions {
			pos := strings.LastIndex(ref, "/")
			if pos < 0 {
				return fmt.Errorf("malformed assertion reference %s", ref)
			}
			cid, aid := ref[:pos], ref[pos+1:]
			if !contains(p.Cases, cid) {
				return fmt.Errorf("assertion not in detecting cases")
			}
			c := ids[cid]
			found := false
			for _, a := range c.Assertions {
				if a.ID != aid {
					continue
				}
				found = true
				for _, id := range a.RequirementIDs {
					if !contains(c.Coverage.PrimaryRequirements, id) || !contains(p.Requirements, id) {
						return fmt.Errorf("supporting or unrelated mutation evidence")
					}
					if contains(p.Detected, ref) {
						detectedRequirements[id] = true
					}
				}
			}
			if !found {
				return fmt.Errorf("missing detecting assertion %s", ref)
			}
		}
		for _, ref := range p.Detected {
			if !contains(p.Assertions, ref) {
				return fmt.Errorf("unmapped detected assertion")
			}
		}
		for _, r := range p.Requirements {
			if !detectedRequirements[r] {
				return fmt.Errorf("%s has no detected owned assertion", r)
			}
		}
		categories[p.Category] = true
	}
	for _, category := range requiredMutationCategories {
		if !categories[category] {
			return fmt.Errorf("missing required mutation category %s", category)
		}
	}
	return nil
}
func TestTOSCA13MutationCoverage(t *testing.T) {
	var m mutationCoverage
	loadYAML(t, filepath.Join(repositoryRoot(t), "docs/conformance/tosca-1.3/mutation-coverage.yaml"), &m)
	if err := validateMutationCoverage(m, loadManifest(t, corpusRoot(t)).Cases, func(p string) ([]byte, error) { return os.ReadFile(filepath.Join(repositoryRoot(t), p)) }); err != nil {
		t.Fatal(err)
	}
}
func TestMutationCoverageIntegrity(t *testing.T) {
	makeMap := func() (mutationCoverage, []corpusCase) {
		c := corpusCase{ID: "case", File: "fixture.yaml", Coverage: corpusCoverage{PrimaryRequirements: []string{"R"}}, Assertions: []ownedAssertion{{ID: "A", RequirementIDs: []string{"R"}}}}
		m := mutationCoverage{}
		digest := fmt.Sprintf("%x", sha256.Sum256([]byte("bytes")))
		for _, category := range requiredMutationCategories {
			m.Paths = append(m.Paths, mutationPath{ID: category, Category: category, Symbols: []string{"symbol"}, Requirements: []string{"R"}, Cases: []string{"case"}, Assertions: []string{"case/A"}, Detected: []string{"case/A"}, Status: "caught", SourceSHA256: digest, FixtureSHA256: map[string]string{"case": digest}})
		}
		return m, []corpusCase{c}
	}
	read := func(string) ([]byte, error) { return []byte("bytes"), nil }
	m, c := makeMap()
	if err := validateMutationCoverage(m, c, read); err != nil {
		t.Fatal(err)
	}
	tests := map[string]func(*mutationCoverage, []corpusCase){
		"missing-category":  func(m *mutationCoverage, c []corpusCase) { m.Paths = m.Paths[1:] },
		"missing-fixture":   func(m *mutationCoverage, c []corpusCase) { m.Paths[0].Cases = []string{"absent"} },
		"missing-assertion": func(m *mutationCoverage, c []corpusCase) { m.Paths[0].Assertions = []string{"case/absent"} },
		"supporting-only": func(m *mutationCoverage, c []corpusCase) {
			c[0].Coverage.PrimaryRequirements = nil
			c[0].Coverage.SupportingRequirements = []string{"R"}
		},
		"unrelated-requirement": func(m *mutationCoverage, c []corpusCase) { m.Paths[0].Requirements = []string{"other"} },
		"missed-mutation":       func(m *mutationCoverage, c []corpusCase) { m.Paths[0].Status = "missed" },
		"uncaught-assertion":    func(m *mutationCoverage, c []corpusCase) { m.Paths[0].Detected = nil },
		"stale-path":            func(m *mutationCoverage, c []corpusCase) { m.Paths[0].SourceSHA256 = "old" },
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			m, c := makeMap()
			mutate(&m, c)
			if validateMutationCoverage(m, c, read) == nil {
				t.Fatal("invalid evidence accepted")
			}
		})
	}
}

func TestMutationReplayProvenance(t *testing.T) {
	var m mutationCoverage
	loadYAML(t, filepath.Join(repositoryRoot(t), "docs/conformance/tosca-1.3/mutation-coverage.yaml"), &m)
	var report struct {
		Results []struct {
			ID           string   `yaml:"id"`
			Status       string   `yaml:"status"`
			Baseline     bool     `yaml:"baseline_passed"`
			ExitCode     int      `yaml:"exit_code"`
			Restored     bool     `yaml:"restored"`
			Detected     []string `yaml:"detected_assertions"`
			SourceSHA256 string   `yaml:"source_sha256"`
			Output       string   `yaml:"failure_output"`
		} `yaml:"results"`
	}
	loadYAML(t, filepath.Join(repositoryRoot(t), "docs/conformance/tosca-1.3/re-audit/evidence-mutation-replay.yaml"), &report)
	for _, p := range m.Paths {
		found := false
		for _, r := range report.Results {
			if r.ID != p.ID {
				continue
			}
			found = true
			if r.Status != "caught" || !r.Baseline || r.ExitCode != 1 || !r.Restored || r.SourceSHA256 != p.SourceSHA256 {
				t.Errorf("invalid replay provenance %s", p.ID)
			}
			for _, a := range p.Detected {
				if !contains(r.Detected, a) || !strings.Contains(r.Output, a+":") {
					t.Errorf("%s lacks observed assertion failure %s", p.ID, a)
				}
			}
		}
		if !found {
			t.Errorf("missing replay %s", p.ID)
		}
	}
}
