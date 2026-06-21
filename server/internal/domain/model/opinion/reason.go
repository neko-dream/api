package opinion

type Reason int32

//go:generate go tool enumer -type=Reason -trimprefix=Reason -transform=kebab -linecomment
const (
	// 不適切な内容
	ReasonInappropriate Reason = iota + 1
	// セッションとは関係のない発言
	ReasonIrrelevant
	// スパム
	ReasonSpam
	// プライバシー
	ReasonPrivacy
	// その他
	ReasonOther = 255
)

// StringJP
func (r Reason) StringJP() string {
	switch r {
	case ReasonInappropriate:
		return "不適切な内容"
	case ReasonIrrelevant:
		return "セッションとは関係のない発言"
	case ReasonSpam:
		return "スパム"
	case ReasonPrivacy:
		return "プライバシー"
	default:
		return "その他"
	}
}
