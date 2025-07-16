package main

import (
	"log/slog"
	"net/http"
	"os"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	todov1connect "github.com/kogamitora/todo/gen/proto/todo/v1/v1connect"
	"github.com/kogamitora/todo/internal/config"
	"github.com/kogamitora/todo/internal/db"
	"github.com/kogamitora/todo/internal/handler"
)

func main() {
	// 基本コンポーネントの初期化 (Logger)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// 設定の読み込みと検証
	if err := config.LoadFromFile(".env"); err != nil {
		logger.Warn("failed to load .env file", "error", err)
	}

	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	if err := cfg.Validate(); err != nil {
		logger.Error("invalid config", "error", err)
		os.Exit(1)
	}

	logger.Info("loaded config",
		"server_host", cfg.Server.Host,
		"server_port", cfg.Server.Port,
		"db_host", cfg.Database.Host,
		"db_port", cfg.Database.Port,
		"db_name", cfg.Database.Database,
		"db_user", cfg.Database.User,
	)

	// 依存サービスの初期化 (データベース)
	database, err := db.NewDB(cfg.GetDSN(), logger)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer database.Close()

	// HTTPハンドラとルーティングの設定 (Mux)
	todoHandler := handler.NewTodoHandler(database, logger)
	path, h := todov1connect.NewTodoServiceHandler(todoHandler)

	mux := http.NewServeMux()
	mux.Handle(path, h)

	// サーバーの起動
	addr := ":" + cfg.Server.Port
	server := &http.Server{
		Addr:    addr,
		Handler: h2c.NewHandler(mux, &http2.Server{}),
	}

	logger.Info("server starting", "addr", addr, "gRPC_path", path)
	if err := server.ListenAndServe(); err != nil {
		logger.Error("server failed", "error", err)
		os.Exit(1)
	}
}
