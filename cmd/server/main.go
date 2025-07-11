/*
cmd/server/main.go: 服务器的启动器。
作用:
初始化日志 (slog)。
连接数据库。
创建 TodoHandler 的实例（把数据库连接传进去）。
使用 connect-go 将 handler 注册到 HTTP 路由器。
启动 HTTP 服务器，监听端口，开始接收请求。
*/
package main

import (
	"log/slog"
	"net/http"
	"os"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	// 替换为你的模块路径
	todov1connect "github.com/kogamitora/todo/gen/proto/todo/v1/v1connect"

	"github.com/kogamitora/todo/internal/db"
	"github.com/kogamitora/todo/internal/handler"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// 从环境变量获取 DSN
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		// 提供一个默认值方便开发
		dsn = "user:password@tcp(127.0.0.1:3307)/todo_db?parseTime=true"
	}

	database, err := db.NewDB(dsn, logger)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer database.Close()

	todoHandler := handler.NewTodoHandler(database, logger)
	path, h := todov1connect.NewTodoServiceHandler(todoHandler)

	mux := http.NewServeMux()
	mux.Handle(path, h)

	addr := ":8080"
	logger.Info("server starting", "addr", addr)

	// 使用 h2c 来支持 HTTP/2 Cleartext (非 TLS)
	err = http.ListenAndServe(
		addr,
		h2c.NewHandler(mux, &http2.Server{}),
	)
	if err != nil {
		logger.Error("server failed", "error", err)
		os.Exit(1)
	}
}
