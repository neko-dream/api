package analysis

type GroupID int

const (
	GroupIDUnknown    GroupID = -1
	GroupIDStrawberry GroupID = 0
	GroupIDLemon      GroupID = 1
	GroupIDGrape      GroupID = 2
	GroupIDWatermelon GroupID = 3
)

var (
	GroupIDNameMap = map[GroupID]string{
		0: "いちご",
		1: "レモン",
		2: "ぶどう",
		3: "すいか",
	}
)

func NewGroupIDFromInt(i int) GroupID {
	switch i {
	case 0:
		return GroupIDStrawberry
	case 1:
		return GroupIDLemon
	case 2:
		return GroupIDGrape
	case 3:
		return GroupIDWatermelon
	default:
		return GroupIDUnknown
	}
}

func (g GroupID) String() string {
	return GroupIDNameMap[g]
}
