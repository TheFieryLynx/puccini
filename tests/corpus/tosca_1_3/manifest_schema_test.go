package tosca_1_3_corpus_test

import (
	"fmt"
	"path/filepath"
	"testing"
)

// Cross-record predicates are enforced by validateEvidence. This gate checks
// the required v2 shape before zero-valued Go fields can hide missing metadata.
func validateManifestShape(m map[string]any) error {
	if m["schema_version"] != 2 {
		return fmt.Errorf("expected manifest schema_version 2")
	}
	cs, ok := m["cases"].([]any)
	if !ok {
		return fmt.Errorf("cases must be an array")
	}
	for _, raw := range cs {
		c, ok := raw.(map[string]any)
		if !ok {
			return fmt.Errorf("case must be a map")
		}
		coverage, ok := c["coverage"].(map[string]any)
		if !ok {
			return fmt.Errorf("case lacks coverage ownership")
		}
		for _, key := range []string{"primary_requirements", "supporting_requirements"} {
			a, ok := coverage[key].([]any)
			if !ok {
				return fmt.Errorf("%s must be an explicit array", key)
			}
			seen := map[string]bool{}
			for _, v := range a {
				s, ok := v.(string)
				if !ok || s == "" || seen[s] {
					return fmt.Errorf("invalid/duplicate requirement")
				}
				seen[s] = true
			}
		}
		if _, ok := c["assertions"].([]any); !ok {
			return fmt.Errorf("assertions must be an explicit array")
		}
	}
	return nil
}
func TestEvidenceManifestSchema(t *testing.T) {
	for _, p := range []string{"manifest.yaml", "non_must/manifest.yaml", "examples/manifest.yaml", "oasis_community/manifest.yaml"} {
		var m map[string]any
		loadYAML(t, filepath.Join(corpusRoot(t), p), &m)
		if err := validateManifestShape(m); err != nil {
			t.Fatalf("%s: %v", p, err)
		}
		delete(m, "schema_version")
		if validateManifestShape(m) == nil {
			t.Fatal("missing version accepted")
		}
	}
}
