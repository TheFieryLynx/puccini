package tosca_1_3_corpus_test

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

type ownedAssertion struct {
	ID             string   `yaml:"id"`
	RequirementIDs []string `yaml:"requirement_ids"`
	Kind           string   `yaml:"kind"`
	Path           string   `yaml:"path"`
	Expected       any      `yaml:"expected"`
}
type relatedEvidence struct {
	NearestValid string `yaml:"nearest_valid"`
	Change       string `yaml:"change"`
	Exception    string `yaml:"exception"`
}
type predicateRegistry struct {
	Requirements []reviewedPredicate `yaml:"requirements"`
}
type reviewedPredicate struct {
	ID               string             `yaml:"id"`
	Predicate        string             `yaml:"predicate"`
	SourceSHA256     string             `yaml:"source_predicate_sha256"`
	EvidenceKind     string             `yaml:"evidence_kind"`
	PositiveRequired bool               `yaml:"positive_required"`
	NegativeRequired bool               `yaml:"negative_required"`
	Bindings         []predicateBinding `yaml:"bindings"`
	MutationPaths    []string           `yaml:"mutation_paths"`
}

// Independently reviewed bindings prevent changing both a case label and its
// assertion label from automatically creating evidence for another predicate.
// Hashes force a fresh review when fixture bytes change. They are provenance,
// not proof that a human interpretation of the specification is infallible.
type predicateBinding struct {
	CaseID        string         `yaml:"case_id"`
	FixtureSHA256 string         `yaml:"fixture_sha256"`
	Assertion     ownedAssertion `yaml:"assertion"`
}

var evidenceKinds = map[string]bool{"structural-acceptance": true, "semantic-value": true, "semantic-resolution": true, "diagnostic": true, "normalization": true}
var assertionKinds = map[string]bool{"accept": true, "reject": true, "value": true, "diagnostic": true, "phase": true, "normalization": true}

func registryIndex(reg predicateRegistry) map[string]reviewedPredicate {
	index := map[string]reviewedPredicate{}
	for _, r := range reg.Requirements {
		index[r.ID] = r
	}
	return index
}
func validateEvidence(cases []corpusCase, reg predicateRegistry, read func(string) ([]byte, error)) error {
	ids := map[string]corpusCase{}
	predicates := registryIndex(reg)
	for _, c := range cases {
		if _, ok := ids[c.ID]; ok {
			return fmt.Errorf("duplicate case %s", c.ID)
		}
		ids[c.ID] = c
	}
	for _, c := range cases {
		if !c.Expected.Accepted {
			if c.Related.NearestValid == "" {
				if c.Related.Exception == "" {
					return fmt.Errorf("%s missing nearest-valid sibling", c.ID)
				}
				// Exceptions require a reviewed predicate binding, not an arbitrary escape hatch.
				if len(c.Coverage.PrimaryRequirements) == 0 {
					return fmt.Errorf("%s unreviewed sibling exception", c.ID)
				}
			} else {
				sibling, ok := ids[c.Related.NearestValid]
				if !ok || !sibling.Expected.Accepted {
					return fmt.Errorf("%s missing or invalid nearest-valid sibling", c.ID)
				}
				if strings.TrimSpace(c.Related.Change) == "" {
					return fmt.Errorf("%s unexplained sibling change", c.ID)
				}
			}
		}
		seen := map[string]bool{}
		owned := map[string]bool{}
		semantic := map[string]bool{}
		for _, a := range c.Assertions {
			if a.ID == "" || seen[a.ID] || !assertionKinds[a.Kind] || a.Path == "" {
				return fmt.Errorf("%s malformed/duplicate assertion %s", c.ID, a.ID)
			}
			seen[a.ID] = true
			if len(a.RequirementIDs) != 1 {
				return fmt.Errorf("%s/%s must own exactly one atomic predicate", c.ID, a.ID)
			}
			id := a.RequirementIDs[0]
			if !contains(c.Coverage.PrimaryRequirements, id) {
				return fmt.Errorf("%s/%s assertion is not owned by a primary requirement", c.ID, a.ID)
			}
			r, ok := predicates[id]
			if !ok || !evidenceKinds[r.EvidenceKind] {
				return fmt.Errorf("%s unreviewed requirement %s", c.ID, id)
			}
			bytes, err := read(c.File)
			if err != nil {
				return err
			}
			sum := sha256.Sum256(bytes)
			digest := hex.EncodeToString(sum[:])
			found := false
			for _, b := range r.Bindings {
				if b.CaseID == c.ID && b.FixtureSHA256 == digest && reflect.DeepEqual(b.Assertion, a) {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("%s/%s unrelated or stale predicate binding %s", c.ID, a.ID, id)
			}
			owned[id] = true
			if a.Kind == "value" || a.Kind == "normalization" {
				semantic[id] = true
			}
		}
		for _, id := range c.Coverage.PrimaryRequirements {
			if contains(c.Coverage.SupportingRequirements, id) {
				return fmt.Errorf("%s primary/supporting overlap %s", c.ID, id)
			}
			if !owned[id] {
				return fmt.Errorf("%s primary requirement %s has no owned assertion", c.ID, id)
			}
			r := predicates[id]
			if c.Expected.Accepted && (strings.HasPrefix(r.EvidenceKind, "semantic-") || r.EvidenceKind == "normalization") && !semantic[id] {
				return fmt.Errorf("%s semantic predicate %s has only acceptance evidence", c.ID, id)
			}
		}
	}
	return nil
}
func assertionResult(a ownedAssertion, r corpusResult) error {
	var got any
	switch a.Kind {
	case "accept", "reject":
		got = r.Err == nil
	case "phase":
		got = string(r.Phase)
	case "diagnostic":
		expected, ok := a.Expected.(string)
		if !ok || !strings.Contains(r.Problems, expected) {
			return fmt.Errorf("diagnostic lacks %q", a.Expected)
		}
		return nil
	case "value", "normalization":
		if r.Err != nil || r.Template == nil {
			return fmt.Errorf("no normalized result")
		}
		raw, err := json.Marshal(r.Template)
		if err != nil {
			return err
		}
		var value any
		if err = json.Unmarshal(raw, &value); err != nil {
			return err
		}
		value, err = evidenceMapView(value)
		if err != nil {
			return err
		}
		got, err = jsonPointer(value, a.Path)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported assertion kind %s", a.Kind)
	}
	actual, err := json.Marshal(got)
	if err != nil {
		return err
	}
	expected, err := json.Marshal(a.Expected)
	if err != nil {
		return err
	}
	if string(actual) != string(expected) {
		return fmt.Errorf("%s: got %s, expected %s", a.Path, actual, expected)
	}
	return nil
}
func jsonPointer(value any, path string) (any, error) {
	if !strings.HasPrefix(path, "/") {
		return nil, fmt.Errorf("not a JSON pointer: %s", path)
	}
	for _, part := range strings.Split(path[1:], "/") {
		part = strings.ReplaceAll(strings.ReplaceAll(part, "~1", "/"), "~0", "~")
		switch v := value.(type) {
		case map[string]any:
			var ok bool
			value, ok = v[part]
			if !ok {
				return nil, fmt.Errorf("missing assertion path %s", path)
			}
		case []any:
			i, err := strconv.Atoi(part)
			if err != nil || i < 0 || i >= len(v) {
				return nil, fmt.Errorf("bad array path %s", path)
			}
			value = v[i]
		default:
			return nil, fmt.Errorf("cannot traverse assertion path %s", path)
		}
	}
	return value, nil
}
func TestTOSCA13EvidenceOwnership(t *testing.T) {
	root := corpusRoot(t)
	cases := loadManifest(t, root).Cases
	// Non-MUST corpus labels are also supporting unless explicitly reviewed.
	for _, c := range loadNonMUSTManifest(t, root).Cases {
		cases = append(cases, corpusCase{ID: c.ID, File: c.File, Expected: corpusExpected{Accepted: c.Expected.Accepted}, Coverage: c.Coverage, Assertions: c.Assertions, Related: c.Related})
	}
	var reg predicateRegistry
	loadYAML(t, filepath.Join(root, "evidence-predicates.yaml"), &reg)
	if err := validateEvidence(cases, reg, func(p string) ([]byte, error) { return os.ReadFile(filepath.Join(root, p)) }); err != nil {
		t.Fatal(err)
	}
	// Review catalog predicates independently of coverage.yaml status labels.
	var catalog struct {
		Requirements []struct {
			ID          string `yaml:"id"`
			Requirement string `yaml:"requirement"`
		} `yaml:"requirements"`
	}
	loadYAML(t, filepath.Join(repositoryRoot(t), "docs/conformance/tosca-1.3/requirements.yaml"), &catalog)
	source := map[string]string{}
	for _, r := range catalog.Requirements {
		source[r.ID] = r.Requirement
	}
	for _, r := range reg.Requirements {
		sum := sha256.Sum256([]byte(source[r.ID]))
		if source[r.ID] == "" || hex.EncodeToString(sum[:]) != r.SourceSHA256 {
			t.Errorf("predicate source changed or absent: %s", r.ID)
		}
	}
}
