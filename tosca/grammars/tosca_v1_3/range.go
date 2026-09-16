package tosca_v1_3

import (
	"github.com/tliron/go-ard"
	"github.com/tliron/go-puccini/tosca/parsing"
)

// Range is the signed numeric value in sections 3.3.3.1 and 3.3.3.2.
// A nil upper bound represents UNBOUNDED without consuming an integer value.
// Cardinality ranges use a separate, non-negative reader.
type Range struct {
	Lower int64  `json:"lower" yaml:"lower"`
	Upper *int64 `json:"upper" yaml:"upper"`
}

func ReadRange(context *parsing.Context) parsing.EntityPtr {
	self := &Range{}
	if !context.ValidateType(ard.TypeList) {
		return self
	}
	list := context.Data.(ard.List)
	if len(list) != 2 {
		context.ReportValueMalformed("range", "list length not 2")
		return self
	}
	lower := context.ListChild(0, list[0]).ReadInteger()
	if lower != nil {
		self.Lower = *lower
	}
	upperContext := context.ListChild(1, list[1])
	if upperContext.Is(ard.TypeString) {
		if list[1] != "UNBOUNDED" {
			upperContext.ReportValueMalformed("range", "upper bound string not UNBOUNDED")
		}
	} else {
		self.Upper = upperContext.ReadInteger()
	}
	if lower != nil && self.Upper != nil && *self.Upper < self.Lower {
		context.ReportValueMalformed("range", "upper bound lower than lower bound")
	}
	return self
}

func (self *Range) InRange(number int64) bool {
	return number >= self.Lower && (self.Upper == nil || number <= *self.Upper)
}

func (self *Range) Within(parent *Range) bool {
	return self.Lower >= parent.Lower && (parent.Upper == nil || (self.Upper != nil && *self.Upper <= *parent.Upper))
}

// Operands of in_range on a range constrain its two endpoints; they are not
// themselves range values. This follows the same signed domain and ordering.
func rangeConstraintBounds(operands ard.List) (*Range, bool) {
	if len(operands) != 2 {
		return nil, false
	}
	lower, ok := integerOperand(operands[0])
	if !ok {
		return nil, false
	}
	bounds := &Range{Lower: int64(lower)}
	if operands[1] == "UNBOUNDED" {
		return bounds, true
	}
	upper, ok := integerOperand(operands[1])
	if !ok || upper < lower {
		return nil, false
	}
	upperBound := int64(upper)
	bounds.Upper = &upperBound
	return bounds, true
}
