package main

import (
	"context"
	"fmt"
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
	taskStore map[string]*pb.Task
	*pb.UnimplementedTaskmasterServer

	sync.Mutex
}

func (s *server) DeleteTask(ctx context.Context, id *pb.TaskID) (*pb.Empty, error) {
	panic("implement me")
}

func (s *server) DeleteAllTasks(ctx context.Context, empty *pb.Empty) (*pb.Empty, error) {
	panic("implement me")
}

func (s *server) GetAllTasks(ctx context.Context, empty *pb.Empty) (*pb.TaskList, error) {
	panic("implement me")
}

func (s *server) GetTask(ctx context.Context, id *pb.TaskID) (*pb.Task, error) {
	s.Lock()
	task, ok := s.taskStore[id.Value]
	s.Unlock()
	if !ok {
		return nil, status.Errorf(codes.NotFound, "Task does not exist", id.Value)
	}

	return task, status.New(codes.OK, "").Err()
}

func (s *server) CreateTask(ctx context.Context, in *pb.Task) (*pb.TaskID, error) {
	data := TaskObject{
		//ID:          primitive.ObjectID{},
		Description: in.Description,
		Status:      in.Status,
	}

	res, err := taskDb.InsertOne(mongoCtx, data)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v", err))
	}
	taskId := res.InsertedID.(primitive.ObjectID)

	in.Id = taskId.Hex()

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
	db, err := mongo.Connect(mongoCtx, options.Client().ApplyURI("mongodb://localhost:27017"))
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
