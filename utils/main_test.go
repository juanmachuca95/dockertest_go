package main

import (
	"context"
	"database/sql"
	"dockertest_go/domain"
	"dockertest_go/repository"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/ory/dockertest/v3"
)

var dbDocker *sql.DB

func TestMain(m *testing.M) {
	pool, resource := initDB()
	code := m.Run()
	closeDB(pool, resource)
	os.Exit(code)
}

func initDB() (*dockertest.Pool, *dockertest.Resource) {
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
		dbDocker, err = sql.Open("mysql", fmt.Sprintf("root:root@(localhost:%s)/mysql", resource.GetPort("3306/tcp"))) // database name need to be mysql
		if err != nil {
			return err
		}
		return dbDocker.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}
	// Init Migration
	initMigration("./db")

	return pool, resource
}

func closeDB(pool *dockertest.Pool, resource *dockertest.Resource) {
	if err := pool.Purge(resource); err != nil {
		phrase := fmt.Sprintf("Could not purge resource: %s", err)
		log.Fatal(phrase)
	}
}

func initMigration(path string) {
	driver, err := mysql.WithInstance(dbDocker, &mysql.Config{})
	if err != nil {
		panic(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://../db",
		"mysql", // Need to be mysql
		driver,
	)

	if err != nil {
		panic(err)
	}

	m.Steps(2)
}

func TestSaveUser(t *testing.T) {
	repo := repository.NewRepo(dbDocker)

	u := domain.User{
		Nombre:   "Test user",
		Apellido: "Test apellido",
	}

	err := repo.SaveUser(context.TODO(), u)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}
