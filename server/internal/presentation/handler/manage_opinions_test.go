package handler

import (
	"context"
	"errors"
	"testing"
	"time"

	opinionquery "github.com/neko-dream/api/internal/application/query/opinion"
	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/domain/model/admin"
	"github.com/neko-dream/api/internal/domain/model/opinion"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/talksession"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/neko-dream/api/internal/presentation/oas"
)

type fakeManageOpinions struct {
	got opinionquery.ManageOpinionsInput
	out *opinionquery.ManageOpinionsOutput
	err error
}

func (f *fakeManageOpinions) List(_ context.Context, in opinionquery.ManageOpinionsInput) (*opinionquery.ManageOpinionsOutput, error) {
	f.got = in
	return f.out, f.err
}

func adminCtx() context.Context {
	return admin.SetAdminSession(context.Background(), &admin.AdminClaim{Sub: shared.NewUUID[admin.AdminUser]().String()})
}

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

func TestGetTalkSessionOpinionsManage(t *testing.T) {
	talkSessionID := shared.NewUUID[talksession.TalkSession]()
	params := oas.GetTalkSessionOpinionsManageParams{TalkSessionID: talkSessionID.String(), Limit: oas.NewOptInt32(1000)}

	t.Run("管理者でなければ拒否", func(t *testing.T) {
		h := &manageHandler{manageOpinions: &fakeManageOpinions{}}
		if _, err := h.GetTalkSessionOpinionsManage(context.Background(), params); !errors.Is(err, messages.ForbiddenError) {
			t.Fatalf("err = %v, want %v", err, messages.ForbiddenError)
		}
	})

	t.Run("不正なIDは400", func(t *testing.T) {
		h := &manageHandler{manageOpinions: &fakeManageOpinions{}}
		if _, err := h.GetTalkSessionOpinionsManage(adminCtx(), oas.GetTalkSessionOpinionsManageParams{TalkSessionID: "x"}); !errors.Is(err, messages.BadRequestError) {
			t.Fatalf("err = %v, want %v", err, messages.BadRequestError)
		}
	})

	t.Run("存在しないセッションのエラーをそのまま返す", func(t *testing.T) {
		h := &manageHandler{manageOpinions: &fakeManageOpinions{err: messages.TalkSessionNotFound}}
		if _, err := h.GetTalkSessionOpinionsManage(adminCtx(), params); !errors.Is(err, messages.TalkSessionNotFound) {
			t.Fatalf("err = %v, want %v", err, messages.TalkSessionNotFound)
		}
	})

	t.Run("ページングを補正し、レスポンスに変換する", func(t *testing.T) {
		opinionID, parentID := shared.NewUUID[opinion.Opinion](), shared.NewUUID[opinion.Opinion]()
		userID := shared.NewUUID[user.User]()
		created := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
		picture := "https://example.com/a.png"
		query := &fakeManageOpinions{out: &opinionquery.ManageOpinionsOutput{
			TotalCount: 7,
			Opinions: []opinionquery.ManagedOpinion{{
				OpinionID: opinionID, ParentOpinionID: &parentID, Content: "本文", PictureURL: &picture, CreatedAt: created,
				Author:     opinionquery.ManagedOpinionAuthor{UserID: &userID, DisplayID: "id", DisplayName: "名前"},
				AgreeCount: 2, DisagreeCount: 1, PassCount: 3, ReportCount: 2, UnsolvedReportCount: 1, Deleted: true,
			}},
		}}
		h := &manageHandler{manageOpinions: query}
		res, err := h.GetTalkSessionOpinionsManage(adminCtx(), params)
		if err != nil {
			t.Fatal(err)
		}
		if query.got.TalkSessionID != talkSessionID || query.got.Limit != 100 || query.got.Offset != 0 {
			t.Fatalf("query input: %+v", query.got)
		}
		if res.TotalCount != 7 || len(res.Opinions) != 1 {
			t.Fatalf("response: %+v", res)
		}
		got := res.Opinions[0]
		if got.OpinionID != opinionID.String() || got.Content != "本文" || !got.CreatedAt.Equal(created) {
			t.Fatalf("basic fields: %+v", got)
		}
		if got.ParentOpinionID.Value != parentID.String() || got.PictureURL.Value != picture || got.ReferenceURL.Set {
			t.Fatalf("optional fields: %+v", got)
		}
		if got.Author.UserID != userID.String() || got.Author.DisplayName != "名前" || got.Author.DisplayID != "id" {
			t.Fatalf("author: %+v", got.Author)
		}
		if got.AgreeCount != 2 || got.DisagreeCount != 1 || got.PassCount != 3 || got.ReportCount != 2 || got.UnsolvedReportCount != 1 || !got.Deleted {
			t.Fatalf("counts: %+v", got)
		}
	})

	t.Run("ユーザーがいない意見は投稿者IDを空にする", func(t *testing.T) {
		query := &fakeManageOpinions{out: &opinionquery.ManageOpinionsOutput{TotalCount: 1, Opinions: []opinionquery.ManagedOpinion{{
			OpinionID: shared.NewUUID[opinion.Opinion](), Content: "本文",
		}}}}
		res, err := (&manageHandler{manageOpinions: query}).GetTalkSessionOpinionsManage(adminCtx(), params)
		if err != nil {
			t.Fatal(err)
		}
		if res.Opinions[0].Author.UserID != "" {
			t.Fatalf("author: %+v", res.Opinions[0].Author)
		}
	})
}
