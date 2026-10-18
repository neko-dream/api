package opinion_query

import (
	"context"

	opinionquery "github.com/neko-dream/api/internal/application/query/opinion"
	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/domain/model/opinion"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/api/internal/infrastructure/persistence/sqlc/generated"
	"go.opentelemetry.io/otel"
)

type manageOpinionsQuery struct {
	*db.DBManager
}

func NewManageOpinionsQuery(dbManager *db.DBManager) opinionquery.ManageOpinions {
	return &manageOpinionsQuery{DBManager: dbManager}
}

func (q *manageOpinionsQuery) List(ctx context.Context, input opinionquery.ManageOpinionsInput) (*opinionquery.ManageOpinionsOutput, error) {
	ctx, span := otel.Tracer("opinion_query").Start(ctx, "manageOpinionsQuery.List")
	defer span.End()

	queries := q.GetQueries(ctx)
	exists, err := queries.TalkSessionExists(ctx, input.TalkSessionID.UUID())
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, messages.TalkSessionNotFound
	}
	total, err := queries.CountOpinionsForManage(ctx, input.TalkSessionID.UUID())
	if err != nil {
		return nil, err
	}
	rows, err := queries.ListOpinionsForManage(ctx, model.ListOpinionsForManageParams{
		TalkSessionID: input.TalkSessionID.UUID(),
		RowLimit:      input.Limit,
		RowOffset:     input.Offset,
	})
	if err != nil {
		return nil, err
	}

	opinions := make([]opinionquery.ManagedOpinion, 0, len(rows))
	for _, row := range rows {
		opinions = append(opinions, toManagedOpinion(row))
	}
	return &opinionquery.ManageOpinionsOutput{TotalCount: total, Opinions: opinions}, nil
}

func toManagedOpinion(row model.ListOpinionsForManageRow) opinionquery.ManagedOpinion {
	o := row.Opinion
	res := opinionquery.ManagedOpinion{
		OpinionID: shared.UUID[opinion.Opinion](o.OpinionID),
		Content:   o.Content,
		CreatedAt: o.CreatedAt,
		Author: opinionquery.ManagedOpinionAuthor{
			DisplayID:   row.AuthorDisplayID.String,
			DisplayName: row.AuthorDisplayName.String,
			IconURL:     row.AuthorIconUrl.String,
		},
		AgreeCount:          row.AgreeCount,
		DisagreeCount:       row.DisagreeCount,
		PassCount:           row.PassCount,
		ReportCount:         row.ReportCount,
		UnsolvedReportCount: row.UnsolvedReportCount,
		Deleted:             row.Deleted,
	}
	if row.AuthorUserID.Valid {
		userID := shared.UUID[user.User](row.AuthorUserID.UUID)
		res.Author.UserID = &userID
	}
	if o.ParentOpinionID.Valid {
		parent := shared.UUID[opinion.Opinion](o.ParentOpinionID.UUID)
		res.ParentOpinionID = &parent
	}
	if o.PictureUrl.Valid {
		res.PictureURL = &o.PictureUrl.String
	}
	if o.ReferenceUrl.Valid {
		res.ReferenceURL = &o.ReferenceUrl.String
	}
	return res
}
