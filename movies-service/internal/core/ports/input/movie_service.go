package input

import (
	"context"

	"github.com/FranciscoHonorat/movies/movies-service/internal/core/domain"
	"github.com/FranciscoHonorat/movies/movies-service/internal/core/ports/output"
)

type MovieService interface {
	GetMovie(ctx context.Context, id int32) (*domain.Movie, error)
	ListMovies(ctx context.Context, filters output.ListFilters, pagination output.Pagination, sorting output.Sorting) ([]domain.Movie, int32, error)
	CreateMovie(ctx context.Context, movie *domain.Movie) (*domain.Movie, error)
	DeleteMovie(ctx context.Context, id int32) error
}
