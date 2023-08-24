package database

import (
	"embed"
	"github.com/KowalskiPiotr98/gotabase"
	"github.com/KowalskiPiotr98/gotabase/migrations"
)

var (
	//go:embed sql
	migrationFiles embed.FS
)

func RunMigrations(connector gotabase.Connector) error {
	return migrations.Migrate(connector, migrationFiles)
}
