package main

import (
	"dockertest_go/db"
	"dockertest_go/handlers"
	"dockertest_go/repository"

	"github.com/joho/godotenv"
	"github.com/kataras/iris/v12"
)

func main() {
	loadEnv()
	app := iris.New()

	db, err := db.ConnectDB()
	if err != nil {
		panic(err)
	}
	repo := repository.NewRepo(db)
	u := handlers.NewUserHander(repo)
	user := app.Party("/users")
	{
		user.Post("/", u.Create)
	}

	app.Listen(":8080")
}

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
}
