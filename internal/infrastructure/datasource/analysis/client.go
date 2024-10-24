package client

import (
	"context"
	"net/http"

	"github.com/neko-dream/server/internal/domain/model/analysis"
	"github.com/neko-dream/server/internal/domain/model/shared"
	talksession "github.com/neko-dream/server/internal/domain/model/talk_session"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/config"
	"github.com/neko-dream/server/pkg/utils"
)

type analysisService struct {
	conf *config.Config
}

func NewAnalysisService(
	conf *config.Config,
) analysis.AnalysisService {
	return &analysisService{
		conf: conf,
	}
}

// GenerateReport implements analysis.AnalysisService.
func (a *analysisService) GenerateReport(ctx context.Context, talkSessionID shared.UUID[talksession.TalkSession]) error {
	// カスタムHTTPクライアントを作成
	httpClient := &http.Client{
		Transport: &BasicAuthTransport{
			Username: a.conf.ANALYSIS_USER,
			Password: a.conf.ANALYSIS_USER_PASSWORD,
		},
	}

	// クライアントの初期化
	c, err := NewClient(a.conf.ANALYSIS_API_DOMAIN, WithHTTPClient(httpClient))
	if err != nil {
		utils.HandleError(ctx, err, "NewClient")
		return err
	}

	// APIリクエストの実行
	resp, err := c.PostReportsGenerates(ctx, PostReportsGeneratesJSONRequestBody{
		TalkSessionId: talkSessionID.String(),
	})
	if err != nil {
		utils.HandleError(ctx, err, "PostReportsGenerates")
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		utils.HandleError(ctx, err, "PostReportsGenerates")
		return err
	}

	return nil
}

// StartAnalysis 会話分析を開始する
func (a *analysisService) StartAnalysis(ctx context.Context, talkSessionID shared.UUID[talksession.TalkSession], userID shared.UUID[user.User]) error {
	go func() error {
		// カスタムHTTPクライアントを作成
		httpClient := &http.Client{
			Transport: &BasicAuthTransport{
				Username: a.conf.ANALYSIS_USER,
				Password: a.conf.ANALYSIS_USER_PASSWORD,
			},
		}

		// クライアントの初期化
		c, err := NewClient(a.conf.ANALYSIS_API_DOMAIN, WithHTTPClient(httpClient))
		if err != nil {
			return nil
		}

		// APIリクエストの実行
		resp, err := c.PostPredictsGroups(ctx, PostPredictsGroupsJSONRequestBody{
			TalkSessionId: talkSessionID.String(),
			UserId:        userID.String(),
		})
		if err != nil {
			return nil
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil
		}
		return nil
	}()
	return nil
}

// Basic認証用のカスタムTransport
type BasicAuthTransport struct {
	Username string
	Password string
}

func (t *BasicAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(t.Username, t.Password)
	return http.DefaultTransport.RoundTrip(req)
}
