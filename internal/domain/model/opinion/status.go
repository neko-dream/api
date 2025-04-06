package opinion

import "errors"

type Status string

const (
	// StatusUnsolved 未解決
	StatusUnsolved Status = "unsolved"
	// StatusDeleted 削除済み
	StatusDeleted Status = "deleted"
	// StatusHold 保留
	StatusHold Status = "hold"
)

var (
	ErrInvalidStatus = errors.New("invalid status")
)

func NewStatus(status string) (Status, error) {
	switch status {
	case "unsolved":
		return StatusUnsolved, nil
	case "deleted":
		return StatusDeleted, nil
	case "hold":
		return StatusHold, nil
	default:
		return "", ErrInvalidStatus
	}
}
