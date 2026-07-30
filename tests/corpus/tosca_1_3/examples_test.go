package tosca_1_3_corpus_test

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

type exampleManifest struct {
	Cases []struct {
		ID                string   `yaml:"id"`
		File              string   `yaml:"file"`
		Section           string   `yaml:"section"`
		Accepted          bool     `yaml:"accepted"`
		NormativeEvidence bool     `yaml:"normative_evidence"`
		Demonstrates      []string `yaml:"indirectly_demonstrates"`
	} `yaml:"cases"`
}

func TestTOSCA13ExampleCompatibilityCorpus(t *testing.T) {
	root := corpusRoot(t)
	data, err := os.ReadFile(filepath.Join(root, "examples", "manifest.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	var manifest exampleManifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		t.Fatal(err)
	}
	ids := make(map[string]bool)
	for _, testCase := range manifest.Cases {
		if testCase.ID == "" || ids[testCase.ID] {
			t.Fatalf("empty or duplicate example case ID %q", testCase.ID)
		}
		ids[testCase.ID] = true
		if testCase.Section == "" || testCase.NormativeEvidence || len(testCase.Demonstrates) == 0 {
			t.Fatalf("%s lacks non-normative example metadata", testCase.ID)
		}
		_, problems, parseErr := parseCorpusFile(t, filepath.Join(root, filepath.FromSlash(testCase.File)))
		if testCase.Accepted && parseErr != nil {
			t.Fatalf("%s rejected: %v\n%s", testCase.ID, parseErr, problems)
		}
		if !testCase.Accepted && parseErr == nil {
			t.Fatalf("%s was accepted", testCase.ID)
		}
	}
}
