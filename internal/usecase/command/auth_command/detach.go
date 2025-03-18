package auth_command

import (
	"context"

	"github.com/google/uuid"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/server/internal/infrastructure/persistence/sqlc/generated"
	"go.opentelemetry.io/otel"
)

// 開発向け。退会処理を作るまでの代替。Subjectを付け替えることで、一度SSOしても再度SSOさせることができるやつ。
type DetachAccount interface {
	Execute(ctx context.Context, input DetachAccountInput) error
}

type DetachAccountInput struct {
	UserID shared.UUID[user.User]
}

type detachAccountInteractor struct {
	userService user.UserService
	dbm         *db.DBManager
}

func NewDetachAccount(
	userService user.UserService,
	dbm *db.DBManager,
) DetachAccount {
	return &detachAccountInteractor{
		userService: userService,
		dbm:         dbm,
	}
}

// Execute
func (d *detachAccountInteractor) Execute(ctx context.Context, input DetachAccountInput) error {
	ctx, span := otel.Tracer("auth_command").Start(ctx, "detachAccountInteractor.Execute")
	defer span.End()

	d.dbm.GetQueries(ctx).ChangeSubject(ctx, model.ChangeSubjectParams{
		UserID:  input.UserID.UUID(),
		Subject: uuid.New().String(),
	})

	return nil
}
