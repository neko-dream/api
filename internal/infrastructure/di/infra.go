package di

import (
	"github.com/neko-dream/server/internal/infrastructure/auth/jwt"
	"github.com/neko-dream/server/internal/infrastructure/auth/oauth"
	"github.com/neko-dream/server/internal/infrastructure/config"
	"github.com/neko-dream/server/internal/infrastructure/crypto"
	client "github.com/neko-dream/server/internal/infrastructure/external/analysis"
	"github.com/neko-dream/server/internal/infrastructure/external/aws"
	"github.com/neko-dream/server/internal/infrastructure/http/cookie"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	"github.com/neko-dream/server/internal/infrastructure/persistence/postgresql"
	"github.com/neko-dream/server/internal/infrastructure/persistence/repository"
	"github.com/neko-dream/server/internal/infrastructure/telemetry"
)

func infraDeps() []ProvideArg {
	return []ProvideArg{
		{config.LoadConfig, nil},
		{postgresql.Connect, nil},
		{db.NewMigrator, nil},
		{db.NewDBManager, nil},
		{oauth.NewProviderFactory, nil},
		// {telemetry.SentryProvider, nil},
		{telemetry.BaselimeProvider, nil},
		{repository.InitS3Client, nil},
		{repository.NewImageRepository, nil},
		{repository.NewImageStorage, nil},
		{repository.NewSessionRepository, nil},
		{repository.NewUserRepository, nil},
		{repository.NewTalkSessionRepository, nil},
		{repository.NewOpinionRepository, nil},
		{repository.NewVoteRepository, nil},
		{repository.NewConclusionRepository, nil},
		{repository.NewActionItemRepository, nil},
		{repository.NewPolicyRepository, nil},
		{repository.NewConsentRecordRepository, nil},
		{repository.NewReportRepository, nil},
		{jwt.NewTokenManager, nil},
		{db.NewDummyInitializer, nil},
		{client.NewAnalysisService, nil},
		{cookie.NewCookieManager, nil},
		{crypto.NewEncryptor, nil},
		{aws.NewAWSConfig, nil},
		{aws.NewSESClient, nil},
		{repository.NewPasswordAuthRepository, nil},
		{repository.NewOrganizationUserRepository, nil},
		{repository.NewOrganizationRepository, nil},
	}
}
