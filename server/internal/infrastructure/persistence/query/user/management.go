package user_query

import (
	"context"
	"database/sql"
	"errors"

	userquery "github.com/neko-dream/api/internal/application/query/user"
	"github.com/neko-dream/api/internal/domain/model/crypto"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/api/internal/infrastructure/persistence/sqlc/generated"
	"github.com/neko-dream/api/pkg/utils"
	"go.opentelemetry.io/otel"
)

type managementQuery struct {
	db        *db.DBManager
	encryptor crypto.Encryptor
}

func NewManagementQuery(db *db.DBManager, encryptor crypto.Encryptor) userquery.Management {
	return &managementQuery{db: db, encryptor: encryptor}
}

func (q *managementQuery) List(ctx context.Context, input userquery.ManagementListInput) (*userquery.ManagementListOutput, error) {
	ctx, span := otel.Tracer("user_query").Start(ctx, "managementQuery.List")
	defer span.End()
	queries := q.db.GetQueries(ctx)
	total, err := queries.CountManagedUsers(ctx, input.Search)
	if err != nil {
		return nil, err
	}
	rows, err := queries.ListManagedUsers(ctx, model.ListManagedUsersParams{Search: input.Search, PageOffset: input.Offset, PageLimit: input.Limit})
	if err != nil {
		return nil, err
	}
	users := make([]userquery.ManagedUser, 0, len(rows))
	for _, row := range rows {
		u, err := q.mapUser(ctx, row.User, row.Providers)
		if err != nil {
			return nil, err
		}
		users = append(users, *u)
	}
	return &userquery.ManagementListOutput{Users: users, TotalCount: total}, nil
}

func (q *managementQuery) FindByID(ctx context.Context, id shared.UUID[user.User]) (*userquery.ManagedUser, error) {
	ctx, span := otel.Tracer("user_query").Start(ctx, "managementQuery.FindByID")
	defer span.End()
	row, err := q.db.GetQueries(ctx).GetManagedUser(ctx, id.UUID())
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return q.mapUser(ctx, row.User, row.Providers)
}

func (q *managementQuery) mapUser(ctx context.Context, row model.User, providers []string) (*userquery.ManagedUser, error) {
	u := &userquery.ManagedUser{
		CanCreateTalkSession: row.CanCreateTalkSession,
		UserID:               shared.UUID[user.User](row.UserID),
		DisplayID:            row.DisplayID.String,
		DisplayName:          row.DisplayName.String,
		EmailVerified:        row.EmailVerified,
		Providers:            providers,
		CreatedAt:            row.CreatedAt,
		UpdatedAt:            row.UpdatedAt,
	}
	if u.Providers == nil {
		u.Providers = []string{}
	}
	if row.Email.Valid {
		email, err := q.encryptor.DecryptString(ctx, row.Email.String)
		if err != nil {
			utils.HandleError(ctx, err, "メアドの復号に失敗しました。")
			u.EmailVerified = false
		} else {
			u.Email = &email
		}
	}
	if row.WithdrawalDate.Valid {
		u.WithdrawalDate = &row.WithdrawalDate.Time
	}
	return u, nil
}
