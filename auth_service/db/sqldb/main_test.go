package sqldb

import (
	"auth_service/util"
	"log"
	"os"
	"testing"
)

var testStore *SqlDB

func TestMain(m *testing.M) {
	envConf, err := util.LoadConfig("../../")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	f := &MySQLFactory{
		Config: &DBConfig{
			SQLAddr:     envConf.DbSource,
			EnableLog:   false,
			MaxConn:     envConf.MaxConnect,
			IdleConn:    envConf.IdleConnect,
			MaxLifeTime: envConf.MaxLifeTime,
		},
	}

	db, err := f.CreateDB()
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}

	testStore = db

	// Add migration to ensure table exists
	err = db.MigrateDB()
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	os.Exit(m.Run())
}
