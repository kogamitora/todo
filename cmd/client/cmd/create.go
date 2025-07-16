package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/timestamppb"

	todov1 "github.com/kogamitora/todo/gen/proto/todo/v1"
	todov1connect "github.com/kogamitora/todo/gen/proto/todo/v1/v1connect"
)

var (
	title       string
	description string
	dueDate     string
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new TODO item",
	Run: func(cmd *cobra.Command, args []string) {
		client := todov1connect.NewTodoServiceClient(
			http.DefaultClient,
			ServerURL,
		)
		req := &todov1.CreateTodoRequest{
			Title:       title,
			Description: description,
		}
		if dueDate != "" {
			t, err := time.Parse("2006-01-02", dueDate) //参考タイムパッケージのフォーマット
			if err != nil {
				log.Fatalf("Invalid due date format. Use YYYY-MM-DD: %v", err)
			}
			req.DueDate = timestamppb.New(t)
		}
		res, err := client.CreateTodo(context.Background(), connect.NewRequest(req))
		if err != nil {
			log.Fatalf("Failed to create todo: %v", err)
		}
		fmt.Printf("Successfully created TODO item with ID: %d\n", res.Msg.Todo.Id)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVarP(&title, "title", "t", "", "Title of the TODO (required)")
	createCmd.Flags().StringVarP(&description, "description", "d", "", "Description of the TODO")
	createCmd.Flags().StringVar(&dueDate, "due-date", "", "Due date in YYYY-MM-DD format")
	createCmd.MarkFlagRequired("title")
}
