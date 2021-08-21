package main

import (
	"context"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"io"
	"os"
	"strconv"
	pb "taskmaster/client/proto"

	log "github.com/sirupsen/logrus"
)

const address = "localhost:9055"

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
			{
				Name: "delete all",
				Aliases: []string{"dall"},
				Usage: "removes all tasks from the collection",
				Action: func(ctx *cli.Context) error {
					_, err := c.DeleteAllTasks(context.Background(), &pb.Empty{})
					if err != nil {
						log.WithFields(log.Fields{
							"package": "client",
							"method": "DeleteAllTasks",
						}).Fatalf("failed to delete all tasks: %v", err)
					}

					return nil
				},
			},
			{
				Name: "get all",
				Aliases: []string{"gall", "getall", "list"},
				Usage: "lists all tasks",
				Action: func(ctx *cli.Context) error {
					stream, err := c.GetAllTasks(context.Background(), &pb.Empty{})
					if err != nil {
						log.WithFields(log.Fields{
							"package": "client",
							"method": "GetAllTasks",
						}).Fatalf("failed to get the task list: %v", err)
					}

					for {
						task, err := stream.Recv()
						if err == io.EOF {
							log.Println("no more tasks")
						}
						if err != nil {
							log.Fatalf("internal error: %v", err)
						}

						log.Println(task.Text, ": ", task.Status)
					}
				},
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}