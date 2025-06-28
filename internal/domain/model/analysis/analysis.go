package analysis

import (
	"context"
	"time"

	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/talksession"
)

type (
	AnalysisService interface {
		StartAnalysis(context.Context, shared.UUID[talksession.TalkSession]) error
		GenerateReport(context.Context, shared.UUID[talksession.TalkSession]) error
		GenerateImage(ctx context.Context, talkSessionID shared.UUID[talksession.TalkSession]) (*WordCloudResponse, error)
	}

	AnalysisRepository interface {
		FindByTalkSessionID(ctx context.Context, talkSessionID shared.UUID[talksession.TalkSession]) (*AnalysisReport, error)
	}
)
type WordCloudResponse struct {
	Wordcloud string `json:"wordcloud"`
	Tsne      string `json:"tsne"`
}

type AnalysisReport struct {
	Report    *string
	UpdatedAt time.Time
	CreatedAt time.Time
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
