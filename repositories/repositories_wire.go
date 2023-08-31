//go:build wireinject
// +build wireinject

package repositories

import (
	"github.com/KowalskiPiotr98/gotabase"
	"github.com/google/wire"
)

var (
	PlatformRepositorySet = wire.NewSet(NewPlatformRepository)
)

func CreatePlatformRepository() *PlatformRepository {
	wire.Build(gotabase.GetConnection, PlatformRepositorySet)
	return nil
}
