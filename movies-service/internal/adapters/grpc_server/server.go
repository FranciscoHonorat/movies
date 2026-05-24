package grpcserver

import (
	"context"

	"github.com/FranciscoHonorat/movies/movies-service/internal/core/domain"
	"github.com/FranciscoHonorat/movies/movies-service/internal/core/ports/input"
	"github.com/FranciscoHonorat/movies/movies-service/internal/core/ports/output"
	"github.com/FranciscoHonorat/movies/proto"
)

type GrpcServer struct {
	proto.UnimplementedMovieServiceServer
	service input.MovieService
}

func NewGrpcServer(service input.MovieService) *GrpcServer {
	return &GrpcServer{service: service}
}

func (s *GrpcServer) GetMovie(ctx context.Context, req *proto.GetMovieRequest) (*proto.GetMovieResponse, error) {
	movie, err := s.service.GetMovie(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &proto.GetMovieResponse{
		Movie: &proto.Movie{
			Id:    movie.ID,
			Title: movie.Title,
			Year:  movie.Year,
		},
	}, nil
}

func (s *GrpcServer) CreateMovie(ctx context.Context, req *proto.CreateMovieRequest) (*proto.CreateMovieResponse, error) {
	movie, err := domain.NewMovie(req.Title, req.Year)
	if err != nil {
		return nil, err
	}

	created, err := s.service.CreateMovie(ctx, movie)
	if err != nil {
		return nil, err
	}
	return &proto.CreateMovieResponse{
		Movie: &proto.Movie{
			Title: created.Title,
			Year:  created.Year,
		},
	}, nil
}

func (s *GrpcServer) ListMovie(ctx context.Context, req *proto.ListMovieRequest) (*proto.ListMovieResponse, error) {
	filters := output.ListFilters{
		Title: req.Title,
		Year:  req.Year,
	}
	pagination := output.Pagination{
		Page:  req.Page,
		Limit: req.Limit,
	}
	sorting := output.Sorting{
		By: req.SortBy,
	}
	movies, _, err := s.service.ListMovies(ctx, filters, pagination, sorting)
	if err != nil {
		return nil, err
	}

	var protoMovies []*proto.Movie
	for _, m := range movies {
		protoMovies = append(protoMovies, &proto.Movie{
			Id:    m.ID,
			Title: m.Title,
			Year:  m.Year,
		})
	}
	return &proto.ListMovieResponse{Movie: protoMovies}, nil
}

func (s *GrpcServer) DeleteMovie(ctx context.Context, req *proto.DeleteMovieRequest) (*proto.DeleteMovieResponse, error) {
	err := s.service.DeleteMovie(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &proto.DeleteMovieResponse{Success: true}, nil
}
