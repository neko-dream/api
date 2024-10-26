package analysis

type GroupID int

const (
	GroupIDUnknown    GroupID = iota
	GroupIDStrawberry GroupID = iota + 1
	GroupIDLemon
	GroupIDGrape
	GroupIDWatermelon
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
	case 1:
		return GroupIDStrawberry
	case 2:
		return GroupIDLemon
	case 3:
		return GroupIDGrape
	case 4:
		return GroupIDWatermelon
	default:
		return GroupIDUnknown
	}
}

func (g GroupID) String() string {
	return GroupIDNameMap[g]
}
