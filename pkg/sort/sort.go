package sort

type SortKey string

var (
	SortKeyLatest      SortKey = "latest"
	SortKeyNearest     SortKey = "nearest"
	SortKeyMostReplies SortKey = "mostReplies"
	SortKeyOldest      SortKey = "oldest"
)

func (s SortKey) IsValid() bool {
	switch s {
	case SortKeyLatest, SortKeyNearest, SortKeyMostReplies, SortKeyOldest:
		return true
	}

	return false
}

func (s SortKey) String() string {
	if s == "" {
		return "latest"
	}

	return string(s)
}
