package di

import (
	"github.com/neko-dream/server/internal/infrastructure/auth/oauth"
	"github.com/neko-dream/server/internal/infrastructure/auth/session"
	"github.com/neko-dream/server/internal/infrastructure/config"
	"github.com/neko-dream/server/internal/infrastructure/crypto"
	client "github.com/neko-dream/server/internal/infrastructure/external/analysis"
	"github.com/neko-dream/server/internal/infrastructure/external/aws"
	"github.com/neko-dream/server/internal/infrastructure/external/aws/ses"
	"github.com/neko-dream/server/internal/infrastructure/http/cookie"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	"github.com/neko-dream/server/internal/infrastructure/persistence/postgresql"
	"github.com/neko-dream/server/internal/infrastructure/persistence/query/organization"
	"github.com/neko-dream/server/internal/infrastructure/persistence/repository"
	"github.com/neko-dream/server/internal/infrastructure/telemetry"
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
		// {telemetry.SentryProvider, nil},
		{telemetry.BaselimeProvider, nil},
		{repository.InitS3Client, nil},
		{repository.NewImageRepository, nil},
		{repository.NewImageStorage, nil},
		{repository.NewSessionRepository, nil},
		{repository.NewUserRepository, nil},
		{repository.NewUserAuthRepository, nil},
		{repository.NewTalkSessionRepository, nil},
		{repository.NewOpinionRepository, nil},
		{repository.NewVoteRepository, nil},
		{repository.NewConclusionRepository, nil},
		{repository.NewActionItemRepository, nil},
		{repository.NewPolicyRepository, nil},
		{repository.NewConsentRecordRepository, nil},
		{repository.NewReportRepository, nil},
		{repository.NewPasswordAuthRepository, nil},
		{repository.NewOrganizationUserRepository, nil},
		{repository.NewOrganizationRepository, nil},
		{repository.NewOrganizationAliasRepository, nil},
		{repository.NewTalkSessionConsentRepository, nil},
		{repository.NewAnalysisRepository, nil},
		{repository.NewAuthStateRepository, nil},
		{client.NewAnalysisService, nil},
		{aws.NewAWSConfig, nil},
		{aws.NewSESClient, nil},
		{ses.NewSESEmailSender, nil},
		{cookie.NewCookieManager, nil},
		{crypto.NewEncryptor, nil},
		{db.NewDummyInitializer, nil},
		{organization.NewListJoinedOrganizationQuery, nil},
	}
}
