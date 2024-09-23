package repository

import (
	"context"

	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
)

type sessionRepository struct {
}

// Create implements session.SessionRepository.
func (s *sessionRepository) Create(context.Context, session.Session) (*session.Session, error) {
	panic("unimplemented")
}

// FindBySessionID implements session.SessionRepository.
func (s *sessionRepository) FindBySessionID(context.Context, shared.UUID[session.Session]) (*session.Session, error) {
	panic("unimplemented")
}

// FindByUserID implements session.SessionRepository.
func (s *sessionRepository) FindByUserID(context.Context, shared.UUID[user.User]) ([]session.Session, error) {
	panic("unimplemented")
}

// Update implements session.SessionRepository.
func (s *sessionRepository) Update(context.Context, session.Session) (*session.Session, error) {
	panic("unimplemented")
}

func NewSessionRepository() session.SessionRepository {
	return &sessionRepository{}
}
