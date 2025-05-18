package vote

import (
	"errors"

	"github.com/samber/lo"
)

type VoteType int

const (
	UnVoted VoteType = iota
	Agree
	Disagree
	Pass
)

func (v VoteType) Int() int {
	return int(v)
}

func (v VoteType) String() string {
	switch v {
	case Agree:
		return "agree"
	case Disagree:
		return "disagree"
	case Pass:
		return "pass"
	default:
		return ""
	}
}

func VoteTypeFromInt(i int) VoteType {
	switch i {
	case 1:
		return Agree
	case 2:
		return Disagree
	case 3:
		return Pass
	default:
		return UnVoted
	}
}

func VoteFromString(s *string) (*VoteType, error) {
	if s == nil {
		return nil, nil
	}
	switch *s {
	case "agree":
		return lo.ToPtr(Agree), nil
	case "disagree":
		return lo.ToPtr(Disagree), nil
	case "pass":
		return lo.ToPtr(Pass), nil
	default:
		return nil, errors.New("invalid vote type")
	}
}
