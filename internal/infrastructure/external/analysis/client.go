package client

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/neko-dream/server/internal/domain/model/analysis"
	"github.com/neko-dream/server/internal/domain/model/clock"
	"github.com/neko-dream/server/internal/domain/model/image"
	"github.com/neko-dream/server/internal/domain/model/shared"
	talksession "github.com/neko-dream/server/internal/domain/model/talk_session"
	"github.com/neko-dream/server/internal/infrastructure/config"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/server/internal/infrastructure/persistence/sqlc/generated"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

type analysisService struct {
	conf     *config.Config
	imageRep image.ImageRepository
	db.DBManager
}

func NewAnalysisService(
	conf *config.Config,
	imageRep image.ImageRepository,
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
		log.Println("PostPredictsGroups", err)
		return nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Println("PostPredictsGroups", err)
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
	println("GenerateImage")

	// クライアントの初期化
	c, err := NewClient(a.conf.ANALYSIS_API_DOMAIN, WithHTTPClient(httpClient))
	if err != nil {
		log.Println("NewClient", err)
		return nil, err
	}
	println("request")

	// APIリクエストの実行
	resp, err := c.PostReportsWordclouds(ctx, PostReportsWordcloudsJSONRequestBody{
		TalkSessionId: talkSessionID.UUID().String(),
	})
	if err != nil {
		log.Println("PostReportsWordclouds", err)
		return nil, err
	}
	println("response")
	defer resp.Body.Close()

	var wordcloud analysis.WordCloudResponse
	if err := json.NewDecoder(resp.Body).Decode(&wordcloud); err != nil {
		log.Println("json.NewDecoder", err)
		return nil, err
	}
	println("decode end")
	// Base64デコード
	wcData, err := base64.StdEncoding.DecodeString(wordcloud.Wordcloud)
	if err != nil {
		return nil, err
	}
	tsneData, err := base64.StdEncoding.DecodeString(wordcloud.Tsne)
	if err != nil {
		return nil, err
	}

	// 種類-talkSessionID-時間.jpg
	objectPath := "generated/%v-%v-%v.png"
	wcImg := image.NewImage(wcData)
	tsncImg := image.NewImage(tsneData)
	now := clock.Now(ctx)
	wcImgInfo := image.NewImageInfo(
		fmt.Sprintf(objectPath, "wordcloud", talkSessionID.String(), now.UnixNano()),
		"png",
		wcImg,
	)
	tsnCImgInfo := image.NewImageInfo(
		fmt.Sprintf(objectPath, "tsne", talkSessionID.String(), now.UnixNano()),
		"png",
		tsncImg,
	)

	wc, err := a.imageRep.Create(ctx, *wcImgInfo)
	if err != nil {
		return nil, err
	}
	tsne, err := a.imageRep.Create(ctx, *tsnCImgInfo)
	if err != nil {
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
