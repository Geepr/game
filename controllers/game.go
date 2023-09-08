package controllers

import (
	"fmt"
	"github.com/Geepr/game/repositories"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type GameController struct {
	repo *repositories.GameRepository
}

func NewGameController(repo *repositories.GameRepository) *GameController {
	return &GameController{repo: repo}
}

var _ Routable = (*GameController)(nil)

func (g *GameController) Get(c *gin.Context) {
	var query struct {
		Title     string                     `form:"title"`
		SortOrder repositories.GameSortOrder `form:"order"`
		PageIndex int                        `form:"index"`
		PageSize  int                        `form:"size"`
	}
	if err := c.BindQuery(&query); err != nil {
		log.Warnf("Failed to bind game query: %s", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	games, err := g.repo.GetGames(query.Title, query.PageIndex, query.PageSize, query.SortOrder)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, games)
}

func (g *GameController) SetupRoutes(engine *gin.Engine, basePath string) {
	baseUrl := fmt.Sprintf("%s/api/v0/games", basePath)

	engine.GET(baseUrl, g.Get)
}
