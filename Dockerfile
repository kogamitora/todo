# ビルドステージ
FROM golang:1.24.4-alpine AS builder

WORKDIR /app

# 依存関係のファイルをコピーしてダウンロード
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピーしてビルド
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server cmd/server/main.go

# 実行ステージ
FROM alpine:latest

# 実行時に必要な依存関係をインストール
RUN apk --no-cache add ca-certificates netcat-openbsd

WORKDIR /app

# ビルドされたバイナリをコピー
COPY --from=builder /app/server .

# マイグレーションファイルをコピー
COPY --from=builder /app/migrations ./migrations

# ヘルスチェックスクリプトをコピー
COPY --from=builder /app/scripts ./scripts
RUN chmod +x ./scripts/health_check.sh

# ポートを公開
EXPOSE 8080

# ヘルスチェック
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD ./scripts/health_check.sh

# アプリケーションを起動
CMD ["./server"]