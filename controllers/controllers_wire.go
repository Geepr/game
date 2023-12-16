//go:build wireinject
// +build wireinject

package controllers

import (
	"github.com/Geepr/game/repositories"
	"github.com/KowalskiPiotr98/gotabase"
	"github.com/google/wire"
)

var (
	PlatformControllerSet            = wire.NewSet(NewPlatformController, repositories.PlatformRepositorySet)
	GameReleaseControllerSet         = wire.NewSet(NewGameReleaseController, repositories.GameReleaseRepositorySet)
	GameReleasePlatformControllerSet = wire.NewSet(NewGameReleasePlatformController, repositories.GameReleasePlatformRepositorySet)
)

func CreatePlatformController() *PlatformController {
	wire.Build(gotabase.GetConnection, PlatformControllerSet)
	return nil
}

func CreateGameReleaseController() *GameReleaseController {
	wire.Build(gotabase.GetConnection, GameReleaseControllerSet)
	return nil
}

func CreateGameReleasePlatformController() *GameReleasePlatformController {
	wire.Build(gotabase.GetConnection, GameReleasePlatformControllerSet)
	return nil
}
