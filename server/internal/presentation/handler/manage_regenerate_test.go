package handler

import (
	"context"
	"testing"
	"time"

	"github.com/neko-dream/api/internal/domain/model/analysis"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/talksession"
	"github.com/neko-dream/api/internal/presentation/oas"
)

type imageAnalysisService struct {
	analysis.AnalysisService
	proceed chan struct{}
	ctxErr  chan error
}

func (s *imageAnalysisService) GenerateImage(ctx context.Context, _ shared.UUID[talksession.TalkSession]) (*analysis.WordCloudResponse, error) {
	<-s.proceed
	s.ctxErr <- ctx.Err()
	return nil, nil
}

// 画像の生成はレスポンスを返したあとも続くので、リクエストの終了でキャンセルされない ctx で動かす
func TestManageRegenerateImageOutlivesRequest(t *testing.T) {
	svc := &imageAnalysisService{proceed: make(chan struct{}), ctxErr: make(chan error, 1)}
	h := &manageHandler{AnalysisService: svc}
	ctx, cancel := context.WithCancel(adminCtx())

	if _, err := h.ManageRegenerateManage(ctx, &oas.RegenerateRequest{Type: oas.RegenerateRequestTypeImage}, oas.ManageRegenerateManageParams{TalkSessionID: shared.NewUUID[talksession.TalkSession]().String()}); err != nil {
		t.Fatal(err)
	}
	cancel() // レスポンスを返し、リクエストが終わった
	close(svc.proceed)

	select {
	case err := <-svc.ctxErr:
		if err != nil {
			t.Fatalf("image generation context was cancelled with the request: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("image generation did not run")
	}
}
