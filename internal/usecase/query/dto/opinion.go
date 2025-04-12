package dto

import (
	"time"

	"github.com/neko-dream/server/internal/domain/model/opinion"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/talksession"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/domain/model/vote"
	model "github.com/neko-dream/server/internal/infrastructure/persistence/sqlc/generated"
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

func (s *Opinion) ReplaceReported(reports []model.OpinionReport) {
	if len(reports) == 0 {
		return
	}
	report := "この意見は運営により削除されました。\n削除理由:\n"
	for _, r := range reports {
		reason := opinion.Reason(r.Reason)
		report += "・" + reason.StringJP() + "\n"
	}
	s.Content = report
	s.PictureURL = nil
	s.ReferenceURL = nil
	s.Title = nil
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
