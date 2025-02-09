package db

import (
	"database/sql"
	"github.com/SuperMatch/utilities"
	"github.com/SuperMatch/zapLogger"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
)

func AutoMigration(db *sql.DB) error {
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Println("error in getting instance", err)
	}

	rootDirectory := utilities.RootDir()
	log.Println(rootDirectory)

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+rootDirectory+"/migrations",
		"mysql",
		driver,
	)
	if err != nil {
		log.Println("error in getting migrate instance", err)
	}
	err = m.Up()
	if err != nil {
		log.Println("Failed to migrate database", err)
	}

	zapLogger.Logger.Info("Database Migration successful.")

	return nil
}
