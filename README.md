# Go TODO アプリ (コーディングチャレンジ)

これは、「コーディングチャレンジ：Go 言語 TODO アプリ」への回答として、Go 言語で実装された TODO アプリケーションです。このプロジェクトは、`connect-go` を使用して gRPC サービスを構築し、コマンドラインインターフェース（CLI）を介して対話することを核としています。

このプロジェクトは、API 設計、データベース操作、レイヤー化されたアーキテクチャ、および保守性を含む、現代の Go マイクロサービス開発におけるベストプラクティスを示すことを目的としています。

## ✨ 機能

- **TODO の作成**: タイトル、説明、期日を含む新しいタスクを追加します。
- **TODO の更新**: 既存のタスクのステータスやタイトルなどの情報を変更します。
- **TODO の削除**: 指定されたタスクを論理削除（ソフトデリート）します。
- **TODO の取得**: 未削除のすべてのタスクを表示し、ステータスでのフィルタリングと期日でのソートをサポートします。
- **gRPC 通信**: `connect-go` を使用して、効率的で型安全なクライアント-サーバー間通信を実装します。
- **CLI クライアント**: 使いやすいコマンドラインツールを提供し、タスクを管理します。

## 🏛️ プロジェクト構成

プロジェクトは、保守と拡張を容易にするために、明確なレイヤー化アーキテクチャに従っています。

```
todo/
├── cmd/
│   ├── client/           # CLIクライアント
│   └── server/           # gRPCサーバー
├── internal/
│   ├── config/          # 設定管理
│   ├── db/              # データベース接続
│   └── handler/         # ビジネスロジックハンドラ
├── proto/
│   └── todo/v1/         # Protobuf定義
├── gen/
│   └── proto/           # 生成されたProtobufコード
├── models/              # 生成されたORMモデル
├── migrations/          # データベースマイグレーションファイル
├── bin/                 # コンパイル後のバイナリファイル
├── docker-compose.yml   # Docker Compose設定
├── Dockerfile          # アプリケーションのイメージビルドファイル
├── Makefile            # ビルドスクリプト
└── .env                # 環境変数設定
```

## 🛠️ 技術スタック

- **言語**: Go
- **通信プロトコル**: gRPC ([connect-go](https://github.com/connectrpc/connect-go)経由)
- **API 定義**: Protocol Buffers ([buf](https://buf.build/)経由)
- **データベース**: MySQL 8.0
- **ORM**: [sqlboiler](https://github.com/volatiletech/sqlboiler) (コード生成モデル)
- **ロギング**: Go 標準ライブラリ `slog`
- **CLI**: [Cobra](https://github.com/spf13/cobra)
- **環境**: Docker & Docker Compose

## 🚀 クイックスタート

### 前提条件

開始する前に、お使いのシステム（macOS または Linux を推奨）に以下のツールがインストールされていることを確認してください：

- **Go** (最新の安定版を推奨)
- **Docker** & **Docker Compose**
- **Make**
- **Buf CLI**: `brew install bufbuild/buf/buf`
- **SQLBoiler**: `brew install volatiletech/sqlboiler`
- **golang-migrate**: `brew install golang-migrate`

### 実行手順

**リポジトリのクローン**

```bash
git clone https://github.com/kogamitora/todo
cd todo
```

**Go 依存関係のダウンロード**

```bash
go mod tidy
```

**環境設定ファイルの作成**

```bash
cp .env.example .env
```

**サービスの起動**

```bash
make docker-up
```

**クライアントのテスト**

```bash
make test-client
```

### その他のコマンド

```bash
# ProtobufとORMコードの生成
make generate
```

すべてのコマンドについては、プロジェクトのルートディレクトリにある [Makefile](Makefile) を参照してください。

### [クライアントの使用方法](CLIENT_README.md)

## 📐 設計

### 1\. API 設計 (`proto`)

API 定義はプロジェクトの中核です。Protocol Buffers (proto3) を使用してサービス契約を定義し、クライアントとサーバー間の厳密な型付けと後方互換性を保証しています。

- **バージョニング**: API 定義は `proto/todo/v1` ディレクトリに配置され、将来の API 進化（例：v2）のためのスペースを確保しています。
- **明確なメッセージ**: `CreateTodoRequest` や `CreateTodoResponse` のように、リクエストとレスポンスのメッセージが明確に定義されています。
- **オプショナルなフィールド**: `UpdateTodoRequest` では、すべてのフィールドが `optional` としてマークされており、クライアントが一部のフィールドのみを更新できる、一般的な PATCH パターンを実装しています。
- **Enum**: `Status` を定義するために `enum` を使用し、「マジックストリング」の使用を避けています。

### 2\. サーバー (`internal/handler`)

サーバーサイドは、明確な責務分離の原則に従っています。

- **ハンドラ層**: `internal/handler/todo_handler.go` は `TodoService` インターフェースを実装しています。gRPC リクエストの処理、入力の検証、データベースロジックの呼び出し、およびデータ形式の変換（データベースモデルから Protobuf メッセージへ）を担当します。
- **データベース層**: 複雑なリポジトリパターンは実装せず、ハンドラ内で直接 `SQLBoiler` によって生成されたモデルを使用しています。この規模のプロジェクトでは、このアプローチの方が直接的で、不要な抽象化を減らせます。`SQLBoiler` は型安全なデータベースクエリを提供し、手書き SQL に起因するタイプミスや SQL インジェクションのリスクを回避します。
- **エラーハンドリング**: `connect.NewError` を使用して、標準の gRPC エラーコード（例：`CodeNotFound`, `CodeInternal`）を返し、クライアントがエラーを適切に処理できるようにしています。
- **ロギング**: 構造化ロギングライブラリの `slog` を使用し、主要な操作やエラー情報を記録することで、デバッグと監視を容易にしています。

### 3\. データベース (`migrations` & `sqlboiler`)

- **マイグレーション管理**: `golang-migrate` を使用してデータベーススキーマの変更を管理します。これにより、チームでの共同作業やデプロイの自動化がより信頼性の高いものになります。
- **論理削除**: `todos` テーブルには `deleted_at` フィールドが含まれており、削除操作は物理的にデータを削除するのではなく、このフィールドのタイムスタンプを更新します。これはデータを保護し、復旧を容易にする一般的な手法です。
- **ORM の選定**: `SQLBoiler` は「コード生成」型の ORM です。GORM のように大量のリフレクションを使用しないため、パフォーマンスが良く、生成されるコードは型安全であるため、コンパイル時により多くのエラーを検出できます。

### 4\. クライアント (`cmd/client`)

- **CLI フレームワーク**: `Cobra` を使用して、機能が豊富で使いやすいコマンドラインインターフェースを構築しています。サブコマンド、フラグ、および自動生成されるヘルプメッセージをサポートしています。
- **疎結合**: クライアントロジックはサーバーから完全に独立しており、生成された gRPC クライアントを介してのみサーバーと通信します。
