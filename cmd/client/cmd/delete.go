// cmd/client/cmd/delete.go
package cmd

import (
	"bufio" // 导入 bufio 包
	"context"
	"fmt"
	"log"
	"net/http"
	"os" // 导入 os 包
	"strconv"
	"strings" // 导入 strings 包

	"connectrpc.com/connect"
	"github.com/spf13/cobra"

	todov1 "github.com/kogamitora/todo/gen/proto/todo/v1"
	todov1connect "github.com/kogamitora/todo/gen/proto/todo/v1/v1connect"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [ID]",
	Short: "Delete a TODO item",
	Long:  "Logically delete a TODO item by its ID. The item will not be listed anymore but remains in the database.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			log.Fatalf("Invalid ID provided: %v", err)
		}

		// --- 新增确认逻辑 ---
		fmt.Printf("Are you sure you want to delete TODO item with ID %d? (y/N): ", id)

		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.ToLower(strings.TrimSpace(input))

		if input != "y" && input != "yes" {
			fmt.Println("Deletion cancelled.")
			return // 用户取消，直接退出
		}
		// --- 确认逻辑结束 ---

		client := todov1connect.NewTodoServiceClient(
			http.DefaultClient,
			"http://localhost:8080",
		)

		req := &todov1.DeleteTodoRequest{
			Id: id,
		}

		_, err = client.DeleteTodo(context.Background(), connect.NewRequest(req))
		if err != nil {
			log.Fatalf("Failed to delete todo: %v", err)
		}

		fmt.Printf("Successfully deleted TODO item with ID: %d\n", id)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
