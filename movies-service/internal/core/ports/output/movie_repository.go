package output

import (
	"context"

	"github.com/FranciscoHonorat/movies/movies-service/internal/core/domain"
)

type MovieRepository interface {
	GetMovie(ctx context.Context, id int32) (*domain.Movie, error)
	ListMovies(ctx context.Context, filters ListFilters, pagination Pagination, sorting Sorting) ([]domain.Movie, error)
	Count(ctx context.Context, filters ListFilters) (int32, error)
	CreateMovie(ctx context.Context, movie *domain.Movie) (*domain.Movie, error)
	DeleteMovie(ctx context.Context, id int32) error
}

type ListFilters struct {
	Title string
	Year  string
}

type Pagination struct {
	Page  int32
	Limit int32
}

type Sorting struct {
	By string
}
