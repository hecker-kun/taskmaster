package main

import (
	"context"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"log"
	"os"
	pb "taskmaster/client/proto"
)

const address = "localhost:9055"

// TODO: Добавить поле ID, через который можно будет получать задачу методом GetTask()

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewTaskmasterClient(conn)

	app := &cli.App{
		Name: "taskmaster",
		Usage: "Simple task management app",
		Commands: []*cli.Command{
			{
				Name: "create",
				Aliases: []string{"c"},
				Usage: "creates a new task",
				Action: func(ctx *cli.Context) error {
					description := ctx.Args().First()
					status := ctx.Args().Get(1)
					taskid := ctx.Args().Get(2)

					r, err := c.CreateTask(context.Background(), &pb.CreateTaskReq{Task: &pb.Task{
						Description: description,
						Status:      status,
						Taskid:      taskid,
					}})
					if err != nil {
						log.Fatalf("failed to create task: %v", err)
					}

					log.Println(r.String())
					return nil
				},
			},
			{
				Name: "get",
				Aliases: []string{"g"},
				Usage: "returns the task with the specified ID",
				Action: func(ctx *cli.Context) error {
					taskid := ctx.Args().First()
					r, err := c.GetTask(context.Background(), &pb.GetTaskReq{Id: taskid})
					if err != nil {
						log.Fatalf("failed to get the task: %v", err)
					}

					log.Println(r.String())
					return nil
				},
			},
			{
				Name: "deleteall",
				Aliases: []string{"dall"},
				Usage: "removes all tasks",
				Action: func(ctx *cli.Context) error {
					_, err := c.DeleteAllTasks(context.Background(), &pb.Empty{})
					if err != nil {
						log.Fatalf("failed to delete all tasks: %v", err)
					}

					log.Println("all tasks deleted successfully")
					return nil
				},
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

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
	//req := &pb.GetAllTasksReq{}
	//
	//stream, err := c.GetAllTasks(ctx, req)
	//if err != nil {
	//	return
	//}
	//
	//for {
	//	res, err := stream.Recv()
	//	if err == io.EOF {
	//		log.Println("no more tasks")
	//		break
	//	}
	//	if err != nil {
	//		return
	//	}
	//
	//	log.Println(res.GetTask())
	//}
}