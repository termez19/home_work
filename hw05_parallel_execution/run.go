package hw05parallelexecution

import (
	"errors"
	"sync"
)

// ErrErrorsLimitExceeded возвращается, когда количество ошибок достигло m.
var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

// Task — функция-работа, возвращающая ошибку или nil.
type Task func() error

// Run запускает не более n ворутин для выполнения tasks и прерывает работу при m ошибках.
func Run(tasks []Task, n, m int) error {
	// Логика для m <= 0: будем игнорировать ошибки.
	if m <= 0 {
		m = len(tasks) + 1
	}

	jobs := make(chan Task)
	// Буферизованный канал ошибок, чтобы воркеры не блокировались при отправке.
	errs := make(chan error, len(tasks))
	done := make(chan struct{}) // сигнал для остановки диспетчера.

	var wg sync.WaitGroup

	// Запускаем n воркеров.
	for range n {
		wg.Add(1)
		go worker(jobs, errs, &wg)
	}

	// Диспетчер: раздаёт задачи или останавливается по сигналу done.
	go func() {
		defer close(jobs)
		for _, task := range tasks {
			select {
			case <-done:
				return
			case jobs <- task:
			}
		}
	}()

	// Считаем ошибки из воркеров.
	var errCount int
	for range tasks {
		if err := <-errs; err != nil {
			errCount++
			if errCount == m {
				// при достижении лимита ошибок — сигналить диспетчеру и ждать воркеров.
				close(done)
				wg.Wait()
				return ErrErrorsLimitExceeded
			}
		}
	}

	// Ждём завершения оставшихся воркеров перед возвратом.
	close(done)
	wg.Wait()
	return nil
}

// worker забирает задачи из jobs, выполняет и шлёт ошибки в errs.
func worker(jobs <-chan Task, errs chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range jobs {
		errs <- task()
	}
}
