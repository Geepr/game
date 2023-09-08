package controllers

import (
	"fmt"
	"github.com/Geepr/game/models"
	"github.com/Geepr/game/repositories"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type PlatformController struct {
	repo *repositories.PlatformRepository
}

var _ Routable = (*PlatformController)(nil)

func NewPlatformController(repo *repositories.PlatformRepository) *PlatformController {
	return &PlatformController{repo: repo}
}

func (p *PlatformController) Get(c *gin.Context) {
	var query struct {
		Name      string                         `form:"name"`
		SortOrder repositories.PlatformSortOrder `form:"order"`
		PageIndex int                            `form:"index"`
		PageSize  int                            `form:"size"`
	}
	if err := c.MustBindWith(&query, binding.Query); err != nil {
		log.Infof("Failed to bind platform query: %s", err.Error())
		return
	}

	platforms, err := p.repo.GetPlatforms(query.Name, query.PageIndex, query.PageSize, query.SortOrder)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, platforms)
}

func (p *PlatformController) GetById(c *gin.Context) {
	id, err := parseUuidFromParam(c)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	game, err := p.repo.GetPlatformById(id)
	if err != nil {
		abortWithRelevantError(err, c)
		return
	}

	c.JSON(http.StatusOK, game)
}

func (p *PlatformController) Create(c *gin.Context) {
	var createModel struct {
		Name      string `json:"name" binding:"required,max=200"`
		ShortName string `json:"shortName" binding:"required,max=10"`
	}
	if err := c.MustBindWith(&createModel, binding.JSON); err != nil {
		log.Infof("Failed to parse platform creation model: %s", err.Error())
		return
	}

	platform := models.Platform{
		Name:      createModel.Name,
		ShortName: createModel.ShortName,
	}
	if err := p.repo.AddPlatform(&platform); err != nil {
		abortWithRelevantError(err, c)
		return
	}

	c.JSON(http.StatusCreated, &platform)
}

func (p *PlatformController) Update(c *gin.Context) {
	var updateModel struct {
		Name      string `json:"name" binding:"required,max=200"`
		ShortName string `json:"shortName" binding:"required,max=10"`
	}
	if err := c.MustBindWith(&updateModel, binding.JSON); err != nil {
		log.Infof("Failed to parse platform creation model: %s", err.Error())
		return
	}
	id, err := parseUuidFromParam(c)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	platform := models.Platform{
		Name:      updateModel.Name,
		ShortName: updateModel.ShortName,
	}
	if err := p.repo.UpdatePlatform(id, &platform); err != nil {
		abortWithRelevantError(err, c)
		return
	}

	platform.Id = id
	c.JSON(http.StatusOK, &platform)
}

func (p *PlatformController) Delete(c *gin.Context) {
	id, err := parseUuidFromParam(c)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := p.repo.DeletePlatform(id); err != nil {
		abortWithRelevantError(err, c)
		return
	}

	c.Status(http.StatusOK)
}

func (p *PlatformController) SetupRoutes(engine *gin.Engine, basePath string) {
	baseUrl := fmt.Sprintf("%s/api/v0/platforms", basePath)

	engine.GET(baseUrl, p.Get)
	engine.GET(baseUrl+"/:id", p.GetById)
	engine.POST(baseUrl, p.Create)
	engine.PUT(baseUrl+"/:id", p.Update)
}
