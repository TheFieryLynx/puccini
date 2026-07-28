package tosca_v1_3

import (
	"sort"

	"github.com/tliron/go-puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/go-puccini/tosca/parsing"
)

// TOSCA Simple Profile in YAML 1.3 sections 3.6.10.1, 3.6.10.5,
// 3.6.12.1, and 3.6.12.4.
//
// Reflection is intentionally a TOSCA 1.3 effective-type policy. It runs after
// hierarchy/inheritance and before assignments render. The shared v2 model is
// used only as storage and rendering infrastructure; TOSCA 2.0 never invokes
// this hook.
func reflectFilePropertyDefinitions(file *tosca_v2_0.File) {
	for _, capabilityType := range file.CapabilityTypes {
		reflectPropertyDefinitions(capabilityType.PropertyDefinitions, capabilityType.AttributeDefinitions)
	}
	for _, relationshipType := range file.RelationshipTypes {
		reflectPropertyDefinitions(relationshipType.PropertyDefinitions, relationshipType.AttributeDefinitions)
	}
	for _, nodeType := range file.NodeTypes {
		reflectPropertyDefinitions(nodeType.PropertyDefinitions, nodeType.AttributeDefinitions)
		for _, capabilityDefinition := range nodeType.CapabilityDefinitions {
			reflectPropertyDefinitions(capabilityDefinition.PropertyDefinitions, capabilityDefinition.AttributeDefinitions)
		}
	}
}

func reflectPropertyDefinitions(properties tosca_v2_0.PropertyDefinitions, attributes tosca_v2_0.AttributeDefinitions) {
	names := make([]string, 0, len(properties))
	for name := range properties {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		property := properties[name]
		if property == nil {
			continue
		}

		if attribute, exists := attributes[name]; exists {
			if !sameEffectiveDataType(property.DataType, attribute.DataType) {
				attribute.Context.ReportPathf(
					0,
					"attribute %q type %q conflicts with reflected property type %q",
					name,
					effectiveDataTypeName(attribute.DataType, attribute.DataTypeName),
					effectiveDataTypeName(property.DataType, property.DataTypeName),
				)
				continue
			}
			attribute.ReflectedProperty = property
			continue
		}

		attribute := tosca_v2_0.NewAttributeDefinition(property.Context)
		attribute.Name = name
		attribute.Description = cloneString(property.Description)
		attribute.DataTypeName = cloneString(property.DataTypeName)
		attribute.DataType = property.DataType
		attribute.KeySchema = cloneReflectedSchema(property.KeySchema)
		attribute.EntrySchema = cloneReflectedSchema(property.EntrySchema)
		attribute.ReflectedProperty = property
		attributes[name] = attribute
	}
}

func sameEffectiveDataType(propertyType *tosca_v2_0.DataType, attributeType *tosca_v2_0.DataType) bool {
	if propertyType == nil || attributeType == nil {
		return propertyType == attributeType
	}
	return parsing.GetCanonicalName(propertyType) == parsing.GetCanonicalName(attributeType)
}

func effectiveDataTypeName(dataType *tosca_v2_0.DataType, declaredName *string) string {
	if dataType != nil {
		return parsing.GetCanonicalName(dataType)
	}
	if declaredName != nil {
		return *declaredName
	}
	return ""
}

func cloneReflectedSchema(schema *tosca_v2_0.Schema) *tosca_v2_0.Schema {
	if schema == nil {
		return nil
	}

	clone := tosca_v2_0.NewSchema(schema.Context)
	clone.Description = cloneString(schema.Description)
	clone.DataTypeName = cloneString(schema.DataTypeName)
	clone.DataType = schema.DataType
	clone.KeySchema = cloneReflectedSchema(schema.KeySchema)
	clone.EntrySchema = cloneReflectedSchema(schema.EntrySchema)
	return clone
}

func cloneString(value *string) *string {
	if value == nil {
		return nil
	}
	clone := *value
	return &clone
}
