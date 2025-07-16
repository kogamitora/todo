# Todo CLI 客户端操作指南

## 概述

Todo CLI 客户端 (`todocli`) 是一个命令行工具，用于与 Todo gRPC 服务器进行交互。它提供了完整的 Todo 管理功能，包括创建、查看、更新和删除 Todo 项目。

## 安装和配置

### 构建客户端

```bash
make build-client
```

### 服务器连接

客户端默认连接到 `http://localhost:8080`。请确保服务器正在运行

## 基本使用

### 查看帮助信息

```bash
# 查看主帮助
./bin/todocli --help

# 查看具体命令的帮助
./bin/todocli create --help
./bin/todocli get --help
./bin/todocli update --help
./bin/todocli delete --help
```

## 详细命令说明

### 1. 创建 Todo (`create`)

创建新的 Todo 项目。

#### 命令格式

```bash
./bin/todocli create [选项]
```

#### 可用选项

- `-t, --title string`: Todo 标题 (必需)
- `-d, --description string`: Todo 描述 (可选)
- `--due-date string`: 截止日期，格式为 YYYY-MM-DD (可选)

#### 使用示例

**基本创建（仅标题）**

```bash
./bin/todocli create --title "学习 Go 语言"
```

**带描述的创建**

```bash
./bin/todocli create --title "完成项目" --description "完成 Todo 应用的开发和测试"
```

**带截止日期的创建**

```bash
./bin/todocli create --title "提交报告" --due-date "2024-12-31"
```

**完整参数创建**

```bash
./bin/todocli create \
    --title "重要会议" \
    --description "与客户讨论项目进展" \
    --due-date "2024-12-25"
```

**使用短参数**

```bash
./bin/todocli create -t "快速任务" -d "简单的任务描述"
```

#### 成功输出示例

```
Successfully created TODO item with ID: 5
```

#### 错误处理

- 如果缺少必需的 `--title` 参数，会显示错误信息
- 如果日期格式不正确，会显示格式错误信息
- 如果无法连接到服务器，会显示连接错误

---

### 2. 获取 Todo (`get`)

查看所有 Todo 项目，支持过滤和排序。

#### 命令格式

```bash
./bin/todocli get [选项]
```

#### 可用选项

- `-s, --status string`: 按状态过滤 (completed|incomplete)
- `--sort-by-due string`: 按截止日期排序 (asc|desc)

#### 使用示例

**获取所有 Todo**

```bash
./bin/todocli get
```

**按状态过滤**

```bash
# 只显示已完成的 Todo
./bin/todocli get --status completed

# 只显示未完成的 Todo
./bin/todocli get --status incomplete
```

**按截止日期排序**

```bash
# 按截止日期升序排列（最早的在前面）
./bin/todocli get --sort-by-due asc

# 按截止日期降序排列（最晚的在前面）
./bin/todocli get --sort-by-due desc
```

**组合过滤和排序**

```bash
# 显示未完成的 Todo，按截止日期升序排列
./bin/todocli get --status incomplete --sort-by-due asc

# 显示已完成的 Todo，按截止日期降序排列
./bin/todocli get --status completed --sort-by-due desc
```

#### 输出格式

```
ID	Status		Due Date	Title
----------------------------------------------------------
1	INCOMPLETE	2024-12-25	学习 Go 语言
2	COMPLETED	2024-12-30	完成项目
3	INCOMPLETE	N/A		日常任务
```

#### 输出说明

- **ID**: Todo 的唯一标识符
- **Status**: Todo 状态（INCOMPLETE/COMPLETED）
- **Due Date**: 截止日期（如果没有设置则显示 N/A）
- **Title**: Todo 标题

---

### 3. 更新 Todo (`update`)

更新现有的 Todo 项目。

#### 命令格式

```bash
./bin/todocli update [ID] [选项]
```

#### 位置参数

- `ID`: 要更新的 Todo 的 ID（必需）

#### 可用选项

- `-t, --title string`: 新的标题
- `-d, --description string`: 新的描述
- `--due-date string`: 新的截止日期，格式为 YYYY-MM-DD
- `--status string`: 新的状态 (completed|incomplete)

#### 使用示例

**更新标题**

```bash
./bin/todocli update 1 --title "更新后的标题"
```

**更新状态**

```bash
# 标记为已完成
./bin/todocli update 1 --status completed

# 标记为未完成
./bin/todocli update 1 --status incomplete
```

**更新描述**

```bash
./bin/todocli update 1 --description "这是更新后的描述"
```

**更新截止日期**

```bash
./bin/todocli update 1 --due-date "2024-12-30"
```

**同时更新多个字段**

```bash
./bin/todocli update 1 \
    --title "完全更新的标题" \
    --description "完全更新的描述" \
    --due-date "2024-12-28" \
    --status completed
```

**使用短参数**

```bash
./bin/todocli update 1 -t "新标题" -d "新描述"
```

#### 成功输出示例

```
Successfully updated TODO item with ID: 1
Title: 更新后的标题
Status: STATUS_COMPLETED
```

#### 错误处理

- 如果 ID 不存在，会显示 "not found" 错误
- 如果 ID 格式不正确，会显示 "Invalid ID" 错误
- 如果状态值不正确，会显示状态错误信息
- 如果日期格式不正确，会显示格式错误信息

---

### 4. 删除 Todo (`delete`)

删除指定的 Todo 项目（软删除）。

#### 命令格式

```bash
./bin/todocli delete [ID]
```

#### 位置参数

- `ID`: 要删除的 Todo 的 ID（必需）

#### 使用示例

**删除 Todo**

```bash
./bin/todocli delete 1
```

#### 交互式确认

删除操作会要求用户确认：

```
Are you sure you want to delete TODO item with ID 1? (y/N):
```

**确认选项：**

- 输入 `y` 或 `yes`：确认删除
- 输入 `n`、`no` 或直接回车：取消删除

#### 成功输出示例

```
Successfully deleted TODO item with ID: 1
```

#### 取消删除输出

```
Deletion cancelled.
```

#### 错误处理

- 如果 ID 不存在，会显示 "not found" 错误
- 如果 ID 格式不正确，会显示 "Invalid ID" 错误

## 错误处理和故障排除

### 常见错误及解决方法

#### 1. 连接错误

```
Failed to create todo: connect: connection refused
```

**解决方法：**

- 确保服务器正在运行：`make run-server`
- 检查服务器地址和端口是否正确

#### 2. ID 不存在错误

```
Failed to update todo: not found
```

**解决方法：**

- 使用 `./bin/todocli get` 查看可用的 ID
- 确保 ID 是数字且存在

#### 3. 参数错误

```
Error: required flag(s) "title" not set
```

**解决方法：**

- 检查必需参数是否提供
- 使用 `--help` 查看正确的参数格式

#### 4. 日期格式错误

```
Invalid due date format. Use YYYY-MM-DD
```

**解决方法：**

- 确保日期格式为 YYYY-MM-DD
- 例如：2024-12-25

#### 5. 状态值错误

```
Invalid status. Use 'completed' or 'incomplete'
```

**解决方法：**

- 使用正确的状态值：`completed` 或 `incomplete`
