package analysis_test

import (
	"context"
	"testing"
	"time"

	"github.com/neko-dream/server/internal/domain/model/analysis"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestFeedbackType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected analysis.FeedbackType
	}{
		{
			name:     "good を FeedbackTypeGood に変換できる",
			input:    "good",
			expected: analysis.FeedbackTypeGood,
		},
		{
			name:     "bad を FeedbackTypeBad に変換できる",
			input:    "bad",
			expected: analysis.FeedbackTypeBad,
		},
		{
			name:     "unknown を FeedbackTypeUnknown に変換できる",
			input:    "unknown",
			expected: analysis.FeedbackTypeUnknown,
		},
		{
			name:     "空文字列は FeedbackTypeUnknown になる",
			input:    "",
			expected: analysis.FeedbackTypeUnknown,
		},
		{
			name:     "不正な値は FeedbackTypeUnknown になる",
			input:    "invalid",
			expected: analysis.FeedbackTypeUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analysis.NewFeedbackTypeFromString(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAnalysisReport_ShouldReGenerateReport(t *testing.T) {
	t.Run("レポートがnilの場合はtrueを返す", func(t *testing.T) {
		report := &analysis.AnalysisReport{
			AnalysisReportID: shared.MustParseUUID[analysis.AnalysisReport]("00000000-0000-0000-0000-000000000001"),
			Report:           nil,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
			Feedbacks:        []analysis.Feedback{},
		}

		assert.True(t, report.ShouldReGenerateReport())
	})

	t.Run("レポートが存在し、10分以内の場合はfalseを返す", func(t *testing.T) {
		report := &analysis.AnalysisReport{
			AnalysisReportID: shared.MustParseUUID[analysis.AnalysisReport]("00000000-0000-0000-0000-000000000001"),
			Report:           lo.ToPtr("テストレポート"),
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
			Feedbacks:        []analysis.Feedback{},
		}

		assert.False(t, report.ShouldReGenerateReport())
	})

	t.Run("レポートが存在し、10分以上経過している場合はtrueを返す", func(t *testing.T) {
		oldTime := time.Now().Add(-11 * time.Minute)
		report := &analysis.AnalysisReport{
			AnalysisReportID: shared.MustParseUUID[analysis.AnalysisReport]("00000000-0000-0000-0000-000000000001"),
			Report:           lo.ToPtr("テストレポート"),
			CreatedAt:        oldTime,
			UpdatedAt:        oldTime,
			Feedbacks:        []analysis.Feedback{},
		}

		assert.True(t, report.ShouldReGenerateReport())
	})
}

func TestAnalysisReport_ApplyFeedback(t *testing.T) {
	ctx := context.Background()

	t.Run("フィードバックを追加できる", func(t *testing.T) {
		// 固定時刻を設定（clockパッケージの実装を確認する必要があるが、ここでは削除）
		// fixedTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

		report := &analysis.AnalysisReport{
			AnalysisReportID: shared.MustParseUUID[analysis.AnalysisReport]("00000000-0000-0000-0000-000000000001"),
			Report:           lo.ToPtr("テストレポート"),
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
			Feedbacks:        []analysis.Feedback{},
		}

		userID := shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000003")
		feedbackType := analysis.FeedbackTypeGood

		assert.Empty(t, report.Feedbacks)

		report.ApplyFeedback(ctx, feedbackType, userID)

		assert.Len(t, report.Feedbacks, 1)
		assert.Equal(t, feedbackType, report.Feedbacks[0].Type)
		assert.Equal(t, userID, report.Feedbacks[0].UserID)
		assert.NotEmpty(t, report.Feedbacks[0].FeedbackID)
		// UpdatedAtは現在時刻に更新されるはず
		assert.True(t, report.UpdatedAt.After(report.CreatedAt))
	})

	t.Run("複数のフィードバックを追加できる", func(t *testing.T) {
		report := &analysis.AnalysisReport{
			AnalysisReportID: shared.MustParseUUID[analysis.AnalysisReport]("00000000-0000-0000-0000-000000000001"),
			Report:           lo.ToPtr("テストレポート"),
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
			Feedbacks:        []analysis.Feedback{},
		}

		userID1 := shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000003")
		userID2 := shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000004")

		report.ApplyFeedback(ctx, analysis.FeedbackTypeGood, userID1)
		report.ApplyFeedback(ctx, analysis.FeedbackTypeBad, userID2)

		assert.Len(t, report.Feedbacks, 2)
		assert.Equal(t, analysis.FeedbackTypeGood, report.Feedbacks[0].Type)
		assert.Equal(t, userID1, report.Feedbacks[0].UserID)
		assert.Equal(t, analysis.FeedbackTypeBad, report.Feedbacks[1].Type)
		assert.Equal(t, userID2, report.Feedbacks[1].UserID)
	})
}

func TestAnalysisReport_HasReceivedFeedbackFrom(t *testing.T) {
	ctx := context.Background()

	t.Run("フィードバックを受け取っていないユーザーはfalseを返す", func(t *testing.T) {
		report := &analysis.AnalysisReport{
			AnalysisReportID: shared.MustParseUUID[analysis.AnalysisReport]("00000000-0000-0000-0000-000000000001"),
			Report:           lo.ToPtr("テストレポート"),
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
			Feedbacks:        []analysis.Feedback{},
		}

		userID := shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000003")
		assert.False(t, report.HasReceivedFeedbackFrom(ctx, userID))
	})

	t.Run("フィードバックを受け取ったユーザーはtrueを返す", func(t *testing.T) {
		userID1 := shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000003")
		userID2 := shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000004")

		report := &analysis.AnalysisReport{
			AnalysisReportID: shared.MustParseUUID[analysis.AnalysisReport]("00000000-0000-0000-0000-000000000001"),
			Report:           lo.ToPtr("テストレポート"),
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
			Feedbacks: []analysis.Feedback{
				{
					FeedbackID: shared.MustParseUUID[analysis.Feedback]("00000000-0000-0000-0000-000000000005"),
					Type:       analysis.FeedbackTypeGood,
					UserID:     userID1,
				},
			},
		}

		assert.True(t, report.HasReceivedFeedbackFrom(ctx, userID1))
		assert.False(t, report.HasReceivedFeedbackFrom(ctx, userID2))
	})

	t.Run("複数のフィードバックから正しく判定できる", func(t *testing.T) {
		userID1 := shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000003")
		userID2 := shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000004")
		userID3 := shared.MustParseUUID[user.User]("00000000-0000-0000-0000-000000000005")

		report := &analysis.AnalysisReport{
			AnalysisReportID: shared.MustParseUUID[analysis.AnalysisReport]("00000000-0000-0000-0000-000000000001"),
			Report:           lo.ToPtr("テストレポート"),
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
			Feedbacks: []analysis.Feedback{
				{
					FeedbackID: shared.MustParseUUID[analysis.Feedback]("00000000-0000-0000-0000-000000000006"),
					Type:       analysis.FeedbackTypeGood,
					UserID:     userID1,
				},
				{
					FeedbackID: shared.MustParseUUID[analysis.Feedback]("00000000-0000-0000-0000-000000000007"),
					Type:       analysis.FeedbackTypeBad,
					UserID:     userID2,
				},
			},
		}

		assert.True(t, report.HasReceivedFeedbackFrom(ctx, userID1))
		assert.True(t, report.HasReceivedFeedbackFrom(ctx, userID2))
		assert.False(t, report.HasReceivedFeedbackFrom(ctx, userID3))
	})
}
