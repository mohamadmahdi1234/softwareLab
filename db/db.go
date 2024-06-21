package db

import (
	"context"
	"simpleAPI/db/entity"
	"time"
)

type DB interface {
	GetTodoItems(ctx context.Context) (*entity.Hall, error)
	DeleteTodoItem(ctx context.Context, todoID int)
	UpdateTodoItem(ctx context.Context, prevTodoItem, nextTodoItem *entity.Hall)
	AddTodoItem(ctx context.Context, todoItem *entity.Hall) error
	ReserveHall(ctx context.Context, req entity.ReserveReq, date time.Time) error
}
