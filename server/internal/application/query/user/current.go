package user_query

import (
	"context"

	"github.com/neko-dream/api/internal/application/query/dto"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
)

// Current reads only the fields needed by the signed-in UI, without private demographics.
type Current interface {
	FindByID(context.Context, shared.UUID[user.User]) (*CurrentOutput, error)
}

type CurrentOutput struct {
	User                 dto.User
	IsGuest              bool
	CanCreateTalkSession bool
}
