package admin_usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/mail"
	"net/url"
	"strings"
	"unicode/utf8"

	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/domain/model/admin"
	password_auth "github.com/neko-dream/api/internal/domain/model/auth/password"
	"github.com/neko-dream/api/internal/domain/model/clock"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/neko-dream/api/internal/infrastructure/config"
	"github.com/neko-dream/api/internal/infrastructure/email"
	email_template "github.com/neko-dream/api/internal/infrastructure/email/template"
	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
	"github.com/neko-dream/api/pkg/hash"
	"github.com/neko-dream/api/pkg/utils"
	"go.opentelemetry.io/otel"
)

type ManageUserInvitations interface {
	Create(context.Context, InviteUserInput) (*InviteUserOutput, error)
	Resend(context.Context, shared.UUID[user.User]) (*InviteUserOutput, error)
}

type InviteUserInput struct{ DisplayName, Email string }
type InviteUserOutput struct {
	UserID         shared.UUID[user.User]
	InvitationSent bool
}
type invitationTransactionRunner interface {
	ExecTx(context.Context, func(context.Context) error) error
}

type manageUserInvitations struct {
	users     user.UserRepository
	passwords password_auth.PasswordAuthRepository
	tokens    password_auth.SetupTokenRepository
	sender    email.EmailSender
	cfg       *config.Config
	tx        invitationTransactionRunner
}

func NewManageUserInvitations(users user.UserRepository, passwords password_auth.PasswordAuthRepository, tokens password_auth.SetupTokenRepository, sender email.EmailSender, cfg *config.Config, tx *db.DBManager) ManageUserInvitations {
	return &manageUserInvitations{users: users, passwords: passwords, tokens: tokens, sender: sender, cfg: cfg, tx: tx}
}

func (m *manageUserInvitations) Create(ctx context.Context, input InviteUserInput) (*InviteUserOutput, error) {
	ctx, span := otel.Tracer("admin_usecase").Start(ctx, "manageUserInvitations.Create")
	defer span.End()
	claim := admin.GetAdminSession(ctx)
	if claim == nil {
		return nil, messages.ForbiddenError
	}
	input.DisplayName, input.Email = strings.TrimSpace(input.DisplayName), strings.TrimSpace(input.Email)
	address, err := mail.ParseAddress(input.Email)
	if err != nil || strings.ContainsAny(input.Email, "\r\n") || address.Address != input.Email || !utils.IsEmail(input.Email) || len(input.Email) > 254 || utf8.RuneCountInString(input.DisplayName) < 1 || utf8.RuneCountInString(input.DisplayName) > 50 {
		return nil, messages.BadRequestError
	}
	subject, err := hash.HMACHash(input.Email, m.cfg.HASH_PEPPER)
	if err != nil {
		return nil, err
	}
	id := shared.NewUUID[user.User]()
	usr := user.NewUser(id, nil, &input.DisplayName, subject, "PASSWORD", nil)
	// Keep public profile IDs independent of email and organization identifiers.
	compactID := strings.ReplaceAll(id.String(), "-", "")
	if err := usr.SetDisplayID("u" + compactID[len(compactID)-24:]); err != nil {
		return nil, err
	}
	usr.ChangeEmail(input.Email)
	link, token, err := m.newInvitation(ctx, usr)
	if err != nil {
		return nil, err
	}
	err = m.tx.ExecTx(ctx, func(ctx context.Context) error {
		if err := m.tokens.LockEmailSubject(ctx, subject); err != nil {
			return err
		}
		found, err := m.users.FindBySubject(ctx, user.UserSubject(subject))
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		if found != nil {
			return messages.ManagedUserEmailExists
		}
		if err := m.users.Create(ctx, usr); err != nil {
			return err
		}
		return m.tokens.Save(ctx, token)
	})
	if err != nil {
		return nil, err
	}
	slog.InfoContext(ctx, "managed user created", "admin_user_id", claim.Sub, "user_id", id.String())
	return m.send(ctx, claim.Sub, usr, link), nil
}

func (m *manageUserInvitations) Resend(ctx context.Context, id shared.UUID[user.User]) (*InviteUserOutput, error) {
	ctx, span := otel.Tracer("admin_usecase").Start(ctx, "manageUserInvitations.Resend")
	defer span.End()
	claim := admin.GetAdminSession(ctx)
	if claim == nil {
		return nil, messages.ForbiddenError
	}
	var usr *user.User
	var link string
	err := m.tx.ExecTx(ctx, func(ctx context.Context) error {
		if err := m.tokens.LockUser(ctx, id); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return messages.UserNotFoundError
			}
			return err
		}
		pending, err := m.tokens.FindByUserID(ctx, id)
		if errors.Is(err, sql.ErrNoRows) {
			return messages.PasswordSetupNotPending
		}
		if err != nil {
			return err
		}
		usr, err = m.users.FindByID(ctx, id)
		if err != nil {
			return err
		}
		if usr == nil || usr.IsWithdrawn() || usr.Provider() != "PASSWORD" || usr.Email() == nil || usr.Subject() != pending.EmailSubject {
			return messages.PasswordSetupNotPending
		}
		subject, err := hash.HMACHash(*usr.Email(), m.cfg.HASH_PEPPER)
		if err != nil {
			return err
		}
		if subject != pending.EmailSubject {
			return messages.PasswordSetupNotPending
		}
		_, err = m.passwords.GetPasswordAuthByUserID(ctx, id)
		if err == nil {
			return messages.PasswordSetupNotPending
		}
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		if clock.Now(ctx).Before(pending.CreatedAt.Add(password_auth.SetupTokenResendInterval)) {
			return messages.PasswordSetupRateLimit
		}
		var token password_auth.SetupToken
		link, token, err = m.newInvitation(ctx, *usr)
		if err != nil {
			return err
		}
		return m.tokens.Save(ctx, token)
	})
	if err != nil {
		return nil, err
	}
	return m.send(ctx, claim.Sub, *usr, link), nil
}

func (m *manageUserInvitations) newInvitation(ctx context.Context, usr user.User) (string, password_auth.SetupToken, error) {
	u, err := url.Parse(m.cfg.WEBSITE_URL)
	if err != nil || u.Scheme != "https" || u.Host == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" || (u.Path != "" && u.Path != "/") {
		return "", password_auth.SetupToken{}, fmt.Errorf("WEBSITE_URL must be a trusted HTTPS origin for password setup")
	}
	secret, err := hash.GenerateRandomToken(32)
	if err != nil {
		return "", password_auth.SetupToken{}, err
	}
	digest, _ := password_auth.SetupTokenHash(secret)
	u.Path, u.Fragment = "/auth/password/setup", secret
	now := clock.Now(ctx)
	return u.String(), password_auth.SetupToken{UserID: usr.UserID(), TokenHash: digest, EmailSubject: usr.Subject(), CreatedAt: now, ExpiresAt: now.Add(password_auth.SetupTokenTTL)}, nil
}

func (m *manageUserInvitations) send(ctx context.Context, actor string, usr user.User, link string) *InviteUserOutput {
	// Sending is deliberately after commit: a delivery failure leaves a visible,
	// pending account which operators can retry, not an email for a rolled-back user.
	err := m.sender.Send(ctx, *usr.Email(), email_template.PasswordSetupEmailTemplate, map[string]any{"Title": "ことひろのパスワードを設定してください", "RecipientName": *usr.DisplayName(), "SetupURL": link, "ExpiryHours": int(password_auth.SetupTokenTTL.Hours())})
	if err != nil {
		// Do not log transport errors: they may contain addresses or the one-time URL.
		slog.WarnContext(ctx, "password setup email delivery failed", "admin_user_id", actor, "user_id", usr.UserID().String())
	} else {
		slog.InfoContext(ctx, "password setup email sent", "admin_user_id", actor, "user_id", usr.UserID().String())
	}
	return &InviteUserOutput{UserID: usr.UserID(), InvitationSent: err == nil}
}
