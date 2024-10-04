package talk_session_usecase

import (
	"context"

	talksession "github.com/neko-dream/server/internal/domain/model/talk_session"
	"github.com/neko-dream/server/internal/domain/model/user"
)

type (
	ListTalkSessionQuery interface {
		Execute(context.Context, ListTalkSessionInput) (ListTalkSessionOutput, error)
	}

	ListTalkSessionInput struct {
	}

	ListTalkSessionOutput struct {
		TalkSessions []TalkSessionDTO
		OpinionCount int
	}

	TalkSessionDTO struct {
		ID    string
		Theme string
		Owner user.User
	}

	listTalkSessionQueryHandler struct {
		talkSessionRepository talksession.TalkSessionRepository
		userRepository        user.UserRepository
	}
)

func NewListTalkSessionQueryHandler(tr talksession.TalkSessionRepository, ur user.UserRepository) ListTalkSessionQuery {
	return &listTalkSessionQueryHandler{
		talkSessionRepository: tr,
		userRepository:        ur,
	}
}

func (h *listTalkSessionQueryHandler) Execute(ctx context.Context, input ListTalkSessionInput) (ListTalkSessionOutput, error) {
	return ListTalkSessionOutput{}, nil
}
