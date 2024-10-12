package handler

import (
	"context"

	"github.com/neko-dream/server/internal/presentation/oas"
)

type opinionHandler struct {
}

func NewOpinionHandler() oas.OpinionHandler {
	return &opinionHandler{}
}

// GetTopOpinions implements oas.OpinionHandler.
func (o *opinionHandler) GetTopOpinions(ctx context.Context, params oas.GetTopOpinionsParams) (oas.GetTopOpinionsRes, error) {
	panic("unimplemented")
}

// ListOpinions implements oas.OpinionHandler.
func (o *opinionHandler) ListOpinions(ctx context.Context, params oas.ListOpinionsParams) (oas.ListOpinionsRes, error) {
	panic("unimplemented")
}

// OpinionComments implements oas.OpinionHandler.
func (o *opinionHandler) OpinionComments(ctx context.Context, params oas.OpinionCommentsParams) (oas.OpinionCommentsRes, error) {
	panic("unimplemented")
}

// PostOpinionPost implements oas.OpinionHandler.
func (o *opinionHandler) PostOpinionPost(ctx context.Context, req oas.OptPostOpinionPostReq, params oas.PostOpinionPostParams) (oas.PostOpinionPostRes, error) {
	panic("unimplemented")
}
