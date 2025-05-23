package talksession_query

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jinzhu/copier"
	"github.com/neko-dream/server/internal/application/query/dto"
	"github.com/neko-dream/server/internal/application/query/talksession"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/server/internal/infrastructure/persistence/sqlc/generated"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

type BrowseTalkSessionQueryImpl struct {
	*db.DBManager
}

func NewBrowseTalkSessionQueryHandler(
	tm *db.DBManager,
) talksession.BrowseTalkSessionQuery {
	return &BrowseTalkSessionQueryImpl{
		DBManager: tm,
	}
}

func (b *BrowseTalkSessionQueryImpl) Execute(ctx context.Context, in talksession.BrowseTalkSessionQueryInput) (*talksession.BrowseTalkSessionQueryOutput, error) {
	ctx, span := otel.Tracer("talksession_query").Start(ctx, "BrowseTalkSessionQueryImpl.Execute")
	defer span.End()

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
	var theme sql.NullString
	if in.Theme != nil {
		theme = sql.NullString{String: *in.Theme, Valid: true}
	}

	var status sql.NullString
	if in.Status != nil {
		status = sql.NullString{String: string(*in.Status), Valid: true}
	}

	talkSessionRow, err := b.GetQueries(ctx).ListTalkSessions(ctx, model.ListTalkSessionsParams{
		Limit:     int32(*in.Limit),
		Offset:    int32(*in.Offset),
		Theme:     theme,
		Status:    status,
		SortKey:   sql.NullString{String: in.SortKey.String(), Valid: true},
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
		Status: status,
	})
	if err != nil {
		utils.HandleError(ctx, err, "failed to count talk sessions")
		return nil, err
	}

	out.TalkSessions = talkSessions
	out.TotalCount = int(talkSessionCount.TalkSessionCount)
	out.Limit = *in.Limit
	out.Offset = *in.Offset

	return &out, nil
}
