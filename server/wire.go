//go:build wireinject
// +build wireinject

package server

import (
	"github.com/google/wire"
)

func Initialize() (*App, error) {
	panic(wire.Build(
		baseSet,
		repoSet,
		handlerSet,
		controllerSet,
		wire.Struct(new(Controllers), "*"),
		NewApp,
	))
}
