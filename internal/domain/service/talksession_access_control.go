package service

import (
	"context"
	"strings"

	"github.com/neko-dream/server/internal/domain/messages"
	talksession "github.com/neko-dream/server/internal/domain/model/talk_session"
	"github.com/neko-dream/server/internal/domain/model/user"
)

type TalkSessionAccessControl interface {
	// CanUserJoin はユーザーがトークセッションに参加できるかを判定する
	CanUserJoin(ctx context.Context, talkSession *talksession.TalkSession, user *user.User) (bool, error)
}

type talkSessionAccessControl struct {
}

func NewTalkSessionAccessControl() TalkSessionAccessControl {
	return &talkSessionAccessControl{}
}

var (
	// 制限が満たされていない場合のエラーメッセージ
	ErrRestrictionNotSatisfied = messages.APIError{
		Code:       "restriction_not_satisfied",
		StatusCode: 400,
		Message:    "参加条件が満たされていません",
	}
)

func (t *talkSessionAccessControl) CanUserJoin(ctx context.Context, talkSession *talksession.TalkSession, user *user.User) (bool, error) {
	if talkSession == nil {
		return false, &ErrRestrictionNotSatisfied
	}

	// あとあとユーザーが設定されていなくても参加できるような仕様になるかもしれないが、現状はユーザーが設定されていない場合は参加不可とする
	if user == nil {
		return false, &ErrRestrictionNotSatisfied
	}

	// 参加制限がない場合は参加可能
	if len(talkSession.Restrictions()) == 0 {
		return true, nil
	}

	// 参加制限がある場合は、ユーザーが参加可能かを判定し、もし参加制限に引っかかる場合はエラーを返す
	var restrictions []talksession.RestrictionAttribute
	for _, restriction := range talkSession.Restrictions() {
		if !restriction.Fn(*user) {
			restrictions = append(restrictions, *restriction)
		}
	}

	if len(restrictions) > 0 {
		// 必要な項目を,で結合し、エラーメッセージを作成
		var restrictionKeys []string
		for _, restriction := range restrictions {
			restrictionKeys = append(restrictionKeys, string(restriction.Key))
		}

		e := ErrRestrictionNotSatisfied
		e.Message = "このセッションでは、" + strings.Join(restrictionKeys, ",") + "が必要です。"
		return false, &e
	}

	return true, nil
}
