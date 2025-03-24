package opinion

type Reason int32

const (
	// ReasonInappropriate 不適切な内容
	ReasonInappropriate Reason = iota + 1
	// ReasonSpam スパム
	ReasonSpam
	// ReasonHarassment 嫌がらせ
	ReasonHarassment
	// ReasonOther その他
	ReasonOther
)

func NewReason(reason int32) Reason {
	// validation
	if reason < 1 || reason > 4 {
		return ReasonOther
	}
	return Reason(reason)
}
