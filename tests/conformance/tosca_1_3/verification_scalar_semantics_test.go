package tosca_1_3_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/tliron/go-puccini/normal"
	"github.com/tliron/go-puccini/tests/conformance/internal/testsupport"
	"github.com/tliron/go-puccini/tosca/grammars/tosca_v1_3"
	"github.com/tliron/go-puccini/tosca/grammars/tosca_v2_0"
)

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 3.3.2.1 Grammar; 3.3.3.1 Grammar; 3.3.6.1 Grammar;
// 3.3.6.2 Additional requirements
// Requirements: TOSCA13-3.3.2.1-004, TOSCA13-3.3.2.1-005,
// TOSCA13-3.3.3.1-004..006, TOSCA13-3.3.6.1-004..006,
// TOSCA13-3.3.6.2-001..002
// Expected: the TOSCA 1.3 processor accepts every valid boundary form,
// rejects each malformed component during value rendering, and preserves the
// concrete scalar value in normalization.
// Category: positive, negative, boundary, rendering, normalization, direct
func TestVerificationScalarSemantics(t *testing.T) {
	t.Run("TOSCA13-3.3.2.1-004", func(t *testing.T) {
		assertVersionValue(t, "0.1", 0, 1)
		assertScalarRejected(t, "version", `"-1.1"`, "release", "version", "malformed")
	})

	t.Run("TOSCA13-3.3.2.1-005", func(t *testing.T) {
		assertVersionValue(t, "1.0", 1, 0)
		assertScalarRejected(t, "version", `"1"`, "release", "version", "malformed")
		assertScalarRejected(t, "version", `"1.-1"`, "release", "version", "malformed")
	})

	t.Run("TOSCA13-3.3.3.1-004", func(t *testing.T) {
		assertRangeValue(t, "[0, 0]", 0, 0)
		assertScalarRejected(t, "range", `["0", 1]`, "release", "integer")
	})

	t.Run("TOSCA13-3.3.3.1-005", func(t *testing.T) {
		assertRangeValue(t, "[1, 1]", 1, 1)
		assertScalarRejected(t, "range", "[1]", "release", "range", "length")
		assertScalarRejected(t, "range", `[1, "bounded"]`, "release", "UNBOUNDED")
	})

	t.Run("TOSCA13-3.3.3.1-006", func(t *testing.T) {
		assertRangeValue(t, "[1, 1]", 1, 1)
		assertScalarRejected(t, "range", "[2, 1]", "release", "upper", "lower")
	})

	t.Run("TOSCA13-3.3.6.1-004", func(t *testing.T) {
		assertScalarUnitValue(t, "scalar-unit.size", `"1 MB"`, "1000000 bytes")
		assertScalarRejected(t, "scalar-unit.size", `"MB"`, "release", "scalar-unit.size", "malformed")
	})

	t.Run("TOSCA13-3.3.6.1-005", func(t *testing.T) {
		assertScalarUnitValue(t, "scalar-unit.size", `"1 MB"`, "1000000 bytes")
		assertScalarRejected(t, "scalar-unit.size", `"1"`, "release", "scalar-unit.size", "malformed")
	})

	t.Run("TOSCA13-3.3.6.1-006", func(t *testing.T) {
		assertScalarUnitValue(t, "scalar-unit.time", `"1 ms"`, "0.001 seconds")
		assertScalarRejected(t, "scalar-unit.time", `"1 MB"`, "release", "scalar-unit.time", "malformed")
	})

	t.Run("TOSCA13-3.3.6.2-001", func(t *testing.T) {
		for _, value := range []string{`"1MB"`, `"1 MB"`, `"1     MB"`} {
			assertScalarUnitValue(t, "scalar-unit.size", value, "1000000 bytes")
		}
	})

	t.Run("TOSCA13-3.3.6.2-002", func(t *testing.T) {
		assertScalarRejected(t, "scalar-unit.size", `"1"`, "release", "scalar-unit.size", "malformed")
		assertScalarRejected(t, "scalar-unit.size", `"MB"`, "release", "scalar-unit.size", "malformed")
	})
}

func assertVersionValue(t *testing.T, value string, major, minor uint32) {
	t.Helper()
	primitive := parseScalarPrimitive(t, "version", fmt.Sprintf("%q", value))
	version, ok := primitive.(*tosca_v2_0.Version)
	if !ok {
		t.Fatalf("normalized version has type %T", primitive)
	}
	if version.Major != major || version.Minor != minor {
		t.Fatalf("normalized version = %d.%d, want %d.%d", version.Major, version.Minor, major, minor)
	}
}

func assertRangeValue(t *testing.T, value string, lower, upper int64) {
	t.Helper()
	primitive := parseScalarPrimitive(t, "range", value)
	rangeValue, ok := primitive.(*tosca_v1_3.Range)
	if !ok {
		t.Fatalf("normalized range has type %T", primitive)
	}
	if rangeValue.Lower != lower || rangeValue.Upper == nil || *rangeValue.Upper != upper {
		t.Fatalf("normalized range = %#v, want [%d,%d]", rangeValue, lower, upper)
	}
}

func assertScalarUnitValue(t *testing.T, dataType, value, expected string) {
	t.Helper()
	primitive := parseScalarPrimitive(t, dataType, value)
	scalarUnit, ok := primitive.(*tosca_v1_3.ScalarUnit)
	if !ok {
		t.Fatalf("normalized %s has type %T", dataType, primitive)
	}
	if got := scalarUnit.String(); got != expected {
		t.Fatalf("normalized %s = %q, want %q", dataType, got, expected)
	}
}

func parseScalarPrimitive(t *testing.T, dataType, value string) any {
	t.Helper()
	serviceTemplate, problems, err := testsupport.ParseSource(t, scalarVerificationTemplate(dataType, value))
	if err != nil {
		t.Fatalf("valid %s value %s failed: %v\n%s", dataType, value, err, problems)
	}
	node := serviceTemplate.NodeTemplates["node"]
	if node == nil {
		t.Fatal("normalized node is absent")
	}
	primitive, ok := node.Properties["release"].(*normal.Primitive)
	if !ok {
		t.Fatalf("normalized release has type %T", node.Properties["release"])
	}
	return primitive.Primitive
}

func assertScalarRejected(t *testing.T, dataType, value string, expected ...string) {
	t.Helper()
	_, problems, err := testsupport.ParseSource(t, scalarVerificationTemplate(dataType, value))
	if err == nil {
		t.Fatalf("invalid %s value %s was accepted", dataType, value)
	}
	for _, fragment := range expected {
		if !strings.Contains(problems, fragment) {
			t.Fatalf("%s diagnostic does not contain %q:\n%s", dataType, fragment, problems)
		}
	}
}

func scalarVerificationTemplate(dataType, value string) string {
	return fmt.Sprintf(`tosca_definitions_version: tosca_simple_yaml_1_3
node_types:
  ScalarHolder:
    derived_from: tosca.nodes.Root
    properties:
      release:
        type: %s
topology_template:
  node_templates:
    node:
      type: ScalarHolder
      properties:
        release: %s
`, dataType, value)
}
