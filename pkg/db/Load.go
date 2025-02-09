package db

import (
	"database/sql"
	"github.com/SuperMatch/zapLogger"
	"log"
	"os"
	"time"

	"github.com/SuperMatch/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// var DbConn *pg.DB

var GlobalOrm *gorm.DB

func Load(conf *config.Database, appEnv string) (*sql.DB, error) {

	dbConfig := mysql.Config{
		DSN: conf.USER + ":" + conf.PASS + "@tcp(" + conf.HOST + ":" + conf.PORT + ")/" + conf.DB + "?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=true",
	}

	GlobalDB, err := sql.Open("mysql", dbConfig.DSN)

	if err != nil {
		log.Fatal(err)
	}

	pingErr := GlobalDB.Ping()

	if pingErr != nil {
		log.Fatal(pingErr)
	}
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,         // Disable color
		},
	)

	GlobalOrm, err = gorm.Open(mysql.New(mysql.Config{
		Conn: GlobalDB,
	}), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		log.Fatal(pingErr)
	}

	zapLogger.Logger.Info("DB connection successful.")

	err = AutoMigration(GlobalDB)
	if err != nil {
		log.Fatal(err)
	}

	return GlobalDB, nil
}
