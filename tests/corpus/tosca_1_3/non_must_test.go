package tosca_1_3_corpus_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type nonMUSTManifest struct {
	Cases []nonMUSTCase `yaml:"cases"`
}

type nonMUSTCase struct {
	ID              string `yaml:"id"`
	File            string `yaml:"file"`
	ConformanceTier string `yaml:"conformance_tier"`
	Specification   struct {
		Sections              []string `yaml:"sections"`
		NonMUSTRequirementIDs []string `yaml:"non_must_requirement_ids"`
	} `yaml:"specification"`
	Policy struct {
		Support   string `yaml:"support"`
		Rationale string `yaml:"rationale"`
	} `yaml:"policy"`
	Classification struct {
		Kind     string   `yaml:"kind"`
		Category string   `yaml:"category"`
		Features []string `yaml:"features"`
	} `yaml:"classification"`
	Expected struct {
		Accepted   bool              `yaml:"accepted"`
		Phase      string            `yaml:"phase"`
		Diagnostic *corpusDiagnostic `yaml:"diagnostic"`
		Assertions []string          `yaml:"assertions"`
	} `yaml:"expected"`
}

type nonMUSTCatalog struct {
	Requirements []struct {
		ID                  string   `yaml:"id"`
		FixtureIDs          []string `yaml:"fixture_ids"`
		ProcessorApplicable bool     `yaml:"processor_applicable"`
	} `yaml:"requirements"`
}

var allowedNonMUSTTiers = map[string]bool{
	"processor-should":       true,
	"processor-may":          true,
	"processor-default":      true,
	"implementation-defined": true,
	"compatibility":          true,
}

func TestTOSCA13NonMUSTManifestIntegrity(t *testing.T) {
	root := corpusRoot(t)
	manifest := loadNonMUSTManifest(t, root)
	catalogIDs := loadNonMUSTCatalogIDs(t)
	caseIDs := make(map[string]bool)
	files := make(map[string]bool)

	for _, testCase := range manifest.Cases {
		if !strings.HasPrefix(testCase.ID, "TOSCA13-NONMUST-CORPUS-") || caseIDs[testCase.ID] {
			t.Errorf("invalid or duplicate non-MUST case ID %q", testCase.ID)
		}
		caseIDs[testCase.ID] = true
		path, ok := canonicalCorpusPath(testCase.File)
		if !ok || !strings.HasPrefix(path, "non_must/") {
			t.Errorf("%s has unsafe non-MUST fixture path %q", testCase.ID, testCase.File)
			continue
		}
		if files[path] {
			t.Errorf("non-MUST fixture %q is mapped more than once", path)
		}
		files[path] = true
		if info, err := os.Stat(filepath.Join(root, filepath.FromSlash(path))); err != nil || info.IsDir() {
			t.Errorf("%s fixture %q is absent", testCase.ID, path)
		}
		if !allowedNonMUSTTiers[testCase.ConformanceTier] {
			t.Errorf("%s has invalid conformance tier %q", testCase.ID, testCase.ConformanceTier)
		}
		if len(testCase.Specification.Sections) == 0 || len(testCase.Specification.NonMUSTRequirementIDs) == 0 {
			t.Errorf("%s lacks specification traceability", testCase.ID)
		}
		for _, requirementID := range testCase.Specification.NonMUSTRequirementIDs {
			if !catalogIDs[requirementID] {
				t.Errorf("%s references absent non-MUST requirement %q", testCase.ID, requirementID)
			}
		}
		if strings.TrimSpace(testCase.Policy.Rationale) == "" {
			t.Errorf("%s has no policy rationale", testCase.ID)
		}
		if testCase.Classification.Kind == "invalid" && testCase.Expected.Accepted {
			t.Errorf("%s is invalid but expects acceptance", testCase.ID)
		}
		if !testCase.Expected.Accepted && testCase.Expected.Diagnostic == nil {
			t.Errorf("%s rejection has no diagnostic expectation", testCase.ID)
		}
		if testCase.Expected.Accepted && len(testCase.Expected.Assertions) == 0 {
			t.Errorf("%s acceptance has no semantic assertion", testCase.ID)
		}
	}

	for path := range discoverNonMUSTFixtures(t, root) {
		if !files[path] {
			t.Errorf("non-MUST fixture %q has no manifest entry", path)
		}
	}
	for path := range files {
		if !discoverNonMUSTFixtures(t, root)[path] {
			t.Errorf("non-MUST manifest path %q is not a discoverable fixture", path)
		}
	}
}

func TestTOSCA13NonMUSTCorpus(t *testing.T) {
	root := corpusRoot(t)
	manifest := loadNonMUSTManifest(t, root)
	filter := os.Getenv("TOSCA_NON_MUST_CASE")
	for _, testCase := range manifest.Cases {
		testCase := testCase
		if filter != "" && testCase.ID != filter {
			continue
		}
		t.Run(testCase.ID, func(t *testing.T) {
			path := filepath.Join(root, filepath.FromSlash(testCase.File))
			serviceTemplate, firstProblems, firstErr := parseCorpusFile(t, path)
			if testCase.Expected.Accepted {
				if firstErr != nil {
					t.Fatalf("%s rejected at %s: %v\n%s", testCase.ID, testCase.Expected.Phase, firstErr, firstProblems)
				}
				for _, assertion := range testCase.Expected.Assertions {
					assertCorpusResult(t, testCase.ID, serviceTemplate, assertion)
				}
				return
			}
			if firstErr == nil {
				t.Fatalf("%s was accepted, want rejection at %s", testCase.ID, testCase.Expected.Phase)
			}
			_, secondProblems, secondErr := parseCorpusFile(t, path)
			if secondErr == nil || firstProblems != secondProblems {
				t.Fatalf("%s has nondeterministic rejection diagnostics", testCase.ID)
			}
			for _, fragment := range []*string{
				testCase.Expected.Diagnostic.Category,
				testCase.Expected.Diagnostic.Path,
			} {
				if fragment == nil || !strings.Contains(firstProblems, *fragment) {
					t.Fatalf("%s diagnostic lacks %q:\n%s", testCase.ID, valueOrEmpty(fragment), firstProblems)
				}
			}
		})
	}
}

func loadNonMUSTManifest(t *testing.T, root string) nonMUSTManifest {
	t.Helper()
	var manifest nonMUSTManifest
	loadYAML(t, filepath.Join(root, "non_must/manifest.yaml"), &manifest)
	return manifest
}

func loadNonMUSTCatalogIDs(t *testing.T) map[string]bool {
	t.Helper()
	var catalog nonMUSTCatalog
	loadYAML(t, filepath.Join(repositoryRoot(t), "docs/conformance/tosca-1.3/non-must-requirements.yaml"), &catalog)
	ids := make(map[string]bool)
	for _, requirement := range catalog.Requirements {
		if requirement.ProcessorApplicable {
			ids[requirement.ID] = true
		}
	}
	return ids
}

func discoverNonMUSTFixtures(t *testing.T, root string) map[string]bool {
	t.Helper()
	fixtures := make(map[string]bool)
	base := filepath.Join(root, "non_must")
	err := filepath.WalkDir(base, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			if strings.HasPrefix(entry.Name(), "_") {
				return filepath.SkipDir
			}
			return nil
		}
		if entry.Name() == "manifest.yaml" || filepath.Ext(entry.Name()) != ".yaml" {
			return nil
		}
		relative, relativeErr := filepath.Rel(root, path)
		if relativeErr != nil {
			return relativeErr
		}
		fixtures[filepath.ToSlash(relative)] = true
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return fixtures
}
