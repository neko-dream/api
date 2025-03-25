package dto

import (
	"context"

	"github.com/neko-dream/server/internal/domain/model/analysis"
	"go.opentelemetry.io/otel"
)

type UserPosition struct {
	PosX           float64
	PosY           float64
	DisplayID      string
	DisplayName    string
	IconURL        *string
	GroupName      string
	GroupID        int
	PerimeterIndex *int
}

func (u *UserPosition) GetGroupName(ctx context.Context) string {
	ctx, span := otel.Tracer("dto").Start(ctx, "UserPosition.GetGroupName")
	defer span.End()

	return analysis.NewGroupIDFromInt(int(u.GroupID)).String()
}

type OpinionGroup struct {
	GroupName string
	GroupID   int
	Opinions  []OpinionWithRepresentative
}

type OpinionGroupRatio struct {
	GroupName     string
	GroupID       int
	AgreeCount    int
	DisagreeCount int
	PassCount     int
}
