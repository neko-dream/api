package talksession_usecase

import (
	"context"
	"database/sql"
	"errors"

	"braces.dev/errtrace"
	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/talksession"
	"github.com/neko-dream/api/internal/domain/model/talksession/talksession_consent"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
	"github.com/neko-dream/api/pkg/utils"
	"go.opentelemetry.io/otel"
)

type (
	JoinClosedSessionUseCase interface {
		Execute(context.Context, JoinClosedSessionInput) error
	}

	JoinClosedSessionInput struct {
		TalkSessionID shared.UUID[talksession.TalkSession]
		UserID        shared.UUID[user.User]
		Passphrase    string
	}

	joinClosedSessionHandler struct {
		talksession.TalkSessionRepository
		talksession_consent.TalkSessionConsentService
		talksession.PassphraseManager
		*db.DBManager
	}
)

func NewJoinClosedSessionUseCase(
	talkSessionRepository talksession.TalkSessionRepository,
	consentService talksession_consent.TalkSessionConsentService,
	passphraseManager talksession.PassphraseManager,
	DBManager *db.DBManager,
) JoinClosedSessionUseCase {
	return &joinClosedSessionHandler{
		TalkSessionRepository:     talkSessionRepository,
		TalkSessionConsentService: consentService,
		PassphraseManager:         passphraseManager,
		DBManager:                 DBManager,
	}
}

func (h *joinClosedSessionHandler) Execute(ctx context.Context, input JoinClosedSessionInput) error {
	ctx, span := otel.Tracer("talksession_command").Start(ctx, "joinClosedSessionHandler.Execute")
	defer span.End()

	ts, err := h.TalkSessionRepository.FindByID(ctx, input.TalkSessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return messages.TalkSessionNotFound
		}
		utils.HandleError(ctx, err, "TalkSessionRepository.FindByID")
		return errtrace.Wrap(err)
	}
	if ts == nil {
		return messages.TalkSessionNotFound
	}

	// 公開セッションは参加不要
	if !ts.IsClosed() {
		return nil
	}

	// オーナーは常にアクセス可能
	if ts.OwnerUserID() == input.UserID {
		return nil
	}

	// 既に同意済みなら合言葉不要
	hasConsented, err := h.TalkSessionConsentService.HasConsented(ctx, input.TalkSessionID, input.UserID)
	if err != nil {
		utils.HandleError(ctx, err, "TalkSessionConsentService.HasConsented")
		return errtrace.Wrap(err)
	}
	if hasConsented {
		return nil
	}

	// 合言葉が空の場合は拒否
	if input.Passphrase == "" {
		return messages.TalkSessionPassphraseRequired
	}

	// 合言葉を検証
	if !ts.VerifyPassphrase(h.PassphraseManager, input.Passphrase) {
		return messages.TalkSessionPassphraseInvalid
	}

	// 同意記録を作成（参加記録として使用）
	if err := h.DBManager.ExecTx(ctx, func(ctx context.Context) error {
		return h.TalkSessionConsentService.TakeConsent(ctx, input.TalkSessionID, input.UserID, nil)
	}); err != nil {
		// 並行joinで既に同意済みになった場合は成功扱い
		if errors.Is(err, messages.TalkSessionAlreadyConsented) {
			return nil
		}
		utils.HandleError(ctx, err, "TalkSessionConsentService.TakeConsent")
		return errtrace.Wrap(err)
	}

	return nil
}
