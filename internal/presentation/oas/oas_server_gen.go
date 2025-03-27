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
	PolicyHandler
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
	// AuthAccountDetach implements authAccountDetach operation.
	//
	// そのアカウントには再度ログインできなくなります。ログインしたければ言ってね！.
	//
	// DELETE /auth/dev/detach
	AuthAccountDetach(ctx context.Context) (AuthAccountDetachRes, error)
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
	// GetOpinionAnalysis implements getOpinionAnalysis operation.
	//
	// 意見に投票したグループごとの割合.
	//
	// GET /opinions/{opinionID}/analysis
	GetOpinionAnalysis(ctx context.Context, params GetOpinionAnalysisParams) (GetOpinionAnalysisRes, error)
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
	// GetOpinionReportReasons implements getOpinionReportReasons operation.
	//
	// 意見への通報理由一覧.
	//
	// GET /opinions/report_reasons
	GetOpinionReportReasons(ctx context.Context) (GetOpinionReportReasonsRes, error)
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
	// 意見に対するリプライ意見一覧.
	//
	// GET /opinions/{opinionID}/replies
	OpinionComments2(ctx context.Context, params OpinionComments2Params) (OpinionComments2Res, error)
	// PostOpinionPost implements postOpinionPost operation.
	//
	// ParentOpinionIDがなければルートの意見として投稿される.
	//
	// POST /talksessions/{talkSessionID}/opinions
	PostOpinionPost(ctx context.Context, req OptPostOpinionPostReq, params PostOpinionPostParams) (PostOpinionPostRes, error)
	// PostOpinionPost2 implements postOpinionPost2 operation.
	//
	// ParentOpinionIDがなければルートの意見として投稿される
	// parentOpinionIDがない場合はtalkSessionIDが必須.
	//
	// POST /opinions
	PostOpinionPost2(ctx context.Context, req OptPostOpinionPost2Req) (PostOpinionPost2Res, error)
	// ReportOpinion implements reportOpinion operation.
	//
	// 意見通報API.
	//
	// POST /opinions/{opinionID}/report
	ReportOpinion(ctx context.Context, req OptReportOpinionReq, params ReportOpinionParams) (ReportOpinionRes, error)
	// SwipeOpinions implements swipe_opinions operation.
	//
	// セッションの中からまだ投票していない意見をランダムに取得する
	// remainingCountは取得した意見を含めてスワイプできる意見の総数を返す.
	//
	// GET /talksessions/{talkSessionID}/swipe_opinions
	SwipeOpinions(ctx context.Context, params SwipeOpinionsParams) (SwipeOpinionsRes, error)
}

// PolicyHandler handles operations described by OpenAPI v3 specification.
//
// x-ogen-operation-group: Policy
type PolicyHandler interface {
	// GetPolicyConsentStatus implements getPolicyConsentStatus operation.
	//
	// 最新のポリシーに同意したかを取得.
	//
	// GET /policy/consent
	GetPolicyConsentStatus(ctx context.Context) (GetPolicyConsentStatusRes, error)
	// PolicyConsent implements policyConsent operation.
	//
	// 最新のポリシーに同意する.
	//
	// POST /policy/consent
	PolicyConsent(ctx context.Context, req OptPolicyConsentReq) (PolicyConsentRes, error)
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
	// GetTalkSessionRestrictionSatisfied implements getTalkSessionRestrictionSatisfied operation.
	//
	// 特定のセッションで満たしていない条件があれば返す.
	//
	// GET /talksessions/{talkSessionID}/restrictions
	GetTalkSessionRestrictionSatisfied(ctx context.Context, params GetTalkSessionRestrictionSatisfiedParams) (GetTalkSessionRestrictionSatisfiedRes, error)
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
	// Vote2 implements vote2 operation.
	//
	// 意思表明API.
	//
	// POST /opinions/{opinionID}/votes
	Vote2(ctx context.Context, req OptVote2Req, params Vote2Params) (Vote2Res, error)
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
