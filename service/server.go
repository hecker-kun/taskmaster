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

	"taskmaster/service/proto"
)

const port = ":9055"

var (
	mongoCtx context.Context
	taskdb   *mongo.Collection
	db       *mongo.Client

	taskUniqID int32 = 1
)

type server struct {
	*proto.UnimplementedTaskmasterServer
}

func (s *server) CreateTask(ctx context.Context, taskParams *proto.AddTask) (*proto.Task, error) {
	task := &TaskObject{
		Text:   taskParams.Text,
		Status: taskParams.Status,
		TaskID: taskUniqID,
	}

	res, err := taskdb.InsertOne(mongoCtx, task)
	if err != nil {
		status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v", err))
	}
	task.ID = res.InsertedID.(primitive.ObjectID)

	// We increase the ID counter by one, so the tasks will have IDs from 1 to N
	defer incrementID(&taskUniqID)

	return &proto.Task{
		Text:   task.Text,
		Status: task.Status,
		Id:     taskUniqID,
	}, nil
}

func (s *server) DeleteTask(ctx context.Context, params *proto.DeleteParams) (*proto.Empty, error) {
	_, err := taskdb.DeleteOne(ctx, bson.M{"task_id": params.Id})
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("internal error: %v"), err)
	}

	return &proto.Empty{}, nil
}

func (s *server) DeleteAllTasks(ctx context.Context, empty *proto.Empty) (*proto.Empty, error) {
	_, err := taskdb.DeleteMany(ctx, bson.M{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("internal error: %v", err))
	}

	return &proto.Empty{}, nil
}

func (s *server) GetAllTasks(empty *proto.Empty, stream proto.Taskmaster_GetAllTasksServer) error {
	data := &TaskObject{}

	cursor, err := taskdb.Find(context.Background(), bson.M{})
	if err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintf("internal error: %v", err))
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		err := cursor.Decode(data)
		if err != nil {
			return status.Errorf(codes.Unavailable, fmt.Sprintf("could not decode data: %v", err))
		}

		stream.Send(&proto.Task{
			Text:   data.Text,
			Status: data.Status,
			Id:     data.TaskID,
		})
	}
	
	if err := cursor.Err(); err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintf("internal error: %v", err))
	}

	return nil
}

type TaskObject struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Text   string             `bson:"description,omitempty"`
	Status bool               `bson:"status,omitempty"`
	TaskID int32              `bson:"task_id,omitempty"`
}

func incrementID(id *int32) {
	*id++
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	proto.RegisterTaskmasterServer(srv, &server{})

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
	taskdb = db.Database("taskmaster-beta").Collection("tasks")

	log.Printf("Starting gRPC listener on port %s", port)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
