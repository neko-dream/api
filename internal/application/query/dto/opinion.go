package dto

import (
	"time"

	"github.com/neko-dream/server/internal/domain/model/opinion"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/talksession"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/domain/model/vote"
	model "github.com/neko-dream/server/internal/infrastructure/persistence/sqlc/generated"
	"github.com/neko-dream/server/internal/presentation/oas"
	"github.com/neko-dream/server/pkg/utils"
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
	IsDeleted       bool
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

// SwipeOpinionも、Replaceされる
func (s *SwipeOpinion) Mask(reports []model.OpinionReport) {
	s.User = User{
		DisplayID:   "",
		DisplayName: "",
		IconURL:     nil,
	}
	s.ParentVoteType = 0
	s.CurrentVoteType = 0
	s.Opinion.Mask(reports)
}
func (s *Opinion) Mask(reports []model.OpinionReport) {
	if len(reports) == 0 {
		return
	}
	report := "この意見は運営により削除されました。\n削除理由:\n"
	for _, r := range reports {
		reason := opinion.Reason(r.Reason)
		report += "・" + reason.StringJP() + "\n"
	}
	s.UserID = shared.UUID[user.User](shared.NilUUID)
	s.Content = report
	s.PictureURL = nil
	s.ReferenceURL = nil
	s.Title = nil
	s.IsDeleted = true
}

type OpinionWithRepresentative struct {
	Opinion
	User
	RepresentativeOpinion
	ReplyCount int
}

func (o *OpinionWithRepresentative) Mask(reports []model.OpinionReport) {
	o.User = User{
		DisplayID:   "",
		DisplayName: "",
		IconURL:     nil,
	}
	o.RepresentativeOpinion = RepresentativeOpinion{}
	o.ReplyCount = 0
	o.Opinion.Mask(reports)
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

func (o *Opinion) ToResponse() oas.Opinion {
	var parentID oas.OptString
	if o.ParentOpinionID != nil {
		parentID = oas.OptString{
			Value: o.ParentOpinionID.String(),
			Set:   true,
		}
	}

	return oas.Opinion{
		ID:           o.OpinionID.String(),
		Title:        utils.ToOpt[oas.OptString](o.Title),
		Content:      o.Content,
		ParentID:     parentID,
		PictureURL:   utils.ToOptNil[oas.OptNilString](o.PictureURL),
		ReferenceURL: utils.ToOpt[oas.OptString](o.ReferenceURL),
		PostedAt:     o.CreatedAt.Format(time.RFC3339),
		IsDeleted:    o.IsDeleted,
	}
}
