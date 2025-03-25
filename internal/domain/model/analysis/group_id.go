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
		0:  "A",
		1:  "B",
		2:  "C",
		3:  "D",
		4:  "E",
		5:  "F",
		6:  "G",
		7:  "H",
		8:  "I",
		9:  "J",
		10: "K",
		11: "L",
		12: "M",
		13: "N",
		14: "O",
		15: "P",
		16: "Q",
		17: "R",
		18: "S",
	}
)

func NewGroupIDFromInt(i int) GroupID {
	return GroupID(i)
}

func (g GroupID) String() string {
	return GroupIDNameMap[g]
}
