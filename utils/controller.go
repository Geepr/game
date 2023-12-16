package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type id struct {
	Id string `form:"id" uri:"id" binding:"required,uuid"`
}

func ParseUuidFromParam(c *gin.Context) (uuid.UUID, error) {
	var id id
	err := c.BindUri(&id)
	if err != nil {
		log.Infof("Failed to parse uuid: %s", err.Error())
		return uuid.Nil, err
	}
	return uuid.FromString(id.Id)
}

func GetNilIfDefault[T comparable](value T) *T {
	var defaultValue T
	if value == defaultValue {
		return nil
	}
	return &value
}

func AbortWithRelevantError(err error, c *gin.Context) {
	if errors.Is(err, DuplicateDataErr) {
		c.AbortWithStatus(http.StatusBadRequest)
	} else if errors.Is(err, DataNotFoundErr) {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func GetPagesFromItems(totalItems int, pageSize int) int {
	if totalItems%pageSize == 0 {
		return totalItems / pageSize
	}
	return (totalItems / pageSize) + 1
}
