package timelineactions

type ActionStatus string

// CHECK (status IN ('未着手', '進行中', '完了', '保留', '中止')),
const (
	NotStarted ActionStatus = "未着手"
	InProgress ActionStatus = "進行中"
	Completed  ActionStatus = "完了"
	Pending    ActionStatus = "保留"
	Canceled   ActionStatus = "中止"
)

func (s ActionStatus) Valid() bool {
	switch s {
	case NotStarted, InProgress, Completed, Pending, Canceled:
		return true
	}
	return false
}

// String は、ActionStatusを文字列に変換します。
func (s ActionStatus) String() string {
	return string(s)
}
