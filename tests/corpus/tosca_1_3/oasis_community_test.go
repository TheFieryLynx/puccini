package tosca_1_3_corpus_test

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 5.3 through 5.11; 8.5
// Expected: every adapted OASIS community profile case preserves complete
// provenance and resolves all exact normative parent type names through the
// ordinary TOSCA 1.3 processor hierarchy path.
// Category: positive, independent corpus, provenance, hierarchy, inheritance
func TestTOSCA13OASISCommunityCorpus(t *testing.T) {
	root := corpusRoot(t)
	manifest := loadOASISCommunityManifest(t, root)
	sourceCommit, err := os.ReadFile(filepath.Join(
		repositoryRoot(t),
		"third_party/oasis/tosca-simple-profile/1.3/SOURCE_COMMIT",
	))
	if err != nil {
		t.Fatal(err)
	}
	wantCommit := strings.TrimSpace(string(sourceCommit))
	if manifest.UpstreamCommit != wantCommit {
		t.Fatalf("manifest upstream commit = %q, want %q", manifest.UpstreamCommit, wantCommit)
	}

	caseIDs := make(map[string]bool)
	files := make(map[string]bool)
	coveredTypes := make(map[string]bool)
	for _, testCase := range manifest.Cases {
		testCase := testCase
		if testCase.ID == "" || caseIDs[testCase.ID] {
			t.Errorf("empty or duplicate case ID %q", testCase.ID)
		}
		caseIDs[testCase.ID] = true
		if testCase.Upstream.Repository != manifest.UpstreamRepository ||
			testCase.Upstream.Commit != wantCommit ||
			testCase.Upstream.OriginalPath == "" ||
			testCase.Upstream.PinnedPath == "" ||
			testCase.Upstream.OriginalLicense != "Apache-2.0" ||
			testCase.Upstream.LicenseFile == "" {
			t.Errorf("%s has incomplete or inconsistent provenance: %#v", testCase.ID, testCase.Upstream)
		}
		if testCase.Adaptation.Status != "adapted" ||
			strings.TrimSpace(testCase.Adaptation.Explanation) == "" ||
			strings.TrimSpace(testCase.Adaptation.DiffSummary) == "" {
			t.Errorf("%s has incomplete adaptation metadata", testCase.ID)
		}
		if len(testCase.Specification.Sections) == 0 ||
			len(testCase.Specification.NonMUSTRequirementIDs) != len(testCase.TypesCovered) ||
			testCase.Specification.FrozenRequirementIDs == nil ||
			strings.TrimSpace(testCase.Specification.FrozenRequirementRationale) == "" {
			t.Errorf("%s has incomplete specification traceability", testCase.ID)
		}
		if testCase.Expected.ProcessorResult != "accept" ||
			testCase.Expected.Phase != "hierarchy" ||
			len(testCase.Expected.Assertions) == 0 {
			t.Errorf("%s has incomplete expected processor result", testCase.ID)
		}
		for _, path := range []string{
			testCase.File,
			testCase.Upstream.PinnedPath,
			testCase.Upstream.LicenseFile,
		} {
			if path == "" || filepath.IsAbs(path) || strings.Contains(path, "..") {
				t.Errorf("%s has unsafe path %q", testCase.ID, path)
				continue
			}
			if _, statErr := os.Stat(filepath.Join(repositoryRoot(t), filepath.FromSlash(path))); statErr != nil {
				if path == testCase.File {
					_, statErr = os.Stat(filepath.Join(
						root, "oasis_community", filepath.FromSlash(path)))
				}
				if statErr != nil {
					t.Errorf("%s path %q is absent: %v", testCase.ID, path, statErr)
				}
			}
		}
		if files[testCase.File] {
			t.Errorf("fixture %q is used more than once", testCase.File)
		}
		files[testCase.File] = true
		for _, name := range testCase.TypesCovered {
			if coveredTypes[name] {
				t.Errorf("type %q is covered by more than one OASIS case", name)
			}
			coveredTypes[name] = true
		}

		t.Run(testCase.ID, func(t *testing.T) {
			path := filepath.Join(root, "oasis_community", filepath.FromSlash(testCase.File))
			_, problems, parseErr := parseCorpusFile(t, path)
			if parseErr != nil {
				t.Fatalf("adapted OASIS community case rejected: %v\n%s", parseErr, problems)
			}
		})
	}
	if len(manifest.Cases) != 8 {
		t.Fatalf("OASIS community case count = %d, want 8", len(manifest.Cases))
	}
	if len(coveredTypes) != 66 {
		names := make([]string, 0, len(coveredTypes))
		for name := range coveredTypes {
			names = append(names, name)
		}
		sort.Strings(names)
		t.Fatalf("OASIS community covered type count = %d, want 66: %v", len(names), names)
	}
}

type oasisCommunityManifest struct {
	SchemaVersion      int                  `yaml:"schema_version"`
	Layer              string               `yaml:"layer"`
	UpstreamRepository string               `yaml:"upstream_repository"`
	UpstreamCommit     string               `yaml:"upstream_commit"`
	PinnedProfileRoot  string               `yaml:"pinned_profile_root"`
	Cases              []oasisCommunityCase `yaml:"cases"`
}

type oasisCommunityCase struct {
	ID       string `yaml:"id"`
	File     string `yaml:"file"`
	Upstream struct {
		Repository      string `yaml:"repository"`
		Commit          string `yaml:"commit"`
		OriginalPath    string `yaml:"original_path"`
		PinnedPath      string `yaml:"pinned_path"`
		OriginalLicense string `yaml:"original_license"`
		LicenseFile     string `yaml:"license_file"`
	} `yaml:"upstream"`
	Adaptation struct {
		Status      string `yaml:"status"`
		Explanation string `yaml:"explanation"`
		DiffSummary string `yaml:"diff_summary"`
	} `yaml:"adaptation"`
	Specification struct {
		Version                    string   `yaml:"version"`
		Sections                   []string `yaml:"sections"`
		FrozenRequirementIDs       []string `yaml:"frozen_requirement_ids"`
		FrozenRequirementRationale string   `yaml:"frozen_requirement_rationale"`
		NonMUSTRequirementIDs      []string `yaml:"non_must_requirement_ids"`
	} `yaml:"specification"`
	TypesCovered []string `yaml:"types_covered"`
	Expected     struct {
		ProcessorResult string   `yaml:"processor_result"`
		Phase           string   `yaml:"phase"`
		Assertions      []string `yaml:"assertions"`
	} `yaml:"expected"`
}

func loadOASISCommunityManifest(t *testing.T, root string) oasisCommunityManifest {
	t.Helper()
	var manifest oasisCommunityManifest
	loadYAML(t, filepath.Join(root, "oasis_community/manifest.yaml"), &manifest)
	if manifest.SchemaVersion != 1 ||
		manifest.Layer == "" ||
		manifest.UpstreamRepository == "" ||
		manifest.PinnedProfileRoot == "" {
		t.Fatalf("OASIS community manifest header is incomplete")
	}
	return manifest
}
