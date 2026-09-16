package tosca_1_3_corpus_test

import (
	"path/filepath"
	"strings"
	"testing"
)

type sectionInventory struct {
	Sections []struct {
		Section               string   `yaml:"section"`
		ApplicableToProcessor bool     `yaml:"applicable_to_processor"`
		Fixtures              []string `yaml:"fixtures"`
		Status                string   `yaml:"status"`
		Exclusions            []string `yaml:"exclusions"`
	} `yaml:"sections"`
}

type pairwiseInventory struct {
	Pairs []struct {
		Axes                  []string `yaml:"axes"`
		Fixtures              []string `yaml:"fixtures"`
		UncoveredCombinations []struct {
			Combination []string `yaml:"combination"`
			Reason      string   `yaml:"reason"`
		} `yaml:"uncovered_combinations"`
	} `yaml:"pairs"`
}

func TestTOSCA13HistoricalLabelInventory(t *testing.T) {
	var coverage coverageCatalog
	loadYAML(t, filepath.Join(repositoryRoot(t), "docs/conformance/tosca-1.3/coverage.yaml"), &coverage)
	frozenIDs := make(map[string]bool)
	for _, requirement := range coverage.Requirements {
		if requirement.IncludedInMUSTDenomator {
			frozenIDs[requirement.RequirementID] = true
		}
	}
	manifest := loadManifest(t, corpusRoot(t))
	for _, testCase := range manifest.Cases {
		for _, requirementID := range testCase.Specification.RequirementIDs {
			delete(frozenIDs, requirementID)
		}
	}
	if len(frozenIDs) != 0 {
		t.Fatalf("frozen MUST requirements absent from corpus manifest: %v", sortedKeys(frozenIDs))
	}
}

func TestTOSCA13CorpusCoverageGate(t *testing.T) {
	var sections sectionInventory
	loadYAML(t, filepath.Join(repositoryRoot(t), "docs/conformance/tosca-1.3/template-corpus-sections.yaml"), &sections)
	if len(sections.Sections) == 0 {
		t.Fatal("section inventory is empty")
	}
	for _, section := range sections.Sections {
		if section.Status != "insufficient-evidence" {
			t.Errorf("section %s status = %s", section.Section, section.Status)
		}

	}

	var pairwise pairwiseInventory
	loadYAML(t, filepath.Join(repositoryRoot(t), "docs/conformance/tosca-1.3/template-corpus-pairwise.yaml"), &pairwise)
	for _, pair := range pairwise.Pairs {
		if len(pair.Axes) != 2 || len(pair.Fixtures) == 0 {
			t.Errorf("invalid pairwise record: %#v", pair)
		}
		for _, gap := range pair.UncoveredCombinations {
			if len(gap.Combination) != 2 || strings.TrimSpace(gap.Reason) == "" {
				t.Errorf("unexplained pairwise gap: %#v", gap)
			}
		}
	}
}
