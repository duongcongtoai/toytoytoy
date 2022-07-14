package migration

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	migratemysql "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Up(dburl string) {
	var (
		init   bool
		err    error
		db     *sql.DB
		driver database.Driver
	)
	for i := 0; i < 100; i++ {
		db, err = sql.Open("mysql", dburl)
		if err != nil {
			time.Sleep(4 * time.Second)
			fmt.Printf("fail to connect to db %s, retrying\n", err)
			continue
		}

		driver, err = migratemysql.WithInstance(db, &migratemysql.Config{})
		if err != nil {
			time.Sleep(4 * time.Second)
			fmt.Printf("fail to connect to db %s, retrying\n", err)
			continue
		}
		init = true

		defer db.Close()
		break
	}
	if !init {
		panic(fmt.Sprintf("failed to connect db %s", err.Error()))
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
