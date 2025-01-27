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
		return false
	}

	switch string(s) {
	case "latest", "nearest", "mostReplies", "oldest":
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
