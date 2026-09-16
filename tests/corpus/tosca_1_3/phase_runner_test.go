package tosca_1_3_corpus_test

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

// Exercise the actual runner in an isolated test process. Merely testing the
// comparison helper would not catch a runner that stops calling that helper.
func TestEvidenceRunnerRejectsWrongPhase(t *testing.T) {
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	for _, phase := range []string{"namespaces", "hierarchy", "inheritance", "rendering", "normalization"} {
		t.Run(phase, func(t *testing.T) {
			command := exec.Command(executable, "-test.run=^TestEvidenceWrongPhaseChild$", "-test.v")
			command.Env = append(os.Environ(), "TOSCA_EVIDENCE_WRONG_PHASE_CHILD="+phase)
			output, err := command.CombinedOutput()
			if err == nil || !strings.Contains(string(output), "wrong phase: expected "+phase+", actual read") {
				t.Fatalf("runner credited wrong-phase fixture: %v\n%s", err, output)
			}
		})
	}
}
func TestEvidenceWrongPhaseChild(t *testing.T) {
	expected := os.Getenv("TOSCA_EVIDENCE_WRONG_PHASE_CHILD")
	if expected == "" {
		return
	}
	root := corpusRoot(t)
	for _, c := range loadManifest(t, root).Cases {
		if c.ID == "TOSCA13-CORPUS-SERVICE-TEMPLATE-V20-KEY-INVALID" {
			c.Expected.Phase = expected
			runCorpusCase(t, root, c)
			return
		}
	}
	t.Fatal("wrong-phase seed absent")
}
