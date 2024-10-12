package opinion

type VoteType int

const (
	UnVoted VoteType = iota
	Agreed
	Disagreed
	Pass
)

func VoteFromString(s *string) VoteType {
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
