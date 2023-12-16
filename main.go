package main

import (
	"github.com/Geepr/game/controllers"
	"github.com/Geepr/game/database"
	"github.com/Geepr/game/game"
	"github.com/Geepr/game/services"
	"github.com/KowalskiPiotr98/gotabase"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, DisableColors: true})
	log.SetOutput(os.Stdout)
	//todo: configurable
	log.SetLevel(log.InfoLevel)
}

func main() {
	err := gotabase.InitialiseConnection("user=postgres dbname=geepr password=postgres sslmode=disable", "postgres")
	if err != nil {
		panic(err)
	}
	defer gotabase.CloseConnection()
	if err = database.RunMigrations(gotabase.GetConnection()); err != nil {
		panic(err)
	}

	router := setupEngine()

	//todo: configurable address
	if err := router.Run("localhost:5500"); err != nil {
		log.Panicf("Server failed while listening: %s", err.Error())
	}
}

func setupEngine() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(services.GetGinLogger())
	//todo: trusted proxies

	//todo: base path config
	basePath := ""
	game.SetupRoutes(router, basePath)
	platformController := controllers.CreatePlatformController()
	platformController.SetupRoutes(router, basePath)
	releaseController := controllers.CreateGameReleaseController()
	releaseController.SetupRoutes(router, basePath)
	releasePlatformController := controllers.CreateGameReleasePlatformController()
	releasePlatformController.SetupRoutes(router, basePath)

	return router
}
