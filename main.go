package main

import (
	"github.com/Geepr/game/database"
	"github.com/KowalskiPiotr98/gotabase"
	_ "github.com/lib/pq"
)

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
