package dto

import (
	"time"

	"github.com/neko-dream/server/internal/domain/model/opinion"
	"github.com/neko-dream/server/internal/domain/model/shared"
	talksession "github.com/neko-dream/server/internal/domain/model/talk_session"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/domain/model/vote"
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

func (s *SwipeOpinion) GetMyVoteType() string {
	return vote.VoteTypeFromInt(s.CurrentVoteType).String()
}

func (s *SwipeOpinion) GetParentVoteType() string {
	return vote.VoteTypeFromInt(s.ParentVoteType).String()
}
