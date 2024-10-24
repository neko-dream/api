// Code generated by ogen, DO NOT EDIT.

package oas

import (
	"context"
)

// Handler handles operations described by OpenAPI v3 specification.
type Handler interface {
	AuthHandler
	OpinionHandler
	TalkSessionHandler
	TestHandler
	UserHandler
	VoteHandler
}

// AuthHandler handles operations described by OpenAPI v3 specification.
//
// x-ogen-operation-group: Auth
type AuthHandler interface {
	// Authorize implements authorize operation.
	//
	// 認証画面を表示する.
	//
	// GET /auth/{provider}/login
	Authorize(ctx context.Context, params AuthorizeParams) (*AuthorizeFound, error)
	// OAuthCallback implements oauth_callback operation.
	//
	// Auth Callback.
	//
	// GET /auth/{provider}/callback
	OAuthCallback(ctx context.Context, params OAuthCallbackParams) (*OAuthCallbackFound, error)
	// OAuthRevoke implements oauth_revoke operation.
	//
	// アクセストークンを失効.
	//
	// POST /auth/revoke
	OAuthRevoke(ctx context.Context) (OAuthRevokeRes, error)
	// OAuthTokenInfo implements oauth_token_info operation.
	//
	// JWTの内容を返してくれる.
	//
	// GET /auth/token/info
	OAuthTokenInfo(ctx context.Context) (OAuthTokenInfoRes, error)
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
	// OpinionComments2 implements opinionComments2 operation.
	//
	// 意見に対するリプライ意見一覧 Copy.
	//
	// GET /talksessions/{talkSessionID}/opinions/{opinionID}/replies2
	OpinionComments2(ctx context.Context, params OpinionComments2Params) (OpinionComments2Res, error)
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
	// セッション作成.
	//
	// POST /talksessions
	CreateTalkSession(ctx context.Context, req OptCreateTalkSessionReq) (CreateTalkSessionRes, error)
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
	// 🚧 セッションレポートを返す.
	//
	// GET /talksessions/{talkSessionId}/report
	GetTalkSessionReport(ctx context.Context, params GetTalkSessionReportParams) (GetTalkSessionReportRes, error)
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
