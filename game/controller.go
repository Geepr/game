package game

import (
	"fmt"
	"github.com/Geepr/game/utils"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func getRoute(c *gin.Context) {
	var query struct {
		Title     string    `form:"title"`
		SortOrder SortOrder `form:"order"`
		PageIndex int       `form:"page"`
		PageSize  int       `form:"size"`
	}
	if err := c.BindQuery(&query); err != nil {
		log.Infof("Failed to bind game query: %s", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	games, totalItems, err := getGames(query.Title, query.PageIndex, query.PageSize, query.SortOrder)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	response := struct {
		Games      []*Game `json:"games"`
		Page       int     `json:"page"`
		PageSize   int     `json:"pageSize"`
		TotalPages int     `json:"totalPages"`
	}{
		Games:      games,
		Page:       query.PageIndex,
		PageSize:   query.PageSize,
		TotalPages: utils.GetPagesFromItems(totalItems, query.PageSize),
	}
	c.JSON(http.StatusOK, response)
}

func getByIdRoute(c *gin.Context) {
	lookupUuid, err := utils.ParseUuidFromParam(c)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	game, err := getGameById(lookupUuid)
	if err != nil {
		utils.AbortWithRelevantError(err, c)
		return
	}

	c.JSON(http.StatusOK, game)
}

func createRoute(c *gin.Context) {
	var createModel struct {
		Title       string `json:"title" binding:"required,max=200"`
		Description string `json:"description" binding:"max=2000"`
	}
	if err := c.BindJSON(&createModel); err != nil {
		log.Infof("Failed to parse game creation model: %s", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	game := Game{
		Title:       createModel.Title,
		Description: utils.GetNilIfDefault(createModel.Description),
		Archived:    false,
	}
	if err := addGame(&game); err != nil {
		utils.AbortWithRelevantError(err, c)
		return
	}

	c.JSON(http.StatusCreated, &game)
}

func updateRoute(c *gin.Context) {
	var updateModel struct {
		Title       string `json:"title" binding:"required,max=200"`
		Description string `json:"description" binding:"max=2000"`
		Archived    bool   `json:"archived"`
	}
	if err := c.BindJSON(&updateModel); err != nil {
		log.Infof("Failed to parse game updateRoute model: %s", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	id, err := utils.ParseUuidFromParam(c)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	game := Game{
		Title:       updateModel.Title,
		Description: utils.GetNilIfDefault(updateModel.Description),
		Archived:    updateModel.Archived,
	}
	if err := updateGame(id, &game); err != nil {
		utils.AbortWithRelevantError(err, c)
		return
	}

	// this is here just so that it displays properly when returned to the user
	game.Id = id
	c.JSON(http.StatusOK, &game)
}

func deleteRoute(c *gin.Context) {
	id, err := utils.ParseUuidFromParam(c)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = deleteGame(id)
	if err != nil {
		utils.AbortWithRelevantError(err, c)
		return
	}

	c.Status(http.StatusOK)
}

func SetupRoutes(engine *gin.Engine, basePath string) {
	baseUrl := fmt.Sprintf("%s/api/v0/games", basePath)

	engine.GET(baseUrl, getRoute)
	engine.GET(baseUrl+"/:id", getByIdRoute)
	engine.POST(baseUrl, createRoute)
	engine.PUT(baseUrl+"/:id", updateRoute)
	engine.DELETE(baseUrl+"/:id", deleteRoute)
}
