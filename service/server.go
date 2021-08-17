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
	taskdb   *mongo.Collection
	db       *mongo.Client
)

type server struct {
	*pb.UnimplementedTaskmasterServer

	sync.Mutex
}

func (s *server) CreateTask(ctx context.Context, req *pb.CreateTaskReq) (*pb.CreateTaskRes, error) {
	task := req.GetTask()
	data := TaskObject{
		Description: task.GetDescription(),
		Status:      task.GetStatus(),
	}

	result, err := taskdb.InsertOne(mongoCtx, data)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v"),
			err,
		)
	}

	tid := result.InsertedID.(primitive.ObjectID)
	task.Id = tid.Hex()

	return &pb.CreateTaskRes{Task: task}, nil
}

func (s *server) DeleteTask(ctx context.Context, req *pb.DeleteTaskReq) (*pb.DeleteTaskRes, error) {
	tid, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Could not convert to ObjectID"),
		)
	}

	_, err = taskdb.DeleteOne(ctx, bson.M{"_id": tid})
	if err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Task with id %s not found in database", tid.String()),
		)
	}

	return &pb.DeleteTaskRes{Success: true}, nil
}

func (s *server) GetTask(ctx context.Context, req *pb.GetTaskReq) (*pb.GetTaskRes, error) {
	tid, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Could not convert to ObjectID"),
			err,
		)
	}

	result := taskdb.FindOne(ctx, bson.M{"_id": tid})
	data := TaskObject{}

	if err := result.Decode(&data); err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Task not found: %v"), err)
	}

	resp := &pb.GetTaskRes{
		Task: &pb.Task{
			Description: data.Description,
			Status:      data.Status,
			Id:          tid.Hex(),
		},
	}

	return resp, nil
}

func (s *server) DeleteAllTasks(ctx context.Context, empty *pb.Empty) (*pb.DeleteAllTasksRes, error) {
	_, err := taskdb.DeleteMany(ctx, bson.M{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("can not delete all tasks of collection"))
	}

	return &pb.DeleteAllTasksRes{Success: true}, nil
}

func (s *server) GetAllTasks(req *pb.GetAllTasksReq, stream pb.Taskmaster_GetAllTasksServer) error {
	panic("implement me")
}

type TaskObject struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
	Description string `bson:"description,omitempty"`
	Status string `bson:"status,omitempty"`
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
	taskdb = db.Database("taskmaster").Collection("tasks")

	log.Printf("Starting gRPC listener on port %s", port)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
