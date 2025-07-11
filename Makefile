# Makefile
.PHONY: all generate proto db-up db-down migrate-up migrate-down run-server build-client

# 替换为你的模块路径
MODULE_NAME=github.com/yourusername/todo01

# 数据库 DSN
DB_URL=mysql://user:password@tcp(127.0.0.1:3306)/todo_db?multiStatements=true

all: generate

# 生成所有代码
generate: proto sqlboiler

# 生成 proto 代码
proto:
	@echo ">> generating protobuf code..."
	buf mod update
	buf generate

# 生成 sqlboiler ORM 代码
sqlboiler:
	@echo ">> generating sqlboiler code..."
	sqlboiler mysql

# 启动数据库容器
db-up:
	@echo ">> starting database container..."
	docker-compose up -d

# 停止并移除数据库容器
db-down:
	@echo ">> stopping database container..."
	docker-compose down

# 数据库迁移
migrate-up:
	@echo ">> applying database migrations..."
	migrate -database "${DB_URL}" -path migrations up

# 数据库回滚
migrate-down:
	@echo ">> rolling back database migrations..."
	migrate -database "${DB_URL}" -path migrations down

# 运行 gRPC 服务器
run-server:
	@echo ">> running server..."
	go run ./cmd/server/main.go

# 构建 CLI 客户端
build-client:
	@echo ">> building client..."
	go build -o ./bin/todocli ./cmd/client