package db

import (
	"context"
	"testing"
)

func TestAfterCommit_RunsImmediatelyOutsideTransaction(t *testing.T) {
	called := false
	AfterCommit(context.Background(), func() { called = true })
	if !called {
		t.Fatal("トランザクション外ではすぐに実行される")
	}
}

func TestAfterCommit_DefersInsideTransaction(t *testing.T) {
	var hooks []func()
	ctx := context.WithValue(context.Background(), afterCommitKey{}, &hooks)

	called := false
	AfterCommit(ctx, func() { called = true })
	if called {
		t.Fatal("トランザクション中はコミットまで実行されない")
	}
	if len(hooks) != 1 {
		t.Fatalf("hooks = %d, want 1", len(hooks))
	}
	hooks[0]()
	if !called {
		t.Fatal("登録した関数が実行される")
	}
}
