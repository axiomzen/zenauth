package data

import (
	"database/sql"
	"errors"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/axiomzen/zenauth/config"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database"
	"github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"
)

func Migrate(conf *config.ZENAUTHConfig) {

	var err = errors.New("temp")
	var numtries uint16
	var driver database.Driver

	log.Infof("Retries: %d\n", conf.PostgreSQLRetryNumTimes)

	for err != nil && numtries < conf.PostgreSQLRetryNumTimes {

		sslMode := "require"
		if !*conf.PostgreSQLSSL {
			sslMode = "disable"
		}
		pgURL := fmt.Sprintf("postgres://%s:%s@%s:%v/%s?sslmode=%s",
			conf.PostgreSQLUsername,
			conf.PostgreSQLPassword,
			conf.PostgreSQLHost,
			conf.PostgreSQLPort,
			conf.PostgreSQLDatabase,
			sslMode)
		log.Infoln(pgURL)
		var db *sql.DB
		db, err = sql.Open("postgres", pgURL)
		log.Infoln(err)
		if err != nil {
			log.Error("Migration DB Connection failed. Retrying ...")
			continue
		}
		defer db.Close()

		driver, err = postgres.WithInstance(db, &postgres.Config{})
		log.Infof("Num Tries: %d, Err: %v\n", numtries, err)
	}

	if err != nil {
		log.Fatal(err.Error())
	}

	// migrate db
	m, migrateErr := migrate.NewWithDatabaseInstance("file://data/migrations", "postgres", driver)
	defer m.Close()
	if migrateErr != nil {
		log.Fatal(migrateErr.Error())
	} else if upErr := m.Up(); err != nil {
		log.Fatal(upErr.Error())
	} else if sErr, dbErr := m.Close(); sErr != nil {
		log.Fatal(sErr.Error())
	} else if dbErr != nil {
		log.Fatal(dbErr.Error())
	}

}
