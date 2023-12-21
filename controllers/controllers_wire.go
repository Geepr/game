//go:build wireinject
// +build wireinject

package controllers

import (
	"github.com/Geepr/game/repositories"
	"github.com/KowalskiPiotr98/gotabase"
	"github.com/google/wire"
)

var (
	GameReleaseControllerSet         = wire.NewSet(NewGameReleaseController, repositories.GameReleaseRepositorySet)
	GameReleasePlatformControllerSet = wire.NewSet(NewGameReleasePlatformController, repositories.GameReleasePlatformRepositorySet)
)

func CreateGameReleaseController() *GameReleaseController {
	wire.Build(gotabase.GetConnection, GameReleaseControllerSet)
	return nil
}

func CreateGameReleasePlatformController() *GameReleasePlatformController {
	wire.Build(gotabase.GetConnection, GameReleasePlatformControllerSet)
	return nil
}
