package handler

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/neko-dream/api/internal/domain/messages"
	model "github.com/neko-dream/api/internal/infrastructure/persistence/sqlc/generated"
	"github.com/neko-dream/api/internal/presentation/oas"
)

func TestManageOpinionPage(t *testing.T) {
	for name, tc := range map[string]struct {
		limit, offset         oas.OptInt32
		wantLimit, wantOffset int32
	}{
		"未指定":     {oas.OptInt32{}, oas.OptInt32{}, 20, 0},
		"指定":      {oas.NewOptInt32(50), oas.NewOptInt32(40), 50, 40},
		"上限を超える":  {oas.NewOptInt32(1000), oas.OptInt32{}, 100, 0},
		"0以下と負の値": {oas.NewOptInt32(0), oas.NewOptInt32(-5), 20, 0},
	} {
		t.Run(name, func(t *testing.T) {
			limit, offset := manageOpinionPage(tc.limit, tc.offset)
			if limit != tc.wantLimit || offset != tc.wantOffset {
				t.Fatalf("got (%d, %d), want (%d, %d)", limit, offset, tc.wantLimit, tc.wantOffset)
			}
		})
	}
}

func TestOpinionForManage(t *testing.T) {
	opinionID, parentID, userID := uuid.New(), uuid.New(), uuid.New()
	created := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
	got := opinionForManage(model.ListOpinionsForManageRow{
		Opinion: model.Opinion{
			OpinionID: opinionID, UserID: userID, Content: "本文", CreatedAt: created,
			ParentOpinionID: uuid.NullUUID{UUID: parentID, Valid: true},
			PictureUrl:      sql.NullString{String: "https://example.com/a.png", Valid: true},
		},
		User:       model.User{UserID: userID, DisplayID: sql.NullString{String: "id", Valid: true}, DisplayName: sql.NullString{String: "名前", Valid: true}},
		AgreeCount: 2, DisagreeCount: 1, PassCount: 3,
		ReportCount: 2, UnsolvedReportCount: 1, Deleted: true,
	})
	if got.OpinionID != opinionID.String() || got.Content != "本文" || !got.CreatedAt.Equal(created) {
		t.Fatalf("basic fields: %+v", got)
	}
	if got.ParentOpinionID.Value != parentID.String() || !got.ParentOpinionID.Set {
		t.Fatalf("parent: %+v", got.ParentOpinionID)
	}
	if got.PictureURL.Value != "https://example.com/a.png" || got.ReferenceURL.Set {
		t.Fatalf("urls: %+v %+v", got.PictureURL, got.ReferenceURL)
	}
	if got.Author.UserID != userID.String() || got.Author.DisplayName != "名前" || got.Author.DisplayID != "id" {
		t.Fatalf("author: %+v", got.Author)
	}
	if got.AgreeCount != 2 || got.DisagreeCount != 1 || got.PassCount != 3 || got.ReportCount != 2 || got.UnsolvedReportCount != 1 || !got.Deleted {
		t.Fatalf("counts: %+v", got)
	}
}

func TestGetTalkSessionOpinionsManageRequiresAdmin(t *testing.T) {
	h := &manageHandler{}
	if _, err := h.GetTalkSessionOpinionsManage(context.Background(), oas.GetTalkSessionOpinionsManageParams{TalkSessionID: uuid.NewString()}); !errors.Is(err, messages.ForbiddenError) {
		t.Fatalf("err = %v, want %v", err, messages.ForbiddenError)
	}
}
