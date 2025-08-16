package di

import (
	"log"

	"braces.dev/errtrace"

	"go.uber.org/dig"
)

var (
	deps = []ProvideArg{}
)

func AddProvider(arg ProvideArg) {
	deps = append(deps, arg)
}

func BuildContainer() *dig.Container {
	container := ProvideDependencies(
		infraDeps(),
		domainDeps(),
		useCaseDeps(),
		presentationDeps(),
	)
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

// InvokeWithError コンテナに登録したプロバイダの型をTにわたすとそのインスタンスを得られる（エラーを返す）
func InvokeWithError[T any](container *dig.Container, opts ...dig.InvokeOption) (T, error) {
	var res T

	if err := container.Invoke(func(t T) {
		res = t
	}, opts...); err != nil {
		return res, errtrace.Wrap(err)
	}
	return res, nil
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
