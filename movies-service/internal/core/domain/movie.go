package domain

import (
	"strings"
)

type Movie struct {
	ID    int32
	Title string
	Year  string
}

func (m *Movie) Validate() error {
	if strings.TrimSpace(m.Title) == "" {
		return ErrInvalidMovieData
	}
	if strings.TrimSpace(m.Year) == "" {
		return ErrInvalidMovieData
	}
	return nil
}

func NewMovie(title, year string) (*Movie, error) {
	movie := &Movie{
		Title: title,
		Year:  year,
	}

	if err := movie.Validate(); err != nil {
		return nil, err
	}

	return movie, nil
}
