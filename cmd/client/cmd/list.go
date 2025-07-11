// cmd/client/cmd/list.go
package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"connectrpc.com/connect"
	"github.com/spf13/cobra"

	todov1 "github.com/kogamitora/todo/gen/proto/todo/v1"
	todov1connect "github.com/kogamitora/todo/gen/proto/todo/v1/v1connect"
)

var (
	statusFilter  string
	sortByDueDate string // 新增：用于接收排序参数
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all TODO items",
	Run: func(cmd *cobra.Command, args []string) {
		client := todov1connect.NewTodoServiceClient(
			http.DefaultClient,
			"http://localhost:8080",
		)
		req := &todov1.ListTodosRequest{}

		// 处理状态过滤 (逻辑不变)
		if statusFilter != "" {
			var status todov1.Status
			switch statusFilter {
			case "completed":
				status = todov1.Status_STATUS_COMPLETED
			case "incomplete":
				status = todov1.Status_STATUS_INCOMPLETE
			default:
				log.Fatalf("Invalid status filter. Use 'completed' or 'incomplete'")
			}
			req.StatusFilter = &status
		}

		// --- 新增逻辑开始 ---

		// 处理按截止日期排序
		if sortByDueDate != "" {
			var sortOrder todov1.SortOrder
			switch strings.ToLower(sortByDueDate) {
			case "asc":
				sortOrder = todov1.SortOrder_SORT_ORDER_ASC
			case "desc":
				sortOrder = todov1.SortOrder_SORT_ORDER_DESC
			default:
				log.Fatalf("Invalid sort order. Use 'asc' or 'desc'.")
			}
			req.SortByDueDate = &sortOrder
		}

		// --- 新增逻辑结束 ---

		res, err := client.ListTodos(context.Background(), connect.NewRequest(req))
		if err != nil {
			log.Fatalf("Failed to list todos: %v", err)
		}

		fmt.Println("ID\tStatus\t\tDue Date\tTitle")
		fmt.Println("----------------------------------------------------------")
		for _, todo := range res.Msg.Todos {
			dueDateStr := "N/A"
			if todo.DueDate != nil && todo.DueDate.IsValid() {
				dueDateStr = todo.DueDate.AsTime().Format("2006-01-02")
			}
			// 修正一下状态的显示，去掉前缀
			statusStr := strings.Replace(todo.Status.String(), "STATUS_", "", 1)
			fmt.Printf("%d\t%-10s\t%s\t%s\n", todo.Id, statusStr, dueDateStr, todo.Title)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVar(&statusFilter, "status", "", "Filter by status (completed|incomplete)")
	// 新增一个 flag
	listCmd.Flags().StringVar(&sortByDueDate, "sort-by-due", "", "Sort by due date (asc|desc)")
}
