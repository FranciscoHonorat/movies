package main

import (
	"log"
	"os"

	"github.com/FranciscoHonorat/movies/api-gateway/internal/handlers"
	"github.com/FranciscoHonorat/movies/proto"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient(os.Getenv("GRPC_SERVER_URL"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	client := proto.NewMovieServiceClient(conn)

	movieHandler := handlers.NewMovieHandler(client)

	r := gin.Default()

	v1 := r.Group("/api/v1")

	v1.GET("/movies/:id", movieHandler.GetMovie)
	v1.GET("/movies", movieHandler.ListMovie)
	v1.POST("/movies", movieHandler.CreateMovie)
	v1.DELETE("/movies/:id", movieHandler.DeleteMovie)

	r.Run(":8080")
}
