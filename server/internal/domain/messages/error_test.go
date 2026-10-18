package messages

import (
	"errors"
	"fmt"
	"testing"
)

func TestWithMessageReturnsCopy(t *testing.T) {
	base := &APIError{StatusCode: 400, Code: "TEST-0001", Message: "元の文言"}
	changed := base.WithMessage("差し替えた文言")

	if base.Message != "元の文言" {
		t.Fatalf("original was modified: %q", base.Message)
	}
	if changed.Message != "差し替えた文言" || changed.Code != base.Code || changed.StatusCode != base.StatusCode {
		t.Fatalf("changed = %+v", changed)
	}
	if !errors.Is(changed, base) {
		t.Fatal("a copy with a different message must still match the original by code")
	}
}

func TestWithCauseKeepsOriginalError(t *testing.T) {
	base := &APIError{StatusCode: 500, Code: "TEST-0002", Message: "失敗しました"}
	cause := errors.New("db down")
	wrapped := base.WithCause(cause)

	if !errors.Is(wrapped, base) || !errors.Is(wrapped, cause) {
		t.Fatal("wrapped error must match both the API error and its cause")
	}
	if base.Cause != nil {
		t.Fatal("original was modified")
	}
	if wrapped.Error() != "失敗しました: db down" {
		t.Fatalf("Error() = %q", wrapped.Error())
	}
	if base.WithCause(nil) != base {
		t.Fatal("WithCause(nil) must return the same error")
	}

	// 既に APIError を包んでいる場合は、そちらを優先する（二重に包まない）
	inner := &APIError{StatusCode: 404, Code: "TEST-0003", Message: "見つかりません"}
	if got := base.WithCause(fmt.Errorf("context: %w", inner)); got != inner {
		t.Fatalf("WithCause over an existing APIError = %+v, want the inner error", got)
	}
}

func TestErrorWithoutCauseIsMessage(t *testing.T) {
	if got := (&APIError{Code: "TEST-0004", Message: "文言"}).Error(); got != "文言" {
		t.Fatalf("Error() = %q", got)
	}
}

func TestIsComparesCode(t *testing.T) {
	a := &APIError{Code: "TEST-0005", Message: "a"}
	if errors.Is(a, &APIError{Code: "TEST-0006", Message: "a"}) {
		t.Fatal("different codes must not match")
	}
	if errors.Is(a, errors.New("a")) {
		t.Fatal("non-APIError must not match")
	}
}
