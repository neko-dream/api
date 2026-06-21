package repository

import (
	"context"
	"database/sql"
	"encoding/json"

	"braces.dev/errtrace"
	"github.com/google/uuid"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
	model "github.com/neko-dream/api/internal/infrastructure/persistence/sqlc/generated"
	"github.com/neko-dream/api/pkg/utils"
	"github.com/sqlc-dev/pqtype"
	"go.opentelemetry.io/otel"
)

type userStatusChangeLogRepository struct {
	*db.DBManager
}

func NewUserStatusChangeLogRepository(dbManager *db.DBManager) user.UserStatusChangeLogRepository {
	return &userStatusChangeLogRepository{
		DBManager: dbManager,
	}
}

func (r *userStatusChangeLogRepository) Create(ctx context.Context, log user.UserStatusChangeLog) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "userStatusChangeLogRepository.Create")
	defer span.End()

	// additionalDataをJSON形式に変換
	var additionalData pqtype.NullRawMessage
	if log.AdditionalData() != nil && len(log.AdditionalData()) > 0 {
		jsonData, err := json.Marshal(log.AdditionalData())
		if err != nil {
			utils.HandleError(ctx, err, "json.Marshal")
			return errtrace.Wrap(err)
		}
		additionalData = pqtype.NullRawMessage{
			RawMessage: jsonData,
			Valid:      true,
		}
	}

	// IPアドレスの変換
	var ipAddress pqtype.Inet
	if log.IPAddress() != nil {
		ipAddress = pqtype.Inet{
			IPNet: utils.ParseIPNet(*log.IPAddress()),
			Valid: true,
		}
	}

	// UserAgentの変換
	var userAgent sql.NullString
	if log.UserAgent() != nil {
		userAgent = sql.NullString{
			String: *log.UserAgent(),
			Valid:  true,
		}
	}

	// Reasonの変換
	var reason sql.NullString
	if log.Reason() != nil {
		reason = sql.NullString{
			String: *log.Reason(),
			Valid:  true,
		}
	}

	err := r.DBManager.GetQueries(ctx).CreateUserStatusChangeLog(ctx, model.CreateUserStatusChangeLogParams{
		UserStatusChangeLogsID: uuid.UUID(log.ID()),
		UserID:                 uuid.UUID(log.UserID()),
		Status:                 string(log.Status()),
		Reason:                 reason,
		ChangedAt:              log.ChangedAt(),
		ChangedBy:              string(log.ChangedBy()),
		IpAddress:              ipAddress,
		UserAgent:              userAgent,
		AdditionalData:         additionalData,
	})
	if err != nil {
		utils.HandleError(ctx, err, "queries.CreateUserStatusChangeLog")
		return errtrace.Wrap(err)
	}

	return nil
}

func (r *userStatusChangeLogRepository) FindByUserID(ctx context.Context, userID shared.UUID[user.User]) ([]user.UserStatusChangeLog, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "userStatusChangeLogRepository.FindByUserID")
	defer span.End()

	logs, err := r.DBManager.GetQueries(ctx).FindUserStatusChangeLogsByUserID(ctx, uuid.UUID(userID))
	if err != nil {
		utils.HandleError(ctx, err, "queries.FindUserStatusChangeLogsByUserID")
		return nil, errtrace.Wrap(err)
	}

	var result []user.UserStatusChangeLog
	for _, log := range logs {
		// IPアドレスの変換
		var ipAddress *string
		if log.IpAddress.Valid {
			ip := log.IpAddress.IPNet.IP.String()
			ipAddress = &ip
		}

		// UserAgentの変換
		var userAgent *string
		if log.UserAgent.Valid {
			userAgent = &log.UserAgent.String
		}

		// Reasonの変換
		var reason *string
		if log.Reason.Valid {
			reason = &log.Reason.String
		}

		// additionalDataの変換
		var additionalData map[string]interface{}
		if log.AdditionalData.Valid {
			if err := json.Unmarshal(log.AdditionalData.RawMessage, &additionalData); err != nil {
				utils.HandleError(ctx, err, "json.Unmarshal")
				// エラーでも続行
			}
		}

		changeLog := user.NewUserStatusChangeLogWithID(
			shared.UUID[user.UserStatusChangeLog](log.UserStatusChangeLogsID),
			shared.UUID[user.User](log.UserID),
			user.UserStatus(log.Status),
			reason,
			log.ChangedAt,
			user.ChangedBy(log.ChangedBy),
			ipAddress,
			userAgent,
			additionalData,
		)

		result = append(result, changeLog)
	}

	return result, nil
}
