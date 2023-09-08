package controllers

import "github.com/gin-gonic/gin"

// Routable defines a struct used for route handling. BasePath should be used as a global path for the entire application. It will be prepended to the default controller path.
type Routable interface {
	SetupRoutes(engine *gin.Engine, basePath string)
}
