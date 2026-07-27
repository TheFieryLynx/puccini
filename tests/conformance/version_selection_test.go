package conformance_test

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/tliron/go-puccini/tests/conformance/internal/testsupport"
	"github.com/tliron/go-puccini/tosca/grammars"
)

var removedToscaDefinitionsVersions = []string{
	"tosca_simple_yaml_1_0",
	"tosca_simple_yaml_1_1",
	"tosca_simple_yaml_1_2",
	"tosca_simple_profile_for_nfv_1_0",
	"cloudify_dsl_1_3",
}

var removedHotVersions = []string{
	"wallaby",
	"train",
	"stein",
	"rocky",
	"queens",
	"pike",
	"newton",
	"ocata",
	"2021-04-16",
	"2018-08-31",
	"2018-03-02",
	"2017-09-01",
	"2017-02-24",
	"2016-10-14",
	"2016-04-08",
	"2015-10-15",
	"2015-04-30",
	"2014-10-16",
	"2013-05-23",
}

// Specification: TOSCA Simple Profile in YAML 1.3 section 3.10.3.1,
// and TOSCA 2.0 section 6.2, TOSCA Definitions Version
// Expected: the public supported-version API is closed to the two retained values
// Category: positive, public API, regression
func TestSupportedDefinitionsVersions(t *testing.T) {
	got := grammars.SupportedDefinitionsVersions()
	want := []string{"tosca_simple_yaml_1_3", "tosca_2_0"}
	if !slices.Equal(got, want) {
		t.Fatalf("supported versions = %v, want %v", got, want)
	}

	got[0] = "mutated"
	if slices.Equal(got, grammars.SupportedDefinitionsVersions()) {
		t.Fatal("supported-version API exposed mutable registry state")
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3 section 3.10.3.1,
// and TOSCA 2.0 section 6.2, TOSCA Definitions Version
// Expected: every formerly registered tosca_definitions_version is rejected in the read phase
// Category: negative, removed versions, regression
func TestRemovedToscaDefinitionsVersionsAreRejected(t *testing.T) {
	for _, version := range removedToscaDefinitionsVersions {
		t.Run(version, func(t *testing.T) {
			_, problems, err := testsupport.ParseSource(t, fmt.Sprintf("tosca_definitions_version: %s\n", version))
			if err == nil {
				t.Fatalf("removed version %q was accepted", version)
			}
			assertUnsupportedDefinitionsVersion(t, problems, version)
		})
	}
}

// Specification: TOSCA 1.3 section 3.10.1 and TOSCA 2.0 section 6.1
// Expected: every former non-TOSCA HOT selector is rejected as missing the mandatory TOSCA selector
// Category: negative, removed dialect, regression
func TestRemovedHotVersionsAreRejected(t *testing.T) {
	for _, version := range removedHotVersions {
		t.Run(version, func(t *testing.T) {
			_, problems, err := testsupport.ParseSource(t, fmt.Sprintf("heat_template_version: %s\nresources: {}\n", version))
			if err == nil {
				t.Fatalf("removed HOT version %q was accepted", version)
			}
			assertMissingDefinitionsVersion(t, problems)
		})
	}
}

// Specification: TOSCA 1.3 section 3.10.3.1 and TOSCA 2.0 section 6.2
// Expected: unknown values and legacy URI aliases are rejected in the read phase
// Category: negative, unknown version, aliases
func TestUnknownAndAliasVersionsAreRejected(t *testing.T) {
	for _, version := range []string{
		"tosca_future_9_9",
		"http://docs.oasis-open.org/tosca/ns/simple/yaml/1.3",
		"https://docs.oasis-open.org/tosca/ns/simple/yaml/1.3",
	} {
		t.Run(version, func(t *testing.T) {
			_, problems, err := testsupport.ParseSource(t, fmt.Sprintf("tosca_definitions_version: %s\n", version))
			if err == nil {
				t.Fatalf("unknown or alias version %q was accepted", version)
			}
			assertUnsupportedDefinitionsVersion(t, problems, version)
		})
	}
}

// Specification: TOSCA 1.3 section 3.10.1 and TOSCA 2.0 sections 6.1-6.2
// Expected: absence of the mandatory selector is rejected in the read phase
// Category: negative, missing required keyname
func TestMissingDefinitionsVersionIsRejected(t *testing.T) {
	_, problems, err := testsupport.ParseSource(t, "description: missing version\n")
	if err == nil {
		t.Fatal("template without tosca_definitions_version was accepted")
	}
	assertMissingDefinitionsVersion(t, problems)
}

// Specification: TOSCA 1.3 section 3.10.1 and TOSCA 2.0 sections 6.1-6.2
// Expected: a non-string selector is rejected for its YAML type, not as an unknown version
// Category: negative, boundary, wrong YAML type
func TestDefinitionsVersionWrongTypeIsRejected(t *testing.T) {
	_, problems, err := testsupport.ParseSource(t, "tosca_definitions_version: 2\n")
	if err == nil {
		t.Fatal("numeric tosca_definitions_version was accepted")
	}
	if !strings.Contains(problems, `"integer" instead of "string"`) || strings.Contains(problems, "unsupported TOSCA definitions version") {
		t.Fatalf("wrong rejection reason:\n%s", problems)
	}
}

func assertUnsupportedDefinitionsVersion(t *testing.T, problems string, version string) {
	t.Helper()
	for _, expected := range []string{
		"unsupported TOSCA definitions version",
		version,
		"tosca_simple_yaml_1_3",
		"tosca_2_0",
	} {
		if !strings.Contains(problems, expected) {
			t.Fatalf("missing %q in rejection:\n%s", expected, problems)
		}
	}
}

func assertMissingDefinitionsVersion(t *testing.T, problems string) {
	t.Helper()
	if !strings.Contains(problems, "tosca_definitions_version") || !strings.Contains(problems, "missing required keyname") {
		t.Fatalf("wrong rejection reason:\n%s", problems)
	}
}
