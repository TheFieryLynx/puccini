package tosca_v1_3

import (
	"strings"

	"github.com/tliron/go-puccini/normal"
	"github.com/tliron/go-puccini/tosca/parsing"
)

const (
	operatingSystemCapabilityType           = "tosca::OperatingSystem"
	operatingSystemLowercaseConverterName   = "tosca.converter.tosca_1_3.operating_system_lowercase"
	operatingSystemLowercaseConverterScript = `
exports.convert = function(value) {
	if (typeof value !== "string") {
		throw new Error("OperatingSystem lowercase conversion requires a string");
	}
	return value.toLowerCase();
};
`
)

var operatingSystemLowercaseProperties = map[string]struct{}{
	"architecture": {},
	"type":         {},
	"distribution": {},
}

// normalizeOperatingSystemCapabilities applies TOSCA 1.3 section 5.5.12.3
// after rendering has supplied inherited definitions, defaults, and explicit
// assignments. Static string values are normalized immediately. Intrinsic
// function objects are preserved and receive a converter that runs only after
// evaluation.
func normalizeOperatingSystemCapabilities(serviceTemplate *normal.ServiceTemplate) {
	for _, nodeTemplate := range serviceTemplate.NodeTemplates {
		for _, capability := range nodeTemplate.Capabilities {
			if _, ok := capability.Types[operatingSystemCapabilityType]; !ok {
				continue
			}

			for propertyName := range operatingSystemLowercaseProperties {
				value, ok := capability.Properties[propertyName]
				if !ok {
					continue
				}

				switch value := value.(type) {
				case *normal.Primitive:
					if stringValue, ok := value.Primitive.(string); ok {
						value.Primitive = strings.ToLower(stringValue)
					}

				case *normal.FunctionCall:
					if value.ValueMeta == nil {
						value.ValueMeta = normal.NewValueMeta()
					}
					functionCall := value.FunctionCall
					value.ValueMeta.Converter = normal.NewFunctionCall(parsing.NewFunctionCall(
						operatingSystemLowercaseConverterName,
						nil,
						functionCall.URL,
						functionCall.Row,
						functionCall.Column,
						functionCall.Path,
					))
				}
			}
		}
	}
}
