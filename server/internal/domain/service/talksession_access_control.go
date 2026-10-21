package service

import (
	"context"
	"strings"

	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/talksession"
	"github.com/neko-dream/api/internal/domain/model/talksession/survey"
	"github.com/neko-dream/api/internal/domain/model/talksession/talksession_consent"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/neko-dream/api/pkg/utils"
	"go.opentelemetry.io/otel"
)

type TalkSessionAccessControl interface {
	// CanUserJoin はユーザーがトークセッションに参加できるかを判定する
	CanUserJoin(ctx context.Context, talkSessionID shared.UUID[talksession.TalkSession], userID *shared.UUID[user.User]) (bool, error)
	UserSatisfiesRestriction(ctx context.Context, talkSessionID shared.UUID[talksession.TalkSession], userID shared.UUID[user.User]) ([]talksession.RestrictionAttribute, error)
}

type talkSessionAccessControl struct {
	talksession.TalkSessionRepository
	user.UserRepository
	talksession_consent.TalkSessionConsentService
	surveyRepository survey.Repository
}

func NewTalkSessionAccessControl(
	talkSessionRepository talksession.TalkSessionRepository,
	userRepository user.UserRepository,
	talkSessionConsentService talksession_consent.TalkSessionConsentService,
	surveyRepository survey.Repository,
) TalkSessionAccessControl {
	return &talkSessionAccessControl{
		TalkSessionRepository:     talkSessionRepository,
		UserRepository:            userRepository,
		TalkSessionConsentService: talkSessionConsentService,
		surveyRepository:          surveyRepository,
	}
}

var (
	// 制限が満たされていない場合のエラーメッセージ
	// 共有のエラー。文言は書き換えず WithMessage で複製して返す
	ErrRestrictionNotSatisfied = &messages.APIError{
		Code:       "restriction_not_satisfied",
		StatusCode: 400,
		Message:    "参加条件が満たされていません",
	}
)

func (t *talkSessionAccessControl) CanUserJoin(ctx context.Context, talkSessionID shared.UUID[talksession.TalkSession], userID *shared.UUID[user.User]) (bool, error) {
	ctx, span := otel.Tracer("service").Start(ctx, "talkSessionAccessControl.CanUserJoin")
	defer span.End()

	// talksessionが存在するか確認
	talkSession, err := t.TalkSessionRepository.FindByID(ctx, talkSessionID)
	if err != nil || talkSession == nil {
		utils.HandleError(ctx, err, "TalkSessionRepository.FindByID")
		return false, messages.TalkSessionNotFound
	}

	// セッション作成者なら参加可能
	if userID != nil && talkSession.OwnerUserID() == *userID {
		return true, nil
	}

	// userの存在確認
	var user *user.User
	if userID != nil {
		user, err = t.UserRepository.FindByID(ctx, *userID)
		if err != nil || user == nil {
			utils.HandleError(ctx, err, "UserRepository.FindByID")
			return false, messages.UserNotFoundError
		}
	}

	// ゲストユーザーかつ制限付きセッションの場合は参加不可
	if user != nil && user.IsGuest() && len(talkSession.Restrictions()) > 0 {
		return false, messages.GuestCannotJoinRestrictedSessionError
	}

	// 未認証ユーザーは制限付きセッションには参加できない（user/userID 参照前のガード）
	if userID == nil && len(talkSession.Restrictions()) > 0 {
		return false, ErrRestrictionNotSatisfied.WithMessage("このセッションは認証が必要です。")
	}

	// 参加制限がある場合は、ユーザーが参加可能かを判定し、もし参加制限に引っかかる場合はエラーを返す
	var restrictions []talksession.RestrictionAttribute
	for _, restriction := range talkSession.Restrictions() {
		if restriction.Key == talksession.SurveyAnswered {
			submitted, err := t.surveyRepository.HasSubmittedBySessionID(ctx, talkSessionID, *userID)
			if err != nil {
				utils.HandleError(ctx, err, "surveyRepository.HasSubmittedBySessionID")
				return false, err
			}
			if !submitted {
				restrictions = append(restrictions, *restriction)
			}
			continue
		}
		if !restriction.IsSatisfied(*user) {
			restrictions = append(restrictions, *restriction)
		}
	}

	if len(restrictions) > 0 {
		// 必要な項目を,で結合し、エラーメッセージを作成
		var restrictionKeys []string
		for _, restriction := range restrictions {
			restrictionKeys = append(restrictionKeys, string(restriction.Key))
		}

		return false, ErrRestrictionNotSatisfied.WithMessage("このセッションでは、" + strings.Join(restrictionKeys, ",") + "が必要です。")
	}

	// 同意ししていなければ参加できない
	consent, err := t.TalkSessionConsentService.HasConsented(ctx, talkSessionID, *userID)
	if err != nil {
		utils.HandleError(ctx, err, "TalkSessionConsentService.HasConsented")
		return false, err
	}

	if !consent {
		return false, ErrRestrictionNotSatisfied.WithMessage("このセッションでは、参加するために同意が必要です。")
	}

	return true, nil
}

// UserSatisfiesRestriction implements TalkSessionAccessControl.
func (t *talkSessionAccessControl) UserSatisfiesRestriction(ctx context.Context, talkSessionID shared.UUID[talksession.TalkSession], userID shared.UUID[user.User]) ([]talksession.RestrictionAttribute, error) {
	ctx, span := otel.Tracer("service").Start(ctx, "talkSessionAccessControl.UserSatisfiesRestriction")
	defer span.End()

	// talksessionが存在するか確認
	talkSession, err := t.TalkSessionRepository.FindByID(ctx, talkSessionID)
	if err != nil || talkSession == nil {
		utils.HandleError(ctx, err, "TalkSessionRepository.FindByID")
		return nil, messages.TalkSessionNotFound
	}

	// userの存在確認
	u, err := t.UserRepository.FindByID(ctx, userID)
	if err != nil || u == nil {
		utils.HandleError(ctx, err, "UserRepository.FindByID")
		return nil, messages.UserNotFoundError
	}

	// セッション作成者なら参加可能
	if talkSession.OwnerUserID() == userID {
		return nil, nil
	}

	// 参加制限がない場合は参加可能
	if len(talkSession.Restrictions()) == 0 {
		return nil, nil
	}

	// 参加制限がある場合は、ユーザーが参加可能かを判定し、もし参加制限に引っかかる場合はエラーを返す
	var restrictions []talksession.RestrictionAttribute
	for _, restriction := range talkSession.Restrictions() {
		if restriction.Key == talksession.SurveyAnswered {
			submitted, err := t.surveyRepository.HasSubmittedBySessionID(ctx, talkSessionID, userID)
			if err != nil {
				utils.HandleError(ctx, err, "surveyRepository.HasSubmittedBySessionID")
				return nil, err
			}
			if !submitted {
				restrictions = append(restrictions, *restriction)
			}
			continue
		}
		if !restriction.IsSatisfied(*u) {
			restrictions = append(restrictions, *restriction)
		}
	}

	return restrictions, nil
}
