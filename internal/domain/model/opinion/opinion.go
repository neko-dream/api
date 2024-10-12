package opinion

import (
	"context"
	"os/user"
	"time"

	"github.com/neko-dream/server/internal/domain/model/shared"
	talksession "github.com/neko-dream/server/internal/domain/model/talk_session"
	"github.com/neko-dream/server/internal/domain/model/vote"
)

type (
	OpinionRepository interface {
		Create(context.Context, Opinion) error
		FindByTalkSessionID(context.Context, shared.UUID[talksession.TalkSession]) ([]Opinion, error)
		FindByParentID(context.Context, shared.UUID[Opinion]) ([]Opinion, error)
		// FindByTalkSessionWithoutVote まだユーザーが投票していない意見をランダムに取得
		FindByTalkSessionWithoutVote(
			ctx context.Context,
			userID shared.UUID[user.User],
			talkSessionID shared.UUID[talksession.TalkSession],
			limit int,
		) ([]Opinion, error)
	}

	Opinion struct {
		opinionID       shared.UUID[Opinion]
		talkSessionID   shared.UUID[talksession.TalkSession]
		userID          shared.UUID[user.User]
		parentOpinionID *shared.UUID[Opinion]
		content         string
		createdAt       time.Time
		opinions        []Opinion
		voteStatus      *vote.VoteStatus
	}
)

func NewOpinion(
	opinionID shared.UUID[Opinion],
	talkSessionID shared.UUID[talksession.TalkSession],
	userID shared.UUID[user.User],
	parentOpinionID *shared.UUID[Opinion],
	content string,
	createdAt time.Time,
	voteStatus *vote.VoteStatus,
) Opinion {
	return Opinion{
		opinionID:       opinionID,
		talkSessionID:   talkSessionID,
		userID:          userID,
		parentOpinionID: parentOpinionID,
		content:         content,
		createdAt:       createdAt,
		voteStatus:      voteStatus,
		opinions:        []Opinion{},
	}
}

func (o *Opinion) Reply(opinion Opinion) {
	o.opinions = append(o.opinions, opinion)
}

func (o *Opinion) Count() int {
	return len(o.opinions)
}

func (o *Opinion) IsVoted() bool {
	return o.voteStatus != nil
}

func (o *Opinion) Vote(voteStatus vote.VoteStatus) {
	o.voteStatus = &voteStatus
}

func (o *Opinion) OpinionID() shared.UUID[Opinion] {
	return o.opinionID
}

func (o *Opinion) TalkSessionID() shared.UUID[talksession.TalkSession] {
	return o.talkSessionID
}

func (o *Opinion) UserID() shared.UUID[user.User] {
	return o.userID
}

func (o *Opinion) ParentOpinionID() *shared.UUID[Opinion] {
	return o.parentOpinionID
}

func (o *Opinion) Content() string {
	return o.content
}

func (o *Opinion) CreatedAt() time.Time {
	return o.createdAt
}

func (o *Opinion) VoteStatus() vote.VoteStatus {
	if o.voteStatus == nil {
		return vote.UnVoted
	}
	return *o.voteStatus
}

func (o *Opinion) Opinions() []Opinion {
	return o.opinions
}
