package repository

import (
	"context"

	"github.com/neko-dream/api/internal/domain/model/analysis"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/talksession"
	"go.opentelemetry.io/otel"
)

// AnalysisServiceMock is a mock implementation of AnalysisService for testing
type AnalysisServiceMock struct{}

// NewAnalysisServiceMock creates a new mock analysis service
func NewAnalysisServiceMock() analysis.AnalysisService {
	return &AnalysisServiceMock{}
}

// StartAnalysis starts analysis for a talk session
func (m *AnalysisServiceMock) StartAnalysis(ctx context.Context, talkSessionID shared.UUID[talksession.TalkSession]) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "AnalysisServiceMock.StartAnalysis")
	defer span.End()

	_ = ctx

	return nil
}

// GenerateReport generates a report for a talk session
func (m *AnalysisServiceMock) GenerateReport(ctx context.Context, talkSessionID shared.UUID[talksession.TalkSession]) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "AnalysisServiceMock.GenerateReport")
	defer span.End()

	_ = ctx

	return nil
}

// GenerateImage generates images for a talk session
func (m *AnalysisServiceMock) GenerateImage(ctx context.Context, talkSessionID shared.UUID[talksession.TalkSession]) (*analysis.WordCloudResponse, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "AnalysisServiceMock.GenerateImage")
	defer span.End()

	_ = ctx

	return nil, nil
}
