package platform

import (
	"fmt"
	"github.com/Geepr/game/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func getRoute(c *gin.Context) {
	var query struct {
		Name      string    `form:"name"`
		SortOrder SortOrder `form:"order"`
		PageIndex int       `form:"index"`
		PageSize  int       `form:"size"`
	}
	if err := c.MustBindWith(&query, binding.Query); err != nil {
		log.Infof("Failed to bind platform query: %s", err.Error())
		return
	}

	platforms, err := getPlatforms(query.Name, query.PageIndex, query.PageSize, query.SortOrder)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, platforms)
}

func getByIdRoute(c *gin.Context) {
	id, err := utils.ParseUuidFromParam(c)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	game, err := getPlatformById(id)
	if err != nil {
		utils.AbortWithRelevantError(err, c)
		return
	}

	c.JSON(http.StatusOK, game)
}

func createRoute(c *gin.Context) {
	var createModel struct {
		Name      string `json:"name" binding:"required,max=200"`
		ShortName string `json:"shortName" binding:"required,max=10"`
	}
	if err := c.MustBindWith(&createModel, binding.JSON); err != nil {
		log.Infof("Failed to parse platform creation model: %s", err.Error())
		return
	}

	platform := Platform{
		Name:      createModel.Name,
		ShortName: createModel.ShortName,
	}
	if err := addPlatform(&platform); err != nil {
		utils.AbortWithRelevantError(err, c)
		return
	}

	c.JSON(http.StatusCreated, &platform)
}

func updateRoute(c *gin.Context) {
	var updateModel struct {
		Name      string `json:"name" binding:"required,max=200"`
		ShortName string `json:"shortName" binding:"required,max=10"`
	}
	if err := c.MustBindWith(&updateModel, binding.JSON); err != nil {
		log.Infof("Failed to parse platform creation model: %s", err.Error())
		return
	}
	id, err := utils.ParseUuidFromParam(c)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	platform := Platform{
		Name:      updateModel.Name,
		ShortName: updateModel.ShortName,
	}
	if err := updatePlatform(id, &platform); err != nil {
		utils.AbortWithRelevantError(err, c)
		return
	}

	platform.Id = id
	c.JSON(http.StatusOK, &platform)
}

func DeleteRoute(c *gin.Context) {
	id, err := utils.ParseUuidFromParam(c)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := deletePlatform(id); err != nil {
		utils.AbortWithRelevantError(err, c)
		return
	}

	c.Status(http.StatusOK)
}

func SetupRoutes(engine *gin.Engine, basePath string) {
	baseUrl := fmt.Sprintf("%s/api/v0/platforms", basePath)

	engine.GET(baseUrl, getRoute)
	engine.GET(baseUrl+"/:id", getByIdRoute)
	engine.POST(baseUrl, createRoute)
	engine.PUT(baseUrl+"/:id", updateRoute)
}
