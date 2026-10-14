package user_query

import (
	"context"
	"database/sql"
	"errors"

	"github.com/neko-dream/api/internal/application/query/dto"
	userquery "github.com/neko-dream/api/internal/application/query/user"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
	"go.opentelemetry.io/otel"
)

type currentQuery struct{ db *db.DBManager }

func NewCurrentQuery(dbm *db.DBManager) userquery.Current { return &currentQuery{db: dbm} }

func (q *currentQuery) FindByID(ctx context.Context, id shared.UUID[user.User]) (*userquery.CurrentOutput, error) {
	ctx, span := otel.Tracer("user_query").Start(ctx, "currentQuery.FindByID")
	defer span.End()
	r, err := q.db.GetQueries(ctx).GetCurrentUser(ctx, id.UUID())
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	out := &userquery.CurrentOutput{
		User:                 dto.User{DisplayID: r.DisplayID.String, DisplayName: r.DisplayName.String},
		IsGuest:              r.IsGuest,
		CanCreateTalkSession: user.CanCreateTalkSession(r.CanCreateTalkSession, r.IsGuest, r.WithdrawalDate.Valid),
	}
	if r.IconUrl.Valid {
		out.User.IconURL = &r.IconUrl.String
	}
	if r.WithdrawalDate.Valid {
		out.User.WithdrawalDate = &r.WithdrawalDate.Time
	}
	return out, nil
}
