# Todo CLI クライアント 操作ガイド

## 概要

Todo CLI クライアント (`todocli`) は、Todo gRPC サーバーと対話するためのコマンドラインツールです。Todo アイテムの作成、表示、更新、削除といった完全な管理機能を提供します。

## インストールと設定

### クライアントのビルド

```bash
make build-client
```

### サーバー接続

クライアントはデフォルトで `http://localhost:8080` に接続します。サーバーが実行中であることを確認してください。

## 基本的な使い方

### ヘルプ情報の表示

```bash
# メインのヘルプを表示
./bin/todocli --help

# 特定のコマンドのヘルプを表示
./bin/todocli create --help
./bin/todocli get --help
./bin/todocli update --help
./bin/todocli delete --help
```

## コマンド詳細

### 1\. Todo の作成 (`create`)

新しい Todo アイテムを作成します。

#### コマンド形式

```bash
./bin/todocli create [オプション]
```

#### 利用可能なオプション

- `-t, --title string`: Todo のタイトル (必須)
- `-d, --description string`: Todo の説明 (任意)
- `--due-date string`: 期日、YYYY-MM-DD 形式 (任意)

#### 使用例

**基本的な作成（タイトルのみ）**

```bash
./bin/todocli create --title "Go言語を学習する"
```

**説明付きで作成**

```bash
./bin/todocli create --title "プロジェクトを完了させる" --description "Todoアプリの開発とテストを終える"
```

**期日付きで作成**

```bash
./bin/todocli create --title "レポートを提出する" --due-date "2024-12-31"
```

**すべてのパラメータを指定して作成**

```bash
./bin/todocli create \
    --title "重要会議" \
    --description "クライアントとプロジェクトの進捗について協議" \
    --due-date "2024-12-25"
```

**短縮オプションを使用**

```bash
./bin/todocli create -t "クイックタスク" -d "簡単なタスクの説明"
```

#### 成功時の出力例

```
Successfully created TODO item with ID: 5
```

#### エラーハンドリング

- 必須の `--title` パラメータが欠けている場合、エラーメッセージが表示されます。
- 日付の形式が正しくない場合、形式エラーメッセージが表示されます。
- サーバーに接続できない場合、接続エラーが表示されます。

---

### 2\. Todo の取得 (`get`)

すべての Todo アイテムを表示します。フィルタリングとソートをサポートしています。

#### コマンド形式

```bash
./bin/todocli get [オプション]
```

#### 利用可能なオプション

- `-s, --status string`: ステータスでフィルタリング (completed|incomplete)
- `--sort-by-due string`: 期日でソート (asc|desc)

#### 使用例

**すべての Todo を取得**

```bash
./bin/todocli get
```

**ステータスによるフィルタリング**

```bash
# 完了済みのTodoのみ表示
./bin/todocli get --status completed

# 未完了のTodoのみ表示
./bin/todocli get --status incomplete
```

**期日によるソート**

```bash
# 期日で昇順（古いものが先）
./bin/todocli get --sort-by-due asc

# 期日で降順（新しいものが先）
./bin/todocli get --sort-by-due desc
```

**フィルタリングとソートの組み合わせ**

```bash
# 未完了のTodoを、期日で昇順に表示
./bin/todocli get --status incomplete --sort-by-due asc

# 完了済みのTodoを、期日で降順に表示
./bin/todocli get --status completed --sort-by-due desc
```

#### 出力形式

```
ID  Status      Due Date    Title
----------------------------------------------------------
1   INCOMPLETE  2024-12-25  Go言語を学習する
2   COMPLETED   2024-12-30  プロジェクトを完了させる
3   INCOMPLETE  N/A         日常タスク
```

#### 出力の説明

- **ID**: Todo の一意の識別子
- **Status**: Todo のステータス（INCOMPLETE/COMPLETED）
- **Due Date**: 期日（未設定の場合は N/A）
- **Title**: Todo のタイトル

---

### 3\. Todo の更新 (`update`)

既存の Todo アイテムを更新します。

#### コマンド形式

```bash
./bin/todocli update [ID] [オプション]
```

#### 位置引数

- `ID`: 更新対象の Todo の ID（必須）

#### 利用可能なオプション

- `-t, --title string`: 新しいタイトル
- `-d, --description string`: 新しい説明
- `--due-date string`: 新しい期日、YYYY-MM-DD 形式
- `--status string`: 新しいステータス (completed|incomplete)

#### 使用例

**タイトルの更新**

```bash
./bin/todocli update 1 --title "更新後のタイトル"
```

**ステータスの更新**

```bash
# 完了としてマーク
./bin/todocli update 1 --status completed

# 未完了としてマーク
./bin/todocli update 1 --status incomplete
```

**説明の更新**

```bash
./bin/todocli update 1 --description "これは更新後の説明です"
```

**期日の更新**

```bash
./bin/todocli update 1 --due-date "2024-12-30"
```

**複数フィールドの同時更新**

```bash
./bin/todocli update 1 \
    --title "完全に更新されたタイトル" \
    --description "完全に更新された説明" \
    --due-date "2024-12-28" \
    --status completed
```

**短縮オプションを使用**

```bash
./bin/todocli update 1 -t "新しいタイトル" -d "新しい説明"
```

#### 成功時の出力例

```
Successfully updated TODO item with ID: 1
Title: 更新後のタイトル
Status: STATUS_COMPLETED
```

#### エラーハンドリング

- ID が存在しない場合、"not found" エラーが表示されます。
- ID の形式が正しくない場合、"Invalid ID" エラーが表示されます。
- ステータス値が不正な場合、エラーメッセージが表示されます。
- 日付の形式が正しくない場合、形式エラーメッセージが表示されます。

---

### 4\. Todo の削除 (`delete`)

指定された Todo アイテムを削除します（論理削除）。

#### コマンド形式

```bash
./bin/todocli delete [ID]
```

#### 位置引数

- `ID`: 削除対象の Todo の ID（必須）

#### 使用例

**Todo の削除**

```bash
./bin/todocli delete 1
```

#### 対話式の確認

削除操作では、ユーザーに確認を求めます：

```
Are you sure you want to delete TODO item with ID 1? (y/N):
```

**確認オプション：**

- `y` または `yes` を入力：削除を実行
- `n`、`no` または単に Enter キーを押す：削除をキャンセル

#### 成功時の出力例

```
Successfully deleted TODO item with ID: 1
```

#### キャンセル時の出力

```
Deletion cancelled.
```

#### エラーハンドリング

- ID が存在しない場合、"not found" エラーが表示されます。
- ID の形式が正しくない場合、"Invalid ID" エラーが表示されます。

## エラーハンドリングとトラブルシューティング

### よくあるエラーと解決策

#### 1\. 接続エラー

```
Failed to create todo: connect: connection refused
```

**解決策：**

- サーバーが実行中であることを確認してください： `make run-server`
- サーバーのアドレスとポートが正しいか確認してください。

#### 2\. ID が存在しないエラー

```
Failed to update todo: not found
```

**解決策：**

- `./bin/todocli get` を使用して利用可能な ID を確認してください。
- ID が数字であり、存在することを確認してください。

#### 3\. 引数エラー

```
Error: required flag(s) "title" not set
```

**解決策：**

- 必須の引数が提供されているか確認してください。
- `--help` を使用して正しい引数の形式を確認してください。

#### 4\. 日付形式エラー

```
Invalid due date format. Use YYYY-MM-DD
```

**解決策：**

- 日付形式が YYYY-MM-DD であることを確認してください。
- 例：2024-12-25

#### 5\. ステータス値エラー

```
Invalid status. Use 'completed' or 'incomplete'
```

**解決策：**

- 正しいステータス値 `completed` または `incomplete` を使用してください。
