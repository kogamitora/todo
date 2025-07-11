// cmd/client/cmd/update.go
package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"connectrpc.com/connect"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/timestamppb"

	todov1 "github.com/kogamitora/todo/gen/proto/todo/v1"
	todov1connect "github.com/kogamitora/todo/gen/proto/todo/v1/v1connect"
)

var (
	updateTitle       string
	updateDescription string
	updateDueDate     string
	updateStatus      string
)

var updateCmd = &cobra.Command{
	Use:   "update [ID]",
	Short: "Update a TODO item",
	Long:  "Update a TODO item by its ID. You can update its title, description, due date, or status.",
	Args:  cobra.ExactArgs(1), // 规定必须且只能有一个参数，即 TODO 的 ID
	Run: func(cmd *cobra.Command, args []string) {
		// 1. 解析并验证 ID
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			log.Fatalf("Invalid ID provided: %v", err)
		}

		// 2. 创建 gRPC 客户端
		client := todov1connect.NewTodoServiceClient(
			http.DefaultClient,
			"http://localhost:8080",
		)

		// 3. 构建请求体
		req := &todov1.UpdateTodoRequest{
			Id: id,
		}

		// 4. 检查每个 flag 是否被用户设置，并相应地填充请求
		// cobra 提供了 IsChanged 方法来判断一个 flag 是否在命令行中被用户显式设置
		if cmd.Flags().Changed("title") {
			req.Title = &updateTitle
		}
		if cmd.Flags().Changed("description") {
			req.Description = &updateDescription
		}
		if cmd.Flags().Changed("due-date") {
			t, err := time.Parse("2006-01-02", updateDueDate)
			if err != nil {
				log.Fatalf("Invalid due date format. Use YYYY-MM-DD: %v", err)
			}
			req.DueDate = timestamppb.New(t)
		}
		if cmd.Flags().Changed("status") {
			var status todov1.Status
			switch updateStatus {
			case "completed":
				status = todov1.Status_STATUS_COMPLETED
			case "incomplete":
				status = todov1.Status_STATUS_INCOMPLETE
			default:
				log.Fatalf("Invalid status. Use 'completed' or 'incomplete'.")
			}
			req.Status = &status
		}

		// 5. 发送请求
		res, err := client.UpdateTodo(context.Background(), connect.NewRequest(req))
		if err != nil {
			log.Fatalf("Failed to update todo: %v", err)
		}

		fmt.Printf("Successfully updated TODO item with ID: %d\n", res.Msg.Todo.Id)
		fmt.Printf("Title: %s\nStatus: %s\n", res.Msg.Todo.Title, res.Msg.Todo.Status)
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
	// 定义所有可更新字段的 flags
	updateCmd.Flags().StringVarP(&updateTitle, "title", "t", "", "New title for the TODO")
	updateCmd.Flags().StringVarP(&updateDescription, "description", "d", "", "New description for the TODO")
	updateCmd.Flags().StringVar(&updateDueDate, "due-date", "", "New due date in YYYY-MM-DD format")
	updateCmd.Flags().StringVar(&updateStatus, "status", "", "New status (completed|incomplete)")
}
