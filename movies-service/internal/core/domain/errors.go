package domain

import "errors"

var (
	//
	ErrInvalidMovieData = errors.New("Invalid movie data")

	//
	ErrMovieNotFound = errors.New("Movie not found")

	//
	ErrUnauthorized = errors.New("Unauthorized")

	//
	ErrInternalServer = errors.New("Internal server error")
)
