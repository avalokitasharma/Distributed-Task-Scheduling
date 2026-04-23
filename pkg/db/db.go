package db

import (
	"database/sql"
	"sync"

	"github.com/avalokitasharma/job-scheduler/pkg/config"
	_ "github.com/lib/pq"
)

var (
	db   *sql.DB
	once sync.Once
)

func Get() *sql.DB {
	once.Do(func() {
		var err error
		db, err = sql.Open("postgres", config.App.DBURL)
		if err != nil {
			panic(err)
		}
	})
	return db
}
