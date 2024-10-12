package vote

type VoteStatus int

const (
	UnVoted VoteStatus = iota
	Agreed
	Disagreed
	Pass
)
