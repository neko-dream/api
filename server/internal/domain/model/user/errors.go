package user

import "errors"

var (
	// ErrNotWithdrawn ユーザーが退会していない
	ErrNotWithdrawn = errors.New("退会していません")

	// ErrReactivationPeriodExpired 復活可能期間を過ぎている
	ErrReactivationPeriodExpired = errors.New("復活可能期間（30日）を過ぎています")

	// ErrAlreadyWithdrawn すでに退会済み
	ErrAlreadyWithdrawn = errors.New("既に退会済みです")
)
