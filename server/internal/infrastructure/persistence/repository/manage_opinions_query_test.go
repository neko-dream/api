package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/neko-dream/api/internal/domain/model/clock"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/talksession"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/neko-dream/api/internal/infrastructure/di"
	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
	"github.com/neko-dream/api/internal/infrastructure/persistence/repository"
	model "github.com/neko-dream/api/internal/infrastructure/persistence/sqlc/generated"
	"github.com/neko-dream/api/internal/test/txtest"
)

func TestListOpinionsForManage(t *testing.T) {
	container := di.BuildContainer()
	dbManager := di.Invoke[*db.DBManager](container)
	type TestData struct{}

	talkSessionID := shared.NewUUID[talksession.TalkSession]()
	users := make([]uuid.UUID, 4)
	for i := range users {
		users[i] = uuid.New()
	}
	parentID, replyID, deletedID := uuid.New(), uuid.New(), uuid.New()
	base := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)

	testCases := []*txtest.TransactionalTestCase[TestData]{
		{
			Name: "削除済み・返信も含め、投票と通報を意見ごとに集計して新しい順に返す",
			SetupFn: func(ctx context.Context, _ *TestData) error {
				q := dbManager.GetQueries(ctx)
				for i, id := range users {
					if err := q.CreateUser(ctx, model.CreateUserParams{
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
					{OpinionID: deletedID, UserID: users[2], Content: "削除された意見"},
				} {
					op.TalkSessionID = talkSessionID.UUID()
					op.CreatedAt = base.Add(time.Duration(i) * time.Hour)
					if err := q.CreateOpinion(ctx, op); err != nil {
						return err
					}
				}
				// 親の意見に 賛成2・反対1・保留1
				for i, voteType := range []int16{1, 1, 2, 3} {
					if err := q.CreateVote(ctx, model.CreateVoteParams{
						VoteID: uuid.New(), OpinionID: parentID, TalkSessionID: talkSessionID.UUID(),
						UserID: users[i], VoteType: voteType, CreatedAt: base,
					}); err != nil {
						return err
					}
				}
				for _, report := range []model.CreateReportParams{
					{OpinionID: deletedID, ReporterID: users[0], Status: "deleted"},
					{OpinionID: deletedID, ReporterID: users[1], Status: "unsolved"},
					{OpinionID: parentID, ReporterID: users[3], Status: "hold"},
				} {
					report.OpinionReportID = uuid.New()
					report.TalkSessionID = talkSessionID.UUID()
					report.Reason = 1
					if err := q.CreateReport(ctx, report); err != nil {
						return err
					}
				}
				return nil
			},
			TestFn: func(ctx context.Context, _ *TestData) error {
				q := dbManager.GetQueries(ctx)
				total, err := q.CountOpinionsForManage(ctx, talkSessionID.UUID())
				if err != nil {
					return err
				}
				if total != 3 {
					return fmt.Errorf("total = %d, want 3", total)
				}
				rows, err := q.ListOpinionsForManage(ctx, model.ListOpinionsForManageParams{TalkSessionID: talkSessionID.UUID(), RowLimit: 10})
				if err != nil {
					return err
				}
				if len(rows) != 3 || rows[0].Opinion.OpinionID != deletedID || rows[1].Opinion.OpinionID != replyID || rows[2].Opinion.OpinionID != parentID {
					return errors.New("opinions are not ordered newest first")
				}
				deleted, reply, parent := rows[0], rows[1], rows[2]
				if !deleted.Deleted || deleted.ReportCount != 2 || deleted.UnsolvedReportCount != 1 {
					return fmt.Errorf("deleted opinion: %+v", deleted)
				}
				if !reply.Opinion.ParentOpinionID.Valid || reply.Opinion.ParentOpinionID.UUID != parentID || reply.Deleted {
					return fmt.Errorf("reply: %+v", reply)
				}
				if parent.AgreeCount != 2 || parent.DisagreeCount != 1 || parent.PassCount != 1 || parent.ReportCount != 1 || parent.UnsolvedReportCount != 0 || parent.Deleted {
					return fmt.Errorf("parent counts: %+v", parent)
				}
				if parent.User.DisplayName.String != "ユーザー0" {
					return fmt.Errorf("author: %+v", parent.User)
				}

				page, err := q.ListOpinionsForManage(ctx, model.ListOpinionsForManageParams{TalkSessionID: talkSessionID.UUID(), RowLimit: 1, RowOffset: 1})
				if err != nil {
					return err
				}
				if len(page) != 1 || page[0].Opinion.OpinionID != replyID {
					return errors.New("offset/limit not applied")
				}
				return nil
			},
		},
	}

	txtest.RunTransactionalTests(t, dbManager, &TestData{}, testCases)
}
