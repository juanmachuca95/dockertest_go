package repository

import (
	"context"
	"database/sql"
	"dockertest_go/domain"
	"dockertest_go/utils"
	"os"
	"testing"
)

var db *sql.DB
var path string = "../db"

func TestMain(m *testing.M) {
	pool, resource, dbDocker := utils.InitDB(path)
	db = dbDocker
	code := m.Run()
	utils.CloseDB(pool, resource)
	os.Exit(code)
}

func TestSaveUser(t *testing.T) {
	repo := NewRepo(db)

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
