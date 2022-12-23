package repository

import (
	"context"
	"database/sql"
	"dockertest_go/domain"
	"errors"
)

type Repo interface {
	SaveUser(ctx context.Context, u domain.User) error
}

type repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) Repo {
	return &repo{db: db}
}

func (r *repo) SaveUser(ctx context.Context, u domain.User) error {
	stmt, err := r.db.Prepare("INSERT INTO users (nombre, apellido) VALUES(?,?);")
	if err != nil {
		return err
		/*
			 		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
						Title("cannot prepare create user into db").DetailErr(err))
		*/
	}

	result, err := stmt.Exec(u.Nombre, u.Apellido)
	if err != nil {
		return err
		/* ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
		Title("cannot execute create user into db").DetailErr(err)) */
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected <= 0 {
		return errors.New("failed to rowsAffected")
	}
	return err
}
