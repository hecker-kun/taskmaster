package main

import (
	"context"
	"google.golang.org/grpc"
	"io"
	"log"
	pb "taskmaster/client/proto"
	"time"
)

const address = "localhost:9055"

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewTaskmasterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//task := pb.Task{
	//	Description: "First",
	//	Status:      "Canceled",
	//	Id:          "001",
	//}
	//r, err := c.CreateTask(ctx, &pb.CreateTaskReq{Task: &task})
	//if err != nil {
	//	log.Fatalf("failed to create the task: %v", err)
	//}
	//log.Printf(r.String())

	// Test call GetTask()
	//t, err := c.GetTask(ctx, &pb.GetTaskReq{Id: r.Task.GetId()})
	//if err != nil {
	//	log.Fatalf("failed to get the task: %v", err)
	//}
	//log.Printf(t.String())

	// Test call DeleteTask()
	//_, err = c.DeleteTask(ctx, &pb.DeleteTaskReq{Id: r.Task.GetId()})
	//if err != nil {
	//	log.Fatalf("failed to delete the task: %v", err)
	//}
	//log.Printf("the task was deleted successfully")

	// Test call DeleteAllTasks()
	//_, err = c.DeleteAllTasks(ctx, &pb.Empty{})
	//if err != nil {
	//	log.Fatalf("failed to delete all tasks")
	//}
	//log.Printf("all tasks have been deleted")

	// Test call GetAllTasks()
	req := &pb.GetAllTasksReq{}

	stream, err := c.GetAllTasks(ctx, req)
	if err != nil {
		return
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			log.Println("no more tasks")
			break
		}
		if err != nil {
			return
		}

		log.Println(res.GetTask())
	}
}