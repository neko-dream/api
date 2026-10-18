package opinion_query_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	opinionquery "github.com/neko-dream/api/internal/application/query/opinion"
	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/domain/model/clock"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/talksession"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/neko-dream/api/internal/infrastructure/di"
	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
	opinion_query "github.com/neko-dream/api/internal/infrastructure/persistence/query/opinion"
	"github.com/neko-dream/api/internal/infrastructure/persistence/repository"
	model "github.com/neko-dream/api/internal/infrastructure/persistence/sqlc/generated"
	"github.com/neko-dream/api/internal/test/txtest"
)

func TestManageOpinions(t *testing.T) {
	container := di.BuildContainer()
	dbManager := di.Invoke[*db.DBManager](container)
	q := opinion_query.NewManageOpinionsQuery(dbManager)
	type TestData struct{}

	talkSessionID := shared.NewUUID[talksession.TalkSession]()
	users := make([]uuid.UUID, 4)
	for i := range users {
		users[i] = uuid.New()
	}
	parentID, replyID, mixedID, hiddenID := uuid.New(), uuid.New(), uuid.New(), uuid.New()
	base := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)

	setup := func(ctx context.Context, _ *TestData) error {
		queries := dbManager.GetQueries(ctx)
		for i, id := range users {
			if err := queries.CreateUser(ctx, model.CreateUserParams{
				UserID: id, CreatedAt: base,
				DisplayID:   sql.NullString{String: fmt.Sprintf("manage-op-%d-%s", i, id.String()[:8]), Valid: true},
				DisplayName: sql.NullString{String: fmt.Sprintf("ユーザー%d", i), Valid: true},
			}); err != nil {
				return err
			}
		}
		if err := repository.NewTalkSessionRepository(dbManager, repository.NewEventStoreMock()).Create(ctx, talksession.NewTalkSession(
			talkSessionID, "テーマ", nil, nil, shared.UUID[user.User](users[0]),
			clock.Now(ctx), clock.Now(ctx).Add(24*time.Hour), nil, nil, nil, false, nil, nil, nil,
		)); err != nil {
			return err
		}
		for i, op := range []model.CreateOpinionParams{
			{OpinionID: parentID, UserID: users[0], Content: "親の意見"},
			{OpinionID: replyID, UserID: users[1], Content: "返信", ParentOpinionID: uuid.NullUUID{UUID: parentID, Valid: true}},
			{OpinionID: mixedID, UserID: users[2], Content: "削除後にまた通報された意見"},
			{OpinionID: hiddenID, UserID: users[3], Content: "削除された意見"},
		} {
			op.TalkSessionID = talkSessionID.UUID()
			op.CreatedAt = base.Add(time.Duration(i) * time.Hour)
			if err := queries.CreateOpinion(ctx, op); err != nil {
				return err
			}
		}
		// 親の意見に 賛成2・反対1・保留1
		for i, voteType := range []int16{1, 1, 2, 3} {
			if err := queries.CreateVote(ctx, model.CreateVoteParams{
				VoteID: uuid.New(), OpinionID: parentID, TalkSessionID: talkSessionID.UUID(),
				UserID: users[i], VoteType: voteType, CreatedAt: base,
			}); err != nil {
				return err
			}
		}
		for _, report := range []model.CreateReportParams{
			{OpinionID: hiddenID, ReporterID: users[0], Status: "deleted"},
			{OpinionID: hiddenID, ReporterID: users[1], Status: "deleted"},
			// 削除済みの通報と未対応の通報が混ざると、一般の画面では表示される
			{OpinionID: mixedID, ReporterID: users[0], Status: "deleted"},
			{OpinionID: mixedID, ReporterID: users[1], Status: "unsolved"},
			{OpinionID: parentID, ReporterID: users[3], Status: "hold"},
		} {
			report.OpinionReportID = uuid.New()
			report.TalkSessionID = talkSessionID.UUID()
			report.Reason = 1
			if err := queries.CreateReport(ctx, report); err != nil {
				return err
			}
		}
		return nil
	}

	testCases := []*txtest.TransactionalTestCase[TestData]{
		{
			Name:    "削除済み・返信も含め、投票と通報を意見ごとに集計して新しい順に返す",
			SetupFn: setup,
			TestFn: func(ctx context.Context, _ *TestData) error {
				out, err := q.List(ctx, opinionquery.ManageOpinionsInput{TalkSessionID: talkSessionID, Limit: 10})
				if err != nil {
					return err
				}
				if out.TotalCount != 4 || len(out.Opinions) != 4 {
					return fmt.Errorf("total=%d len=%d, want 4", out.TotalCount, len(out.Opinions))
				}
				order := []uuid.UUID{hiddenID, mixedID, replyID, parentID}
				for i, want := range order {
					if out.Opinions[i].OpinionID.UUID() != want {
						return errors.New("opinions are not ordered newest first")
					}
				}
				hidden, mixed, reply, parent := out.Opinions[0], out.Opinions[1], out.Opinions[2], out.Opinions[3]
				if !hidden.Hidden || hidden.ReportCount != 2 || hidden.UnsolvedReportCount != 0 {
					return fmt.Errorf("hidden opinion: %+v", hidden)
				}
				if mixed.Hidden || mixed.ReportCount != 2 || mixed.UnsolvedReportCount != 1 {
					return fmt.Errorf("mixed reports must be visible like the public list: %+v", mixed)
				}
				if reply.ParentOpinionID == nil || reply.ParentOpinionID.UUID() != parentID || reply.Hidden {
					return fmt.Errorf("reply: %+v", reply)
				}
				if parent.AgreeCount != 2 || parent.DisagreeCount != 1 || parent.PassCount != 1 || parent.ReportCount != 1 || parent.Hidden {
					return fmt.Errorf("parent counts: %+v", parent)
				}
				if parent.Author.UserID.UUID() != users[0] || parent.Author.DisplayName != "ユーザー0" {
					return fmt.Errorf("author: %+v", parent.Author)
				}

				page, err := q.List(ctx, opinionquery.ManageOpinionsInput{TalkSessionID: talkSessionID, Limit: 1, Offset: 2})
				if err != nil {
					return err
				}
				if page.TotalCount != 4 || len(page.Opinions) != 1 || page.Opinions[0].OpinionID.UUID() != replyID {
					return errors.New("offset/limit not applied")
				}
				return nil
			},
		},
		{
			Name: "存在しないセッションは TalkSessionNotFound",
			TestFn: func(ctx context.Context, _ *TestData) error {
				_, err := q.List(ctx, opinionquery.ManageOpinionsInput{TalkSessionID: shared.NewUUID[talksession.TalkSession](), Limit: 10})
				if !errors.Is(err, messages.TalkSessionNotFound) {
					return fmt.Errorf("err = %v, want %v", err, messages.TalkSessionNotFound)
				}
				return nil
			},
		},
	}

	txtest.RunTransactionalTests(t, dbManager, &TestData{}, testCases)
}
