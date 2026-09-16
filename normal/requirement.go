package normal

import (
	"encoding/json"

	"github.com/fxamacker/cbor/v2"
	"github.com/tliron/go-ard"
)

//
// Requirement
//

type Requirement struct {
	SourceNodeTemplate *NodeTemplate

	Name                           string
	CapabilityTypeName             *string
	CapabilityName                 *string
	NodeTypeName                   *string
	NodeTemplate                   *NodeTemplate
	NodeTemplatePropertyValidation FunctionCallMap
	CapabilityPropertyValidation   FunctionCallMapMap
	Relationship                   *Relationship
	Directives                     []string
	Optional                       bool
	Occurrences                    *OccurrenceRange
	Location                       *Location
}

// OccurrenceRange retains a declarative multiplicity interval. A nil Upper
// means unbounded. Absence of the whole field preserves the existing model
// used by grammars that expand an exact count into individual requirements.
type OccurrenceRange struct {
	Lower uint64  `json:"lower" yaml:"lower"`
	Upper *uint64 `json:"upper" yaml:"upper"`
}

func (self *NodeTemplate) NewRequirement(name string, location *Location) *Requirement {
	requirement := &Requirement{
		SourceNodeTemplate:             self,
		Name:                           name,
		NodeTemplatePropertyValidation: make(FunctionCallMap),
		CapabilityPropertyValidation:   make(FunctionCallMapMap),
		Location:                       location,
	}
	self.Requirements = append(self.Requirements, requirement)
	return requirement
}

type MarshalableRequirement struct {
	Name                           string             `json:"name" yaml:"name"`
	CapabilityTypeName             string             `json:"capabilityTypeName" yaml:"capabilityTypeName"`
	CapabilityName                 string             `json:"capabilityName" yaml:"capabilityName"`
	NodeTypeName                   string             `json:"nodeTypeName" yaml:"nodeTypeName"`
	NodeTemplateName               string             `json:"nodeTemplateName" yaml:"nodeTemplateName"`
	NodeTemplatePropertyValidators FunctionCallMap    `json:"nodeTemplatePropertyValidators" yaml:"nodeTemplatePropertyValidators"`
	CapabilityPropertyValidators   FunctionCallMapMap `json:"capabilityPropertyValidators" yaml:"capabilityPropertyValidators"`
	Relationship                   *Relationship      `json:"relationship" yaml:"relationship"`
	Directives                     []string           `json:"directives" yaml:"directives"`
	Optional                       bool               `json:"optional" yaml:"optional"`
	Occurrences                    *OccurrenceRange   `json:"occurrences,omitempty" yaml:"occurrences,omitempty"`
	Location                       *Location          `json:"location" yaml:"location"`
}

func (self *Requirement) Marshalable() any {
	var capabilityTypeName string
	var capabilityName string
	var nodeTypeName string
	var nodeTemplateName string
	if self.CapabilityTypeName != nil {
		capabilityTypeName = *self.CapabilityTypeName
	}
	if self.CapabilityName != nil {
		capabilityName = *self.CapabilityName
	}
	if self.NodeTypeName != nil {
		nodeTypeName = *self.NodeTypeName
	}
	if self.NodeTemplate != nil {
		nodeTemplateName = self.NodeTemplate.Name
	}

	return &MarshalableRequirement{
		Name:                           self.Name,
		CapabilityTypeName:             capabilityTypeName,
		CapabilityName:                 capabilityName,
		NodeTypeName:                   nodeTypeName,
		NodeTemplateName:               nodeTemplateName,
		NodeTemplatePropertyValidators: self.NodeTemplatePropertyValidation,
		CapabilityPropertyValidators:   self.CapabilityPropertyValidation,
		Relationship:                   self.Relationship,
		Directives:                     self.Directives,
		Optional:                       self.Optional,
		Occurrences:                    self.Occurrences,
		Location:                       self.Location,
	}
}

// ([json.Marshaler] interface)
func (self *Requirement) MarshalJSON() ([]byte, error) {
	return json.Marshal(self.Marshalable())
}

// ([yaml.Marshaler] interface)
func (self *Requirement) MarshalYAML() (any, error) {
	return self.Marshalable(), nil
}

// ([cbor.Marshaler] interface)
func (self *Requirement) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal(self.Marshalable())
}

// ([msgpack.Marshaler] interface)
func (self *Requirement) MarshalMsgpack() ([]byte, error) {
	return ard.MarshalMessagePack(self.Marshalable())
}

// ([ard.ToARD] interface)
func (self *Requirement) ToARD(reflector *ard.Reflector) (any, error) {
	return reflector.Unpack(self.Marshalable())
}

//
// Requirements
//

type Requirements []*Requirement
