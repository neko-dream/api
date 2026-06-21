package repository

import (
	"context"

	"github.com/neko-dream/api/internal/application/usecase/analysis_usecase"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/talksession"
)

// AnalysisSchedulerMock はテスト用のAnalysisSchedulerモック実装
type AnalysisSchedulerMock struct{}

// NewAnalysisSchedulerMock creates a new mock analysis scheduler
func NewAnalysisSchedulerMock() analysis_usecase.AnalysisScheduler {
	return &AnalysisSchedulerMock{}
}

func (m *AnalysisSchedulerMock) ScheduleAnalysis(_ context.Context, _ shared.UUID[talksession.TalkSession], _ analysis_usecase.ScheduleMode) {
}

func (m *AnalysisSchedulerMock) Start(_ context.Context) {}

func (m *AnalysisSchedulerMock) Stop() {}
