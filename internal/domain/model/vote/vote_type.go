package vote

type VoteType int

const (
	UnVoted VoteType = iota
	Agreed
	Disagreed
	Pass
)

func (v VoteType) Int() int {
	return int(v)
}

func (v VoteType) String() string {
	switch v {
	case Agreed:
		return "agree"
	case Disagreed:
		return "disagree"
	case Pass:
		return "pass"
	default:
		return "unvote"
	}
}

func VoteTypeFromInt(i int) VoteType {
	switch i {
	case 1:
		return Agreed
	case 2:
		return Disagreed
	case 3:
		return Pass
	default:
		return UnVoted
	}
}

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
