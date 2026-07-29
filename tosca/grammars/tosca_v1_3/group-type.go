package tosca_v1_3

import (
	"fmt"

	"github.com/tliron/go-ard"
	"github.com/tliron/go-puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/go-puccini/tosca/parsing"
)

func validateGroupTypeMembers(entityPtr parsing.EntityPtr) {
	groupType := entityPtr.(*tosca_v2_0.GroupType)

	data, ok := groupType.Context.Data.(ard.Map)
	if !ok {
		return
	}
	memberData, declared := data["members"]
	if !declared || (groupType.MemberNodeTypeNames == nil) {
		return
	}

	memberNames := *groupType.MemberNodeTypeNames
	memberTypes := groupType.MemberNodeTypes
	if (len(memberTypes) < 2) || (len(memberTypes) != len(memberNames)) {
		return
	}

	baseType := memberTypes[0]
	context := groupType.Context.FieldChild("members", memberData)
	for index := 1; index < len(memberTypes); index++ {
		if !groupType.Context.Hierarchy.IsInSameHierarchy(baseType, memberTypes[index]) {
			context.ListChild(index, memberNames[index]).ReportValueMalformed(
				"member type",
				fmt.Sprintf(
					"must derive from the same type hierarchy as %q",
					parsing.GetCanonicalName(baseType),
				),
			)
		}
	}
}
