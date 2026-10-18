package opinion_query

import (
	"context"
	"time"

	"github.com/neko-dream/api/internal/domain/model/opinion"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/talksession"
	"github.com/neko-dream/api/internal/domain/model/user"
)

// ManageOpinions は管理画面用に、セッションの意見を削除済み・返信も含めて読む。
// 投票は意見ごとの集計だけを返し、誰がどう投票したかは返さない。
type ManageOpinions interface {
	// List はセッションが存在しなければ messages.TalkSessionNotFound を返す。
	List(context.Context, ManageOpinionsInput) (*ManageOpinionsOutput, error)
}

type ManageOpinionsInput struct {
	TalkSessionID shared.UUID[talksession.TalkSession]
	Limit, Offset int32
}

type ManageOpinionsOutput struct {
	TotalCount int32
	Opinions   []ManagedOpinion
}

type ManagedOpinion struct {
	OpinionID       shared.UUID[opinion.Opinion]
	ParentOpinionID *shared.UUID[opinion.Opinion]
	Content         string
	PictureURL      *string
	ReferenceURL    *string
	CreatedAt       time.Time
	Author          ManagedOpinionAuthor
	AgreeCount      int32
	DisagreeCount   int32
	PassCount       int32
	// ReportCount は通報の件数（すべての状態）、UnsolvedReportCount はそのうち未対応の件数。
	ReportCount         int32
	UnsolvedReportCount int32
	// Deleted は通報の対応で削除されたか。削除済みの通報が1件でもあると、一般の意見一覧・返信・詳細では
	// 「運営により削除されました」と表示される（ProcessReportedOpinions と同じ判定）。
	Deleted bool
}

// ManagedOpinionAuthor は投稿者。opinions.user_id に外部キーはないため、ユーザーがいなければ UserID は nil。
type ManagedOpinionAuthor struct {
	UserID      *shared.UUID[user.User]
	DisplayID   string
	DisplayName string
	IconURL     string
}
