package dto

import (
	"time"

	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/talksession"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/presentation/oas"
	"github.com/neko-dream/server/pkg/utils"
)

type TalkSession struct {
	TalkSessionID    shared.UUID[talksession.TalkSession]
	OwnerID          shared.UUID[user.User]
	Theme            string
	ScheduledEndTime time.Time
	CreatedAt        time.Time
	ThumbnailURL     *string
	Description      *string
	City             *string
	Prefecture       *string
	Restrictions     []string
	HideReport       bool
}

type TalkSessionWithDetail struct {
	TalkSession
	User         User
	OpinionCount int
	Latitude     *float64
	Longitude    *float64
}

// Latitude, Longitudeはnull, または0の場合はfalseを返す
func (t *TalkSessionWithDetail) HasLocation() bool {
	return t.Latitude != nil && t.Longitude != nil && *t.Latitude != 0 && *t.Longitude != 0
}

func (t *TalkSessionWithDetail) ToResponse() oas.TalkSession {
	var location oas.OptTalkSessionLocation
	if t.HasLocation() {
		location = oas.OptTalkSessionLocation{
			Value: oas.TalkSessionLocation{
				Latitude:  utils.ToOpt[oas.OptFloat64](t.Latitude),
				Longitude: utils.ToOpt[oas.OptFloat64](t.Longitude),
			},
			Set: true,
		}
	}

	restrictions := make([]oas.Restriction, 0, len(t.Restrictions))
	for _, restriction := range t.Restrictions {
		res := talksession.RestrictionAttributeKey(restriction)
		attr := res.RestrictionAttribute()
		restrictions = append(restrictions, oas.Restriction{
			Key:         string(attr.Key),
			Description: attr.Description,
		})
	}

	return oas.TalkSession{
		ID:               t.TalkSessionID.String(),
		Theme:            t.Theme,
		Description:      utils.ToOptNil[oas.OptNilString](t.Description),
		Owner:            oas.TalkSessionOwner(t.User.ToResponse()),
		CreatedAt:        t.CreatedAt.Format(time.RFC3339),
		ScheduledEndTime: t.ScheduledEndTime.Format(time.RFC3339),
		Location:         location,
		City:             utils.ToOptNil[oas.OptNilString](t.City),
		Prefecture:       utils.ToOptNil[oas.OptNilString](t.Prefecture),
		ThumbnailURL:     utils.ToOptNil[oas.OptNilString](t.ThumbnailURL),
		Restrictions:     restrictions,
		HideReport:       t.HideReport,
	}
}
