package main

import (
	"context"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"os"
	"strconv"
	pb "taskmaster/client/proto"

	log "github.com/sirupsen/logrus"
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
					text := ctx.Args().First()

					res, err := c.CreateTask(context.Background(), &pb.AddTask{
						Text:   text,
						Status: true,
					})
					if err != nil {
						log.WithFields(log.Fields{
							"package": "client",
							"method": "CreateTask",
						}).Fatalf("failed to create task: %v", err)
					}

					log.Println(res.String())

					return nil
				},
			},
			{
				Name: "delete",
				Aliases: []string{"d", "del"},
				Usage: "deletes the task with the specified ID",
				Action: func(ctx *cli.Context) error {
					id := ctx.Args().First()
					nid, _ := strconv.Atoi(id)

					_, err := c.DeleteTask(context.Background(), &pb.DeleteParams{Id: int32(nid)})
					if err != nil {
						log.WithFields(log.Fields{
							"package": "client",
							"method": "DeleteTask",
						}).Fatalf("failed to delete the task: %v", err)
					}

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