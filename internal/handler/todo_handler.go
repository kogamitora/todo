/*
internal/handler/todo_handler.go: 业务逻辑的核心。
作用: 这个文件里的 TodoHandler 结构体，实现了由 buf 在 gen/ 目录中生成的 TodoServiceHandler 接口。对于每个 API 调用（如 CreateTodo），它会：
接收 protobuf 格式的请求数据。
将其转换为 sqlboiler 生成的数据库模型。
调用 sqlboiler 的方法与数据库交互。
将数据库返回的模型转换回 protobuf 格式的响应数据。
返回响应或错误。
*/

package handler

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"connectrpc.com/connect"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"google.golang.org/protobuf/types/known/timestamppb"

	todov1 "github.com/kogamitora/todo/gen/proto/todo/v1"
	"github.com/kogamitora/todo/gen/proto/todo/v1/v1connect"
	"github.com/kogamitora/todo/models"
)

// TodoHandler 实现了 TodoService
type TodoHandler struct {
	db     *sql.DB
	logger *slog.Logger
}

// 确保 TodoHandler 实现了接口
var _ v1connect.TodoServiceHandler = (*TodoHandler)(nil)

func NewTodoHandler(db *sql.DB, logger *slog.Logger) *TodoHandler {
	return &TodoHandler{
		db:     db,
		logger: logger,
	}
}

// modelToProto 将数据库模型转换为 protobuf 消息
func modelToProto(t *models.Todo) *todov1.Todo {
	todo := &todov1.Todo{
		Id:        t.ID,
		Title:     t.Title,
		CreatedAt: timestamppb.New(t.CreatedAt),
		UpdatedAt: timestamppb.New(t.UpdatedAt),
	}
	if t.Description.Valid {
		todo.Description = t.Description.String
	}
	if t.DueDate.Valid {
		todo.DueDate = timestamppb.New(t.DueDate.Time)
	}
	// 将数据库 ENUM 字符串转换为 protobuf ENUM
	switch t.Status {
	case "TODO_STATUS_INCOMPLETE":
		todo.Status = todov1.Status_STATUS_INCOMPLETE
	case "TODO_STATUS_COMPLETED":
		todo.Status = todov1.Status_STATUS_COMPLETED
	default:
		todo.Status = todov1.Status_STATUS_UNSPECIFIED
	}
	return todo
}

func (h *TodoHandler) CreateTodo(ctx context.Context, req *connect.Request[todov1.CreateTodoRequest]) (*connect.Response[todov1.CreateTodoResponse], error) {
	h.logger.Info("CreateTodo called", "title", req.Msg.Title)

	newTodo := &models.Todo{
		Title: req.Msg.Title,
	}
	if req.Msg.Description != "" {
		newTodo.Description.String = req.Msg.Description
		newTodo.Description.Valid = true
	}
	if req.Msg.DueDate != nil && req.Msg.DueDate.IsValid() {
		newTodo.DueDate.Time = req.Msg.DueDate.AsTime()
		newTodo.DueDate.Valid = true
	}

	err := newTodo.Insert(ctx, h.db, boil.Infer())
	if err != nil {
		h.logger.Error("failed to insert todo", "error", err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&todov1.CreateTodoResponse{
		Todo: modelToProto(newTodo),
	}), nil
}

func (h *TodoHandler) GetTodo(ctx context.Context, req *connect.Request[todov1.GetTodoRequest]) (*connect.Response[todov1.GetTodoResponse], error) {
	h.logger.Info("GetTodo called", "id", req.Msg.Id)

	todo, err := models.Todos(models.TodoWhere.ID.EQ(req.Msg.Id)).One(ctx, h.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}
		h.logger.Error("failed to get todo", "error", err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&todov1.GetTodoResponse{
		Todo: modelToProto(todo),
	}), nil
}

func (h *TodoHandler) UpdateTodo(ctx context.Context, req *connect.Request[todov1.UpdateTodoRequest]) (*connect.Response[todov1.UpdateTodoResponse], error) {
	h.logger.Info("UpdateTodo called", "id", req.Msg.Id)

	todo, err := models.FindTodo(ctx, h.db, req.Msg.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}
		h.logger.Error("failed to find todo for update", "error", err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if req.Msg.Title != nil {
		todo.Title = *req.Msg.Title
	}
	if req.Msg.Description != nil {
		todo.Description.String = *req.Msg.Description
		todo.Description.Valid = true
	}
	if req.Msg.DueDate != nil {
		todo.DueDate.Time = req.Msg.DueDate.AsTime()
		todo.DueDate.Valid = true
	}
	if req.Msg.Status != nil {
		switch *req.Msg.Status {
		case todov1.Status_STATUS_INCOMPLETE:
			todo.Status = "TODO_STATUS_INCOMPLETE"
		case todov1.Status_STATUS_COMPLETED:
			todo.Status = "TODO_STATUS_COMPLETED"
		}
	}

	_, err = todo.Update(ctx, h.db, boil.Infer())
	if err != nil {
		h.logger.Error("failed to update todo", "error", err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&todov1.UpdateTodoResponse{
		Todo: modelToProto(todo),
	}), nil
}

func (h *TodoHandler) DeleteTodo(ctx context.Context, req *connect.Request[todov1.DeleteTodoRequest]) (*connect.Response[todov1.DeleteTodoResponse], error) {
	h.logger.Info("DeleteTodo called", "id", req.Msg.Id)

	todo, err := models.FindTodo(ctx, h.db, req.Msg.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}
		h.logger.Error("failed to find todo for delete", "error", err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// 软删除
	todo.DeletedAt.Time = time.Now()
	todo.DeletedAt.Valid = true
	_, err = todo.Update(ctx, h.db, boil.Whitelist(models.TodoColumns.DeletedAt))
	if err != nil {
		h.logger.Error("failed to soft delete todo", "error", err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&todov1.DeleteTodoResponse{}), nil
}

func (h *TodoHandler) ListTodos(ctx context.Context, req *connect.Request[todov1.ListTodosRequest]) (*connect.Response[todov1.ListTodosResponse], error) {
	h.logger.Info("ListTodos called")

	queryMods := []qm.QueryMod{
		models.TodoWhere.DeletedAt.IsNull(),
	}

	// 默认按创建时间降序排序
	orderByClause := models.TodoColumns.CreatedAt + " DESC"

	if req.Msg.SortByDueDate != nil {
		dueDateColumn := models.TodoColumns.DueDate
		switch *req.Msg.SortByDueDate {
		case todov1.SortOrder_SORT_ORDER_ASC:
			orderByClause = fmt.Sprintf("CASE WHEN %s IS NULL THEN 1 ELSE 0 END, %s ASC", dueDateColumn, dueDateColumn)

		case todov1.SortOrder_SORT_ORDER_DESC:
			orderByClause = fmt.Sprintf("CASE WHEN %s IS NULL THEN 0 ELSE 1 END, %s DESC", dueDateColumn, dueDateColumn)
		}
	}

	queryMods = append(queryMods, qm.OrderBy(orderByClause))

	if req.Msg.StatusFilter != nil {
		var statusStr string
		switch *req.Msg.StatusFilter {
		case todov1.Status_STATUS_INCOMPLETE:
			statusStr = "TODO_STATUS_INCOMPLETE"
		case todov1.Status_STATUS_COMPLETED:
			statusStr = "TODO_STATUS_COMPLETED"
		}
		if statusStr != "" {
			queryMods = append(queryMods, models.TodoWhere.Status.EQ(statusStr))
		}
	}

	// 调试：输出生成的 SQL
	h.logger.Info("Generated ORDER BY clause", "clause", orderByClause)

	// ... 剩下的代码与之前相同 ...
	todos, err := models.Todos(queryMods...).All(ctx, h.db)
	if err != nil {
		h.logger.Error("failed to list todos", "error", err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	protoTodos := make([]*todov1.Todo, len(todos))
	for i, t := range todos {
		protoTodos[i] = modelToProto(t)
	}

	return connect.NewResponse(&todov1.ListTodosResponse{
		Todos: protoTodos,
	}), nil
}
