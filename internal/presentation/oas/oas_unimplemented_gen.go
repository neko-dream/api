// Code generated by ogen, DO NOT EDIT.

package oas

import (
	"context"

	ht "github.com/ogen-go/ogen/http"
)

// UnimplementedHandler is no-op Handler which returns http.ErrNotImplemented.
type UnimplementedHandler struct{}

var _ Handler = UnimplementedHandler{}

// AuthAccountDetach implements authAccountDetach operation.
//
// そのアカウントには再度ログインできなくなります。ログインしたければ言ってね！.
//
// DELETE /auth/dev/detach
func (UnimplementedHandler) AuthAccountDetach(ctx context.Context) (r AuthAccountDetachRes, _ error) {
	return r, ht.ErrNotImplemented
}

// Authorize implements authorize operation.
//
// ログイン.
//
// GET /auth/{provider}/login
func (UnimplementedHandler) Authorize(ctx context.Context, params AuthorizeParams) (r AuthorizeRes, _ error) {
	return r, ht.ErrNotImplemented
}

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
func (UnimplementedHandler) CreateTalkSession(ctx context.Context, req OptCreateTalkSessionReq) (r CreateTalkSessionRes, _ error) {
	return r, ht.ErrNotImplemented
}

// DevAuthorize implements devAuthorize operation.
//
// 開発用登録/ログイン.
//
// GET /auth/dev/login
func (UnimplementedHandler) DevAuthorize(ctx context.Context, params DevAuthorizeParams) (r DevAuthorizeRes, _ error) {
	return r, ht.ErrNotImplemented
}

// DummiInit implements dummiInit operation.
//
// Mudai.
//
// POST /test/dummy
func (UnimplementedHandler) DummiInit(ctx context.Context) (r DummiInitRes, _ error) {
	return r, ht.ErrNotImplemented
}

// EditTalkSession implements editTalkSession operation.
//
// セッション編集.
//
// PUT /talksessions/{talkSessionID}
func (UnimplementedHandler) EditTalkSession(ctx context.Context, req OptEditTalkSessionReq, params EditTalkSessionParams) (r EditTalkSessionRes, _ error) {
	return r, ht.ErrNotImplemented
}

// EditTimeLine implements editTimeLine operation.
//
// タイムライン編集.
//
// PUT /talksessions/{talkSessionID}/timelines/{actionItemID}
func (UnimplementedHandler) EditTimeLine(ctx context.Context, req OptEditTimeLineReq, params EditTimeLineParams) (r EditTimeLineRes, _ error) {
	return r, ht.ErrNotImplemented
}

// EditUserProfile implements editUserProfile operation.
//
// ユーザー情報の変更.
//
// PUT /user
func (UnimplementedHandler) EditUserProfile(ctx context.Context, req OptEditUserProfileReq) (r EditUserProfileRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetConclusion implements getConclusion operation.
//
// 結論取得.
//
// GET /talksessions/{talkSessionID}/conclusion
func (UnimplementedHandler) GetConclusion(ctx context.Context, params GetConclusionParams) (r GetConclusionRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetOpenedTalkSession implements getOpenedTalkSession operation.
//
// 自分が開いたセッション一覧.
//
// GET /talksessions/opened
func (UnimplementedHandler) GetOpenedTalkSession(ctx context.Context, params GetOpenedTalkSessionParams) (r GetOpenedTalkSessionRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetOpinionAnalysis implements getOpinionAnalysis operation.
//
// 意見に投票したグループごとの割合.
//
// GET /opinions/{opinionID}/analysis
func (UnimplementedHandler) GetOpinionAnalysis(ctx context.Context, params GetOpinionAnalysisParams) (r GetOpinionAnalysisRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetOpinionDetail implements getOpinionDetail operation.
//
// 意見の詳細.
//
// GET /talksessions/{talkSessionID}/opinions/{opinionID}
func (UnimplementedHandler) GetOpinionDetail(ctx context.Context, params GetOpinionDetailParams) (r GetOpinionDetailRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetOpinionDetail2 implements getOpinionDetail2 operation.
//
// 意見詳細.
//
// GET /opinions/{opinionID}
func (UnimplementedHandler) GetOpinionDetail2(ctx context.Context, params GetOpinionDetail2Params) (r GetOpinionDetail2Res, _ error) {
	return r, ht.ErrNotImplemented
}

// GetOpinionReportReasons implements getOpinionReportReasons operation.
//
// 意見への通報理由一覧.
//
// GET /opinions/report_reasons
func (UnimplementedHandler) GetOpinionReportReasons(ctx context.Context) (r GetOpinionReportReasonsRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetOpinionsForTalkSession implements getOpinionsForTalkSession operation.
//
// セッションに対する意見一覧.
//
// GET /talksessions/{talkSessionID}/opinions
func (UnimplementedHandler) GetOpinionsForTalkSession(ctx context.Context, params GetOpinionsForTalkSessionParams) (r GetOpinionsForTalkSessionRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetPolicyConsentStatus implements getPolicyConsentStatus operation.
//
// 最新のポリシーに同意したかを取得.
//
// GET /policy/consent
func (UnimplementedHandler) GetPolicyConsentStatus(ctx context.Context) (r GetPolicyConsentStatusRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetReportsForTalkSession implements getReportsForTalkSession operation.
//
// 通報一覧.
//
// GET /talksessions/{talkSessionID}/reports
func (UnimplementedHandler) GetReportsForTalkSession(ctx context.Context, params GetReportsForTalkSessionParams) (r GetReportsForTalkSessionRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetTalkSessionDetail implements getTalkSessionDetail operation.
//
// トークセッションの詳細.
//
// GET /talksessions/{talkSessionID}
func (UnimplementedHandler) GetTalkSessionDetail(ctx context.Context, params GetTalkSessionDetailParams) (r GetTalkSessionDetailRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetTalkSessionList implements getTalkSessionList operation.
//
// セッション一覧.
//
// GET /talksessions
func (UnimplementedHandler) GetTalkSessionList(ctx context.Context, params GetTalkSessionListParams) (r GetTalkSessionListRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetTalkSessionReport implements getTalkSessionReport operation.
//
// セッションレポートを返す.
//
// GET /talksessions/{talkSessionID}/report
func (UnimplementedHandler) GetTalkSessionReport(ctx context.Context, params GetTalkSessionReportParams) (r GetTalkSessionReportRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetTalkSessionRestrictionKeys implements getTalkSessionRestrictionKeys operation.
//
// セッションの投稿制限に使用できるキーの一覧を返す.
//
// GET /talksessions/restrictions
func (UnimplementedHandler) GetTalkSessionRestrictionKeys(ctx context.Context) (r GetTalkSessionRestrictionKeysRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetTalkSessionRestrictionSatisfied implements getTalkSessionRestrictionSatisfied operation.
//
// 特定のセッションで満たしていない条件があれば返す.
//
// GET /talksessions/{talkSessionID}/restrictions
func (UnimplementedHandler) GetTalkSessionRestrictionSatisfied(ctx context.Context, params GetTalkSessionRestrictionSatisfiedParams) (r GetTalkSessionRestrictionSatisfiedRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetTimeLine implements getTimeLine operation.
//
// タイムラインはセッション終了後にセッション作成者が設定できるその後の予定を知らせるもの.
//
// GET /talksessions/{talkSessionID}/timelines
func (UnimplementedHandler) GetTimeLine(ctx context.Context, params GetTimeLineParams) (r GetTimeLineRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetUserInfo implements get_user_info operation.
//
// ユーザー情報の取得.
//
// GET /user
func (UnimplementedHandler) GetUserInfo(ctx context.Context) (r GetUserInfoRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ManageIndex implements manageIndex operation.
//
// GET /manage
func (UnimplementedHandler) ManageIndex(ctx context.Context) (r ManageIndexOK, _ error) {
	return r, ht.ErrNotImplemented
}

// ManageRegenerate implements manageRegenerate operation.
//
// Analysisを再生成する。enum: [report, group, image].
//
// POST /manage/regenerate
func (UnimplementedHandler) ManageRegenerate(ctx context.Context, req OptManageRegenerateReq) (r *ManageRegenerateOK, _ error) {
	return r, ht.ErrNotImplemented
}

// OAuthCallback implements oauth_callback operation.
//
// Auth Callback.
//
// GET /auth/{provider}/callback
func (UnimplementedHandler) OAuthCallback(ctx context.Context, params OAuthCallbackParams) (r OAuthCallbackRes, _ error) {
	return r, ht.ErrNotImplemented
}

// OAuthTokenInfo implements oauth_token_info operation.
//
// JWTの内容を返してくれる.
//
// GET /auth/token/info
func (UnimplementedHandler) OAuthTokenInfo(ctx context.Context) (r OAuthTokenInfoRes, _ error) {
	return r, ht.ErrNotImplemented
}

// OAuthTokenRevoke implements oauth_token_revoke operation.
//
// トークンを失効（ログアウト）.
//
// POST /auth/revoke
func (UnimplementedHandler) OAuthTokenRevoke(ctx context.Context) (r OAuthTokenRevokeRes, _ error) {
	return r, ht.ErrNotImplemented
}

// OpinionComments implements opinionComments operation.
//
// 意見に対するリプライ意見一覧.
//
// GET /talksessions/{talkSessionID}/opinions/{opinionID}/replies
func (UnimplementedHandler) OpinionComments(ctx context.Context, params OpinionCommentsParams) (r OpinionCommentsRes, _ error) {
	return r, ht.ErrNotImplemented
}

// OpinionComments2 implements opinionComments2 operation.
//
// 意見に対するリプライ意見一覧.
//
// GET /opinions/{opinionID}/replies
func (UnimplementedHandler) OpinionComments2(ctx context.Context, params OpinionComments2Params) (r OpinionComments2Res, _ error) {
	return r, ht.ErrNotImplemented
}

// OpinionsHistory implements opinionsHistory operation.
//
// 今までに投稿した異見.
//
// GET /opinions/histories
func (UnimplementedHandler) OpinionsHistory(ctx context.Context, params OpinionsHistoryParams) (r OpinionsHistoryRes, _ error) {
	return r, ht.ErrNotImplemented
}

// PolicyConsent implements policyConsent operation.
//
// 最新のポリシーに同意する.
//
// POST /policy/consent
func (UnimplementedHandler) PolicyConsent(ctx context.Context, req OptPolicyConsentReq) (r PolicyConsentRes, _ error) {
	return r, ht.ErrNotImplemented
}

// PostConclusion implements postConclusion operation.
//
// 結論（conclusion）はセッションが終了した後にセッっションの作成者が投稿できる文章。
// セッションの流れやグループの分かれ方などに対するセッション作成者の感想やそれらの意見を受け、これからの方向性などを記入する。.
//
// POST /talksessions/{talkSessionID}/conclusion
func (UnimplementedHandler) PostConclusion(ctx context.Context, req OptPostConclusionReq, params PostConclusionParams) (r PostConclusionRes, _ error) {
	return r, ht.ErrNotImplemented
}

// PostImage implements postImage operation.
//
// 画像を投稿してURLを返すAPI.
//
// POST /images
func (UnimplementedHandler) PostImage(ctx context.Context, req OptPostImageReq) (r PostImageRes, _ error) {
	return r, ht.ErrNotImplemented
}

// PostOpinionPost implements postOpinionPost operation.
//
// ParentOpinionIDがなければルートの意見として投稿される.
//
// POST /talksessions/{talkSessionID}/opinions
func (UnimplementedHandler) PostOpinionPost(ctx context.Context, req OptPostOpinionPostReq, params PostOpinionPostParams) (r PostOpinionPostRes, _ error) {
	return r, ht.ErrNotImplemented
}

// PostOpinionPost2 implements postOpinionPost2 operation.
//
// ParentOpinionIDがなければルートの意見として投稿される
// parentOpinionIDがない場合はtalkSessionIDが必須.
//
// POST /opinions
func (UnimplementedHandler) PostOpinionPost2(ctx context.Context, req OptPostOpinionPost2Req) (r PostOpinionPost2Res, _ error) {
	return r, ht.ErrNotImplemented
}

// PostTimeLineItem implements postTimeLineItem operation.
//
// タイムラインアイテム追加.
//
// POST /talksessions/{talkSessionID}/timeline
func (UnimplementedHandler) PostTimeLineItem(ctx context.Context, req OptPostTimeLineItemReq, params PostTimeLineItemParams) (r PostTimeLineItemRes, _ error) {
	return r, ht.ErrNotImplemented
}

// RegisterUser implements registerUser operation.
//
// ユーザー作成.
//
// POST /user
func (UnimplementedHandler) RegisterUser(ctx context.Context, req OptRegisterUserReq) (r RegisterUserRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ReportOpinion implements reportOpinion operation.
//
// 意見通報API.
//
// POST /opinions/{opinionID}/report
func (UnimplementedHandler) ReportOpinion(ctx context.Context, req OptReportOpinionReq, params ReportOpinionParams) (r ReportOpinionRes, _ error) {
	return r, ht.ErrNotImplemented
}

// SessionsHistory implements sessionsHistory operation.
//
// リアクション済みのセッション一覧.
//
// GET /talksessions/histories
func (UnimplementedHandler) SessionsHistory(ctx context.Context, params SessionsHistoryParams) (r SessionsHistoryRes, _ error) {
	return r, ht.ErrNotImplemented
}

// SwipeOpinions implements swipe_opinions operation.
//
// セッションの中からまだ投票していない意見をランダムに取得する
// remainingCountは取得した意見を含めてスワイプできる意見の総数を返す.
//
// GET /talksessions/{talkSessionID}/swipe_opinions
func (UnimplementedHandler) SwipeOpinions(ctx context.Context, params SwipeOpinionsParams) (r SwipeOpinionsRes, _ error) {
	return r, ht.ErrNotImplemented
}

// TalkSessionAnalysis implements talkSessionAnalysis operation.
//
// 分析結果一覧.
//
// GET /talksessions/{talkSessionID}/analysis
func (UnimplementedHandler) TalkSessionAnalysis(ctx context.Context, params TalkSessionAnalysisParams) (r TalkSessionAnalysisRes, _ error) {
	return r, ht.ErrNotImplemented
}

// Test implements test operation.
//
// OpenAPIテスト用.
//
// GET /test
func (UnimplementedHandler) Test(ctx context.Context) (r TestRes, _ error) {
	return r, ht.ErrNotImplemented
}

// Vote implements vote operation.
//
// 意思表明API.
//
// POST /talksessions/{talkSessionID}/opinions/{opinionID}/votes
func (UnimplementedHandler) Vote(ctx context.Context, req OptVoteReq, params VoteParams) (r VoteRes, _ error) {
	return r, ht.ErrNotImplemented
}

// Vote2 implements vote2 operation.
//
// 意思表明API.
//
// POST /opinions/{opinionID}/votes
func (UnimplementedHandler) Vote2(ctx context.Context, req OptVote2Req, params Vote2Params) (r Vote2Res, _ error) {
	return r, ht.ErrNotImplemented
}
