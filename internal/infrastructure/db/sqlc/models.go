// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Opinion struct {
	OpinionID       uuid.UUID
	TalkSessionID   uuid.UUID
	UserID          uuid.UUID
	ParentOpinionID uuid.NullUUID
	Title           sql.NullString
	Content         string
	CreatedAt       time.Time
}

type RepresentativeOpinion struct {
	TalkSessionID uuid.UUID
	OpinionID     uuid.UUID
	GroupID       int32
	Rank          int32
	UpdatedAt     time.Time
	CreatedAt     time.Time
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
}

type TalkSessionLocation struct {
	TalkSessionID uuid.UUID
	Location      interface{}
	City          string
	Prefecture    string
}

type User struct {
	UserID      uuid.UUID
	DisplayID   sql.NullString
	DisplayName sql.NullString
	IconUrl     sql.NullString
	CreatedAt   time.Time
	UpdatedAt   time.Time
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
	YearOfBirth        sql.NullInt32
	Occupation         sql.NullInt16
	Gender             int16
	Municipality       sql.NullString
	HouseholdSize      sql.NullInt16
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type UserGroupInfo struct {
	TalkSessionID uuid.UUID
	UserID        uuid.UUID
	GroupID       int32
	PosX          float64
	PosY          float64
	UpdatedAt     time.Time
	CreatedAt     time.Time
}

type Vote struct {
	VoteID        uuid.UUID
	OpinionID     uuid.UUID
	UserID        uuid.UUID
	VoteType      int16
	CreatedAt     time.Time
	TalkSessionID uuid.UUID
}
