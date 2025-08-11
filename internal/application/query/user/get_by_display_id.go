package user_query

import (
	"context"

	"github.com/neko-dream/server/internal/application/query/dto"
)

type (
	GetByDisplayID interface {
		Execute(context.Context, GetByDisplayIDInput) (*GetByDisplayIDOutput, error)
	}

	GetByDisplayIDInput struct {
		DisplayID string
	}

	GetByDisplayIDOutput struct {
		User *dto.User
	}
)
