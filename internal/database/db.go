package db

import (
	"catalogue-app/internal/pkg/log"
	"catalogue-app/internal/pkg/model"
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type MovieClient struct {
	dbClient *gorm.DB
}

func NewClient(dbClient *gorm.DB) *MovieClient {
	return &MovieClient{dbClient: dbClient}
}

type DBCLientIntfc interface {
	GetMovies(ctx context.Context) ([]model.MovieInfo, error)
	GetMovieByID(ctx context.Context, movieID string) (model.MovieInfo, error)
	CreateMovie(ctx context.Context, movieInfo model.MovieInfo) (model.MovieInfo, error)
	UpdateMovie(ctx context.Context, movieInfo model.MovieInfo, movieID string) (model.MovieInfo, error)
	DeleteMovie(ctx context.Context, movieID string) error
}

func (movieClient MovieClient) GetMovies(ctx context.Context) ([]model.MovieInfo, error) {
	movies := []model.MovieInfo{}

	if err := movieClient.dbClient.Find(&movies).Error; err != nil {
		log.Errorf(ctx, "unable to retrieve movies from database")

		return nil, err
	}

	if len(movies) == 0 {
		log.Errorf(ctx, "no movies available in the database")

		return nil, errors.New("no movies found")
	}

	return movies, nil
}

func (movieClient MovieClient) GetMovieByID(ctx context.Context, movieID string) (model.MovieInfo, error) {
	var movie model.MovieInfo

	if err := movieClient.dbClient.First(&movie, movieID).Error; err != nil {
		log.Errorf(ctx, "unable to retrieve movie for id: %s", movieID)

		return model.MovieInfo{}, err
	}

	return movie, nil
}

func (movieClient MovieClient) CreateMovie(ctx context.Context, movieInfo model.MovieInfo) (model.MovieInfo, error) {
	movieInfo.CreatedAt = time.Now()

	tx := movieClient.dbClient.Begin()
	if err := tx.Create(&movieInfo).Error; err != nil {
		tx.Rollback()
		log.Errorf(ctx, "error creating movie in database")

		return model.MovieInfo{}, err
	}

	tx.Commit()

	return movieInfo, nil
}

func (movieClient MovieClient) UpdateMovie(ctx context.Context, movieInfo model.MovieInfo, movieID string) (model.MovieInfo, error) {
	var movie model.MovieInfo

	if err := movieClient.dbClient.First(&movie, movieID).Error; err != nil {
		log.Errorf(ctx, "unable to retrieve movie for id: %s", movieID)

		return model.MovieInfo{}, err
	}

	movieInfo.CreatedAt = movie.CreatedAt
	movieInfo.UpdatedAt = time.Now()

	tx := movieClient.dbClient.Begin()
	if err := tx.Save(&movieInfo).Error; err != nil {
		tx.Rollback()
		log.Errorf(ctx, "error updating movie details in database")

		return model.MovieInfo{}, err
	}

	tx.Commit()

	return movieInfo, nil
}

func (movieClient MovieClient) DeleteMovie(ctx context.Context, movieID string) error {
	var movie model.MovieInfo

	if err := movieClient.dbClient.First(&movie, movieID).Error; err != nil {
		log.Errorf(ctx, "unable to retrieve movie for id: %s", movieID)

		return err
	}

	tx := movieClient.dbClient.Begin()
	if err := tx.Delete(&movie).Error; err != nil {
		tx.Rollback()
		log.Errorf(ctx, "error deleting movie in database")

		return err
	}

	tx.Commit()

	return nil
}
