package main

import (
	"context"
	"google.golang.org/grpc"
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

	task := pb.Task{
		Description: "Fourth",
		Status:      "Canceled",
		Id:          "105",
	}
	r, err := c.CreateTask(ctx, &pb.CreateTaskReq{Task: &task})
	if err != nil {
		log.Fatalf("failed to create the task: %v", err)
	}
	log.Printf(r.String())

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
}