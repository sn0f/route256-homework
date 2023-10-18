package wpool

import (
	"context"
)

type Fn[In, Out any] func(ctx context.Context, args In) (Out, error)

// Задача для запуска функции с заданными типами входных и выходных данных
type Task[In, Out any] struct {
	Func Fn[In, Out]
	Args In
}

// Результат выполнения задачи содержит данные и ошибку
type Result[data any] struct {
	Data  data
	Error error
}

// При создании задачи передаем функцию и аргументы
func NewTask[In, Out any](fn Fn[In, Out], args In) *Task[In, Out] {
	return &Task[In, Out]{
		Func: fn,
		Args: args,
	}
}

// Выполняет функцию из задачи и возвращает результат/ошибку
func (task Task[In, Out]) Execute(ctx context.Context) (result Result[Out]) {
	value, err := task.Func(ctx, task.Args)
	if err != nil {
		result.Error = err
		return
	}

	result.Data = value
	return
}
