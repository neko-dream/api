package opinion_command

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/opinion"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/talksession"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockOpinionRepository struct {
	mock.Mock
}

func (m *mockOpinionRepository) FindByID(ctx context.Context, id shared.UUID[opinion.Opinion]) (*opinion.Opinion, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*opinion.Opinion), args.Error(1)
}

func (m *mockOpinionRepository) Create(ctx context.Context, op opinion.Opinion) error {
	args := m.Called(ctx, op)
	return args.Error(0)
}

func (m *mockOpinionRepository) FindByParentID(ctx context.Context, id shared.UUID[opinion.Opinion]) ([]opinion.Opinion, error) {
	args := m.Called(ctx, id)
	return args.Get(0).([]opinion.Opinion), args.Error(1)
}

func (m *mockOpinionRepository) FindByTalkSessionWithoutVote(
	ctx context.Context,
	userID shared.UUID[user.User],
	talkSessionID shared.UUID[talksession.TalkSession],
	limit int,
) ([]opinion.Opinion, error) {
	args := m.Called(ctx, userID, talkSessionID, limit)
	return args.Get(0).([]opinion.Opinion), args.Error(1)
}

type mockReportRepository struct {
	mock.Mock
}

func (m *mockReportRepository) Create(ctx context.Context, report opinion.Report) error {
	args := m.Called(ctx, report)
	return args.Error(0)
}

func (m *mockReportRepository) UpdateStatus(ctx context.Context, id shared.UUID[opinion.Report], status opinion.Status) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func TestReportOpinion_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("正常系: 意見を報告できる", func(t *testing.T) {
		// モックの準備
		opinionRep := new(mockOpinionRepository)
		reportRep := new(mockReportRepository)

		// テストデータ
		opinionID := shared.NewUUID[opinion.Opinion]()
		userID := shared.NewUUID[user.User]()
		talkSessionID := shared.NewUUID[talksession.TalkSession]()
		content := "テスト意見"

		testOpinion, _ := opinion.NewOpinion(
			opinionID,
			talkSessionID,
			userID,
			nil,
			nil,
			content,
			time.Now(),
			nil,
		)

		// モックの振る舞いを設定
		opinionRep.On("FindByID", mock.Anything, opinionID).Return(testOpinion, nil)
		reportRep.On("Create", mock.Anything, mock.AnythingOfType("opinion.Report")).Return(nil)

		// テスト実行
		usecase := NewReportOpinion(opinionRep, reportRep)
		err := usecase.Execute(ctx, ReportOpinionInput{
			ReporterID: shared.NewUUID[user.User](),
			OpinionID:  opinionID,
			Reason:     1,
		})

		// アサーション
		assert.NoError(t, err)
		opinionRep.AssertExpectations(t)
		reportRep.AssertExpectations(t)
	})

	t.Run("異常系: 意見が見つからない", func(t *testing.T) {
		opinionRep := new(mockOpinionRepository)
		reportRep := new(mockReportRepository)

		opinionID := shared.NewUUID[opinion.Opinion]()

		// 意見が見つからない場合のモック
		opinionRep.On("FindByID", mock.Anything, opinionID).Return(nil, sql.ErrNoRows)

		usecase := NewReportOpinion(opinionRep, reportRep)
		err := usecase.Execute(ctx, ReportOpinionInput{
			ReporterID: shared.NewUUID[user.User](),
			OpinionID:  opinionID,
			Reason:     1,
		})

		assert.Error(t, err)
		assert.Equal(t, messages.OpinionNotFound, err)
		opinionRep.AssertExpectations(t)
	})
}
