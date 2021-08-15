package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"sync"

	pb "taskmaster/service/proto"
)

const port = ":9055"

var (
	mongoCtx context.Context
	taskDb *mongo.Collection
	db *mongo.Client
)

type server struct {
	*pb.UnimplementedTaskmasterServer

	sync.Mutex
}

func (s *server) DeleteTask(ctx context.Context, id *pb.TaskID) (*pb.Empty, error) {
	tid, err := primitive.ObjectIDFromHex(id.GetValue())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Could not convert to ObjectID: %v"), err)
	}

	_, err = taskDb.DeleteOne(ctx, bson.M{"_id": tid})
	if err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Could not find/delete task with id %s: %v"),
			id.GetValue(),
			err,
		)
	}

	return &pb.Empty{}, nil
}

func (s *server) DeleteAllTasks(ctx context.Context, empty *pb.Empty) (*pb.Empty, error) {
	panic("implement me")
}

func (s *server) GetAllTasks(ctx context.Context, empty *pb.Empty) (*pb.TaskList, error) {
	panic("implement me")
}

func (s *server) GetTask(ctx context.Context, id *pb.TaskID) (*pb.Task, error) {
	tid, err := primitive.ObjectIDFromHex(id.GetValue())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Could not convert to ObjectID: %v", err))
	}
	s.Lock()
	result := taskDb.FindOne(ctx, bson.M{"_id": tid})
	s.Unlock()

	// Create an empty TaskObject to write our decode result to
	data := TaskObject{}
	// Decode and write to data
	if err := result.Decode(&data); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Could not find task with ObjectID %s: %v"),
			id.GetValue(),
			err,
		)
	}

	response := &pb.Task{
		Description: data.Description,
		Status:      data.Status,
		Id:          tid.Hex(),
	}

	return response, status.New(codes.OK, "").Err()
}

func (s *server) CreateTask(ctx context.Context, in *pb.Task) (*pb.TaskID, error) {
	data := TaskObject{
		//ID:          primitive.ObjectID{},
		Description: in.Description,
		Status:      in.Status,
	}
	s.Lock()
	res, err := taskDb.InsertOne(mongoCtx, data)
	s.Unlock()
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v", err))
	}
	tid := res.InsertedID.(primitive.ObjectID)

	in.Id = tid.Hex()

	return &pb.TaskID{Value: in.Id}, status.New(codes.OK, "").Err()
}

type TaskObject struct {
	ID primitive.ObjectID `bson:"id,omitempty"`
	Description string `bson:"description"`
	Status string `bson:"status"`
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	pb.RegisterTaskmasterServer(srv, &server{})

	// Connecting to MongoDB server
	mongoCtx = context.Background()
	db, err = mongo.Connect(mongoCtx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = db.Ping(mongoCtx, nil)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	} else {
		log.Printf("Connected to MongoDB")
	}
	taskDb = db.Database("taskmaster").Collection("tasks")

	log.Printf("Starting gRPC listener on port %s", port)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
