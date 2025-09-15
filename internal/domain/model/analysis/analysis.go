package analysis

import (
	"context"
	"time"

	"github.com/neko-dream/api/internal/domain/model/clock"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/talksession"
	"github.com/neko-dream/api/internal/domain/model/user"
	"go.opentelemetry.io/otel"
)

type (
	AnalysisService interface {
		StartAnalysis(context.Context, shared.UUID[talksession.TalkSession]) error
		GenerateReport(context.Context, shared.UUID[talksession.TalkSession]) error
		GenerateImage(ctx context.Context, talkSessionID shared.UUID[talksession.TalkSession]) (*WordCloudResponse, error)
	}

	AnalysisRepository interface {
		FindByTalkSessionID(ctx context.Context, talkSessionID shared.UUID[talksession.TalkSession]) (*AnalysisReport, error)
		FindByID(ctx context.Context, analysisReportID shared.UUID[AnalysisReport]) (*AnalysisReport, error)
		SaveReport(ctx context.Context, report *AnalysisReport) error
	}
)

type WordCloudResponse struct {
	Wordcloud string `json:"wordcloud"`
	Tsne      string `json:"tsne"`
}

type FeedbackType int

const (
	FeedbackTypeUnknown FeedbackType = iota
	FeedbackTypeGood
	FeedbackTypeBad
)

func NewFeedbackTypeFromString(s string) FeedbackType {
	switch s {
	case "good":
		return FeedbackTypeGood
	case "bad":
		return FeedbackTypeBad
	default:
		return FeedbackTypeUnknown
	}
}

type Feedback struct {
	FeedbackID shared.UUID[Feedback]
	Type       FeedbackType
	UserID     shared.UUID[user.User]
	CreatedAt  time.Time
}

type AnalysisReport struct {
	AnalysisReportID shared.UUID[AnalysisReport]
	Report           *string
	UpdatedAt        time.Time
	CreatedAt        time.Time

	Feedbacks []Feedback
}

// ShouldReGenerateReport 再生成するかどうかを判定する
func (r *AnalysisReport) ShouldReGenerateReport() bool {
	if r.Report == nil {
		// レポートが存在しない場合は再生成する
		return true
	}
	// 10分以上経過していれば再生成する
	return time.Since(r.UpdatedAt) > 10*time.Minute
}

func (r *AnalysisReport) ApplyFeedback(ctx context.Context, feedbackType FeedbackType, userID shared.UUID[user.User]) {
	ctx, span := otel.Tracer("analysis").Start(ctx, "AnalysisReport.ApplyFeedback")
	defer span.End()

	r.Feedbacks = append(r.Feedbacks, Feedback{
		FeedbackID: shared.NewUUID[Feedback](),
		Type:       feedbackType,
		UserID:     userID,
	})

	// フィードバックが追加されたら更新日時を更新する
	r.UpdatedAt = clock.Now(ctx)
}

func (r *AnalysisReport) HasReceivedFeedbackFrom(ctx context.Context, userID shared.UUID[user.User]) bool {
	ctx, span := otel.Tracer("analysis").Start(ctx, "AnalysisReport.HasReceivedFeedbackFrom")
	defer span.End()

	_ = ctx

	for _, feedback := range r.Feedbacks {
		if feedback.UserID == userID {
			return true
		}
	}

	return false
}
