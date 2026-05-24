package main

import (
	"context"
	"log"
	"net"
	"os"

	grpcserver "github.com/FranciscoHonorat/movies/movies-service/internal/adapters/grpc_server"
	"github.com/FranciscoHonorat/movies/movies-service/internal/adapters/mongodb"
	"github.com/FranciscoHonorat/movies/movies-service/internal/adapters/seed"
	"github.com/FranciscoHonorat/movies/movies-service/internal/core/service"
	"github.com/FranciscoHonorat/movies/proto"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"google.golang.org/grpc"
)

func main() {
	mongoURI := os.Getenv("MONGO_URI")

	ctx := context.Background()

	client, err := mongo.Connect(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database("moviesDB").Collection("movies")

	if err := seed.Seed(ctx, collection, "movies-service/movies.json"); err != nil {
		log.Fatal(err)
	}

	repo := mongodb.NewMongoRepository(collection)
	svc := service.NewMovieService(repo)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcSrv := grpc.NewServer()
	proto.RegisterMovieServiceServer(grpcSrv, grpcserver.NewGrpcServer(svc))
	grpcSrv.Serve(lis)
}
