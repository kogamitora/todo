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

// TodoService
type TodoHandler struct {
	db     *sql.DB
	logger *slog.Logger
}

var _ v1connect.TodoServiceHandler = (*TodoHandler)(nil)

func NewTodoHandler(db *sql.DB, logger *slog.Logger) *TodoHandler {
	return &TodoHandler{
		db:     db,
		logger: logger,
	}
}

// modelToProto converts a Todo model to a protobuf Todo message.
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
	//
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

	todo, err := h.findTodoByID(ctx, req.Msg.Id)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&todov1.GetTodoResponse{
		Todo: modelToProto(todo),
	}), nil
}

func (h *TodoHandler) UpdateTodo(ctx context.Context, req *connect.Request[todov1.UpdateTodoRequest]) (*connect.Response[todov1.UpdateTodoResponse], error) {
	h.logger.Info("UpdateTodo called", "id", req.Msg.Id)

	todo, err := h.findTodoByID(ctx, req.Msg.Id)
	if err != nil {
		return nil, err
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

	todo, err := h.findTodoByID(ctx, req.Msg.Id)
	if err != nil {
		return nil, err
	}

	// soft delete: set DeletedAt to current time
	todo.DeletedAt.Time = time.Now()
	todo.DeletedAt.Valid = true
	_, err = todo.Update(ctx, h.db, boil.Whitelist(models.TodoColumns.DeletedAt))
	if err != nil {
		h.logger.Error("failed to soft delete todo", "error", err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&todov1.DeleteTodoResponse{}), nil
}

func (h *TodoHandler) GetTodos(ctx context.Context, req *connect.Request[todov1.GetTodosRequest]) (*connect.Response[todov1.GetTodosResponse], error) {
	h.logger.Info("GetTodos called")

	queryMods := []qm.QueryMod{
		models.TodoWhere.DeletedAt.IsNull(),
	}

	// sort by created_at DESC by default
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

	h.logger.Info("Generated ORDER BY clause", "clause", orderByClause)

	todos, err := models.Todos(queryMods...).All(ctx, h.db)
	if err != nil {
		h.logger.Error("failed to list todos", "error", err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	protoTodos := make([]*todov1.Todo, len(todos))
	for i, t := range todos {
		protoTodos[i] = modelToProto(t)
	}

	return connect.NewResponse(&todov1.GetTodosResponse{
		Todos: protoTodos,
	}), nil
}

// findTodoByID finds a todo by its ID and handles common errors.
func (h *TodoHandler) findTodoByID(ctx context.Context, id int64) (*models.Todo, error) {
	todo, err := models.FindTodo(ctx, h.db, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("todo with id %d not found", id))
		}
		h.logger.Error("failed to find todo", "id", id, "error", err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return todo, nil
}
