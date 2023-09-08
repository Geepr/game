package controllers

import (
	"fmt"
	"github.com/Geepr/game/models"
	"github.com/Geepr/game/repositories"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type GameReleaseController struct {
	repo *repositories.GameReleaseRepository
}

var _ Routable = (*GameReleaseController)(nil)

func NewGameReleaseController(repo *repositories.GameReleaseRepository) *GameReleaseController {
	return &GameReleaseController{repo: repo}
}

func (g *GameReleaseController) Get(c *gin.Context) {
	var query struct {
		Title     string                            `form:"title"`
		GameId    uuid.UUID                         `form:"gameId"`
		SortOrder repositories.GameReleaseSortOrder `form:"order"`
		PageIndex int                               `form:"index"`
		PageSize  int                               `form:"size"`
	}
	if err := c.MustBindWith(&query, binding.Query); err != nil {
		log.Infof("Failed to bind game release query: %s", err.Error())
		return
	}

	releases, err := g.repo.GetGameReleases(query.Title, query.GameId, query.PageIndex, query.PageSize, query.SortOrder)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, releases)
}

func (g *GameReleaseController) GetById(c *gin.Context) {
	id, err := parseUuidFromParam(c)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	release, err := g.repo.GetGameReleaseById(id)
	if err != nil {
		abortWithRelevantError(err, c)
		return
	}

	c.JSON(http.StatusOK, release)
}

func (g *GameReleaseController) Create(c *gin.Context) {
	var createModel struct {
		GameId             uuid.UUID `json:"gameId" binding:"required"`
		TitleOverride      string    `json:"title" binding:"max=200"`
		Description        string    `json:"description" binding:"max=2000"`
		ReleaseDateUnknown bool      `json:"releaseDateUnknown"`
		ReleaseDate        time.Time `json:"releaseDate"` //in format 2006-01-02T15:04:05Z07:00
	}
	if err := c.MustBindWith(&createModel, binding.JSON); err != nil {
		log.Infof("Failed to parse release creation model: %s", err.Error())
		return
	}

	release := models.GameRelease{
		GameId:             createModel.GameId,
		TitleOverride:      getNilIfDefault(createModel.TitleOverride),
		Description:        getNilIfDefault(createModel.Description),
		ReleaseDate:        getNilIfDefault(createModel.ReleaseDate),
		ReleaseDateUnknown: createModel.ReleaseDateUnknown,
	}
	if err := g.repo.AddGameRelease(&release); err != nil {
		abortWithRelevantError(err, c)
		return
	}

	c.JSON(http.StatusCreated, &release)
}

func (g *GameReleaseController) Update(c *gin.Context) {
	var updateModel struct {
		TitleOverride      string    `json:"title" binding:"max=200"`
		Description        string    `json:"description" binding:"max=2000"`
		ReleaseDateUnknown bool      `json:"releaseDateUnknown"`
		ReleaseDate        time.Time `json:"releaseDate"` //in format 2006-01-02T15:04:05Z07:00
	}
	if err := c.MustBindWith(&updateModel, binding.JSON); err != nil {
		log.Infof("Failed to parse release update model: %s", err.Error())
		return
	}
	id, err := parseUuidFromParam(c)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	release := models.GameRelease{
		TitleOverride:      getNilIfDefault(updateModel.TitleOverride),
		Description:        getNilIfDefault(updateModel.Description),
		ReleaseDate:        getNilIfDefault(updateModel.ReleaseDate),
		ReleaseDateUnknown: updateModel.ReleaseDateUnknown,
	}
	if err := g.repo.UpdateGameRelease(id, &release); err != nil {
		abortWithRelevantError(err, c)
		return
	}

	release.Id = id
	c.JSON(http.StatusOK, &release)
}

func (g *GameReleaseController) Delete(c *gin.Context) {
	id, err := parseUuidFromParam(c)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := g.repo.DeleteGameRelease(id); err != nil {
		abortWithRelevantError(err, c)
		return
	}

	c.Status(http.StatusOK)
}

func (g *GameReleaseController) SetupRoutes(engine *gin.Engine, basePath string) {
	baseUrl := fmt.Sprintf("%s/api/v0/releases", basePath)

	engine.GET(baseUrl, g.Get)
	engine.GET(baseUrl+"/:id", g.GetById)
	engine.POST(baseUrl, g.Create)
	engine.PUT(baseUrl+"/:id", g.Update)
	engine.DELETE(baseUrl+"/:id", g.Delete)
}
