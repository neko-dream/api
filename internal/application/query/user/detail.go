package user_query

import (
	"context"

	"github.com/neko-dream/server/internal/application/query/dto"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
)

type (
	Detail interface {
		Execute(context.Context, DetailInput) (*DetailOutput, error)
	}

	DetailInput struct {
		UserID shared.UUID[user.User]
	}

	DetailOutput struct {
		User dto.UserDetail
	}
)
