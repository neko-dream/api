package di

import (
	"log"

	"braces.dev/errtrace"
	"github.com/neko-dream/server/internal/domain/service"
	"github.com/neko-dream/server/internal/infrastructure/auth/jwt"
	"github.com/neko-dream/server/internal/infrastructure/config"
	client "github.com/neko-dream/server/internal/infrastructure/datasource/analysis"
	"github.com/neko-dream/server/internal/infrastructure/datasource/postgresql"
	"github.com/neko-dream/server/internal/infrastructure/datasource/repository"
	"github.com/neko-dream/server/internal/infrastructure/db"
	opentelemetry "github.com/neko-dream/server/internal/infrastructure/open_telemetry"
	"github.com/neko-dream/server/internal/presentation/handler"
	analysis_usecase "github.com/neko-dream/server/internal/usecase/analysis"
	auth_usecase "github.com/neko-dream/server/internal/usecase/auth"
	opinion_usecase "github.com/neko-dream/server/internal/usecase/opinion"
	talk_session_usecase "github.com/neko-dream/server/internal/usecase/talk_session"
	timeline_usecase "github.com/neko-dream/server/internal/usecase/timeline"
	user_usecase "github.com/neko-dream/server/internal/usecase/user"
	vote_usecase "github.com/neko-dream/server/internal/usecase/vote"
	"go.uber.org/dig"
)

var (
	deps = []ProvideArg{}
)

func AddProvider(arg ProvideArg) {
	deps = append(deps, arg)
}

func BuildContainer() *dig.Container {
	deps := []ProvideArg{
		{config.LoadConfig, nil},
		{postgresql.Connect, nil},
		{db.NewMigrator, nil},
		{db.NewDBManager, nil},
		{opentelemetry.SentryProvider, nil},
		{repository.InitConfig, nil},
		{repository.InitS3Client, nil},
		{repository.NewImageRepository, nil},
		{repository.NewSessionRepository, nil},
		{repository.NewUserRepository, nil},
		{repository.NewTalkSessionRepository, nil},
		{repository.NewOpinionRepository, nil},
		{repository.NewVoteRepository, nil},
		{repository.NewConclusionRepository, nil},
		{repository.NewActionItemRepository, nil},
		{db.NewDummyInitializer, nil},
		{jwt.NewTokenManager, nil},
		{service.NewAuthService, nil},
		{service.NewSessionService, nil},
		{service.NewUserService, nil},
		{service.NewOpinionService, nil},
		{service.NewActionItemService, nil},
		{auth_usecase.NewAuthLoginUseCase, nil},
		{auth_usecase.NewAuthCallbackUseCase, nil},
		{auth_usecase.NewRevokeUseCase, nil},
		{user_usecase.NewRegisterUserUseCase, nil},
		{user_usecase.NewEditUserUseCase, nil},
		{user_usecase.NewGetUserInformationQueryHandler, nil},
		{talk_session_usecase.NewCreateTalkSessionUseCase, nil},
		{talk_session_usecase.NewListTalkSessionQueryHandler, nil},
		{talk_session_usecase.NewGetTalkSessionDetailUseCase, nil},
		{talk_session_usecase.NewGetTalkSessionHistoriesQuery, nil},
		{talk_session_usecase.NewGetTalkSessionByUserQueryHandler, nil},
		{talk_session_usecase.NewCreateTalkSessionConclusionUseCase, nil},
		{talk_session_usecase.NewGetTalkSessionConclusionQuery, nil},
		{opinion_usecase.NewPostOpinionUseCase, nil},
		{opinion_usecase.NewGetOpinionRepliesUseCase, nil},
		{opinion_usecase.NewGetSwipeOpinionsQueryHandler, nil},
		{opinion_usecase.NewGetOpinionDetailUseCase, nil},
		{opinion_usecase.NewGetUserOpinionListQueryHandler, nil},
		{opinion_usecase.NewGetOpinionsByTalkSessionUseCase, nil},
		{analysis_usecase.NewGetAnalysisResultUseCase, nil},
		{analysis_usecase.NewGetReportQueryHandler, nil},
		{timeline_usecase.NewAddTimeLineUseCase, nil},
		{timeline_usecase.NewGetTimeLineUseCase, nil},
		{timeline_usecase.NewEditTimeLineUseCase, nil},
		{vote_usecase.NewPostVoteUseCase, nil},
		{client.NewAnalysisService, nil},
		{handler.NewSecurityHandler, nil},
		{handler.NewAuthHandler, nil},
		{handler.NewUserHandler, nil},
		{handler.NewVoteHandler, nil},
		{handler.NewOpinionHandler, nil},
		{handler.NewTalkSessionHandler, nil},
		{handler.NewHandler, nil},
		{handler.NewTestHandler, nil},
		{handler.NewManageHandler, nil},
		{handler.NewTimelineHandler, nil},
	}

	container := ProvideDependencies(deps)
	return container
}

type ProvideArg struct {
	Constructor any
	Opts        []dig.ProvideOption
}

func (p *ProvideArg) Provide(container *dig.Container) {
	if err := container.Provide(p.Constructor, p.Opts...); err != nil {
		panic(err)
	}
}

// Invoke コンテナに登録したプロバイダの型をTにわたすとそのインスタンスを得られる
func Invoke[T any](container *dig.Container, opts ...dig.InvokeOption) T {
	var res T

	if err := container.Invoke(func(t T) {
		res = t
	}, opts...); err != nil {
		log.Fatalln("INVOKE ERROR: ", err.Error())
		panic(err)
	}
	return res
}

// Provide コンテナにコンストラクタを登録する。Invokeされるとここで登録されたコンストラクタが実行される
func Provide(container *dig.Container, constructor any, opts ...dig.ProvideOption) error {
	return errtrace.Wrap(container.Provide(constructor, opts...))
}

// Decorate Provideで登録したコンストラクタを上書きする
func Decorate(container *dig.Container, constructor any, opts ...dig.DecorateOption) error {
	if len(opts) >= 0 || opts[0] == nil {
		return errtrace.Wrap(container.Decorate(constructor))
	} else {
		return errtrace.Wrap(container.Decorate(constructor, opts...))
	}
}

func ProvideDependencies(providers ...[]ProvideArg) *dig.Container {
	cont := dig.New()

	for _, args := range providers {
		for _, arg := range args {
			arg.Provide(cont)
		}
	}

	return cont
}
