package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/you/monk-coupons/pkg/endpoints"
	"github.com/you/monk-coupons/pkg/repo"
	"github.com/you/monk-coupons/pkg/service"
	transport "github.com/you/monk-coupons/pkg/transport/http"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	_ = godotenv.Load()

	mongoURI := os.Getenv("MONGO_URI")
	db := os.Getenv("DB_NAME")
	col := os.Getenv("COLLECTION_NAME")
	port := os.Getenv("PORT")

	client, _ := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))

	repository := repo.NewMongoRepo(client, db, col)
	svc := service.NewService(repository)
	eps := endpoints.Make(svc)
	handler := transport.MakeHTTPHandler(eps)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("Server running on port", port)
	log.Fatal(server.ListenAndServe())
}
