package controllers

import (
	"fmt"
	"github.com/Geepr/game/repositories"
	"github.com/Geepr/game/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type GameReleasePlatformController struct {
	repo *repositories.GameReleasePlatformRepository
}

var _ Routable = (*GameReleasePlatformController)(nil)

func NewGameReleasePlatformController(repo *repositories.GameReleasePlatformRepository) *GameReleasePlatformController {
	return &GameReleasePlatformController{repo: repo}
}

func (g *GameReleasePlatformController) GetByReleaseId(c *gin.Context) {
	id, err := utils.ParseUuidFromParam(c)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	result, err := g.repo.GetPlatformIdsForRelease(id)
	if err != nil {
		utils.AbortWithRelevantError(err, c)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (g *GameReleasePlatformController) GetByPlatformId(c *gin.Context) {
	id, err := utils.ParseUuidFromParam(c)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	result, err := g.repo.GetReleaseIdsForPlatforms(id)
	if err != nil {
		utils.AbortWithRelevantError(err, c)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (g *GameReleasePlatformController) Create(c *gin.Context) {
	var createModel struct {
		PlatformId uuid.UUID `json:"platformId" binding:"required"`
		ReleaseId  uuid.UUID `json:"releaseId" binding:"required"`
	}
	if err := c.MustBindWith(&createModel, binding.JSON); err != nil {
		log.Infof("Failed to parse release platform create model: %s", err.Error())
		return
	}

	if err := g.repo.AddReleasePlatform(createModel.ReleaseId, createModel.PlatformId); err != nil {
		utils.AbortWithRelevantError(err, c)
		return
	}

	c.Status(http.StatusCreated)
}

func (g *GameReleasePlatformController) Delete(c *gin.Context) {
	var createModel struct {
		PlatformId uuid.UUID `json:"platformId" binding:"required"`
		ReleaseId  uuid.UUID `json:"releaseId" binding:"required"`
	}
	if err := c.MustBindWith(&createModel, binding.JSON); err != nil {
		log.Infof("Failed to parse release platform delete model: %s", err.Error())
		return
	}

	if err := g.repo.RemoveReleasePlatform(createModel.ReleaseId, createModel.PlatformId); err != nil {
		utils.AbortWithRelevantError(err, c)
		return
	}

	c.Status(http.StatusOK)
}

func (g *GameReleasePlatformController) SetupRoutes(engine *gin.Engine, basePath string) {
	baseUrl := fmt.Sprintf("%s/api/v0/releasePlatforms", basePath)

	engine.GET(baseUrl+"/byReleaseId/:id", g.GetByReleaseId)
	engine.GET(baseUrl+"/byPlatformId/:id", g.GetByPlatformId)
	engine.POST(baseUrl, g.Create)
	engine.DELETE(baseUrl, g.Delete)
}
