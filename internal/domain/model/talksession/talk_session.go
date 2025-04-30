package talksession

import (
	"context"
	"errors"
	"time"

	"github.com/neko-dream/server/internal/domain/model/clock"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
)

type (
	TalkSessionRepository interface {
		Create(ctx context.Context, talkSession *TalkSession) error
		Update(ctx context.Context, talkSession *TalkSession) error
		FindByID(ctx context.Context, talkSessionID shared.UUID[TalkSession]) (*TalkSession, error)
	}

	TalkSession struct {
		talkSessionID    shared.UUID[TalkSession]
		ownerUserID      shared.UUID[user.User]
		theme            string
		description      *string
		thumbnailURL     *string
		scheduledEndTime time.Time // 予定終了時間
		createdAt        time.Time // 作成日時
		location         *Location
		city             *string
		prefecture       *string
		restrictions     []*RestrictionAttribute // 参加制限
	}
)

func NewTalkSession(
	talkSessionID shared.UUID[TalkSession],
	theme string,
	description *string,
	thumbnailURL *string,
	ownerUserID shared.UUID[user.User],
	createdAt time.Time,
	scheduledEndTime time.Time,
	location *Location,
	city *string,
	prefecture *string,
) *TalkSession {
	return &TalkSession{
		talkSessionID:    talkSessionID,
		theme:            theme,
		description:      description,
		thumbnailURL:     thumbnailURL,
		ownerUserID:      ownerUserID,
		createdAt:        createdAt,
		scheduledEndTime: scheduledEndTime,
		location:         location,
		city:             city,
		prefecture:       prefecture,
	}
}

func (t *TalkSession) TalkSessionID() shared.UUID[TalkSession] {
	return t.talkSessionID
}

func (t *TalkSession) OwnerUserID() shared.UUID[user.User] {
	return t.ownerUserID
}

func (t *TalkSession) Theme() string {
	return t.theme
}

func (t *TalkSession) Description() *string {
	return t.description
}

func (t *TalkSession) ThumbnailURL() *string {
	return t.thumbnailURL
}

func (t *TalkSession) ScheduledEndTime() time.Time {
	return t.scheduledEndTime
}

func (t *TalkSession) CreatedAt() time.Time {
	return t.createdAt
}

func (t *TalkSession) Location() *Location {
	return t.location
}

func (t *TalkSession) City() *string {
	return t.city
}

func (t *TalkSession) Prefecture() *string {
	return t.prefecture
}

func (t *TalkSession) ChangeTheme(theme string) {
	t.theme = theme
}
func (t *TalkSession) ChangeDescription(description *string) {
	t.description = description
}
func (t *TalkSession) ChangeThumbnailURL(thumbnailURL *string) {
	t.thumbnailURL = thumbnailURL
}
func (t *TalkSession) ChangeScheduledEndTime(scheduledEndTime time.Time) {
	t.scheduledEndTime = scheduledEndTime
}
func (t *TalkSession) ChangeLocation(location *Location) {
	t.location = location
}
func (t *TalkSession) ChangeCity(city *string) {
	t.city = city
}
func (t *TalkSession) ChangePrefecture(prefecture *string) {
	t.prefecture = prefecture
}

func (t *TalkSession) Restrictions() []*RestrictionAttribute {
	return t.restrictions
}
func (t *TalkSession) RestrictionList() Restrictions {
	var restrictions []string
	for _, restriction := range t.restrictions {
		restrictions = append(restrictions, string(restriction.Key))
	}
	return restrictions
}

// 終了しているかを調べる
func (t *TalkSession) IsFinished(ctx context.Context) bool {
	ctx, span := otel.Tracer("talksession").Start(ctx, "TalkSession.IsFinished")
	defer span.End()

	return t.scheduledEndTime.Before(clock.Now(ctx))
}

// 参加制限を全てアップデートする
func (t *TalkSession) UpdateRestrictions(ctx context.Context, restrictions []string) error {
	ctx, span := otel.Tracer("talksession").Start(ctx, "TalkSession.UpdateRestrictions")
	defer span.End()

	_ = ctx

	var attrs []*RestrictionAttribute
	var errs error
	for _, restriction := range restrictions {
		if restriction == "" {
			continue
		}
		attribute := RestrictionAttributeKey(restriction)
		if err := attribute.IsValid(); err != nil {
			errs = errors.Join(errs, err)
			continue
		}

		attrs = append(attrs, lo.ToPtr(attribute.RestrictionAttribute()))
	}
	if errs != nil {
		err := &ErrInvalidRestrictionAttribute
		err.Message = errs.Error()
		return err
	}

	t.restrictions = attrs
	return nil
}
