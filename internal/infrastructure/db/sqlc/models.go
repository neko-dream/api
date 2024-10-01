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
	OpinionContent  string
	ParentOpinionID uuid.NullUUID
	VoteID          uuid.NullUUID
	CreatedAt       time.Time
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
	TalkSessionID uuid.UUID
	OwnerID       uuid.UUID
	Theme         string
	FinishedAt    sql.NullTime
	CreatedAt     time.Time
}

type User struct {
	UserID      uuid.UUID
	DisplayID   sql.NullString
	DisplayName sql.NullString
	Picture     sql.NullString
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

type Vote struct {
	VoteID    uuid.UUID
	OpinionID uuid.UUID
	UserID    uuid.UUID
	CreatedAt time.Time
}
