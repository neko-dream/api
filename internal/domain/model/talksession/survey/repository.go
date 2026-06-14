package survey

import (
	"context"

	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/talksession"
	"github.com/neko-dream/api/internal/domain/model/user"
)

// Repository は Survey 集約の永続化窓口。
type Repository interface {
	// Create は新規 Survey と質問・選択肢をまとめて保存する。
	Create(ctx context.Context, s *Survey) error

	// Save は Survey の編集結果（質問の追加・更新・soft delete、選択肢の追加・更新・soft delete）を反映する。
	// Save は idempotent。choice/question の hard delete は別 API。
	Save(ctx context.Context, s *Survey) error

	// FindBySurveyID はsoft-deleted 含めて取得する。
	FindBySurveyID(ctx context.Context, surveyID shared.UUID[Survey]) (*Survey, error)

	// FindBySessionID は talk_session_id から逆引きする。存在しない場合は (nil, nil)。
	FindBySessionID(ctx context.Context, sessionID shared.UUID[talksession.TalkSession]) (*Survey, error)

	// HardDeleteQuestion は回答ゼロの質問のみ hard delete 可。回答があると error。
	HardDeleteQuestion(ctx context.Context, questionID shared.UUID[Question]) error

	// HardDeleteChoice は選択ゼロの選択肢のみ hard delete 可。
	HardDeleteChoice(ctx context.Context, choiceID shared.UUID[Choice]) error

	// HasAnswersForQuestion は質問に対する回答が 1 件でも存在するかを返す。
	HasAnswersForQuestion(ctx context.Context, questionID shared.UUID[Question]) (bool, error)

	// HasAnswersForChoice は選択肢が選択された回答が 1 件でも存在するかを返す。
	HasAnswersForChoice(ctx context.Context, choiceID shared.UUID[Choice]) (bool, error)

	// CreateSubmission は1ユーザー1送信を作成する。重複は UNIQUE 違反として error。
	CreateSubmission(ctx context.Context, sub *Submission) error

	// HasSubmittedBySessionID はユーザーが該当セッションのアンケートに回答済みか返す。
	// アンケートが存在しないセッションについては (false, nil)。
	HasSubmittedBySessionID(
		ctx context.Context,
		sessionID shared.UUID[talksession.TalkSession],
		userID shared.UUID[user.User],
	) (bool, error)

	// ListSubmissions は survey の回答一覧を pagination 付きで返す。
	ListSubmissions(
		ctx context.Context,
		surveyID shared.UUID[Survey],
		limit, offset int,
	) ([]*Submission, int32, error)
}
