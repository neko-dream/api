package user

import (
	"context"
	"time"

	"github.com/neko-dream/api/internal/domain/model/shared"
)

// UserStatus ユーザーのステータス
type UserStatus string

const (
	UserStatusActive      UserStatus = "active"
	UserStatusWithdrawn   UserStatus = "withdrawn"
	UserStatusReactivated UserStatus = "reactivated"
)

// ChangedBy 変更者の種類
type ChangedBy string

const (
	ChangedByUser   ChangedBy = "user"
	ChangedByAdmin  ChangedBy = "admin"
	ChangedBySystem ChangedBy = "system"
)

// UserStatusChangeLog ユーザーステータス変更ログ
type UserStatusChangeLog struct {
	id             shared.UUID[UserStatusChangeLog]
	userID         shared.UUID[User]
	status         UserStatus
	reason         *string
	changedAt      time.Time
	changedBy      ChangedBy
	ipAddress      *string
	userAgent      *string
	additionalData map[string]interface{}
}

// NewUserStatusChangeLog 新しいユーザーステータス変更ログを作成
func NewUserStatusChangeLog(
	userID shared.UUID[User],
	status UserStatus,
	reason *string,
	changedAt time.Time,
	changedBy ChangedBy,
	ipAddress *string,
	userAgent *string,
) UserStatusChangeLog {
	return UserStatusChangeLog{
		id:        shared.NewUUID[UserStatusChangeLog](),
		userID:    userID,
		status:    status,
		reason:    reason,
		changedAt: changedAt,
		changedBy: changedBy,
		ipAddress: ipAddress,
		userAgent: userAgent,
	}
}

// Getters
func (l *UserStatusChangeLog) ID() shared.UUID[UserStatusChangeLog] {
	return l.id
}

func (l *UserStatusChangeLog) UserID() shared.UUID[User] {
	return l.userID
}

func (l *UserStatusChangeLog) Status() UserStatus {
	return l.status
}

func (l *UserStatusChangeLog) Reason() *string {
	return l.reason
}

func (l *UserStatusChangeLog) ChangedAt() time.Time {
	return l.changedAt
}

func (l *UserStatusChangeLog) ChangedBy() ChangedBy {
	return l.changedBy
}

func (l *UserStatusChangeLog) IPAddress() *string {
	return l.ipAddress
}

func (l *UserStatusChangeLog) UserAgent() *string {
	return l.userAgent
}

func (l *UserStatusChangeLog) AdditionalData() map[string]interface{} {
	return l.additionalData
}

// UserStatusChangeLogRepository ユーザーステータス変更ログのリポジトリ
type UserStatusChangeLogRepository interface {
	Create(ctx context.Context, log UserStatusChangeLog) error
	FindByUserID(ctx context.Context, userID shared.UUID[User]) ([]UserStatusChangeLog, error)
}

// NewUserStatusChangeLogWithID 既存のIDでユーザーステータス変更ログを作成（リポジトリ用）
func NewUserStatusChangeLogWithID(
	id shared.UUID[UserStatusChangeLog],
	userID shared.UUID[User],
	status UserStatus,
	reason *string,
	changedAt time.Time,
	changedBy ChangedBy,
	ipAddress *string,
	userAgent *string,
	additionalData map[string]interface{},
) UserStatusChangeLog {
	return UserStatusChangeLog{
		id:             id,
		userID:         userID,
		status:         status,
		reason:         reason,
		changedAt:      changedAt,
		changedBy:      changedBy,
		ipAddress:      ipAddress,
		userAgent:      userAgent,
		additionalData: additionalData,
	}
}
