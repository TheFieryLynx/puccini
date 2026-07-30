package tosca_1_3_corpus_test

import (
	"archive/zip"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"testing"

	"github.com/tliron/exturl"
	"github.com/tliron/go-puccini/normal"
	"github.com/tliron/go-puccini/tosca/parser"
	"gopkg.in/yaml.v3"
)

type corpusManifest struct {
	Cases []corpusCase `yaml:"cases"`
}

type corpusCase struct {
	ID             string               `yaml:"id"`
	File           string               `yaml:"file"`
	Specification  corpusSpecification  `yaml:"specification"`
	Classification corpusClassification `yaml:"classification"`
	Expected       corpusExpected       `yaml:"expected"`
	Coverage       corpusCoverage       `yaml:"coverage"`
}

type corpusSpecification struct {
	Sections       []string `yaml:"sections"`
	RequirementIDs []string `yaml:"requirement_ids"`
}

type corpusClassification struct {
	Kind     string   `yaml:"kind"`
	Category string   `yaml:"category"`
	Features []string `yaml:"features"`
}

type corpusExpected struct {
	Accepted   bool             `yaml:"accepted"`
	Phase      string           `yaml:"phase"`
	Diagnostic corpusDiagnostic `yaml:"diagnostic"`
	Assertions []string         `yaml:"assertions"`
}

type corpusDiagnostic struct {
	Category *string `yaml:"category"`
	Path     *string `yaml:"path"`
}

type corpusCoverage struct {
	VariationAxes []string `yaml:"variation_axes"`
	RelatedCases  []string `yaml:"related_cases"`
}

type requirementCatalog struct {
	Requirements []struct {
		ID      string `yaml:"id"`
		Section string `yaml:"section"`
	} `yaml:"requirements"`
}

type coverageCatalog struct {
	Requirements []struct {
		RequirementID           string `yaml:"requirement_id"`
		ImplementationStatus    string `yaml:"implementation_status"`
		VerificationStatus      string `yaml:"verification_status"`
		IncludedInMUSTDenomator bool   `yaml:"included_in_must_denominator"`
	} `yaml:"requirements"`
}

var allowedPhases = map[string]bool{
	"read":          true,
	"namespaces":    true,
	"hierarchy":     true,
	"inheritance":   true,
	"rendering":     true,
	"normalization": true,
	"csar":          true,
}

func TestTOSCA13CorpusManifestIntegrity(t *testing.T) {
	root := corpusRoot(t)
	manifest := loadManifest(t, root)
	requirementIDs, sections := loadRequirementIndex(t)

	caseIDs := make(map[string]bool)
	files := make(map[string]bool)
	for _, testCase := range manifest.Cases {
		if testCase.ID == "" || !strings.HasPrefix(testCase.ID, "TOSCA13-CORPUS-") {
			t.Errorf("invalid case ID %q", testCase.ID)
		}
		if caseIDs[testCase.ID] {
			t.Errorf("duplicate case ID %q", testCase.ID)
		}
		caseIDs[testCase.ID] = true

		canonical, ok := canonicalCorpusPath(testCase.File)
		if !ok {
			t.Errorf("%s has non-canonical or unsafe file path %q", testCase.ID, testCase.File)
			continue
		}
		if files[canonical] {
			t.Errorf("fixture %q is mapped by more than one case", canonical)
		}
		files[canonical] = true
		if info, err := os.Stat(filepath.Join(root, filepath.FromSlash(canonical))); err != nil || info.IsDir() {
			t.Errorf("%s fixture %q is absent or not a file", testCase.ID, canonical)
		}

		if len(testCase.Specification.Sections) == 0 {
			t.Errorf("%s has no specification section", testCase.ID)
		}
		for _, section := range testCase.Specification.Sections {
			if !sections[section] {
				t.Errorf("%s references unknown section %q", testCase.ID, section)
			}
		}
		if len(testCase.Specification.RequirementIDs) == 0 {
			t.Errorf("%s has no requirement ID", testCase.ID)
		}
		for _, requirementID := range testCase.Specification.RequirementIDs {
			if !requirementIDs[requirementID] {
				t.Errorf("%s references unknown requirement %q", testCase.ID, requirementID)
			}
		}
		if len(testCase.Coverage.VariationAxes) == 0 {
			t.Errorf("%s has no variation axis", testCase.ID)
		}
		if !allowedPhases[testCase.Expected.Phase] {
			t.Errorf("%s has unsupported phase %q", testCase.ID, testCase.Expected.Phase)
		}
		if testCase.Classification.Kind != "valid" && testCase.Classification.Kind != "invalid" {
			t.Errorf("%s has unsupported kind %q", testCase.ID, testCase.Classification.Kind)
		}
		if testCase.Classification.Kind == "invalid" {
			if testCase.Expected.Accepted {
				t.Errorf("%s is invalid but expects acceptance", testCase.ID)
			}
			if testCase.Expected.Diagnostic.Category == nil || testCase.Expected.Diagnostic.Path == nil {
				t.Errorf("%s invalid case lacks diagnostic category/path", testCase.ID)
			}
		} else if !testCase.Expected.Accepted {
			t.Errorf("%s is valid but expects rejection", testCase.ID)
		} else if len(testCase.Expected.Assertions) == 0 {
			t.Errorf("%s valid case has no semantic assertion", testCase.ID)
		}
	}

	for _, testCase := range manifest.Cases {
		for _, related := range testCase.Coverage.RelatedCases {
			if !caseIDs[related] {
				t.Errorf("%s references absent related case %q", testCase.ID, related)
			}
		}
	}

	fixtures := discoverFixtures(t, root)
	for file := range fixtures {
		if !files[file] {
			t.Errorf("fixture %q has no manifest entry", file)
		}
	}
	for file := range files {
		if !fixtures[file] {
			t.Errorf("manifest entry fixture %q is outside executable fixture roots", file)
		}
	}
}

func TestTOSCA13CorpusFrozenMUSTGate(t *testing.T) {
	var coverage coverageCatalog
	loadYAML(t, filepath.Join(repositoryRoot(t), "docs/conformance/tosca-1.3/coverage.yaml"), &coverage)
	count := 0
	for _, requirement := range coverage.Requirements {
		if !requirement.IncludedInMUSTDenomator {
			continue
		}
		count++
		if requirement.ImplementationStatus != "implemented" ||
			requirement.VerificationStatus != "verified" {
			t.Errorf("%s frozen status = %s/%s",
				requirement.RequirementID,
				requirement.ImplementationStatus,
				requirement.VerificationStatus)
		}
	}
	if count != 223 {
		t.Fatalf("frozen MUST denominator = %d, want 223", count)
	}
}

func TestTOSCA13Corpus(t *testing.T) {
	root := corpusRoot(t)
	manifest := loadManifest(t, root)
	caseFilter := os.Getenv("TOSCA_CORPUS_CASE")
	sectionFilter := os.Getenv("TOSCA_CORPUS_SECTION")
	featureFilter := os.Getenv("TOSCA_CORPUS_FEATURE")

	selected := 0
	for _, testCase := range manifest.Cases {
		testCase := testCase
		if caseFilter != "" && testCase.ID != caseFilter {
			continue
		}
		if sectionFilter != "" && !contains(testCase.Specification.Sections, sectionFilter) {
			continue
		}
		if featureFilter != "" && !contains(testCase.Classification.Features, featureFilter) {
			continue
		}
		selected++
		t.Run(testCase.ID, func(t *testing.T) {
			runCorpusCase(t, root, testCase)
		})
	}
	if selected == 0 {
		t.Fatalf("corpus filters selected no cases: case=%q section=%q feature=%q",
			caseFilter, sectionFilter, featureFilter)
	}
}

func runCorpusCase(t *testing.T, root string, testCase corpusCase) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(testCase.File))
	if testCase.Classification.Category == "csar" {
		path = materializeCSAR(t, path)
	}
	serviceTemplate, firstProblems, firstErr := parseCorpusFile(t, path)
	if testCase.Expected.Accepted {
		if firstErr != nil {
			t.Fatalf("%s expected acceptance at %s: %v\n%s",
				testCase.ID, testCase.Expected.Phase, firstErr, firstProblems)
		}
		for _, assertion := range testCase.Expected.Assertions {
			assertCorpusResult(t, testCase.ID, serviceTemplate, assertion)
		}
		return
	}
	if firstErr == nil {
		t.Fatalf("%s expected rejection at %s", testCase.ID, testCase.Expected.Phase)
	}
	for label, fragment := range map[string]*string{
		"category": testCase.Expected.Diagnostic.Category,
		"path":     testCase.Expected.Diagnostic.Path,
	} {
		if fragment == nil || !strings.Contains(firstProblems, *fragment) {
			t.Fatalf("%s diagnostic lacks %s %q:\n%s",
				testCase.ID, label, valueOrEmpty(fragment), firstProblems)
		}
	}
	_, secondProblems, secondErr := parseCorpusFile(t, path)
	if secondErr == nil || firstProblems != secondProblems {
		t.Fatalf("%s diagnostic is not deterministic:\nfirst:\n%s\nsecond:\n%s",
			testCase.ID, firstProblems, secondProblems)
	}
}

func materializeCSAR(t *testing.T, recipePath string) string {
	t.Helper()
	var recipe struct {
		Entries map[string]string `yaml:"entries"`
	}
	loadYAML(t, recipePath, &recipe)
	if len(recipe.Entries) == 0 {
		t.Fatalf("CSAR recipe %s has no entries", recipePath)
	}
	path := filepath.Join(t.TempDir(), "fixture.csar")
	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("create CSAR: %v", err)
	}
	writer := zip.NewWriter(file)
	names := make([]string, 0, len(recipe.Entries))
	for name := range recipe.Entries {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		entry, createErr := writer.Create(name)
		if createErr != nil {
			t.Fatalf("create CSAR entry %q: %v", name, createErr)
		}
		if _, writeErr := entry.Write([]byte(recipe.Entries[name])); writeErr != nil {
			t.Fatalf("write CSAR entry %q: %v", name, writeErr)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close CSAR writer: %v", err)
	}
	if err := file.Close(); err != nil {
		t.Fatalf("close CSAR file: %v", err)
	}
	return path
}

func assertCorpusResult(t *testing.T, caseID string, serviceTemplate *normal.ServiceTemplate, assertion string) {
	t.Helper()
	switch {
	case assertion == "service_template_non_nil":
		if serviceTemplate == nil {
			t.Fatalf("%s normalized service template is nil", caseID)
		}
	case strings.HasPrefix(assertion, "node_exists:"):
		name := strings.TrimPrefix(assertion, "node_exists:")
		if serviceTemplate == nil || serviceTemplate.NodeTemplates[name] == nil {
			t.Fatalf("%s normalized node %q is absent", caseID, name)
		}
	case strings.HasPrefix(assertion, "input_exists:"):
		name := strings.TrimPrefix(assertion, "input_exists:")
		if serviceTemplate == nil || serviceTemplate.Inputs[name] == nil {
			t.Fatalf("%s normalized input %q is absent", caseID, name)
		}
	case strings.HasPrefix(assertion, "output_exists:"):
		name := strings.TrimPrefix(assertion, "output_exists:")
		if serviceTemplate == nil || serviceTemplate.Outputs[name] == nil {
			t.Fatalf("%s normalized output %q is absent", caseID, name)
		}
	case strings.HasPrefix(assertion, "requirement_target:"):
		parts := strings.Split(assertion, ":")
		if len(parts) != 4 {
			t.Fatalf("%s has malformed requirement target assertion %q", caseID, assertion)
		}
		source := serviceTemplate.NodeTemplates[parts[1]]
		if source == nil {
			t.Fatalf("%s normalized source node %q is absent", caseID, parts[1])
		}
		for _, requirement := range source.Requirements {
			if requirement.Name == parts[2] {
				if requirement.NodeTemplate == nil || requirement.NodeTemplate.Name != parts[3] {
					t.Fatalf("%s requirement %q target = %#v, want %q",
						caseID, parts[2], requirement.NodeTemplate, parts[3])
				}
				return
			}
		}
		t.Fatalf("%s normalized requirement %q is absent from node %q", caseID, parts[2], parts[1])
	default:
		t.Fatalf("%s has unsupported assertion %q", caseID, assertion)
	}
}

func parseCorpusFile(t *testing.T, path string) (*normal.ServiceTemplate, string, error) {
	t.Helper()
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		t.Fatalf("resolve corpus fixture: %v", err)
	}
	urlContext := exturl.NewContext()
	defer urlContext.Release()
	parserContext := parser.NewParser().NewContext()
	parserContext.URL = urlContext.NewFileURL(filepath.ToSlash(absolutePath))
	serviceTemplate, parseErr := parserContext.Parse(context.Background())
	problems := ""
	if parserContext.Root != nil {
		problems = parserContext.GetProblems().ToString(false)
	}
	return serviceTemplate, problems, parseErr
}

func loadManifest(t *testing.T, root string) corpusManifest {
	t.Helper()
	var manifest corpusManifest
	loadYAML(t, filepath.Join(root, "manifest.yaml"), &manifest)
	return manifest
}

func loadRequirementIndex(t *testing.T) (map[string]bool, map[string]bool) {
	t.Helper()
	var catalog requirementCatalog
	loadYAML(t, filepath.Join(repositoryRoot(t), "docs/conformance/tosca-1.3/requirements.yaml"), &catalog)
	requirements := make(map[string]bool)
	sections := make(map[string]bool)
	for _, requirement := range catalog.Requirements {
		requirements[requirement.ID] = true
		sections[requirement.Section] = true
	}
	return requirements, sections
}

func loadYAML(t *testing.T, path string, target any) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if err := yaml.Unmarshal(data, target); err != nil {
		t.Fatalf("decode %s: %v", path, err)
	}
}

func discoverFixtures(t *testing.T, root string) map[string]bool {
	t.Helper()
	fixtures := make(map[string]bool)
	for _, directory := range []string{"sections", "interactions", "version-isolation"} {
		base := filepath.Join(root, directory)
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
			if !strings.HasSuffix(entry.Name(), ".yaml") ||
				strings.HasSuffix(entry.Name(), ".normalized.yaml") {
				return nil
			}
			relative, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			fixtures[filepath.ToSlash(relative)] = true
			return nil
		})
		if err != nil && !os.IsNotExist(err) {
			t.Fatalf("walk corpus fixtures: %v", err)
		}
	}
	return fixtures
}

func canonicalCorpusPath(path string) (string, bool) {
	if path == "" || filepath.IsAbs(path) || strings.Contains(path, `\`) {
		return "", false
	}
	clean := filepath.ToSlash(filepath.Clean(filepath.FromSlash(path)))
	if clean != path || clean == "." || clean == ".." || strings.HasPrefix(clean, "../") {
		return "", false
	}
	return clean, true
}

func corpusRoot(t testing.TB) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("locate corpus test")
	}
	return filepath.Dir(file)
}

func repositoryRoot(t *testing.T) string {
	t.Helper()
	return filepath.Clean(filepath.Join(corpusRoot(t), "../../.."))
}

func contains(values []string, expected string) bool {
	for _, value := range values {
		if value == expected {
			return true
		}
	}
	return false
}

func valueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func sortedKeys(values map[string]bool) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func formatCorpusFilter(caseID string, section string, feature string) string {
	return fmt.Sprintf("case=%q section=%q feature=%q", caseID, section, feature)
}
