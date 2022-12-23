package utils

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/ory/dockertest/v3"
)

var (
	DBDocker *sql.DB
)

func InitDB(path string) (*dockertest.Pool, *dockertest.Resource, *sql.DB) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	resource, err := pool.Run("mysql", "5.7", []string{"MYSQL_ROOT_PASSWORD=root"})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		var err error
		DBDocker, err = sql.Open("mysql", fmt.Sprintf("root:root@(localhost:%s)/mysql", resource.GetPort("3306/tcp"))) // database name need to be mysql
		if err != nil {
			return err
		}
		return DBDocker.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	// Init Migration
	initMigration(path)

	return pool, resource, DBDocker
}

func CloseDB(pool *dockertest.Pool, resource *dockertest.Resource) {
	if err := pool.Purge(resource); err != nil {
		phrase := fmt.Sprintf("Could not purge resource: %s", err)
		log.Fatal(phrase)
	}
}

func initMigration(path string) {
	driver, err := mysql.WithInstance(DBDocker, &mysql.Config{})
	if err != nil {
		panic(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+path, // -> path = ../db
		"mysql",        // Need to be mysql
		driver,
	)

	if err != nil {
		panic(err)
	}

	m.Steps(2)
}
