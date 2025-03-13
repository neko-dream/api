// Code generated by ogen, DO NOT EDIT.

package oas

import (
	"context"
)

// Handler handles operations described by OpenAPI v3 specification.
type Handler interface {
	AuthHandler
	ImageHandler
	ManageHandler
	OpinionHandler
	TalkSessionHandler
	TestHandler
	TimelineHandler
	UserHandler
	VoteHandler
}

// AuthHandler handles operations described by OpenAPI v3 specification.
//
// x-ogen-operation-group: Auth
type AuthHandler interface {
	// Authorize implements authorize operation.
	//
	// ログイン.
	//
	// GET /auth/{provider}/login
	Authorize(ctx context.Context, params AuthorizeParams) (AuthorizeRes, error)
	// DevAuthorize implements devAuthorize operation.
	//
	// 開発用登録/ログイン.
	//
	// GET /auth/dev/login
	DevAuthorize(ctx context.Context, params DevAuthorizeParams) (DevAuthorizeRes, error)
	// OAuthCallback implements oauth_callback operation.
	//
	// Auth Callback.
	//
	// GET /auth/{provider}/callback
	OAuthCallback(ctx context.Context, params OAuthCallbackParams) (OAuthCallbackRes, error)
	// OAuthTokenInfo implements oauth_token_info operation.
	//
	// JWTの内容を返してくれる.
	//
	// GET /auth/token/info
	OAuthTokenInfo(ctx context.Context) (OAuthTokenInfoRes, error)
	// OAuthTokenRevoke implements oauth_token_revoke operation.
	//
	// トークンを失効（ログアウト）.
	//
	// POST /auth/revoke
	OAuthTokenRevoke(ctx context.Context) (OAuthTokenRevokeRes, error)
}

// ImageHandler handles operations described by OpenAPI v3 specification.
//
// x-ogen-operation-group: Image
type ImageHandler interface {
	// PostImage implements postImage operation.
	//
	// 画像を投稿してURLを返すAPI.
	//
	// POST /images
	PostImage(ctx context.Context, req OptPostImageReq) (PostImageRes, error)
}

// ManageHandler handles operations described by OpenAPI v3 specification.
//
// x-ogen-operation-group: Manage
type ManageHandler interface {
	// ManageIndex implements manageIndex operation.
	//
	// GET /manage
	ManageIndex(ctx context.Context) (ManageIndexOK, error)
	// ManageRegenerate implements manageRegenerate operation.
	//
	// Analysisを再生成する。enum: [report, group, image].
	//
	// POST /manage/regenerate
	ManageRegenerate(ctx context.Context, req OptManageRegenerateReq) (*ManageRegenerateOK, error)
}

// OpinionHandler handles operations described by OpenAPI v3 specification.
//
// x-ogen-operation-group: Opinion
type OpinionHandler interface {
	// GetOpinionDetail implements getOpinionDetail operation.
	//
	// 意見の詳細.
	//
	// GET /talksessions/{talkSessionID}/opinions/{opinionID}
	GetOpinionDetail(ctx context.Context, params GetOpinionDetailParams) (GetOpinionDetailRes, error)
	// GetOpinionDetail2 implements getOpinionDetail2 operation.
	//
	// 意見詳細.
	//
	// GET /opinions/{opinionID}
	GetOpinionDetail2(ctx context.Context, params GetOpinionDetail2Params) (GetOpinionDetail2Res, error)
	// GetOpinionsForTalkSession implements getOpinionsForTalkSession operation.
	//
	// セッションに対する意見一覧.
	//
	// GET /talksessions/{talkSessionID}/opinions
	GetOpinionsForTalkSession(ctx context.Context, params GetOpinionsForTalkSessionParams) (GetOpinionsForTalkSessionRes, error)
	// OpinionComments implements opinionComments operation.
	//
	// 意見に対するリプライ意見一覧.
	//
	// GET /talksessions/{talkSessionID}/opinions/{opinionID}/replies
	OpinionComments(ctx context.Context, params OpinionCommentsParams) (OpinionCommentsRes, error)
	// PostOpinionPost implements postOpinionPost operation.
	//
	// ParentOpinionIDがなければルートの意見として投稿される.
	//
	// POST /talksessions/{talkSessionID}/opinions
	PostOpinionPost(ctx context.Context, req OptPostOpinionPostReq, params PostOpinionPostParams) (PostOpinionPostRes, error)
	// SwipeOpinions implements swipe_opinions operation.
	//
	// セッションの中からまだ投票していない意見をランダムに取得する.
	//
	// GET /talksessions/{talkSessionID}/swipe_opinions
	SwipeOpinions(ctx context.Context, params SwipeOpinionsParams) (SwipeOpinionsRes, error)
}

// TalkSessionHandler handles operations described by OpenAPI v3 specification.
//
// x-ogen-operation-group: TalkSession
type TalkSessionHandler interface {
	// CreateTalkSession implements createTalkSession operation.
	//
	// ## サムネイル画像について
	// - `Description中に出てくる画像で一番最初のものを使用`。
	// - 画像自体は`POST /images`でサーバにポストしたものを使用してください。
	// ## 投稿制限のキーについて
	// restrictionsに値を入れると一定のデモグラ情報を登録していない限り、セッションへの投稿が制限されるようにできる。
	// restrictionsには [GET /talksessions/restrictions](https://app.apidog.
	// com/link/project/674502/apis/api-14271260)
	// より取れるkeyをカンマ区切りで入力してください。.
	//
	// POST /talksessions
	CreateTalkSession(ctx context.Context, req OptCreateTalkSessionReq) (CreateTalkSessionRes, error)
	// EditTalkSession implements editTalkSession operation.
	//
	// セッション編集.
	//
	// PUT /talksessions/{talkSessionId}
	EditTalkSession(ctx context.Context, req OptEditTalkSessionReq, params EditTalkSessionParams) (EditTalkSessionRes, error)
	// GetConclusion implements getConclusion operation.
	//
	// 結論取得.
	//
	// GET /talksessions/{talkSessionID}/conclusion
	GetConclusion(ctx context.Context, params GetConclusionParams) (GetConclusionRes, error)
	// GetOpenedTalkSession implements getOpenedTalkSession operation.
	//
	// 自分が開いたセッション一覧.
	//
	// GET /talksessions/opened
	GetOpenedTalkSession(ctx context.Context, params GetOpenedTalkSessionParams) (GetOpenedTalkSessionRes, error)
	// GetTalkSessionDetail implements getTalkSessionDetail operation.
	//
	// トークセッションの詳細.
	//
	// GET /talksessions/{talkSessionId}
	GetTalkSessionDetail(ctx context.Context, params GetTalkSessionDetailParams) (GetTalkSessionDetailRes, error)
	// GetTalkSessionList implements getTalkSessionList operation.
	//
	// セッション一覧.
	//
	// GET /talksessions
	GetTalkSessionList(ctx context.Context, params GetTalkSessionListParams) (GetTalkSessionListRes, error)
	// GetTalkSessionReport implements getTalkSessionReport operation.
	//
	// セッションレポートを返す.
	//
	// GET /talksessions/{talkSessionId}/report
	GetTalkSessionReport(ctx context.Context, params GetTalkSessionReportParams) (GetTalkSessionReportRes, error)
	// GetTalkSessionRestrictionKeys implements getTalkSessionRestrictionKeys operation.
	//
	// セッションの投稿制限に使用できるキーの一覧を返す.
	//
	// GET /talksessions/restrictions
	GetTalkSessionRestrictionKeys(ctx context.Context) (GetTalkSessionRestrictionKeysRes, error)
	// PostConclusion implements postConclusion operation.
	//
	// 結論（conclusion）はセッションが終了した後にセッっションの作成者が投稿できる文章。
	// セッションの流れやグループの分かれ方などに対するセッション作成者の感想やそれらの意見を受け、これからの方向性などを記入する。.
	//
	// POST /talksessions/{talkSessionID}/conclusion
	PostConclusion(ctx context.Context, req OptPostConclusionReq, params PostConclusionParams) (PostConclusionRes, error)
	// TalkSessionAnalysis implements talkSessionAnalysis operation.
	//
	// 分析結果一覧.
	//
	// GET /talksessions/{talkSessionId}/analysis
	TalkSessionAnalysis(ctx context.Context, params TalkSessionAnalysisParams) (TalkSessionAnalysisRes, error)
}

// TestHandler handles operations described by OpenAPI v3 specification.
//
// x-ogen-operation-group: Test
type TestHandler interface {
	// DummiInit implements dummiInit operation.
	//
	// Mudai.
	//
	// POST /test/dummy
	DummiInit(ctx context.Context) (DummiInitRes, error)
	// Test implements test operation.
	//
	// OpenAPIテスト用.
	//
	// GET /test
	Test(ctx context.Context) (TestRes, error)
}

// TimelineHandler handles operations described by OpenAPI v3 specification.
//
// x-ogen-operation-group: Timeline
type TimelineHandler interface {
	// EditTimeLine implements editTimeLine operation.
	//
	// タイムライン編集.
	//
	// PUT /talksessions/{talkSessionID}/timelines/{actionItemID}
	EditTimeLine(ctx context.Context, req OptEditTimeLineReq, params EditTimeLineParams) (EditTimeLineRes, error)
	// GetTimeLine implements getTimeLine operation.
	//
	// タイムラインはセッション終了後にセッション作成者が設定できるその後の予定を知らせるもの.
	//
	// GET /talksessions/{talkSessionID}/timelines
	GetTimeLine(ctx context.Context, params GetTimeLineParams) (GetTimeLineRes, error)
	// PostTimeLineItem implements postTimeLineItem operation.
	//
	// タイムラインアイテム追加.
	//
	// POST /talksessions/{talkSessionID}/timeline
	PostTimeLineItem(ctx context.Context, req OptPostTimeLineItemReq, params PostTimeLineItemParams) (PostTimeLineItemRes, error)
}

// UserHandler handles operations described by OpenAPI v3 specification.
//
// x-ogen-operation-group: User
type UserHandler interface {
	// EditUserProfile implements editUserProfile operation.
	//
	// ユーザー情報の変更.
	//
	// PUT /user
	EditUserProfile(ctx context.Context, req OptEditUserProfileReq) (EditUserProfileRes, error)
	// GetUserInfo implements get_user_info operation.
	//
	// ユーザー情報の取得.
	//
	// GET /user
	GetUserInfo(ctx context.Context) (GetUserInfoRes, error)
	// OpinionsHistory implements opinionsHistory operation.
	//
	// 今までに投稿した異見.
	//
	// GET /opinions/histories
	OpinionsHistory(ctx context.Context, params OpinionsHistoryParams) (OpinionsHistoryRes, error)
	// RegisterUser implements registerUser operation.
	//
	// ユーザー作成.
	//
	// POST /user
	RegisterUser(ctx context.Context, req OptRegisterUserReq) (RegisterUserRes, error)
	// SessionsHistory implements sessionsHistory operation.
	//
	// リアクション済みのセッション一覧.
	//
	// GET /talksessions/histories
	SessionsHistory(ctx context.Context, params SessionsHistoryParams) (SessionsHistoryRes, error)
}

// VoteHandler handles operations described by OpenAPI v3 specification.
//
// x-ogen-operation-group: Vote
type VoteHandler interface {
	// Vote implements vote operation.
	//
	// 意思表明API.
	//
	// POST /talksessions/{talkSessionID}/opinions/{opinionID}/votes
	Vote(ctx context.Context, req OptVoteReq, params VoteParams) (VoteRes, error)
}

// Server implements http server based on OpenAPI v3 specification and
// calls Handler to handle requests.
type Server struct {
	h   Handler
	sec SecurityHandler
	baseServer
}

// NewServer creates new Server.
func NewServer(h Handler, sec SecurityHandler, opts ...ServerOption) (*Server, error) {
	s, err := newServerConfig(opts...).baseServer()
	if err != nil {
		return nil, err
	}
	return &Server{
		h:          h,
		sec:        sec,
		baseServer: s,
	}, nil
}
