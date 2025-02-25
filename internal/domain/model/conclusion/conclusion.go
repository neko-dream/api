package conclusion

import (
	"context"

	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/talksession"
	user "github.com/neko-dream/server/internal/domain/model/user"
)

type (
	ConclusionRepository interface {
		Create(context.Context, Conclusion) error
		FindByTalkSessionID(context.Context, shared.UUID[talksession.TalkSession]) (*Conclusion, error)
	}

	Conclusion struct {
		talkSessionID shared.UUID[talksession.TalkSession]
		conclusion    string
		createdBy     shared.UUID[user.User]
	}
)

func NewConclusion(
	talkSessionID shared.UUID[talksession.TalkSession],
	conclusion string,
	createdBy shared.UUID[user.User],
) *Conclusion {
	return &Conclusion{
		talkSessionID: talkSessionID,
		conclusion:    conclusion,
		createdBy:     createdBy,
	}
}

func (c *Conclusion) TalkSessionID() shared.UUID[talksession.TalkSession] {
	return c.talkSessionID
}

func (c *Conclusion) Conclusion() string {
	return c.conclusion
}

func (c *Conclusion) CreatedBy() shared.UUID[user.User] {
	return c.createdBy
}

func (c *Conclusion) EditConclusion(
	conclusion string,
) {
	c.conclusion = conclusion
}
