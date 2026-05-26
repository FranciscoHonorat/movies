package seed

import (
	"context"
	"encoding/json"
	"os"

	"github.com/FranciscoHonorat/movies/movies-service/internal/core/domain"
	"github.com/FranciscoHonorat/movies/movies-service/internal/core/ports/output"
)

type MovieSeedDynamodb struct {
	ID    int32  `json:"id"`
	Title string `json:"title"`
	Year  string `json:"year"`
}

func SeedMovie(ctx context.Context, repo output.MovieRepository, filePath string) error {
	count, err := repo.Count(ctx, output.ListFilters{})
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var movies []MovieSeedDynamodb
	err = json.Unmarshal(data, &movies)
	if err != nil {
		return err
	}
	for _, movie := range movies {
		_, err := repo.CreateMovie(ctx, &domain.Movie{
			ID:    movie.ID,
			Title: movie.Title,
			Year:  movie.Year,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
