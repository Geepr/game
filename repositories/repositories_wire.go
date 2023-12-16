//go:build wireinject
// +build wireinject

package repositories

import (
	"github.com/KowalskiPiotr98/gotabase"
	"github.com/google/wire"
)

var (
	PlatformRepositorySet            = wire.NewSet(NewPlatformRepository)
	GameReleaseRepositorySet         = wire.NewSet(NewGameReleaseRepository)
	GameReleasePlatformRepositorySet = wire.NewSet(NewGameReleasePlatformRepository)
)

func CreatePlatformRepository() *PlatformRepository {
	wire.Build(gotabase.GetConnection, PlatformRepositorySet)
	return nil
}

func CreateGameReleaseRepository() *GameReleaseRepository {
	wire.Build(gotabase.GetConnection, GameReleaseRepositorySet)
	return nil
}

func CreateGameReleasePlatformRepository() *GameReleasePlatformRepository {
	wire.Build(gotabase.GetConnection, GameReleasePlatformRepositorySet)
	return nil
}
