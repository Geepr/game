//go:build wireinject
// +build wireinject

package controllers

import (
	"github.com/Geepr/game/repositories"
	"github.com/KowalskiPiotr98/gotabase"
	"github.com/google/wire"
)

var (
	GameControllerSet = wire.NewSet(NewGameController, repositories.GameRepositorySet)
)

func CreateGameController() *GameController {
	wire.Build(gotabase.GetConnection, GameControllerSet)
	return nil
}
