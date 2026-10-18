package di

import (
	"github.com/neko-dream/api/internal/domain/model/event"
	adminauth "github.com/neko-dream/api/internal/infrastructure/auth/admin"
	"github.com/neko-dream/api/internal/infrastructure/auth/oauth"
	"github.com/neko-dream/api/internal/infrastructure/auth/session"
	"github.com/neko-dream/api/internal/infrastructure/config"
	"github.com/neko-dream/api/internal/infrastructure/crypto"
	"github.com/neko-dream/api/internal/infrastructure/email/delivery"
	client "github.com/neko-dream/api/internal/infrastructure/external/analysis"
	"github.com/neko-dream/api/internal/infrastructure/external/aws"
	"github.com/neko-dream/api/internal/infrastructure/external/aws/pinpoint"
	"github.com/neko-dream/api/internal/infrastructure/http/cookie"
	"github.com/neko-dream/api/internal/infrastructure/persistence"
	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
	"github.com/neko-dream/api/internal/infrastructure/persistence/postgresql"
	"github.com/neko-dream/api/internal/infrastructure/persistence/query/organization"
	"github.com/neko-dream/api/internal/infrastructure/persistence/repository"
	"github.com/neko-dream/api/internal/infrastructure/telemetry"
)

// このファイルはインフラ層（DB・リポジトリ・外部API等）のコンストラクタを管理します。
// 新しいリポジトリや外部サービスを追加した場合は必ずここに追記してください。

func infraDeps() []ProvideArg {
	return []ProvideArg{
		{config.LoadConfig, nil},
		{postgresql.Connect, nil},
		{db.NewMigrator, nil},
		{db.NewDBManager, nil},
		{oauth.NewProviderFactory, nil},
		{session.NewSessionTokenManager, nil},
		{adminauth.NewAdminTokenManager, nil},
		{adminauth.NewGoogleTokenVerifier, nil},
		{repository.NewAdminUserRepository, nil},
		{repository.NewAdminSessionRepository, nil},
		{telemetry.SetupOtelProvider, nil},
		{repository.InitS3Client, nil},
		{repository.NewImageRepository, nil},
		{repository.NewImageStorage, nil},
		{repository.NewSessionRepository, nil},
		{repository.NewUserRepository, nil},
		{repository.NewSessionCreationPermissionRepository, nil},
		{repository.NewTalkSessionRepository, nil},
		{repository.NewOpinionRepository, nil},
		{repository.NewVoteRepository, nil},
		{repository.NewConclusionRepository, nil},
		{repository.NewActionItemRepository, nil},
		{repository.NewPolicyRepository, nil},
		{repository.NewConsentRecordRepository, nil},
		{repository.NewReportRepository, nil},
		{repository.NewPasswordAuthRepository, nil},
		{repository.NewPasswordLoginChallengeRepository, nil},
		{repository.NewPasswordSetupTokenRepository, nil},
		{repository.NewOrganizationUserRepository, nil},
		{repository.NewOrganizationRepository, nil},
		{repository.NewOrganizationAliasRepository, nil},
		{repository.NewUserStatusChangeLogRepository, nil},
		{repository.NewTalkSessionConsentRepository, nil},
		{repository.NewTalkSessionSurveyRepository, nil},
		{repository.NewAnalysisRepository, nil},
		{repository.NewAuthStateRepository, nil},
		{client.NewAnalysisService, nil},
		{aws.NewAWSConfig, nil},
		{aws.NewSESClient, nil},
		{aws.NewPinpointClient, nil},
		{delivery.NewSender, nil},
		{cookie.NewCookieManager, nil},
		{crypto.NewEncryptor, nil},
		{pinpoint.NewPushNotificationSender, nil},
		{repository.NewDeviceRepository, nil},
		{repository.NewNotificationPreferenceRepository, nil},
		{db.NewDummyInitializer, nil},
		{organization.NewListJoinedOrganizationQuery, nil},
		{event.NewNotifier, nil},
		{persistence.NewEventStore, nil},
	}
}
