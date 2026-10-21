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
	if changed == base || changed.Message != "差し替えた文言" || changed.Code != base.Code || changed.StatusCode != base.StatusCode {
		t.Fatalf("changed = %+v", changed)
	}
	if !errors.Is(changed, base) || !errors.Is(changed.WithMessage("さらに差し替え"), base) {
		t.Fatal("copies must still match the original")
	}
}

func TestWithCauseKeepsReceiverAndCause(t *testing.T) {
	base := &APIError{StatusCode: 500, Code: "TEST-0002", Message: "失敗しました"}
	cause := errors.New("db down")
	wrapped := base.WithCause(cause)

	if wrapped == base || base.Cause != nil {
		t.Fatal("WithCause must return a copy and leave the original unchanged")
	}
	if !errors.Is(wrapped, base) || !errors.Is(wrapped, cause) {
		t.Fatal("wrapped error must match both the API error and its cause")
	}
	if wrapped.Error() != "失敗しました: db down" {
		t.Fatalf("Error() = %q", wrapped.Error())
	}

	// cause に APIError が含まれていても、呼び出し側のエラー（コード・文言）を残す
	inner := &APIError{StatusCode: 403, Code: "TEST-0003", Message: "権限がありません"}
	got := base.WithMessage("画像の保存に失敗しました").WithCause(fmt.Errorf("upload: %w", inner))
	if got.Code != base.Code || got.StatusCode != base.StatusCode || got.Message != "画像の保存に失敗しました" {
		t.Fatalf("receiver was dropped: %+v", got)
	}
	if !errors.Is(got, base) || !errors.Is(got, inner) {
		t.Fatal("must match both the receiver and the inner API error")
	}
	if got == inner {
		t.Fatal("must not return the shared inner error")
	}

	// cause が nil でも、共有のエラーそのものではなく複製を返す
	if nilCause := base.WithCause(nil); nilCause == base || !errors.Is(nilCause, base) {
		t.Fatal("WithCause(nil) must return a copy that still matches the original")
	}
}

func TestErrorDoesNotRepeatTheSameText(t *testing.T) {
	if got := (&APIError{Code: "TEST-0004", Message: "文言"}).Error(); got != "文言" {
		t.Fatalf("Error() = %q", got)
	}
	same := (&APIError{Code: "TEST-0005", Message: "同じ文言"}).WithCause(errors.New("同じ文言"))
	if got := same.Error(); got != "同じ文言" {
		t.Fatalf("Error() = %q", got)
	}
}

// 比較はコードではなく、どの定義済みのエラーから複製したかで行う（コードが重複していても誤判定しない）
func TestIsComparesOrigin(t *testing.T) {
	a := &APIError{Code: "DUP-001", Message: "a"}
	b := &APIError{Code: "DUP-001", Message: "b"}
	if errors.Is(a, b) || errors.Is(a.WithMessage("x"), b) {
		t.Fatal("different errors with the same code must not match")
	}
	if errors.Is(&APIError{}, &APIError{}) {
		t.Fatal("ad-hoc errors with an empty code must not match each other")
	}
	if errors.Is(a, errors.New("a")) {
		t.Fatal("non-APIError must not match")
	}
	var nilTarget *APIError
	if errors.Is(a, nilTarget) {
		t.Fatal("nil target must not match")
	}
	if !errors.Is(a, a) {
		t.Fatal("an error must match itself")
	}
}

func TestDuplicatedCodesDoNotMatch(t *testing.T) {
	// OrganizationAlreadyExists と OpinionContentBadLength はどちらもコードが OPINION-001
	if errors.Is(OrganizationAlreadyExists, OpinionContentBadLength) {
		t.Fatal("errors sharing a code must not be treated as the same error")
	}
}
