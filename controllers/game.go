package controllers

import (
	"fmt"
	"github.com/Geepr/game/models"
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
		log.Infof("Failed to bind game query: %s", err.Error())
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

func (g *GameController) GetById(c *gin.Context) {
	lookupUuid, err := parseUuidFromParam(c)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	game, err := g.repo.GetGameById(lookupUuid)
	if err != nil {
		abortWithRelevantError(err, c)
		return
	}

	c.JSON(http.StatusOK, game)
}

func (g *GameController) Create(c *gin.Context) {
	var createModel struct {
		Title       string `json:"title" binding:"required,max=200"`
		Description string `json:"description" binding:"max=2000"`
	}
	if err := c.BindJSON(&createModel); err != nil {
		log.Infof("Failed to parse game creation model: %s", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	game := models.Game{
		Title:       createModel.Title,
		Description: getNilIfDefault(createModel.Description),
		Archived:    false,
	}
	if err := g.repo.AddGame(&game); err != nil {
		abortWithRelevantError(err, c)
		return
	}

	c.JSON(http.StatusCreated, &game)
}

func (g *GameController) Update(c *gin.Context) {
	var updateModel struct {
		Title       string `json:"title" binding:"required,max=200"`
		Description string `json:"description" binding:"max=2000"`
		Archived    bool   `json:"archived"`
	}
	if err := c.BindJSON(&updateModel); err != nil {
		log.Infof("Failed to parse game update model: %s", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	id, err := parseUuidFromParam(c)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	game := models.Game{
		Title:       updateModel.Title,
		Description: getNilIfDefault(updateModel.Description),
		Archived:    updateModel.Archived,
	}
	if err := g.repo.UpdateGame(id, &game); err != nil {
		abortWithRelevantError(err, c)
		return
	}

	// this is here just so that it displays properly when returned to the user
	game.Id = id
	c.JSON(http.StatusOK, &game)
}

func (g *GameController) Delete(c *gin.Context) {
	id, err := parseUuidFromParam(c)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = g.repo.DeleteGame(id)
	if err != nil {
		abortWithRelevantError(err, c)
		return
	}

	c.Status(http.StatusOK)
}

func (g *GameController) SetupRoutes(engine *gin.Engine, basePath string) {
	baseUrl := fmt.Sprintf("%s/api/v0/games", basePath)

	engine.GET(baseUrl, g.Get)
	engine.GET(baseUrl+"/:id", g.GetById)
	engine.POST(baseUrl, g.Create)
	engine.PUT(baseUrl+"/:id", g.Update)
	engine.DELETE(baseUrl+"/:id", g.Delete)
}
