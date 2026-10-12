package sort

type SortKey string

var (
	SortKeyLatest      SortKey = "latest"
	SortKeyNearest     SortKey = "nearest"
	SortKeyMostReplies SortKey = "mostReplies"
	SortKeyOldest      SortKey = "oldest"
)

func (s SortKey) IsValid() bool {
	if s == "" {
		return true
	}

	return s == SortKeyLatest || s == SortKeyNearest || s == SortKeyMostReplies || s == SortKeyOldest
}

func (s SortKey) String() string {
	if s == "" {
		return "latest"
	}

	return string(s)
}
