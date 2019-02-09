package database

import (
	"fmt"
	"github.com/gobuffalo/packr/v2"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stashapp/stash/internal/logger"
)

var DB *sqlx.DB

func Initialize(databasePath string) {
	runMigrations(databasePath)

	// https://github.com/mattn/go-sqlite3
	conn, err := sqlx.Open("sqlite3", "file:"+databasePath+"?_fk=true")
	conn.SetMaxOpenConns(10)
	if err != nil {
		logger.Fatalf("db.Open(): %q\n", err)
	}
	DB = conn
}

func Reset(databasePath string) {
	_, _ = DB.Exec("PRAGMA writable_schema = 1;")
	_, _ = DB.Exec("delete from sqlite_master where type in ('table', 'index', 'trigger');")
	_, _ = DB.Exec("PRAGMA writable_schema = 0;")
	_, _ = DB.Exec("VACUUM;")
	_, _ = DB.Exec("PRAGMA INTEGRITY_CHECK;")
	runMigrations(databasePath)
}

// Migrate the database
func runMigrations(databasePath string) {
	migrationsBox := packr.New("Migrations Box", "./migrations")
	packrSource := &Packr2Source{
		Box:        migrationsBox,
		Migrations: source.NewMigrations(),
	}
	s, _ := WithInstance(packrSource)
	m, err := migrate.NewWithSourceInstance(
		"packr2",
		s,
		fmt.Sprintf("sqlite3:%s", databasePath),
	)
	if err != nil {
		panic(err.Error())
	}

	err = m.Steps(1)
	if err != nil {
		//panic(err.Error()) // TODO
	}
}