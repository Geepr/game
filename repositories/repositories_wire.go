//go:build wireinject
// +build wireinject

package repositories

import (
	"github.com/KowalskiPiotr98/gotabase"
	"github.com/google/wire"
)

var (
	GameReleaseRepositorySet         = wire.NewSet(NewGameReleaseRepository)
	GameReleasePlatformRepositorySet = wire.NewSet(NewGameReleasePlatformRepository)
)

func CreateGameReleaseRepository() *GameReleaseRepository {
	wire.Build(gotabase.GetConnection, GameReleaseRepositorySet)
	return nil
}

func CreateGameReleasePlatformRepository() *GameReleasePlatformRepository {
	wire.Build(gotabase.GetConnection, GameReleasePlatformRepositorySet)
	return nil
}
