package tosca_1_3_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/tliron/go-puccini/tests/conformance/internal/testsupport"
	"github.com/tliron/go-puccini/tosca/grammars/tosca_v1_3"
	"github.com/tliron/go-puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/go-puccini/tosca/parser"
)

// Specification: TOSCA Simple Profile in YAML 1.3
// Sections: 5.3.6, 5.3.7, 5.3.11, 5.5.7, 5.5.13, 5.7.1,
// 5.7.5, 5.9.1, 5.9.5, 5.9.8, and 5.9.9
// Requirements: TOSCA13-5.3.6.1-005, TOSCA13-5.3.6.1-008,
// TOSCA13-5.3.7.1-002, TOSCA13-5.3.7.1-005,
// TOSCA13-5.3.7.2-005, TOSCA13-5.3.7.2-008,
// TOSCA13-5.3.11.1-002, TOSCA13-5.3.11.2-005,
// TOSCA13-5.5.7.1-003, TOSCA13-5.5.7.2-002,
// TOSCA13-5.5.7.3-005, TOSCA13-5.5.13.1-003,
// TOSCA13-5.5.13.1-007, TOSCA13-5.7.1.1-002,
// TOSCA13-5.7.1.1-005, TOSCA13-5.7.1.1-010,
// TOSCA13-5.7.5.1-003, TOSCA13-5.9.1.2-002,
// TOSCA13-5.9.1.2-005, TOSCA13-5.9.1.2-010,
// TOSCA13-5.9.5.3-001, TOSCA13-5.9.8.1-002,
// TOSCA13-5.9.9.1-002
// Expected: the ordinary processor path loads each required normative field
// into the corresponding effective type definition with the specified kind and
// type, including both WebServer endpoint capabilities.
// Category: positive, normative profile, effective definition, direct
func TestVerificationNormativeProfileFields(t *testing.T) {
	parserContext, release := parseNormativeProfile(t)
	defer release()

	propertyChecks := []struct {
		requirementID string
		entityKind    string
		entityName    string
		fieldName     string
		dataType      string
	}{
		{"TOSCA13-5.3.6.1-005", "data", "tosca.datatypes.Credential", "token_type", "string"},
		{"TOSCA13-5.3.6.1-008", "data", "tosca.datatypes.Credential", "token", "string"},
		{"TOSCA13-5.3.7.1-002", "data", "tosca.datatypes.TimeInterval", "start_time", "timestamp"},
		{"TOSCA13-5.3.7.1-005", "data", "tosca.datatypes.TimeInterval", "end_time", "timestamp"},
		{"TOSCA13-5.3.7.2-005", "data", "tosca.datatypes.TimeInterval", "start_time", "timestamp"},
		{"TOSCA13-5.3.7.2-008", "data", "tosca.datatypes.TimeInterval", "end_time", "timestamp"},
		{"TOSCA13-5.3.11.1-002", "data", "tosca.datatypes.network.PortSpec", "protocol", "string"},
		{"TOSCA13-5.3.11.2-005", "data", "tosca.datatypes.network.PortSpec", "protocol", "string"},
		{"TOSCA13-5.5.7.1-003", "capability", "tosca.capabilities.Endpoint", "protocol", "string"},
		{"TOSCA13-5.5.7.3-005", "capability", "tosca.capabilities.Endpoint", "protocol", "string"},
		{"TOSCA13-5.5.13.1-003", "capability", "tosca.capabilities.Scalable", "min_instances", "integer"},
		{"TOSCA13-5.5.13.1-007", "capability", "tosca.capabilities.Scalable", "max_instances", "integer"},
		{"TOSCA13-5.7.5.1-003", "relationship", "tosca.relationships.AttachesTo", "location", "string"},
		{"TOSCA13-5.9.8.1-002", "node", "tosca.nodes.Database", "name", "string"},
		{"TOSCA13-5.9.9.1-002", "node", "tosca.nodes.Abstract.Storage", "name", "string"},
	}
	for _, check := range propertyChecks {
		check := check
		t.Run(check.requirementID, func(t *testing.T) {
			definition := findProfilePropertyDefinition(
				t, parserContext, check.entityKind, check.entityName, check.fieldName)
			if !definition.IsRequired() {
				t.Fatalf("%s.%s is optional, want required", check.entityName, check.fieldName)
			}
			if definition.DataTypeName == nil || *definition.DataTypeName != check.dataType {
				t.Fatalf("%s.%s type = %v, want %s",
					check.entityName, check.fieldName, definition.DataTypeName, check.dataType)
			}
			assertNormativeRequirednessCannotBeRelaxed(
				t, check.entityKind, check.entityName, check.fieldName)
		})
	}

	attributeChecks := []struct {
		requirementID string
		entityKind    string
		entityName    string
		fieldName     string
	}{
		{"TOSCA13-5.5.7.2-002", "capability", "tosca.capabilities.Endpoint", "ip_address"},
		{"TOSCA13-5.7.1.1-002", "relationship", "tosca.relationships.Root", "tosca_id"},
		{"TOSCA13-5.7.1.1-005", "relationship", "tosca.relationships.Root", "tosca_name"},
		{"TOSCA13-5.7.1.1-010", "relationship", "tosca.relationships.Root", "state"},
		{"TOSCA13-5.9.1.2-002", "node", "tosca.nodes.Root", "tosca_id"},
		{"TOSCA13-5.9.1.2-005", "node", "tosca.nodes.Root", "tosca_name"},
		{"TOSCA13-5.9.1.2-010", "node", "tosca.nodes.Root", "state"},
	}
	for _, check := range attributeChecks {
		check := check
		t.Run(check.requirementID, func(t *testing.T) {
			definition := findProfileAttributeDefinition(
				t, parserContext, check.entityKind, check.entityName, check.fieldName)
			if definition.DataTypeName == nil || *definition.DataTypeName != "string" {
				t.Fatalf("%s.%s type = %v, want string",
					check.entityName, check.fieldName, definition.DataTypeName)
			}
		})
	}

	t.Run("TOSCA13-5.9.5.3-001", func(t *testing.T) {
		webServer := findProfileNodeType(t, parserContext, "tosca.nodes.WebServer")
		for name, wantType := range map[string]string{
			"admin_endpoint": "tosca.capabilities.Endpoint.Admin",
			"data_endpoint":  "tosca.capabilities.Endpoint",
		} {
			definition := webServer.CapabilityDefinitions[name]
			if definition == nil || definition.CapabilityTypeName == nil ||
				*definition.CapabilityTypeName != wantType {
				t.Fatalf("WebServer capability %s = %#v, want %s", name, definition, wantType)
			}
		}
	})
}

func assertNormativeRequirednessCannotBeRelaxed(
	t *testing.T,
	entityKind string,
	parentName string,
	fieldName string,
) {
	t.Helper()
	section := map[string]string{
		"data":         "data_types",
		"capability":   "capability_types",
		"relationship": "relationship_types",
		"node":         "node_types",
	}[entityKind]
	if section == "" {
		t.Fatalf("unsupported normative entity kind %q", entityKind)
	}
	source := fmt.Sprintf(`tosca_definitions_version: tosca_simple_yaml_1_3
%s:
  example.Invalid:
    derived_from: %s
    properties:
      %s:
        required: false
topology_template: {}
`, section, parentName, fieldName)
	_, problems, err := testsupport.ParseSource(t, source)
	if err == nil {
		t.Fatalf("%s.%s requiredness was relaxed", parentName, fieldName)
	}
	for _, expected := range []string{
		fmt.Sprintf(`properties["%s"].required`, fieldName),
		"cannot refine true to false",
	} {
		if !strings.Contains(problems, expected) {
			t.Fatalf("requiredness diagnostic does not contain %q:\n%s", expected, problems)
		}
	}
}

func findProfilePropertyDefinition(
	t *testing.T,
	parserContext *parser.Context,
	entityKind string,
	entityName string,
	fieldName string,
) *tosca_v2_0.PropertyDefinition {
	t.Helper()
	var definition *tosca_v2_0.PropertyDefinition
	switch entityKind {
	case "data":
		definition = findProfileDataType(t, parserContext, entityName).PropertyDefinitions[fieldName]
	case "capability":
		definition = findProfileCapabilityType(t, parserContext, entityName).PropertyDefinitions[fieldName]
	case "relationship":
		definition = findProfileRelationshipType(t, parserContext, entityName).PropertyDefinitions[fieldName]
	case "node":
		definition = findProfileNodeType(t, parserContext, entityName).PropertyDefinitions[fieldName]
	default:
		t.Fatalf("unsupported normative entity kind %q", entityKind)
	}
	if definition == nil {
		t.Fatalf("%s property %q is absent", entityName, fieldName)
	}
	return definition
}

func findProfileAttributeDefinition(
	t *testing.T,
	parserContext *parser.Context,
	entityKind string,
	entityName string,
	fieldName string,
) *tosca_v2_0.AttributeDefinition {
	t.Helper()
	var definition *tosca_v2_0.AttributeDefinition
	switch entityKind {
	case "capability":
		definition = findProfileCapabilityType(t, parserContext, entityName).AttributeDefinitions[fieldName]
	case "relationship":
		definition = findProfileRelationshipType(t, parserContext, entityName).AttributeDefinitions[fieldName]
	case "node":
		definition = findProfileNodeType(t, parserContext, entityName).AttributeDefinitions[fieldName]
	default:
		t.Fatalf("unsupported normative entity kind %q", entityKind)
	}
	if definition == nil {
		t.Fatalf("%s attribute %q is absent", entityName, fieldName)
	}
	return definition
}

func findProfileDataType(t *testing.T, parserContext *parser.Context, name string) *tosca_v2_0.DataType {
	t.Helper()
	for _, file := range parserContext.Files {
		if profile, ok := file.EntityPtr.(*tosca_v1_3.File); ok {
			for _, type_ := range profile.DataTypes {
				if type_.Name == name {
					return type_
				}
			}
		}
	}
	t.Fatalf("normative data type %q is absent", name)
	return nil
}

func findProfileCapabilityType(t *testing.T, parserContext *parser.Context, name string) *tosca_v2_0.CapabilityType {
	t.Helper()
	for _, file := range parserContext.Files {
		if profile, ok := file.EntityPtr.(*tosca_v1_3.File); ok {
			for _, type_ := range profile.CapabilityTypes {
				if type_.Name == name {
					return type_
				}
			}
		}
	}
	t.Fatalf("normative capability type %q is absent", name)
	return nil
}

func findProfileRelationshipType(t *testing.T, parserContext *parser.Context, name string) *tosca_v2_0.RelationshipType {
	t.Helper()
	for _, file := range parserContext.Files {
		if profile, ok := file.EntityPtr.(*tosca_v1_3.File); ok {
			for _, type_ := range profile.RelationshipTypes {
				if type_.Name == name {
					return type_
				}
			}
		}
	}
	t.Fatalf("normative relationship type %q is absent", name)
	return nil
}
