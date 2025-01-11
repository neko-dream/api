package queryimpl

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jinzhu/copier"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/server/internal/infrastructure/persistence/sqlc/generated"
	"github.com/neko-dream/server/internal/usecase/query/dto"
	"github.com/neko-dream/server/internal/usecase/query/talksession"
	"github.com/neko-dream/server/pkg/utils"
)

type BrowseTalkSessionQueryImpl struct {
	*db.DBManager
}

func NewBrowseTalkSessionQuery(
	tm *db.DBManager,
) talksession.BrowseTalkSessionQuery {
	return &BrowseTalkSessionQueryImpl{
		DBManager: tm,
	}
}

func (b *BrowseTalkSessionQueryImpl) Execute(ctx context.Context, in talksession.BrowseTalkSessionQueryInput) (*talksession.BrowseTalkSessionQueryOutput, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	var latitude, longitude sql.NullFloat64
	if in.Latitude != nil {
		latitude = sql.NullFloat64{Float64: *in.Latitude, Valid: true}
	}
	if in.Longitude != nil {
		longitude = sql.NullFloat64{Float64: *in.Longitude, Valid: true}
	}
	theme := utils.IfThenElse(
		in.Theme != nil,
		sql.NullString{String: *in.Theme, Valid: true},
		sql.NullString{},
	)

	talkSessionRow, err := b.GetQueries(ctx).ListTalkSessions(ctx, model.ListTalkSessionsParams{
		Limit:     int32(in.Limit),
		Offset:    int32(in.Offset),
		Theme:     theme,
		Status:    sql.NullString{String: in.Status, Valid: true},
		SortKey:   sql.NullString{String: string(*in.SortKey), Valid: true},
		Latitude:  latitude,
		Longitude: longitude,
	})

	var out talksession.BrowseTalkSessionQueryOutput
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			out.TalkSessions = make([]dto.TalkSessionWithDetail, 0)
			return &out, nil
		}
		utils.HandleError(ctx, err, "failed to list talk sessions")
		return nil, err
	}
	if len(talkSessionRow) <= 0 {
		out.TalkSessions = make([]dto.TalkSessionWithDetail, 0)
		return &out, nil
	}

	var talkSessions []dto.TalkSessionWithDetail
	if err := copier.CopyWithOption(&talkSessions, talkSessionRow, copier.Option{
		DeepCopy: true,
	}); err != nil {
		utils.HandleError(ctx, err, "failed to copy talk session")
		return nil, err
	}

	talkSessionCount, err := b.GetQueries(ctx).CountTalkSessions(ctx, model.CountTalkSessionsParams{
		Theme:  theme,
		Status: sql.NullString{String: in.Status, Valid: true},
	})
	if err != nil {
		utils.HandleError(ctx, err, "failed to count talk sessions")
		return nil, err
	}

	out.TalkSessions = talkSessions
	out.TotalCount = int(talkSessionCount.TalkSessionCount)

	return &out, nil
}
