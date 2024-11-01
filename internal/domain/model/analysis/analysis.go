package analysis

import (
	"context"

	"github.com/neko-dream/server/internal/domain/model/shared"
	talksession "github.com/neko-dream/server/internal/domain/model/talk_session"
)

type (
	AnalysisService interface {
		StartAnalysis(context.Context, shared.UUID[talksession.TalkSession]) error
		GenerateReport(context.Context, shared.UUID[talksession.TalkSession]) error
		GenerateImage(ctx context.Context, talkSessionID shared.UUID[talksession.TalkSession]) (*WordCloudResponse, error)
	}
)
type WordCloudResponse struct {
	Wordcloud string `json:"wordcloud"`
	Tsne      string `json:"tsne"`
}
