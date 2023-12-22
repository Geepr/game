package relationships

import (
	"fmt"
	"github.com/Geepr/game/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func getByReleaseIdRoute(c *gin.Context) {
	id, err := utils.ParseUuidFromParam(c)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	result, err := getPlatformIdsForRelease(id)
	if err != nil {
		utils.AbortWithRelevantError(err, c)
		return
	}

	c.JSON(http.StatusOK, result)
}

func getByPlatformIdRoute(c *gin.Context) {
	id, err := utils.ParseUuidFromParam(c)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	result, err := getReleaseIdsForPlatforms(id)
	if err != nil {
		utils.AbortWithRelevantError(err, c)
		return
	}

	c.JSON(http.StatusOK, result)
}

func createRoute(c *gin.Context) {
	var createModel struct {
		PlatformId uuid.UUID `json:"platformId" binding:"required"`
		ReleaseId  uuid.UUID `json:"releaseId" binding:"required"`
	}
	if err := c.MustBindWith(&createModel, binding.JSON); err != nil {
		log.Infof("Failed to parse release platform create model: %s", err.Error())
		return
	}

	if err := addReleasePlatform(createModel.ReleaseId, createModel.PlatformId); err != nil {
		utils.AbortWithRelevantError(err, c)
		return
	}

	c.Status(http.StatusCreated)
}

func deleteRoute(c *gin.Context) {
	var createModel struct {
		PlatformId uuid.UUID `json:"platformId" binding:"required"`
		ReleaseId  uuid.UUID `json:"releaseId" binding:"required"`
	}
	if err := c.MustBindWith(&createModel, binding.JSON); err != nil {
		log.Infof("Failed to parse release platform delete model: %s", err.Error())
		return
	}

	if err := removeReleasePlatform(createModel.ReleaseId, createModel.PlatformId); err != nil {
		utils.AbortWithRelevantError(err, c)
		return
	}

	c.Status(http.StatusOK)
}

func SetupRoutes(engine *gin.Engine, basePath string) {
	baseUrl := fmt.Sprintf("%s/api/v0/releasePlatforms", basePath)

	engine.GET(baseUrl+"/byReleaseId/:id", getByReleaseIdRoute)
	engine.GET(baseUrl+"/byPlatformId/:id", getByPlatformIdRoute)
	engine.POST(baseUrl, createRoute)
	engine.DELETE(baseUrl, deleteRoute)
}
