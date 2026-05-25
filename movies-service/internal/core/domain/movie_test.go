package domain

import (
	"testing"
)

func TestNewMovieEmptyTitle(t *testing.T) {
	movie, err := NewMovie("", "2020")
	if movie != nil {
		t.Error("Expected nil, got a movie")
	}
	if err != ErrInvalidMovieData {
		t.Errorf("Expected error '%v', got '%v'", ErrInvalidMovieData, err)
	}
}

func TestNewMovie_EmptyYear(t *testing.T) {
	movie, err := NewMovie("Inception", "")
	if movie != nil {
		t.Error("Expected nil, got a movie")
	}
	if err != ErrInvalidMovieData {
		t.Errorf("Expected error '%v', got '%v'", ErrInvalidMovieData, err)
	}
}

func TestNewMovie_ValidData(t *testing.T) {
	movie, err := NewMovie("Inception", "2010")
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}
	if movie == nil {
		t.Error("Expected a movie, got nil")
	}
}
