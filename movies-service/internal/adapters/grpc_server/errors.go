package grpcserver

import (
	"errors"

	"github.com/FranciscoHonorat/movies/movies-service/internal/core/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func toGRPCError(err error) error {
	if errors.Is(err, domain.ErrMovieNotFound) {
		return status.Error(codes.NotFound, err.Error())
	}

	if errors.Is(err, domain.ErrInvalidMovieData) {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	if errors.Is(err, domain.ErrInternalServer) {
		return status.Error(codes.Internal, err.Error())
	}

	return status.Error(codes.Internal, err.Error())
}
