package release

import (
	"fmt"
	"github.com/Geepr/game/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func getRoute(c *gin.Context) {
	var query struct {
		Title     string    `form:"title"`
		GameId    string    `form:"gameId"`
		SortOrder SortOrder `form:"order"`
		PageIndex int       `form:"index"`
		PageSize  int       `form:"size"`
	}
	var err error
	var gameId uuid.UUID
	if err, gameId = c.MustBindWith(&query, binding.Query), uuid.FromStringOrNil(query.GameId); err != nil || gameId == uuid.Nil {
		log.Infof("Failed to bind game release query: %s", err.Error())
		return
	}

	releases, totalItems, err := getGameReleases(query.Title, gameId, query.PageIndex, query.PageSize, query.SortOrder)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	response := struct {
		Games      []*GameRelease `json:"releases"`
		Page       int            `json:"page"`
		PageSize   int            `json:"pageSize"`
		TotalPages int            `json:"totalPages"`
	}{
		Games:      releases,
		Page:       query.PageIndex,
		PageSize:   query.PageSize,
		TotalPages: utils.GetPagesFromItems(totalItems, query.PageSize),
	}
	c.JSON(http.StatusOK, response)
}

func getByIdRoute(c *gin.Context) {
	id, err := utils.ParseUuidFromParam(c)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	release, err := getGameReleaseById(id)
	if err != nil {
		utils.AbortWithRelevantError(err, c)
		return
	}

	c.JSON(http.StatusOK, release)
}

func createRoute(c *gin.Context) {
	var createModel struct {
		GameId             uuid.UUID   `json:"gameId" binding:"required"`
		TitleOverride      string      `json:"title" binding:"max=200"`
		Description        string      `json:"description" binding:"max=2000"`
		ReleaseDateUnknown bool        `json:"releaseDateUnknown"`
		ReleaseDate        time.Time   `json:"releaseDate"` //in format 2006-01-02T15:04:05Z07:00
		PlatformIds        []uuid.UUID `json:"platformIds" binding:"required"`
	}
	if err := c.MustBindWith(&createModel, binding.JSON); err != nil {
		log.Infof("Failed to parse release creation model: %s", err.Error())
		return
	}

	release := GameRelease{
		GameId:             createModel.GameId,
		TitleOverride:      utils.GetNilIfDefault(createModel.TitleOverride),
		Description:        utils.GetNilIfDefault(createModel.Description),
		ReleaseDate:        utils.GetNilIfDefault(createModel.ReleaseDate),
		ReleaseDateUnknown: createModel.ReleaseDateUnknown,
		PlatformIds:        createModel.PlatformIds,
	}
	if err := addGameRelease(&release); err != nil {
		utils.AbortWithRelevantError(err, c)
		return
	}

	c.JSON(http.StatusCreated, &release)
}

func updateRoute(c *gin.Context) {
	var updateModel struct {
		TitleOverride      string      `json:"title" binding:"max=200"`
		Description        string      `json:"description" binding:"max=2000"`
		ReleaseDateUnknown bool        `json:"releaseDateUnknown"`
		ReleaseDate        time.Time   `json:"releaseDate"` //in format 2006-01-02T15:04:05Z07:00
		PlatformIds        []uuid.UUID `json:"platformIds" binding:"required"`
	}
	if err := c.MustBindWith(&updateModel, binding.JSON); err != nil {
		log.Infof("Failed to parse release update model: %s", err.Error())
		return
	}
	id, err := utils.ParseUuidFromParam(c)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	release := GameRelease{
		Id:                 id,
		TitleOverride:      utils.GetNilIfDefault(updateModel.TitleOverride),
		Description:        utils.GetNilIfDefault(updateModel.Description),
		ReleaseDate:        utils.GetNilIfDefault(updateModel.ReleaseDate),
		ReleaseDateUnknown: updateModel.ReleaseDateUnknown,
		PlatformIds:        updateModel.PlatformIds,
	}
	if err := updateGameRelease(id, &release); err != nil {
		utils.AbortWithRelevantError(err, c)
		return
	}

	release.Id = id
	c.JSON(http.StatusOK, &release)
}

func deleteRoute(c *gin.Context) {
	id, err := utils.ParseUuidFromParam(c)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := deleteGameRelease(id); err != nil {
		utils.AbortWithRelevantError(err, c)
		return
	}

	c.Status(http.StatusOK)
}

func SetupRoutes(engine *gin.Engine, basePath string) {
	baseUrl := fmt.Sprintf("%s/api/v0/releases", basePath)

	engine.GET(baseUrl, getRoute)
	engine.GET(baseUrl+"/:id", getByIdRoute)
	engine.POST(baseUrl, createRoute)
	engine.PUT(baseUrl+"/:id", updateRoute)
	engine.DELETE(baseUrl+"/:id", deleteRoute)
}
