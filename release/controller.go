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
		GameId    uuid.UUID `form:"gameId"`
		SortOrder SortOrder `form:"order"`
		PageIndex int       `form:"index"`
		PageSize  int       `form:"size"`
	}
	if err := c.MustBindWith(&query, binding.Query); err != nil {
		log.Infof("Failed to bind game release query: %s", err.Error())
		return
	}

	releases, err := getGameReleases(query.Title, query.GameId, query.PageIndex, query.PageSize, query.SortOrder)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, releases)
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

	release := GameRelease{
		GameId:             createModel.GameId,
		TitleOverride:      utils.GetNilIfDefault(createModel.TitleOverride),
		Description:        utils.GetNilIfDefault(createModel.Description),
		ReleaseDate:        utils.GetNilIfDefault(createModel.ReleaseDate),
		ReleaseDateUnknown: createModel.ReleaseDateUnknown,
	}
	if err := addGameRelease(&release); err != nil {
		utils.AbortWithRelevantError(err, c)
		return
	}

	c.JSON(http.StatusCreated, &release)
}

func updateRoute(c *gin.Context) {
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
	id, err := utils.ParseUuidFromParam(c)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	release := GameRelease{
		TitleOverride:      utils.GetNilIfDefault(updateModel.TitleOverride),
		Description:        utils.GetNilIfDefault(updateModel.Description),
		ReleaseDate:        utils.GetNilIfDefault(updateModel.ReleaseDate),
		ReleaseDateUnknown: updateModel.ReleaseDateUnknown,
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
