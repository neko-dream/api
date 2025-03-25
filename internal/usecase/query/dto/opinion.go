package dto

import (
	"time"

	"github.com/neko-dream/server/internal/domain/model/opinion"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/talksession"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/domain/model/vote"
	"github.com/samber/lo"
)

type Opinion struct {
	OpinionID       shared.UUID[opinion.Opinion]
	TalkSessionID   shared.UUID[talksession.TalkSession]
	UserID          shared.UUID[user.User]
	ParentOpinionID *shared.UUID[opinion.Opinion]
	Title           *string
	Content         string
	CreatedAt       time.Time
	PictureURL      *string
	ReferenceURL    *string
}

type SwipeOpinion struct {
	Opinion         Opinion
	User            User
	CurrentVoteType int
	ReplyCount      int
	ParentVoteType  int
}

func (s *SwipeOpinion) GetMyVoteType() *string {
	if s.CurrentVoteType == 0 {
		return nil
	}
	return lo.ToPtr(vote.VoteTypeFromInt(s.CurrentVoteType).String())
}

func (s *SwipeOpinion) GetParentVoteType() *string {
	if s.ParentVoteType == 0 {
		return nil
	}
	return lo.ToPtr(vote.VoteTypeFromInt(s.ParentVoteType).String())
}

type OpinionWithRepresentative struct {
	Opinion
	User
	RepresentativeOpinion
	ReplyCount int
}

type RepresentativeOpinion struct {
	TalkSessionID shared.UUID[talksession.TalkSession]
	OpinionID     shared.UUID[opinion.Opinion]
	GroupID       int
	AgreeCount    int
	DisagreeCount int
	PassCount     int
}

type ReportReason struct {
	ReasonID int
	Reason   string
}
