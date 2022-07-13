package migration

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	migratemysql "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Up(dburl string) {
	var (
		init bool
		err  error
		db   *sql.DB
	)
	for i := 0; i < 10; i++ {
		db, err = sql.Open("mysql", dburl)
		if err != nil {
			time.Sleep(3 * time.Second)
			fmt.Printf("fail to connect to db %s, retrying", err)
			continue
		}
		init = true
		break
	}
	if !init {
		panic(fmt.Sprintf("failed to connect db %s", err.Error()))
	}

	driver, err := migratemysql.WithInstance(db, &migratemysql.Config{})
	if err != nil {
		panic(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file:///migrations",
		"mysql", driver)
	if err != nil {
		panic(err)
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		panic(err)
	}
}
