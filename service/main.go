package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"log"
	"net"
	pb "taskmaster/service/proto"
)

var (
	collection *mongo.Collection
	ctx = context.TODO()
)

func main() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("could not connet to MongoDB server: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf(err)
	}

	collection = client.Database("taskmaster").Collection("tasks")

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	pb.RegisterTaskmasterServer(srv, &server{})

	log.Printf("Sarting gRPC listener on port %s", port)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
