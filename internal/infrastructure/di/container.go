package di

import (
	"log"

	"braces.dev/errtrace"
	"github.com/neko-dream/server/internal/domain/model/session"
	"github.com/neko-dream/server/internal/domain/service"
	"github.com/neko-dream/server/internal/infrastructure/auth"
	"github.com/neko-dream/server/internal/infrastructure/config"
	"github.com/neko-dream/server/internal/infrastructure/datasource/postgresql"
	"github.com/neko-dream/server/internal/infrastructure/datasource/repository"
	"github.com/neko-dream/server/internal/infrastructure/db"
	"github.com/neko-dream/server/internal/presentation/handler"
	auth_usecase "github.com/neko-dream/server/internal/usecase/auth"
	talk_session_usecase "github.com/neko-dream/server/internal/usecase/talk_session"
	user_usecase "github.com/neko-dream/server/internal/usecase/user"
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
		{db.NewMigrator, nil},
		{postgresql.Connect, nil},
		{db.NewDBManager, nil},
		{repository.NewSessionRepository, nil},
		{repository.NewUserRepository, nil},
		{repository.NewTalkSessionRepository, nil},

		{func() session.TokenManager {
			return auth.NewTokenManager("")
		}, nil},

		{service.NewAuthService, nil},
		{service.NewSessionService, nil},
		{auth_usecase.NewAuthLoginUseCase, nil},
		{auth_usecase.NewAuthCallbackUseCase, nil},
		{talk_session_usecase.NewCreateTalkSessionUseCase, nil},
		{user_usecase.NewRegisterUserUseCase, nil},
		{handler.NewSecurityHandler, nil},
		{handler.NewAuthHandler, nil},
		{handler.NewUserHandler, nil},
		{handler.NewIntentionHandler, nil},
		{handler.NewOpinionHandler, nil},
		{handler.NewTalkSessionHandler, nil},
		{handler.NewHandler, nil},
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
