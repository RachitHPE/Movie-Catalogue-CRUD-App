package controller

import (
	db "catalogue-app/internal/database"
	"catalogue-app/internal/pkg/model"
	"context"
)

type MovieController struct {
	dbClient db.DBCLientIntfc
}

type ControllerIntfc interface {
	GetMovies(ctx context.Context) ([]model.MovieInfo, error)
	GetMovieByID(ctx context.Context, movieID string) (model.MovieInfo, error)
	CreateMovie(ctx context.Context, movieInfo model.MovieInfo) (model.MovieInfo, error)
	UpdateMovie(ctx context.Context, movieInfo model.MovieInfo, movieID string) (model.MovieInfo, error)
	DeleteMovie(ctx context.Context, movieID string) error
}

func NewMovieController(dbClient db.DBCLientIntfc) *MovieController {
	return &MovieController{dbClient: dbClient}
}

func (movieController MovieController) GetMovies(ctx context.Context) ([]model.MovieInfo, error) {
	movies, err := movieController.dbClient.GetMovies(ctx)
	if err != nil {
		return nil, err
	}

	return movies, nil
}

func (movieController MovieController) GetMovieByID(ctx context.Context, movieID string) (model.MovieInfo, error) {
	movie, err := movieController.dbClient.GetMovieByID(ctx, movieID)
	if err != nil {
		return model.MovieInfo{}, err
	}

	return movie, nil
}

func (movieController MovieController) CreateMovie(ctx context.Context, movieInfo model.MovieInfo) (model.MovieInfo, error) {
	movie, err := movieController.dbClient.CreateMovie(ctx, movieInfo)
	if err != nil {
		return model.MovieInfo{}, err
	}

	return movie, nil
}

func (movieController MovieController) UpdateMovie(ctx context.Context, movieInfo model.MovieInfo, movieID string) (model.MovieInfo, error) {
	movie, err := movieController.dbClient.UpdateMovie(ctx, movieInfo, movieID)
	if err != nil {
		return model.MovieInfo{}, err
	}

	return movie, nil
}

func (movieController MovieController) DeleteMovie(ctx context.Context, movieID string) error {
	err := movieController.dbClient.DeleteMovie(ctx, movieID)
	if err != nil {
		return err
	}

	return nil
}
