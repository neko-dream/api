package vote

type VoteStatus int

const (
	UnVoted VoteStatus = iota
	Agreed
	Disagreed
	Pass
)

func FromString(s *string) VoteStatus {
	if s == nil {
		return UnVoted
	}
	switch *s {
	case "agreed":
		return Agreed
	case "disagreed":
		return Disagreed
	case "pass":
		return Pass
	default:
		return UnVoted
	}
}
