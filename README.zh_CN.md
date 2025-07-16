# Go TODO 应用 (编码挑战)

这是一个使用 Go 语言实现的 TODO 应用程序，作为对 "编码挑战：Go 语言 TODO 应用" 的响应。项目核心是使用 `connect-go` 构建一个 gRPC 服务，并通过命令行界面 (CLI) 与之交互。

此项目旨在展示现代 Go 微服务开发的最佳实践，包括 API 设计、数据库操作、分层架构和可维护性。

## ✨ 功能特性

- **创建 TODO**: 添加新的待办事项，包括标题、描述和截止日期。
- **更新 TODO**: 修改已存在的待办事项的状态、标题等信息。
- **删除 TODO**: 逻辑删除（软删除）指定的待办事项。
- **获取 TODO**: 展示所有未删除的待办事项，并支持按状态过滤和按截止日期排序。
- **gRPC 通信**: 使用 `connect-go` 实现高效、类型安全的客户端-服务器通信。
- **CLI 客户端**: 提供一个易于使用的命令行工具来管理待办事项。

## 🏛️ 项目结构

项目遵循了清晰的分层架构，便于维护和扩展。

```
todo/
├── cmd/
│   ├── client/           # CLI 客户端
│   └── server/           # gRPC 服务器
├── internal/
│   ├── config/          # 配置管理
│   ├── db/              # 数据库连接
│   └── handler/         # 业务逻辑处理器
├── proto/
│   └── todo/v1/         # Protobuf 定义
├── gen/
│   └── proto/           # 生成的 protobuf 代码
├── models/              # 生成的 ORM 模型
├── migrations/          # 数据库迁移文件
├── bin/                 # 编译后的二进制文件
├── docker-compose.yml   # Docker Compose 配置
├── Dockerfile          # 应用镜像构建文件
├── Makefile            # 构建脚本
└── .env                # 环境变量配置
```

## 🛠️ 技术栈

- **语言**: Go
- **通信协议**: gRPC (通过 [connect-go](https://github.com/connectrpc/connect-go))
- **API 定义**: Protocol Buffers (通过 [buf](https://buf.build/))
- **数据库**: MySQL 8.0
- **ORM**: [sqlboiler](https://github.com/volatiletech/sqlboiler) (代码生成模式)
- **日志**: Go 标准库 `slog`
- **CLI**: [Cobra](https://github.com/spf13/cobra)
- **环境**: Docker & Docker Compose

## 🚀 快速开始

### 环境要求

在开始之前，请确保你的系统 (推荐 macOS 或 Linux) 上安装了以下工具：

- **Go** (推荐最新稳定版)
- **Docker** & **Docker Compose**
- **Make**
- **Buf CLI**: `brew install bufbuild/buf/buf`
- **SQLBoiler**: `brew install volatiletech/sqlboiler`
- **golang-migrate**: `brew install golang-migrate`

### 运行步骤

**克隆仓库**

```bash
git clone https://github.com/kogamitora/todo
cd todo
```

**下载 GO 依赖**

```bash
go mod tidy
```

**创建环境配置文件**

```bash
cp .env.example .env
```

**启动服务**

```bash
make docker-up
```

**测试客户端**

```bash
make test-client
```

### 其他命令

```bash
# 生成 protobuf 和 ORM 代码
make generate
```

全部命令请查看项目根目录的 [Makefile](Makefile) 文件

### [客户端使用说明](CLIENT_README.md)

## 📐 设计说明

### 1. API 设计 (`proto`)

API 定义是项目的核心。我们使用 Protocol Buffers (proto3) 来定义服务契约，确保了前后端的强类型和向后兼容性。

- **版本化**: API 定义在 `proto/todo/v1` 目录下，为未来的 API 演进（如 v2）预留了空间。
- **清晰的消息体**: 请求和响应消息被明确定义，例如 `CreateTodoRequest` 和 `CreateTodoResponse`。
- **可选字段**: 在 `UpdateTodoRequest` 中，所有字段都标记为 `optional`，允许客户端只更新部分字段，这是一种常见的 PATCH 模式。
- **枚举**: 使用 `enum` 来定义 `Status`，避免了使用 "魔术字符串"。

### 2. 服务器 (`internal/handler`)

服务器端遵循了清晰的职责分离原则。

- **Handler 层**: `internal/handler/todo_handler.go` 实现了 `TodoService` 接口。它负责处理 gRPC 请求，验证输入，调用数据库逻辑，并转换数据格式（从数据库模型到 Protobuf 消息）。
- **数据库层**: 我们没有实现复杂的 Repository 模式，而是直接在 Handler 中使用了 `SQLBoiler` 生成的模型。对于这个规模的项目，这样做更直接，减少了不必要的抽象。`SQLBoiler` 提供了类型安全的数据库查询，避免了手写 SQL 带来的拼写错误和 SQL 注入风险。
- **错误处理**: 使用 `connect.NewError` 来返回标准的 gRPC 错误码（如 `CodeNotFound`, `CodeInternal`），使客户端能更好地处理错误。
- **日志**: 使用结构化日志 `slog`，记录关键操作和错误信息，便于调试和监控。

### 3. 数据库 (`migrations` & `sqlboiler`)

- **迁移管理**: 使用 `golang-migrate` 管理数据库 schema 的演变。这使得团队协作和部署自动化变得更加可靠。
- **软删除**: `todos` 表中包含 `deleted_at` 字段，删除操作实际上是更新这个字段的时间戳，而不是物理删除数据。这是一种保护数据、便于恢复的常见做法。
- **ORM 选择**: `SQLBoiler` 是一个 "代码生成" 型 ORM。它不会像 GORM 那样使用大量反射，性能更好，并且生成的代码是类型安全的，可以在编译时捕获更多错误。

### 4. 客户端 (`cmd/client`)

- **CLI 框架**: 使用 `Cobra` 构建了一个功能丰富且用户友好的命令行界面。它支持子命令、标志（flags）和自动生成的帮助信息。
- **解耦**: 客户端逻辑完全独立于服务器，只通过生成的 gRPC 客户端与服务器通信。
