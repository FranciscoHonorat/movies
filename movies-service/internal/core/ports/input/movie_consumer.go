package input

import (
	"context"

	"github.com/FranciscoHonorat/movies/movies-service/internal/core/domain"
)

type MovieConsumer interface {
	Consume(ctx context.Context, handler func(movie domain.Movie) error) error
}
