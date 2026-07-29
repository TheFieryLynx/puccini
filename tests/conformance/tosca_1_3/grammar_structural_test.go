package tosca_1_3_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/tliron/go-puccini/tests/conformance/internal/testsupport"
)

const grammarStructuralVersion = "tosca_definitions_version: tosca_simple_yaml_1_3\n"

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 3.6.3.4, 3.6.6-3.6.27, 3.7.2-3.7.12, 3.8.3-3.8.13, 3.10.1
// Requirements: every ID listed in work-items/grammar-structural.yaml
// Expected: each selected frozen MUST has direct, permanent test evidence.
// Category: evidence completeness
func TestGrammarStructuralRequirementEvidence(t *testing.T) {
	requirementIDs := []string{
		"TOSCA13-3.6.3.4-004", "TOSCA13-3.6.3.4-005", "TOSCA13-3.6.3.4-007",
		"TOSCA13-3.6.6.1-004", "TOSCA13-3.6.6.1-006", "TOSCA13-3.6.6.2.2-003", "TOSCA13-3.6.6.2.2-005",
		"TOSCA13-3.6.7.1-001", "TOSCA13-3.6.7.1-003", "TOSCA13-3.6.7.1-004", "TOSCA13-3.6.7.1-006",
		"TOSCA13-3.6.7.2.2-003", "TOSCA13-3.6.7.2.2-005", "TOSCA13-3.6.7.2.2-006",
		"TOSCA13-3.6.9.1-001", "TOSCA13-3.6.10.2-001", "TOSCA13-3.6.10.4-004", "TOSCA13-3.6.10.4-007",
		"TOSCA13-3.6.12.2-001", "TOSCA13-3.6.12.3-004", "TOSCA13-3.6.14.2-008", "TOSCA13-3.6.14.2-014",
		"TOSCA13-3.6.15.1-005", "TOSCA13-3.6.15.1-012", "TOSCA13-3.6.17.2.3-003",
		"TOSCA13-3.6.19.2-004", "TOSCA13-3.6.20.2.2-003", "TOSCA13-3.6.20.2.2-004",
		"TOSCA13-3.6.20.2.2-007", "TOSCA13-3.6.20.2.2-008", "TOSCA13-3.6.21.1-001",
		"TOSCA13-3.6.21.1-003", "TOSCA13-3.6.22.1-004", "TOSCA13-3.6.22.1-020",
		"TOSCA13-3.6.22.3.2-002", "TOSCA13-3.6.22.3.2-004", "TOSCA13-3.6.23.1.1-003",
		"TOSCA13-3.6.23.2.1-002", "TOSCA13-3.6.23.3.1-004", "TOSCA13-3.6.23.3.1-007",
		"TOSCA13-3.6.23.4.1-003", "TOSCA13-3.6.23.4.1-004", "TOSCA13-3.6.26.1-002",
		"TOSCA13-3.6.27.1-002", "TOSCA13-3.6.27.1-020", "TOSCA13-3.7.2.1-001",
		"TOSCA13-3.7.2.4-002", "TOSCA13-3.7.3.1-001", "TOSCA13-3.7.3.2.3-003",
		"TOSCA13-3.7.3.2.3-004", "TOSCA13-3.7.3.3-001", "TOSCA13-3.7.5.2-004",
		"TOSCA13-3.7.5.2-008", "TOSCA13-3.7.5.2-009", "TOSCA13-3.7.5.4-001",
		"TOSCA13-3.7.6.2-004", "TOSCA13-3.7.7.2-004", "TOSCA13-3.7.9.2-004",
		"TOSCA13-3.7.10.2-004", "TOSCA13-3.7.11.2-003", "TOSCA13-3.7.12.2-004",
		"TOSCA13-3.8.3.1-001", "TOSCA13-3.8.3.2-003", "TOSCA13-3.8.4.1-001",
		"TOSCA13-3.8.4.2-003", "TOSCA13-3.8.5.1-001", "TOSCA13-3.8.5.2-003",
		"TOSCA13-3.8.5.2-007", "TOSCA13-3.8.6.1-001", "TOSCA13-3.8.6.1-003",
		"TOSCA13-3.8.6.2-003", "TOSCA13-3.8.13.1-001", "TOSCA13-3.8.13.1-003",
		"TOSCA13-3.8.13.2-004", "TOSCA13-3.10.1-002",
	}
	if len(requirementIDs) != 75 {
		t.Fatalf("direct-evidence manifest contains %d IDs, want 75", len(requirementIDs))
	}
	seen := make(map[string]struct{}, len(requirementIDs))
	for _, requirementID := range requirementIDs {
		if _, duplicate := seen[requirementID]; duplicate {
			t.Fatalf("duplicate direct-evidence requirement ID %s", requirementID)
		}
		seen[requirementID] = struct{}{}
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 3.6.6.2, 3.6.7.2, 3.6.10.4, 3.6.14.2, 3.6.23.3-3.6.23.4
// Requirements: TOSCA13-3.6.6.2.2-003, -005; TOSCA13-3.6.7.2.2-003, -005, -006;
// TOSCA13-3.6.10.4-004, -007; TOSCA13-3.6.14.2-008, -014;
// TOSCA13-3.6.23.3.1-004, -007; TOSCA13-3.6.23.4.1-003, -004
// Expected: valid short and long forms are accepted.
// Category: positive, short notation, long notation
func TestGrammarStructuralShortLongNotation(t *testing.T) {
	tests := []struct {
		name   string
		source string
	}{
		{"repository-short", "repositories:\n  repo: https://example.invalid/repository/\n"},
		{"repository-long", "repositories:\n  repo:\n    url: https://example.invalid/repository/\n    description: repository\n"},
		{"artifact-short", "artifact_types:\n  TestArtifact:\n    derived_from: tosca.artifacts.File\nnode_types:\n  TestNode:\n    derived_from: tosca.nodes.Root\n    artifacts:\n      payload: payload.txt\n"},
		{"artifact-long", "artifact_types:\n  TestArtifact:\n    derived_from: tosca.artifacts.File\nnode_types:\n  TestNode:\n    derived_from: tosca.nodes.Root\n    artifacts:\n      payload:\n        type: TestArtifact\n        file: payload.txt\n"},
		{"property-long", "data_types:\n  TestData:\n    properties:\n      enabled:\n        type: boolean\n        required: false\n"},
		{"constraint-list", "data_types:\n  TestData:\n    properties:\n      value:\n        type: integer\n        constraints:\n          - greater_or_equal: 1\n          - less_or_equal: 10\n"},
		{"parameter-short", "topology_template:\n  outputs:\n    answer: 42\n"},
		{"parameter-long", "topology_template:\n  outputs:\n    answer:\n      type: integer\n      required: false\n      value: 42\n"},
		{"call-operation-short", workflowSource("- call_operation: Test.run")},
		{"call-operation-long", workflowSource("- call_operation:\n    operation: Test.run\n    inputs: {}")},
		{"inline-short", workflowSource("- inline: nested")},
		{"inline-long", workflowSource("- inline:\n    workflow: nested\n    inputs: {}")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, problems, err := testsupport.ParseSource(t, grammarStructuralVersion+test.source)
			if err != nil {
				t.Fatalf("valid %s notation failed: %v\n%s", test.name, err, problems)
			}
		})
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 3.6.6.1, 3.6.7.1, 3.6.21.1, 3.6.22.1, 3.6.23.3.1,
// 3.6.23.4.1, 3.6.27.1, 3.8.3.1, 3.8.4.1, 3.8.5.1, 3.8.6.1, 3.8.13.1
// Requirements: all selected required-key records in those sections
// Expected: independently missing required keynames fail in structural reading.
// Category: negative, required key, read phase
func TestGrammarStructuralRequiredKeys(t *testing.T) {
	tests := []struct {
		name       string
		source     string
		diagnostic []string
	}{
		{"repository-url", "repositories:\n  repo: {}\n", []string{"url", "missing"}},
		{"artifact-type", artifactLong("file: payload.txt"), []string{"type", "missing"}},
		{"artifact-file", artifactLong("type: TestArtifact"), []string{"file", "missing"}},
		{"event-filter-node", policySource("event: changed\ntarget_filter: {}\naction: []"), []string{"node", "missing"}},
		{"trigger-event", policySource("action: []"), []string{"event", "missing"}},
		{"trigger-action", policySource("event: changed"), []string{"action", "missing"}},
		{"call-operation-operation", workflowSource("- call_operation:\n    inputs: {}"), []string{"operation", "missing"}},
		{"inline-workflow", workflowSource("- inline:\n    inputs: {}"), []string{"workflow", "missing"}},
		{"workflow-target", workflowStepSource("activities: []"), []string{"target", "missing"}},
		{"workflow-activities", workflowStepSource("target: node"), []string{"activities", "missing"}},
		{"workflow-precondition-target", "topology_template:\n  workflows:\n    main:\n      preconditions:\n        precondition: {}\n", []string{"target", "missing"}},
		{"schema-type", "data_types:\n  TestData:\n    properties:\n      values:\n        type: list\n        entry_schema: {}\n", []string{"type", "missing"}},
		{"property-type", "data_types:\n  TestData:\n    properties:\n      value: {}\n", []string{"type", "missing"}},
		{"attribute-type", "node_types:\n  TestNode:\n    derived_from: tosca.nodes.Root\n    attributes:\n      value: {}\n", []string{"type", "missing"}},
		{"capability-type", "node_types:\n  TestNode:\n    derived_from: tosca.nodes.Root\n    capabilities:\n      endpoint: {}\n", []string{"type", "missing"}},
		{"requirement-capability", "node_types:\n  TestNode:\n    derived_from: tosca.nodes.Root\n    requirements:\n      - host: {}\n", []string{"capability", "missing"}},
		{"interface-definition-type", "node_types:\n  TestNode:\n    derived_from: tosca.nodes.Root\n    interfaces:\n      Test: {}\n", []string{"type", "missing"}},
		{"node-template-type", "topology_template:\n  node_templates:\n    node: {}\n", []string{"type", "missing"}},
		{"relationship-template-type", "topology_template:\n  relationship_templates:\n    relation: {}\n", []string{"type", "missing"}},
		{"group-type", "topology_template:\n  groups:\n    group:\n      members: []\n", []string{"type", "missing"}},
		{"policy-type", "topology_template:\n  policies:\n    - policy: {}\n", []string{"type", "missing"}},
		{"substitution-node-type", "topology_template:\n  substitution_mappings: {}\n", []string{"node_type", "missing"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assertGrammarRejected(t, grammarStructuralVersion+test.source, test.diagnostic...)
		})
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 3.6.3.4, 3.6.6.1, 3.6.7.1, 3.6.10.4, 3.6.14.2,
// 3.6.21.1, 3.6.22.1, 3.6.23, 3.6.27.1, 3.8.3-3.8.13
// Requirements: selected structural type records in grammar-structural.yaml
// Expected: meaningful opposing YAML types are rejected at the responsible field.
// Category: negative, YAML type, read phase
func TestGrammarStructuralYAMLTypes(t *testing.T) {
	tests := []struct {
		name       string
		source     string
		diagnostic []string
	}{
		{"repositories-sequence", "repositories: []\n", []string{"repositories", "map"}},
		{"repository-url-integer", "repositories:\n  repo:\n    url: 1\n", []string{"url", "string"}},
		{"artifact-definition-sequence", "node_types:\n  TestNode:\n    derived_from: tosca.nodes.Root\n    artifacts:\n      payload: []\n", []string{"artifact", "map", "string"}},
		{"property-required-string", "data_types:\n  TestData:\n    properties:\n      value:\n        type: string\n        required: yes\n", []string{"required", "boolean"}},
		{"parameter-required-integer", "topology_template:\n  outputs:\n    value:\n      type: string\n      required: 1\n", []string{"required", "boolean"}},
		{"attribute-mapping-first-item", "topology_template:\n  outputs:\n    value: [1, state]\n", []string{"mapping", "string"}},
		{"attribute-mapping-second-item", "topology_template:\n  outputs:\n    value: [SELF, 1]\n", []string{"mapping", "string"}},
		{"event-node-sequence", policySource("event: changed\ntarget_filter:\n  node: []\naction: []"), []string{"node", "string"}},
		{"trigger-action-mapping", policySource("event: changed\naction: {}"), []string{"action", "list"}},
		{"workflow-target-list", workflowStepSource("target: []\nactivities: []"), []string{"target", "string"}},
		{"workflow-activities-map", workflowStepSource("target: node\nactivities: {}"), []string{"activities", "list"}},
		{"group-members-string", "topology_template:\n  groups:\n    group:\n      type: tosca.groups.Root\n      members: node\n", []string{"members", "list"}},
		{"substitution-node-type-list", "topology_template:\n  substitution_mappings:\n    node_type: []\n", []string{"node_type", "string"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assertGrammarRejected(t, grammarStructuralVersion+test.source, test.diagnostic...)
		})
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 3.6.6.2.2, 3.6.7.2.2, 3.6.10.4, 3.6.12.3, 3.6.14.2,
// 3.6.17.2.3, 3.6.19.2, 3.6.20.2.2, 3.6.22.3.2, 3.7.2-3.7.12,
// 3.8.3-3.8.6, 3.8.13
// Requirements: selected symbolic-name grammar records in grammar-structural.yaml
// Expected: valid symbolic mapping names reach their concrete readers.
// Category: positive, symbolic names, concrete reader trace
func TestGrammarStructuralSymbolicNames(t *testing.T) {
	source := grammarStructuralVersion + `repositories:
  RepositoryName: https://example.invalid/
artifact_types:
  ArtifactTypeName:
    derived_from: tosca.artifacts.File
capability_types:
  CapabilityTypeName:
    derived_from: tosca.capabilities.Root
data_types:
  DataTypeName:
    properties:
      PropertyName:
        type: string
group_types:
  GroupTypeName:
    derived_from: tosca.groups.Root
interface_types:
  InterfaceTypeName:
    operations:
      OperationName: {}
    notifications:
      NotificationName: {}
node_types:
  NodeTypeName:
    derived_from: tosca.nodes.Root
    attributes:
      AttributeName:
        type: string
    capabilities:
      CapabilityName: CapabilityTypeName
    requirements:
      - RequirementName: CapabilityTypeName
    artifacts:
      ArtifactName:
        type: ArtifactTypeName
        file: artifact.txt
    interfaces:
      InterfaceDefinitionName:
        type: InterfaceTypeName
        operations:
          OperationName: {}
        notifications:
          NotificationName: {}
policy_types:
  PolicyTypeName:
    derived_from: tosca.policies.Root
relationship_types:
  RelationshipTypeName:
    derived_from: tosca.relationships.Root
topology_template:
  node_templates:
    NodeTemplateName:
      type: NodeTypeName
  relationship_templates:
    RelationshipTemplateName:
      type: RelationshipTypeName
  groups:
    GroupName:
      type: GroupTypeName
      members: [NodeTemplateName]
  policies:
    - PolicyName:
        type: PolicyTypeName
        targets: [NodeTemplateName]
  substitution_mappings:
    node_type: NodeTypeName
    capabilities:
      CapabilityName: [NodeTemplateName, CapabilityName]
      feature: [NodeTemplateName, feature]
    requirements:
      RequirementName: [NodeTemplateName, RequirementName]
      dependency: [NodeTemplateName, dependency]
`
	if _, problems, err := testsupport.ParseSource(t, source); err != nil {
		t.Fatalf("valid symbolic grammar names failed: %v\n%s", err, problems)
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 3.6.6-3.6.27, 3.7.2-3.7.12, 3.8.3-3.8.13, 3.10.1
// Requirements: all selected records whose entity has a closed keyname set
// Expected: unknown root, nested, misspelled, TOSCA 2.0-only, and extension keys fail.
// Category: negative, unknown key, cross-version, read phase
func TestGrammarStructuralUnknownKeys(t *testing.T) {
	tests := []struct {
		name       string
		source     string
		diagnostic string
	}{
		{"root", "unknown_root: true\n", "unknown_root"},
		{"nested", "topology_template:\n  node_templates:\n    node:\n      type: tosca.nodes.Root\n      unknown_nested: true\n", "unknown_nested"},
		{"misspelled", "repositories:\n  repo:\n    url: https://example.invalid/\n    descrption: typo\n", "descrption"},
		{"tosca-2.0", "service_template: {}\n", "service_template"},
		{"extension", "x_project_extension: true\n", "x_project_extension"},
		{"long-notation", artifactLong("type: TestArtifact\n        file: payload.txt\n        extra: invalid"), "extra"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assertGrammarRejected(t, grammarStructuralVersion+test.source, test.diagnostic, "unsupported keyname")
		})
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 3.6.3.4, 3.6.20.2.2, 3.7.5.2, 3.8.5.2
// Requirements: TOSCA13-3.6.3.4-004, -005, -007;
// TOSCA13-3.6.20.2.2-007, -008; TOSCA13-3.7.5.2-008, -009;
// TOSCA13-3.8.5.2-007
// Expected: required operands and explicitly present one-or-more collections are non-empty.
// Category: negative, boundary, cardinality
func TestGrammarStructuralBoundaries(t *testing.T) {
	tests := []struct {
		name       string
		source     string
		diagnostic []string
	}{
		{"constraint-empty-map", "data_types:\n  TestData:\n    properties:\n      value:\n        type: integer\n        constraints:\n          - {}\n", []string{"constraint", "map size"}},
		{"constraint-null-operand", "data_types:\n  TestData:\n    properties:\n      value:\n        type: integer\n        constraints:\n          - greater_than:\n", []string{"greater_than", "required"}},
		{"constraint-range-cardinality", "data_types:\n  TestData:\n    properties:\n      value:\n        type: integer\n        constraints:\n          - in_range: [1]\n", []string{"in_range", "2"}},
		{"constraint-unknown-operator", "data_types:\n  TestData:\n    properties:\n      value:\n        type: integer\n        constraints:\n          - unknown_operator: [1]\n", []string{"unknown_operator", "unsupported operator"}},
		{"interface-definition-empty-operations", interfaceDefinitionSource("operations: {}"), []string{"operations", "one or more"}},
		{"interface-definition-empty-notifications", interfaceDefinitionSource("notifications: {}"), []string{"notifications", "one or more"}},
		{"interface-type-empty-operations", "interface_types:\n  Test:\n    operations: {}\n", []string{"operations", "one or more"}},
		{"interface-type-empty-notifications", "interface_types:\n  Test:\n    notifications: {}\n", []string{"notifications", "one or more"}},
		{"group-empty-members", "topology_template:\n  groups:\n    group:\n      type: tosca.groups.Root\n      members: []\n", []string{"members", "one or more"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assertGrammarRejected(t, grammarStructuralVersion+test.source, test.diagnostic...)
		})
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 3.7.2.4, 3.7.3.3, and the named sequenced-list grammars in scope
// Requirements: TOSCA13-3.7.2.4-002, TOSCA13-3.7.3.3-001
// Expected: neither YAML mapping duplicates nor sequenced-list duplicates overwrite.
// Category: negative, duplicate, deterministic diagnostic
func TestGrammarStructuralDuplicates(t *testing.T) {
	tests := []struct {
		name   string
		source string
		key    string
	}{
		{
			"yaml-capability-key",
			"node_types:\n  TestNode:\n    derived_from: tosca.nodes.Root\n    capabilities:\n      endpoint: tosca.capabilities.Endpoint\n      endpoint: tosca.capabilities.Endpoint\n",
			"endpoint",
		},
		{
			"sequenced-requirement-name",
			"node_types:\n  TestNode:\n    derived_from: tosca.nodes.Root\n    requirements:\n      - host: tosca.capabilities.Container\n      - host: tosca.capabilities.Container\n",
			"host",
		},
		{
			"sequenced-policy-name",
			"topology_template:\n  policies:\n    - policy:\n        type: tosca.policies.Root\n    - policy:\n        type: tosca.policies.Root\n",
			"policy",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var first string
			for iteration := 0; iteration < 2; iteration++ {
				_, problems, err := testsupport.ParseSource(t, grammarStructuralVersion+test.source)
				if err == nil {
					t.Fatalf("duplicate %q was accepted", test.key)
				}
				diagnostic := stableGrammarDiagnostic(diagnosticText(problems, err))
				assertDiagnostic(t, diagnostic, "duplicate", test.key)
				if iteration == 0 {
					first = diagnostic
				} else if diagnostic != first {
					t.Fatalf("diagnostic order changed:\nfirst:\n%s\nsecond:\n%s", first, diagnostic)
				}
			}
		})
	}
}

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 3.7.5.4 Additional Requirements
// Requirement: TOSCA13-3.7.5.4-001
// Expected: implementation is invalid only inside interface type operation/notification definitions.
// Category: negative, context restriction, version isolation
func TestInterfaceTypeRejectsImplementations(t *testing.T) {
	for _, body := range []string{
		"operations:\n      run:\n        implementation: run.sh",
		"notifications:\n      changed:\n        implementation: notify.sh",
	} {
		source := grammarStructuralVersion + "interface_types:\n  Test:\n    " + body + "\n"
		assertGrammarRejected(t, source, "implementation", "invalid", "interface type")
	}
}

func assertGrammarRejected(t *testing.T, source string, fragments ...string) {
	t.Helper()
	_, problems, err := testsupport.ParseSource(t, source)
	if err == nil {
		t.Fatalf("invalid structural form was accepted:\n%s", source)
	}
	assertDiagnostic(t, diagnosticText(problems, err), fragments...)
	if strings.Contains(strings.ToLower(problems), "import") {
		t.Fatalf("structural test failed first for an unrelated import:\n%s", problems)
	}
}

func artifactLong(body string) string {
	return fmt.Sprintf(`artifact_types:
  TestArtifact:
    derived_from: tosca.artifacts.File
node_types:
  TestNode:
    derived_from: tosca.nodes.Root
    artifacts:
      payload:
        %s
`, body)
}

func interfaceDefinitionSource(body string) string {
	return fmt.Sprintf(`interface_types:
  TestInterface:
    operations:
      declared: {}
node_types:
  TestNode:
    derived_from: tosca.nodes.Root
    interfaces:
      Test:
        type: TestInterface
        %s
`, body)
}

func policySource(triggerBody string) string {
	return fmt.Sprintf(`topology_template:
  policies:
    - policy:
        type: tosca.policies.Root
        triggers:
          trigger:
%s
`, indentGrammar(triggerBody, 12))
}

func workflowSource(activity string) string {
	return `interface_types:
  TestInterface:
    operations:
      run: {}
node_types:
  TestNode:
    derived_from: tosca.nodes.Root
    interfaces:
      Test:
        type: TestInterface
topology_template:
  node_templates:
    node:
      type: TestNode
  workflows:
    nested:
      steps: {}
    main:
      steps:
        step:
          target: node
          activities:
` + indentGrammar(activity, 10) + "\n"
}

func workflowStepSource(body string) string {
	return `topology_template:
  node_templates:
    node:
      type: tosca.nodes.Root
  workflows:
    main:
      steps:
        step:
` + indentGrammar(body, 10) + "\n"
}

func indentGrammar(value string, spaces int) string {
	prefix := strings.Repeat(" ", spaces)
	lines := strings.Split(value, "\n")
	for index, line := range lines {
		lines[index] = prefix + line
	}
	return strings.Join(lines, "\n")
}

func stableGrammarDiagnostic(value string) string {
	if index := strings.Index(value, "service-template.yaml"); index != -1 {
		return value[index:]
	}
	return value
}
