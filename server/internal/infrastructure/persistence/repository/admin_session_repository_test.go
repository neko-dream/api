package repository

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/neko-dream/api/internal/domain/model/admin"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
	"os"
	"sync"
	"testing"
	"time"
)

func TestAdminRefreshRotationIsAtomic(t *testing.T) {
	url := os.Getenv("ADMIN_TEST_DATABASE_URL")
	if url == "" {
		t.Skip("ADMIN_TEST_DATABASE_URL is required")
	}
	conn, err := sql.Open("pgx", url)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { conn.Close() })
	repo := NewAdminSessionRepository(db.NewDBManager(conn))
	ctx := context.Background()
	id := shared.NewUUID[admin.AdminSession]()
	old := id.String()
	if err := repo.Create(ctx, admin.AdminSession{AdminSessionID: id, AdminUserID: shared.NewUUID[admin.AdminUser](), TokenHash: old, ExpiresAt: time.Now().Add(time.Hour)}); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_, err := conn.Exec("DELETE FROM admin_sessions WHERE admin_session_id=$1", id.UUID())
		if err != nil {
			t.Error(err)
		}
	})
	var wg sync.WaitGroup
	start := make(chan struct{})
	successes := make(chan bool, 2)
	failures := make(chan error, 2)
	for _, hash := range []string{old + "-a", old + "-b"} {
		wg.Add(1)
		go func(hash string) {
			defer wg.Done()
			<-start
			ok, err := repo.Rotate(ctx, id, old, hash, time.Now().Add(time.Hour), time.Now())
			successes <- ok
			failures <- err
		}(hash)
	}
	close(start)
	wg.Wait()
	close(successes)
	close(failures)
	count := 0
	for ok := range successes {
		if ok {
			count++
		}
	}
	for err := range failures {
		if err != nil {
			t.Fatal(err)
		}
	}
	if count != 1 {
		t.Fatalf("successful rotations=%d, want 1", count)
	}
	// 直前のトークンでも見つかるが、現在のトークンではなく、猶予期間を過ぎると使えない
	oldSession, err := repo.FindByTokenHash(ctx, old)
	if err != nil || oldSession == nil || oldSession.TokenHash == old {
		t.Fatalf("previous token lookup: %+v %v", oldSession, err)
	}
	if !oldSession.IsPreviousTokenInGracePeriod(old, time.Now()) {
		t.Fatal("previous token must be accepted within the grace period")
	}
	if oldSession.IsPreviousTokenInGracePeriod(old, time.Now().Add(admin.RotationGracePeriod+time.Second)) {
		t.Fatal("previous token must be rejected after the grace period")
	}
	// 直前のトークンでログアウトしても、セッションごと消える
	if err := repo.DeleteByTokenHash(ctx, old); err != nil {
		t.Fatal(err)
	}
	if current, err := repo.FindByTokenHash(ctx, oldSession.TokenHash); err != nil || current != nil {
		t.Fatalf("session remains after logout with previous token: %v", err)
	}
}
