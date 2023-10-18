package wpool

import (
	"context"
	"sync"
)

// Воркер пул для запуска задач с заданными типами входных и выходных данных.
// Работа с пулом:
// 1. Создание пула.
// 2. Запуск задач. После выполнения задач пул закрывается.
type WorkerPool[In, Out any] interface {
	// Запуск задач в пуле
	StartTasks(ctx context.Context, tasks []Task[In, Out])
	// Получения результатов
	Results() <-chan Result[Out]
}

type workerPool[In, Out any] struct {
	// Число рабочих горутин
	workerCount int
	// Канал для задач
	tasks chan Task[In, Out]
	// Канал для результатов
	results chan Result[Out]
}

func NewWorkerPool[In, Out any](workerCount int) WorkerPool[In, Out] {
	return &workerPool[In, Out]{
		workerCount: workerCount,
		tasks:       make(chan Task[In, Out], workerCount),
		results:     make(chan Result[Out], workerCount),
	}
}

// Запуск выполнения задач в пулае
func (wp *workerPool[In, Out]) StartTasks(ctx context.Context, tasks []Task[In, Out]) {
	go wp.AddTasks(tasks)

	var wg sync.WaitGroup

	// запускаем всех воркеров
	for i := 0; i < wp.workerCount; i++ {
		wg.Add(1)
		go worker(ctx, &wg, wp.tasks, wp.results)
	}

	// ждем окончания работы воркеров и закрываем канал с результатами
	wg.Wait()
	close(wp.results)
}

// Получение результатов из канала (сам канал не торчит наружу)
func (wp *workerPool[In, Out]) Results() <-chan Result[Out] {
	return wp.results
}

// Добавляем пачку задач и закрываем входной канал
func (wp *workerPool[In, Out]) AddTasks(tasks []Task[In, Out]) {
	for i := range tasks {
		wp.tasks <- tasks[i]
	}
	close(wp.tasks)
}

// Обработка задачи одним воркером
func worker[In, Out any](ctx context.Context, wg *sync.WaitGroup, tasks <-chan Task[In, Out], results chan<- Result[Out]) {
	defer wg.Done()

	for {
		select {
		case task, ok := <-tasks:
			// Получаем задачу и запускаем
			if !ok {
				return
			}
			results <- task.Execute(ctx)
		case <-ctx.Done():
			// При закрытии контекста возвращаем ошибку
			results <- Result[Out]{
				Error: ctx.Err(),
			}
			return
		}
	}
}
