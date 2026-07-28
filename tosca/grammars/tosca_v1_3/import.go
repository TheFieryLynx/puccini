package tosca_v1_3

import (
	"strings"

	"github.com/tliron/go-ard"
	"github.com/tliron/go-puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/go-puccini/tosca/parsing"
)

//
// Import
//
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.8
//

// ([parsing.Reader] signature)
func ReadImport(context *parsing.Context) parsing.EntityPtr {
	context.SetReadTag("URL", "file")
	context.SetReadTag("Namespace", "namespace_prefix")
	context.SetReadTag("NamespaceURI", "namespace_uri")

	self := tosca_v2_0.NewImport(context)

	if context.Is(ard.TypeMap) {
		if context.HasQuirk(parsing.QuirkImportsSequencedList) {
			map_ := context.Data.(ard.Map)
			if len(map_) == 1 {
				for _, data := range map_ {
					if data_, ok := data.(ard.Map); ok {
						context.Data = data_
					}
					break
				}
			}
		}

		// Long notation
		context.ValidateUnsupportedFields(context.ReadFields(self))
		if _, ok := context.GetFieldChild("file"); !ok {
			context.FieldChild("file", nil).ReportKeynameMissing()
		} else if self.URL != nil && strings.TrimSpace(*self.URL) == "" {
			context.FieldChild("file", *self.URL).ReportValueInvalid("import file", "must not be empty")
		}
	} else if context.ValidateType(ard.TypeMap, ard.TypeString) {
		// Short notation
		self.URL = context.FieldChild("file", context.Data).ReadString()
		if self.URL != nil && strings.TrimSpace(*self.URL) == "" {
			context.FieldChild("file", *self.URL).ReportValueInvalid("import file", "must not be empty")
		}
	}

	return self
}
