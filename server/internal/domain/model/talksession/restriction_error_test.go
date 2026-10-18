package talksession

import (
	"context"
	"errors"
	"testing"

	"github.com/neko-dream/api/internal/domain/messages"
)

// 共有のエラーを書き換えると、以後のすべてのリクエストに前の文言が残る
func TestUpdateRestrictionsDoesNotModifySharedError(t *testing.T) {
	ctx := context.Background()
	original := ErrInvalidRestrictionAttribute.Message

	errA := (&TalkSession{}).UpdateRestrictions(ctx, []string{"unknown.a"})
	// 不正なキーの数を変え、呼び出しごとに文言が変わるようにする
	errB := (&TalkSession{}).UpdateRestrictions(ctx, []string{"unknown.a", "unknown.b"})

	var a, b *messages.APIError
	if !errors.As(errA, &a) || !errors.As(errB, &b) {
		t.Fatalf("errors are not APIError: %v / %v", errA, errB)
	}
	if a == b {
		t.Fatal("each call must return its own error, not the shared one")
	}
	if a.Message == b.Message {
		t.Fatalf("messages must reflect each call: %q", a.Message)
	}
	if ErrInvalidRestrictionAttribute.Message != original {
		t.Fatalf("shared error was modified: %q", ErrInvalidRestrictionAttribute.Message)
	}
	if !errors.Is(errA, ErrInvalidRestrictionAttribute) {
		t.Fatal("the returned error must still match the shared error by code")
	}
}
