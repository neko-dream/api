package user_query

import (
	"context"
	"time"

	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
)

// Management reads a restricted projection for operators, not public user profiles.
type Management interface {
	List(context.Context, ManagementListInput) (*ManagementListOutput, error)
	FindByID(context.Context, shared.UUID[user.User]) (*ManagedUser, error)
}
type ManagementListInput struct {
	Offset, Limit int32
	Search        string
}
type ManagementListOutput struct {
	Users      []ManagedUser
	TotalCount int64
}
type ManagedUser struct {
	UserID                 shared.UUID[user.User]
	DisplayID, DisplayName string
	Email                  *string
	EmailVerified          bool
	Providers              []string
	CreatedAt, UpdatedAt   time.Time
	WithdrawalDate         *time.Time
}
