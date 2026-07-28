package tosca_1_3_test

import (
	"strings"
	"testing"

	"github.com/tliron/go-puccini/tests/conformance/internal/testsupport"
)

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.7.5.4, Additional Requirements
// Requirement: TOSCA13-3.7.5.4-002
// Expected: an ordinary operation name is accepted.
// Category: positive, read phase, direct
func TestInterfaceTypeOperationNameAccepted(t *testing.T) {
	_, parseProblems, err := testsupport.ParseSource(t, interfaceType13WithOperation("run"))
	if err != nil {
		t.Fatalf("ordinary interface operation name was rejected: %v\n%s", err, parseProblems)
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.7.5.4, Additional Requirements
// Requirement: TOSCA13-3.7.5.4-002
// Expected: the exact reserved operation name inputs is rejected in the read phase.
// Category: negative, reserved name, read phase, direct
func TestInterfaceTypeRejectsReservedInputsOperationName(t *testing.T) {
	problems := rejectedReservedOperation(t)
	for _, fragment := range []string{"operations", "inputs", "reserved"} {
		if !strings.Contains(problems, fragment) {
			t.Fatalf("reserved-name diagnostic does not contain %q:\n%s", fragment, problems)
		}
	}
	if strings.Contains(problems, "unsupported keyname") {
		t.Fatalf("reserved name was rejected for the wrong structural reason:\n%s", problems)
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.7.5.4, Additional Requirements; 5.2.1, name case sensitivity
// Requirement: TOSCA13-3.7.5.4-002
// Expected: the reservation applies to the exact lowercase keyname.
// Category: positive, boundary, case sensitivity, direct
func TestInterfaceTypeReservedOperationNameIsCaseSensitive(t *testing.T) {
	_, parseProblems, err := testsupport.ParseSource(t, interfaceType13WithOperation("Inputs"))
	if err != nil {
		t.Fatalf("non-reserved case variant was rejected: %v\n%s", err, parseProblems)
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.7.5.4, Additional Requirements
// Requirement: TOSCA13-3.7.5.4-002
// Expected: repeated reserved-name rejection has a stable diagnostic body.
// Category: negative, determinism, direct
func TestInterfaceTypeReservedOperationDiagnosticDeterministic(t *testing.T) {
	first := interfaceDiagnosticBody(rejectedReservedOperation(t))
	second := interfaceDiagnosticBody(rejectedReservedOperation(t))
	if first != second {
		t.Fatalf("reserved-name diagnostic is nondeterministic:\nfirst:\n%s\nsecond:\n%s", first, second)
	}
}

// Specification: TOSCA Version 2.0
// Section: interface type grammar (cross-version isolation)
// Expected: the TOSCA 1.3 reserved-name policy does not alter TOSCA 2.0.
// Category: positive, cross-version regression
func TestTosca20InterfaceOperationNamedInputsUnchanged(t *testing.T) {
	source := `tosca_definitions_version: tosca_2_0
interface_types:
  Test:
    operations:
      inputs: {}
`
	if _, parseProblems, err := testsupport.ParseSource(t, source); err != nil {
		t.Fatalf("TOSCA 2.0 interface operation behavior changed: %v\n%s", err, parseProblems)
	}
}

func rejectedReservedOperation(t *testing.T) string {
	t.Helper()
	_, parseProblems, err := testsupport.ParseSource(t, interfaceType13WithOperation("inputs"))
	if err == nil {
		t.Fatal("reserved interface operation name inputs was accepted")
	}
	return parseProblems
}

func interfaceType13WithOperation(name string) string {
	return `tosca_definitions_version: tosca_simple_yaml_1_3
interface_types:
  Test:
    operations:
      ` + name + `: {}
`
}

func interfaceDiagnosticBody(problems string) string {
	if index := strings.Index(problems, "interface_types"); index >= 0 {
		return problems[index:]
	}
	return problems
}
