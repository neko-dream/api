package client

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/neko-dream/server/internal/domain/model/analysis"
	"github.com/neko-dream/server/internal/domain/model/image"
	"github.com/neko-dream/server/internal/domain/model/image/meta"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/talksession"
	"github.com/neko-dream/server/internal/infrastructure/config"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/server/internal/infrastructure/persistence/sqlc/generated"
	http_utils "github.com/neko-dream/server/pkg/http"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

type analysisService struct {
	conf     *config.Config
	imageRep image.ImageStorage
	db.DBManager
}

func NewAnalysisService(
	conf *config.Config,
	imageRep image.ImageStorage,
	dbm *db.DBManager,
) analysis.AnalysisService {
	return &analysisService{
		conf:      conf,
		imageRep:  imageRep,
		DBManager: *dbm,
	}
}

// GenerateReport implements analysis.AnalysisService.
func (a *analysisService) GenerateReport(ctx context.Context, talkSessionID shared.UUID[talksession.TalkSession]) error {
	ctx, span := otel.Tracer("client").Start(ctx, "analysisService.GenerateReport")
	defer span.End()

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
func (a *analysisService) StartAnalysis(ctx context.Context, talkSessionID shared.UUID[talksession.TalkSession]) error {
	ctx, span := otel.Tracer("client").Start(ctx, "analysisService.StartAnalysis")
	defer span.End()

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
		UserId:        "0",
	})
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil
	}
	return nil
}

func (a *analysisService) GenerateImage(ctx context.Context, talkSessionID shared.UUID[talksession.TalkSession]) (*analysis.WordCloudResponse, error) {
	ctx, span := otel.Tracer("client").Start(ctx, "analysisService.GenerateImage")
	defer span.End()

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
		return nil, err
	}

	// APIリクエストの実行
	resp, err := c.PostReportsWordclouds(ctx, PostReportsWordcloudsJSONRequestBody{
		TalkSessionId: talkSessionID.UUID().String(),
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var wordcloud analysis.WordCloudResponse
	if err := json.NewDecoder(resp.Body).Decode(&wordcloud); err != nil {
		utils.HandleError(ctx, err, "json.NewDecoder.Decode")
		return nil, err
	}

	// Base64デコード
	wcData, err := base64.StdEncoding.DecodeString(wordcloud.Wordcloud)
	if err != nil {
		utils.HandleError(ctx, err, "base64.StdEncoding.DecodeString")
		return nil, err
	}
	tsneData, err := base64.StdEncoding.DecodeString(wordcloud.Tsne)
	if err != nil {
		utils.HandleError(ctx, err, "base64.StdEncoding.DecodeString")
		return nil, err
	}
	wcBuf := bytes.NewBuffer(wcData)
	wcImgInfo, err := meta.NewImageForAnalysis(ctx, wcBuf)
	if err != nil {
		utils.HandleError(ctx, err, "meta.NewImageForAnalysis")
		return nil, err
	}
	tsncBuf := bytes.NewBuffer(tsneData)
	tsnCImgInfo, err := meta.NewImageForAnalysis(ctx, tsncBuf)
	if err != nil {
		utils.HandleError(ctx, err, "meta.NewImageForAnalysis")
		return nil, err
	}

	// *multipart.FileHeaderを作成
	wcFile, err := http_utils.CreateFileHeader(ctx, wcBuf, "wordcloud.png")
	if err != nil {
		utils.HandleError(ctx, err, "http_utils.CreateFileHeader")
		return nil, err
	}
	tsneFile, err := http_utils.CreateFileHeader(ctx, tsncBuf, "tsne.png")
	if err != nil {
		utils.HandleError(ctx, err, "http_utils.CreateFileHeader")
		return nil, err
	}

	wc, err := a.imageRep.Upload(ctx, *wcImgInfo, wcFile)
	if err != nil {
		utils.HandleError(ctx, err, "imageRep.Upload")
		return nil, err
	}
	tsne, err := a.imageRep.Upload(ctx, *tsnCImgInfo, tsneFile)
	if err != nil {
		utils.HandleError(ctx, err, "imageRep.Upload")
		return nil, err
	}

	// 画像情報をDBに保存
	if err := a.DBManager.GetQueries(ctx).AddGeneratedImages(ctx, model.AddGeneratedImagesParams{
		TalkSessionID: talkSessionID.UUID(),
		WordmapUrl:    *wc,
		TsncUrl:       *tsne,
	}); err != nil {
		return nil, err
	}

	return &wordcloud, nil
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
