package opinion

import "errors"

type Status string

const (
	// StatusUnconfirmed 未確認
	StatusUnconfirmed Status = "unconfirmed"
	// StatusConfirmed 確認済み
	StatusConfirmed Status = "confirmed"
	// StatusIgnored 無視
	StatusIgnored Status = "ignored"
)

var (
	ErrInvalidStatus = errors.New("invalid status")
)

func NewStatus(status string) (Status, error) {
	switch status {
	case "unconfirmed":
		return StatusUnconfirmed, nil
	case "confirmed":
		return StatusConfirmed, nil
	case "ignored":
		return StatusIgnored, nil
	default:
		return "", ErrInvalidStatus
	}
}
