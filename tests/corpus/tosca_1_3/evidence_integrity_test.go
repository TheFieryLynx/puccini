package tosca_1_3_corpus_test

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TOSCA 1.3 §14.3 Processor conformance: a rejection is evidence only for
// the predicate actually reached. Phase identities are API boundaries.
func TestEvidenceExactPhase(t *testing.T) {
	phases := []processingPhase{phaseRead, phaseNamespaces, phaseHierarchy, phaseInheritance, phaseRendering, phaseNormalization}
	for _, actual := range phases {
		for _, expected := range phases {
			r := corpusResult{Phase: actual, Err: fmt.Errorf("same opaque diagnostic")}
			err := checkExpectedPhase(string(expected), r)
			if (err == nil) != (actual == expected) {
				t.Fatalf("actual=%s expected=%s: %v", actual, expected, err)
			}
		}
	}
	if checkExpectedPhase("read", corpusResult{Phase: phaseRead}) == nil {
		t.Fatal("acceptance credited as rejection")
	}
}

func TestEvidenceEarlyFailureCannotProveLaterRule(t *testing.T) {
	root := corpusRoot(t)
	tests := []struct {
		name, file string
		phase      processingPhase
	}{
		{"structural", "version-isolation/tosca-2.0-service-template-key-invalid.yaml", phaseRead},
		{"namespace", "sections/3.1-namespaces/imports/unknown-qualified-invalid.yaml", phaseNamespaces},
		{"hierarchy", "sections/3.7-type-definitions/hierarchy/cycle-invalid.yaml", phaseHierarchy},
		{"refinement", "sections/3.7-type-definitions/refinement/required-widening-invalid.yaml", phaseInheritance},
		{"assignment", "sections/3.6-grammar-definitions/constraints/operand-invalid.yaml", phaseRendering},
		{"resolved-function", "sections/4-functions/functions/reference-invalid.yaml", phaseRendering},
		{"csar-metadata", "sections/6-csar/archives/entry-missing-invalid.csar.yaml", phaseRead},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(root, tc.file)
			if strings.HasSuffix(path, ".csar.yaml") {
				path = materializeCSAR(t, path)
			}
			r := parseCorpusPhases(t, path)
			if err := checkExpectedPhase(string(tc.phase), r); err != nil {
				t.Fatalf("%v\n%s", err, r.Problems)
			}
			// Read is earliest, including archive opening: there is no earlier API phase.
			if tc.phase == phaseRead {
				if checkExpectedPhase("rendering", r) == nil {
					t.Fatal("read evidence credited as rendering")
				}
				return
			}
			data, err := os.ReadFile(filepath.Join(root, tc.file))
			if err != nil {
				t.Fatal(err)
			}
			evil := filepath.Join(t.TempDir(), "evil.yaml")
			// Unknown root key fails before imports or any semantic predicate.
			if err := os.WriteFile(evil, append(data, []byte("\nevidence_unknown_key: true\n")...), 0600); err != nil {
				t.Fatal(err)
			}
			r = parseCorpusPhases(t, evil)
			if r.Phase != phaseRead || checkExpectedPhase(string(tc.phase), r) == nil {
				t.Fatalf("early sibling credited: %+v", r)
			}
		})
	}
	// No normative normalization rejection category is invented: normalization
	// is checked by values and the phase-comparison unit matrix above.
}

func TestEvidenceOwnershipIntegrity(t *testing.T) {
	makeEvidence := func() ([]corpusCase, predicateRegistry) {
		a := ownedAssertion{ID: "value", RequirementIDs: []string{"R"}, Kind: "value", Path: "/description", Expected: "wanted"}
		c := corpusCase{ID: "valid", File: "valid.yaml", Expected: corpusExpected{Accepted: true}, Coverage: corpusCoverage{PrimaryRequirements: []string{"R"}}, Assertions: []ownedAssertion{a}}
		sum := sha256.Sum256([]byte("fixture"))
		r := reviewedPredicate{ID: "R", EvidenceKind: "semantic-value", Bindings: []predicateBinding{{CaseID: c.ID, FixtureSHA256: fmt.Sprintf("%x", sum), Assertion: a}}}
		return []corpusCase{c}, predicateRegistry{Requirements: []reviewedPredicate{r}}
	}
	read := func(string) ([]byte, error) { return []byte("fixture"), nil }
	cs, reg := makeEvidence()
	if err := validateEvidence(cs, reg, read); err != nil {
		t.Fatal(err)
	}
	tests := map[string]func([]corpusCase, *predicateRegistry){
		"primary-without-assertion": func(c []corpusCase, r *predicateRegistry) { c[0].Assertions = nil },
		"unowned-assertion":         func(c []corpusCase, r *predicateRegistry) { c[0].Assertions[0].RequirementIDs = nil },
		"supporting-only": func(c []corpusCase, r *predicateRegistry) {
			c[0].Coverage.PrimaryRequirements = nil
			c[0].Coverage.SupportingRequirements = []string{"R"}
		},
		"unrelated-predicate": func(c []corpusCase, r *predicateRegistry) {
			c[0].Coverage.PrimaryRequirements = []string{"S"}
			c[0].Assertions[0].RequirementIDs = []string{"S"}
		},
		"unrelated-fixture": func(c []corpusCase, r *predicateRegistry) { c[0].ID = "another" },
		"semantic-accept-only": func(c []corpusCase, r *predicateRegistry) {
			c[0].Assertions[0].Kind = "accept"
			r.Requirements[0].Bindings[0].Assertion = c[0].Assertions[0]
		},
		"negative-without-sibling": func(c []corpusCase, r *predicateRegistry) { c[0].Expected.Accepted = false },
		"negative-with-absent-sibling": func(c []corpusCase, r *predicateRegistry) {
			c[0].Expected.Accepted = false
			c[0].Related.NearestValid = "absent"
		},
		"stale-fixture": func(c []corpusCase, r *predicateRegistry) { r.Requirements[0].Bindings[0].FixtureSHA256 = "old" },
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			c, r := makeEvidence()
			mutate(c, &r)
			if validateEvidence(c, r, read) == nil {
				t.Fatal("invalid evidence accepted")
			}
		})
	}
	// Supporting labels are legal background; they carry no owned assertion.
	cs, reg = makeEvidence()
	cs[0].Coverage.PrimaryRequirements = nil
	cs[0].Coverage.SupportingRequirements = []string{"R"}
	cs[0].Assertions = nil
	if err := validateEvidence(cs, reg, read); err != nil {
		t.Fatal(err)
	}
}
