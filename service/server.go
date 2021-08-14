package main

import (
	"context"
	"github.com/gofrs/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"sync"

	pb "taskmaster/service/proto"
)

const port = ":9055"

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
	out, err := uuid.NewV4()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error while generating Task ID", err)
	}
	in.Id = out.String()

	s.Lock()
	if s.taskStore == nil {
		s.taskStore = make(map[string]*pb.Task)
	}

	s.taskStore[in.Id] = in
	s.Unlock()

	return &pb.TaskID{Value: in.Id}, status.New(codes.OK, "").Err()
}

func main() {
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
