package main

import (
	"github.com/Geepr/game/database"
	"github.com/KowalskiPiotr98/gotabase"
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
}
