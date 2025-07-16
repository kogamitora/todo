# Makefile
.PHONY: all generate proto sqlboiler \
		run-server build-client \
		docker-build docker-up docker-down docker-logs docker-rebuild docker-clean docker-migrate \
		docker-exec docker-exec-db test-client test

# .env ファイルから環境変数をデフォルトで読み込みます
# これにより、sqlboiler のようなローカルスクリプトもこれらの変数を使用できます
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# デフォルトターゲット
all: generate

# ====== コード生成 ======
generate: proto sqlboiler

# proto コード (connect-go) を生成
proto:
	@echo ">> generating protobuf code..."
	buf mod update
	buf generate

# ORM コード (sqlboiler) を生成
sqlboiler:
	@echo ">> generating sqlboiler ORM code..."
	export MYSQL_HOST="${DB_HOST}" && \
	export MYSQL_PORT="${DB_PORT}" && \
	export MYSQL_USER="${DB_USER}" && \
	export MYSQL_PASS="${DB_PASSWORD}"; \
	sqlboiler mysql

# ====== ローカル開発 ======
# ローカルでサーバーを実行
run-server:
	@echo ">> running server..."
	go run ./cmd/server/main.go

# CLI クライアントをビルド
build-client:
	@echo ">> building client..."
	go build -o ./bin/todocli ./cmd/client

# ====== Docker 管理 ======
# Docker イメージをビルド
docker-build:
	@echo ">> building docker images..."
	docker-compose build

# すべてのサービスを起動
docker-up:
	@echo ">> starting all services via docker..."
	docker-compose up -d

# すべてのサービスを停止して削除
docker-down:
	@echo ">> stopping all services..."
	docker-compose down

# すべてのサービスのログを表示
docker-logs:
	@echo ">> tailing all service logs..."
	docker-compose logs -f

# データベースマイグレーションを実行 (コンテナ内) - 現在はアプリケーションが自動でマイグレーションを実行するため、このコマンドは予備のオプションです
docker-migrate:
	@echo ">> running database migrations inside docker..."
	# docker-compose up migrate

# データベースマイグレーションをロールバック (コンテナ内)
docker-migrate-down:
	@echo ">> rolling back database migrations inside docker..."
	docker-compose run --rm migrate down 1 # 例: 1つのバージョンをロールバック

# すべてのサービスを再ビルドして再起動
docker-rebuild:
	@echo ">> rebuilding and restarting all services..."
	docker-compose down
	docker-compose build --no-cache
	docker-compose up -d

# Docker リソース (ボリュームを含む) をクリーンアップ
docker-clean:
	@echo ">> cleaning all docker resources (containers, volumes, networks)..."
	docker-compose down -v --remove-orphans

# ====== ユーティリティ ======
# app コンテナのシェルに入る
docker-exec:
	@echo ">> entering app container..."
	docker-compose exec app /bin/sh

# db コンテナの mysql-cli に入る
docker-exec-db:
	@echo ">> entering db container mysql cli..."
	docker-compose exec db mysql -u${DB_USER} -p${DB_PASSWORD} ${DB_NAME}

# クライアントのスモークテストを実行
test-client: build-client
	@echo ">> running client smoke tests..."
	@echo "\n--- Creating a todo ---"
	./bin/todocli create --title="学習 Go" --description="完成 connect-go 挑战"
	@echo "\n--- Listing todos ---"
	./bin/todocli get
	@echo "\n--- Test completed! ---"

