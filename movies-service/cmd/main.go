package main

import (
	"context"
	"log"
	"net"
	"os"

	dynamodbadapter "github.com/FranciscoHonorat/movies/movies-service/internal/adapters/dynamodb"
	grpcserver "github.com/FranciscoHonorat/movies/movies-service/internal/adapters/grpc_server"
	"github.com/FranciscoHonorat/movies/movies-service/internal/adapters/rabbitmq"
	"github.com/FranciscoHonorat/movies/movies-service/internal/adapters/seed"

	"github.com/FranciscoHonorat/movies/movies-service/internal/core/domain"
	"github.com/FranciscoHonorat/movies/movies-service/internal/core/service"
	"github.com/FranciscoHonorat/movies/proto"

	//"go.mongodb.org/mongo-driver/v2/mongo"
	//"go.mongodb.org/mongo-driver/v2/mongo/options"
	"google.golang.org/grpc"
)

func main() {
	//mongoURI := os.Getenv("MONGO_URI")

	ctx := context.Background()

	//client, err := mongo.Connect(options.Client().ApplyURI(mongoURI))
	//if err != nil {
	//	log.Fatal(err)
	//}

	//err = client.Ping(ctx, nil)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//collection := client.Database("moviesDB").Collection("movies")

	//if err := seed.Seed(ctx, collection, "movies-service/movies.json"); err != nil {
	//	log.Fatal(err)
	//}

	//repo := mongodb.NewMongoRepository(collection)

	dynamoRepo, err := dynamodbadapter.NewDynamoRepository(ctx, os.Getenv("DYNAMODB_ENDPOINT"))
	if err != nil {
		log.Fatal(err)
	}

	if err := seed.SeedMovie(ctx, dynamoRepo, "movies-service/movies.json"); err != nil {
		log.Fatal(err)
	}
	svc := service.NewMovieService(dynamoRepo)

	rabbitmqConsumer, err := rabbitmq.NewRabbitMQConsumer(os.Getenv("RABBITMQ_URI"), "movies_queue")
	if err != nil {
		log.Fatal(err)
	}
	go rabbitmqConsumer.Consume(ctx, func(movie domain.Movie) error {
		_, err := svc.CreateMovie(ctx, &movie)
		return err
	})

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcSrv := grpc.NewServer()
	proto.RegisterMovieServiceServer(grpcSrv, grpcserver.NewGrpcServer(svc))
	grpcSrv.Serve(lis)
}
