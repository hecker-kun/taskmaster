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

	// Creating a test Task
	task := pb.Task{
		Id: "109",
		Description: "Implement MongoDB",
		Status:      "Done",
	}

	r, err := c.CreateTask(ctx, &task)
	if err != nil {
		log.Fatalf("could not create a task: %v", err)
	}
	log.Printf("TaskID: %s created succesfully", r.Value)

	t, err := c.GetTask(ctx, r)
	if err != nil {
		log.Fatalf("task non found: %v", err)
	}
	log.Println(t.String())
}