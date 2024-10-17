package analysis

var GroupNameMap = []string{
	// ランダムな動物の名前
	"ねこ",
	"いぬ",
	"ひつじ",
	"うさぎ",
	"ぞう",
	"ねずみ",
	"へび",
	"とら",
	"さる",
	"ひよこ",
}

type GroupName string

func NewGroupName(num int) GroupName {
	return GroupName(GroupNameMap[num])
}
