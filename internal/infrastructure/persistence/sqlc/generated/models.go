// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/neko-dream/server/internal/domain/model/talksession"
)

type ActionItem struct {
	ActionItemID  uuid.UUID
	TalkSessionID uuid.UUID
	Sequence      int32
	Content       string
	Status        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Opinion struct {
	OpinionID       uuid.UUID
	TalkSessionID   uuid.UUID
	UserID          uuid.UUID
	ParentOpinionID uuid.NullUUID
	Title           sql.NullString
	Content         string
	CreatedAt       time.Time
	PictureUrl      sql.NullString
	ReferenceUrl    sql.NullString
}

type PolicyConsent struct {
	PolicyConsentID uuid.UUID
	UserID          uuid.UUID
	PolicyVersion   string
	ConsentedAt     time.Time
	IpAddress       string
	UserAgent       string
}

type PolicyVersion struct {
	Version   string
	CreatedAt time.Time
}

type RepresentativeOpinion struct {
	TalkSessionID uuid.UUID
	OpinionID     uuid.UUID
	GroupID       int32
	Rank          int32
	UpdatedAt     time.Time
	CreatedAt     time.Time
	AgreeCount    int32
	DisagreeCount int32
	PassCount     int32
}

type Session struct {
	SessionID      uuid.UUID
	UserID         uuid.UUID
	Provider       string
	SessionStatus  int32
	ExpiresAt      time.Time
	CreatedAt      time.Time
	LastActivityAt time.Time
}

type TalkSession struct {
	TalkSessionID    uuid.UUID
	OwnerID          uuid.UUID
	Theme            string
	ScheduledEndTime time.Time
	CreatedAt        time.Time
	City             sql.NullString
	Prefecture       sql.NullString
	Description      sql.NullString
	ThumbnailUrl     sql.NullString
	Restrictions     talksession.Restrictions
	UpdatedAt        time.Time
}

type TalkSessionConclusion struct {
	TalkSessionID uuid.UUID
	Content       string
	CreatedBy     uuid.UUID
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type TalkSessionGeneratedImage struct {
	TalkSessionID uuid.UUID
	WordmapUrl    string
	TsncUrl       string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type TalkSessionLocation struct {
	TalkSessionID uuid.UUID
	Location      interface{}
}

type TalkSessionReport struct {
	TalkSessionID uuid.UUID
	Report        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type User struct {
	UserID        uuid.UUID
	DisplayID     sql.NullString
	DisplayName   sql.NullString
	IconUrl       sql.NullString
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Email         sql.NullString
	EmailVerified bool
}

type UserAuth struct {
	UserAuthID uuid.UUID
	UserID     uuid.UUID
	Provider   string
	Subject    string
	IsVerified bool
	CreatedAt  time.Time
}

type UserDemographic struct {
	UserDemographicsID uuid.UUID
	UserID             uuid.UUID
	YearOfBirth        sql.NullString
	Gender             sql.NullString
	City               sql.NullString
	CreatedAt          time.Time
	UpdatedAt          time.Time
	Prefecture         sql.NullString
}

type UserGroupInfo struct {
	TalkSessionID  uuid.UUID
	UserID         uuid.UUID
	GroupID        int32
	PosX           float64
	PosY           float64
	UpdatedAt      time.Time
	CreatedAt      time.Time
	PerimeterIndex sql.NullInt32
}

type UserImage struct {
	UserImagesID uuid.UUID
	UserID       uuid.UUID
	Key          string
	Width        int32
	Height       int32
	Extension    string
	Archived     bool
	Url          string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Vote struct {
	VoteID        uuid.UUID
	OpinionID     uuid.UUID
	UserID        uuid.UUID
	VoteType      int16
	CreatedAt     time.Time
	TalkSessionID uuid.UUID
}
