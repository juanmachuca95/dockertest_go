package handlers

import (
	"dockertest_go/domain"
	"dockertest_go/repository"

	"github.com/kataras/iris/v12"
)

type UserHandler struct {
	repo repository.Repo
}

func NewUserHander(repo repository.Repo) *UserHandler {
	return &UserHandler{repo: repo}
}

func (us *UserHandler) Create(ctx iris.Context) {
	var u domain.User
	if err := ctx.ReadJSON(&u); err != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("User creation failure").DetailErr(err))
		return
	}

	err := us.repo.SaveUser(ctx, u)
	if err != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("Failed to connect db").DetailErr(err))
	}

	ctx.StatusCode(iris.StatusCreated)
}
