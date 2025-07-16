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

// filter and sort flags
var (
	statusFilter  string
	sortByDueDate string
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get all TODO items",
	Run: func(cmd *cobra.Command, args []string) {
		client := todov1connect.NewTodoServiceClient(
			http.DefaultClient,
			ServerURL,
		)
		req := &todov1.GetTodosRequest{}

		// filter by status
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

		// sort by due date
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

		res, err := client.GetTodos(context.Background(), connect.NewRequest(req))
		if err != nil {
			log.Fatalf("Failed to get todos: %v", err)
		}

		fmt.Println("ID\tStatus\t\tDue Date\tTitle")
		fmt.Println("----------------------------------------------------------")
		for _, todo := range res.Msg.Todos {
			dueDateStr := "N/A"
			if todo.DueDate != nil && todo.DueDate.IsValid() {
				dueDateStr = todo.DueDate.AsTime().Format("2006-01-02")
			}
			// Convert status to string without the prefix
			// e.g., STATUS_COMPLETED -> COMPLETED
			statusStr := strings.Replace(todo.Status.String(), "STATUS_", "", 1)
			fmt.Printf("%d\t%-10s\t%s\t%s\n", todo.Id, statusStr, dueDateStr, todo.Title)
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().StringVar(&statusFilter, "status", "", "Filter by status (completed|incomplete)")
	getCmd.Flags().StringVar(&sortByDueDate, "sort-by-due", "", "Sort by due date (asc|desc)")
}
