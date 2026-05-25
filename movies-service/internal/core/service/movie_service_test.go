package service

import (
	"context"
	"testing"

	"github.com/FranciscoHonorat/movies/movies-service/internal/core/domain"
	"github.com/FranciscoHonorat/movies/movies-service/internal/core/ports/output"
)

type MockMovieRepository struct {
	GetMovieFn    func(ctx context.Context, id int32) (*domain.Movie, error)
	ListMoviesFn  func(ctx context.Context, filters output.ListFilters, pagination output.Pagination, sorting output.Sorting) ([]domain.Movie, error)
	CountFn       func(ctx context.Context, filters output.ListFilters) (int32, error)
	CreateMovieFn func(ctx context.Context, movie *domain.Movie) (*domain.Movie, error)
	DeleteMovieFn func(ctx context.Context, id int32) error
}

func (m *MockMovieRepository) GetMovie(ctx context.Context, id int32) (*domain.Movie, error) {
	if m.GetMovieFn != nil {
		return m.GetMovieFn(ctx, id)
	}
	return nil, nil
}

func (m *MockMovieRepository) ListMovies(ctx context.Context, filters output.ListFilters, pagination output.Pagination, sorting output.Sorting) ([]domain.Movie, error) {
	if m.ListMoviesFn != nil {
		return m.ListMoviesFn(ctx, filters, pagination, sorting)
	}
	return []domain.Movie{}, nil
}

func (m *MockMovieRepository) Count(ctx context.Context, filters output.ListFilters) (int32, error) {
	if m.CountFn != nil {
		return m.CountFn(ctx, filters)
	}
	return 0, nil
}

func (m *MockMovieRepository) CreateMovie(ctx context.Context, movie *domain.Movie) (*domain.Movie, error) {
	if m.CreateMovieFn != nil {
		return m.CreateMovieFn(ctx, movie)
	}
	return movie, nil
}

func (m *MockMovieRepository) DeleteMovie(ctx context.Context, id int32) error {
	if m.DeleteMovieFn != nil {
		return m.DeleteMovieFn(ctx, id)
	}
	return nil
}

func TestGetMovie_InvalidID(t *testing.T) {
	//Mock repository
	repo := &MockMovieRepository{}
	service := NewMovieService(repo)
	_, err := service.GetMovie(context.Background(), -1)
	if err != domain.ErrInvalidMovieData {
		t.Errorf("Expected error for invalid ID, got %v", err)
		return
	}
}

func TestCreateMovie_InvalidData(t *testing.T) {
	//Mock repository
	repo := &MockMovieRepository{}
	service := NewMovieService(repo)
	_, err := service.CreateMovie(context.Background(), nil)

	if err != domain.ErrInvalidMovieData {
		t.Errorf("Expected error for nil movie, got %v", err)
		return
	}
}

func TestListMovies_Success(t *testing.T) {
	//Mock repository
	repo := &MockMovieRepository{
		ListMoviesFn: func(ctx context.Context, filters output.ListFilters, pagination output.Pagination, sorting output.Sorting) ([]domain.Movie, error) {
			return []domain.Movie{
				{ID: 1, Title: "Movie 1", Year: "2020"},
				{ID: 2, Title: "Movie 2", Year: "2021"},
			}, nil
		},
		CountFn: func(ctx context.Context, filters output.ListFilters) (int32, error) {
			return 2, nil
		},
	}
	service := NewMovieService(repo)
	movies, total, err := service.ListMovies(context.Background(), output.ListFilters{}, output.Pagination{Page: 1, Limit: 10}, output.Sorting{By: "title"})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}
	if len(movies) != 2 {
		t.Errorf("Expected 2 movies, got %d", len(movies))
		return
	}
	if total != 2 {
		t.Errorf("Expected total of 2 movies, got %d", total)
		return
	}
}

func TestDeleteMovie_InvalidID(t *testing.T) {
	//Mock repository
	repo := &MockMovieRepository{}
	service := NewMovieService(repo)
	err := service.DeleteMovie(context.Background(), -1)
	if err != domain.ErrMovieNotFound {
		t.Errorf("Expected error for invalid ID, got %v", err)
		return
	}
}

func TestDeleteMovie_Success(t *testing.T) {
	//Mock repository
	repo := &MockMovieRepository{
		DeleteMovieFn: func(ctx context.Context, id int32) error {
			if id == 1 {
				return nil
			}
			return domain.ErrMovieNotFound
		},
	}
	service := NewMovieService(repo)
	err := service.DeleteMovie(context.Background(), 1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}
}

func TestDeleteMovie_NotFound(t *testing.T) {
	//Mock repository
	repo := &MockMovieRepository{
		DeleteMovieFn: func(ctx context.Context, id int32) error {
			return domain.ErrMovieNotFound
		},
	}
	service := NewMovieService(repo)
	err := service.DeleteMovie(context.Background(), 2)
	if err != domain.ErrMovieNotFound {
		t.Errorf("Expected error for not found movie, got %v", err)
		return
	}
}

func TestCreateMovie_Success(t *testing.T) {
	//Mock repository
	repo := &MockMovieRepository{
		CreateMovieFn: func(ctx context.Context, movie *domain.Movie) (*domain.Movie, error) {
			movie.ID = 1
			return movie, nil
		},
	}
	service := NewMovieService(repo)
	_, err := service.CreateMovie(context.Background(), &domain.Movie{Title: "Inception", Year: "2010"})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}
}

func TestGetMovie_Success(t *testing.T) {
	//Mock repository
	repo := &MockMovieRepository{
		GetMovieFn: func(ctx context.Context, id int32) (*domain.Movie, error) {
			if id == 1 {
				return &domain.Movie{ID: 1, Title: "Inception", Year: "2010"}, nil
			}
			return nil, domain.ErrMovieNotFound
		},
	}
	service := NewMovieService(repo)
	movie, err := service.GetMovie(context.Background(), 1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}
	if movie == nil {
		t.Error("Expected a movie, got nil")
		return
	}
}

func TestGetMovie_NotFound(t *testing.T) {
	//Mock repository
	repo := &MockMovieRepository{
		GetMovieFn: func(ctx context.Context, id int32) (*domain.Movie, error) {
			return nil, domain.ErrMovieNotFound
		},
	}
	service := NewMovieService(repo)
	_, err := service.GetMovie(context.Background(), 2)
	if err != domain.ErrMovieNotFound {
		t.Errorf("Expected error for not found movie, got %v", err)
		return
	}
}

func TestListMovies_Empty(t *testing.T) {
	//Mock repository
	repo := &MockMovieRepository{
		ListMoviesFn: func(ctx context.Context, filters output.ListFilters, pagination output.Pagination, sorting output.Sorting) ([]domain.Movie, error) {
			return []domain.Movie{}, nil
		},
		CountFn: func(ctx context.Context, filters output.ListFilters) (int32, error) {
			return 0, nil
		},
	}
	service := NewMovieService(repo)
	movies, total, err := service.ListMovies(context.Background(), output.ListFilters{}, output.Pagination{Page: 1, Limit: 10}, output.Sorting{By: "title"})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}
	if len(movies) != 0 {
		t.Errorf("Expected 0 movies, got %d", len(movies))
		return
	}
	if total != 0 {
		t.Errorf("Expected total of 0 movies, got %d", total)
		return
	}
}

func TestListMovies_Error(t *testing.T) {
	//Mock repository
	repo := &MockMovieRepository{
		ListMoviesFn: func(ctx context.Context, filters output.ListFilters, pagination output.Pagination, sorting output.Sorting) ([]domain.Movie, error) {
			return nil, domain.ErrInvalidMovieData
		},
		CountFn: func(ctx context.Context, filters output.ListFilters) (int32, error) {
			return 0, nil
		},
	}
	service := NewMovieService(repo)
	_, _, err := service.ListMovies(context.Background(), output.ListFilters{}, output.Pagination{Page: 1, Limit: 10}, output.Sorting{By: "title"})
	if err != domain.ErrInvalidMovieData {
		t.Errorf("Expected error for invalid movie data, got %v", err)
		return
	}
}

func TestCreateMovie_Error(t *testing.T) {
	//Mock repository
	repo := &MockMovieRepository{
		CreateMovieFn: func(ctx context.Context, movie *domain.Movie) (*domain.Movie, error) {
			return nil, domain.ErrInvalidMovieData
		},
	}
	service := NewMovieService(repo)
	_, err := service.CreateMovie(context.Background(), &domain.Movie{Title: "Inception", Year: "2010"})
	if err != domain.ErrInvalidMovieData {
		t.Errorf("Expected error for invalid movie data, got %v", err)
		return
	}
}

func TestDeleteMovie_Error(t *testing.T) {
	//Mock repository
	repo := &MockMovieRepository{
		DeleteMovieFn: func(ctx context.Context, id int32) error {
			return domain.ErrMovieNotFound
		},
	}
	service := NewMovieService(repo)
	err := service.DeleteMovie(context.Background(), 1)
	if err != domain.ErrMovieNotFound {
		t.Errorf("Expected error for not found movie, got %v", err)
		return
	}
}
