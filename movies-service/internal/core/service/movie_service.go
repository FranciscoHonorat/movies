package service

import (
	"context"

	"github.com/FranciscoHonorat/movies/movies-service/internal/core/domain"
	"github.com/FranciscoHonorat/movies/movies-service/internal/core/ports/output"
)

type MovieService struct {
	repo output.MovieRepository
}

func NewMovieService(repo output.MovieRepository) *MovieService {
	return &MovieService{
		repo: repo,
	}
}

func (m *MovieService) GetMovie(ctx context.Context, id int32) (*domain.Movie, error) {
	if id <= 0 {
		return nil, domain.ErrInvalidMovieData
	}

	movie, err := m.repo.GetMovie(ctx, id)

	if err != nil {
		return nil, err
	}

	return movie, nil
}

func (m *MovieService) CreateMovie(ctx context.Context, movie *domain.Movie) (*domain.Movie, error) {
	//validate
	if movie == nil {
		return nil, domain.ErrInvalidMovieData
	}
	if err := movie.Validate(); err != nil {
		return nil, domain.ErrInvalidMovieData
	}

	createdMovie, err := m.repo.CreateMovie(ctx, movie)
	if err != nil {
		return nil, err
	}

	return createdMovie, nil
}

func (m *MovieService) ListMovies(ctx context.Context, filters output.ListFilters, pagination output.Pagination, sorting output.Sorting) ([]domain.Movie, int32, error) {
	movies, err := m.repo.ListMovies(ctx, filters, pagination, sorting)
	if err != nil {
		return nil, 0, err
	}

	total, err := m.repo.Count(ctx, filters)
	if err != nil {
		return nil, 0, err
	}

	return movies, total, nil
}

func (m *MovieService) DeleteMovie(ctx context.Context, id int32) error {
	if id <= 0 {
		return domain.ErrMovieNotFound
	}
	if err := m.repo.DeleteMovie(ctx, id); err != nil {
		return err
	}
	return nil
}
