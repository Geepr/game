package controllers

import (
	"errors"
	"github.com/Geepr/game/repositories"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type id struct {
	Id string `form:"id" uri:"id" binding:"required,uuid"`
}

func parseUuidFromParam(c *gin.Context) (uuid.UUID, error) {
	var id id
	err := c.BindUri(&id)
	if err != nil {
		log.Infof("Failed to parse uuid: %s", err.Error())
		return uuid.Nil, err
	}
	return uuid.FromString(id.Id)
}

func getNilIfDefault[T comparable](value T) *T {
	var defaultValue T
	if value == defaultValue {
		return nil
	}
	return &value
}

func abortWithRelevantError(err error, c *gin.Context) {
	if errors.Is(err, repositories.DuplicateDataErr) {
		c.AbortWithStatus(http.StatusBadRequest)
	} else if errors.Is(err, repositories.DataNotFoundErr) {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}
