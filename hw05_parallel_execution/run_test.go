package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := range tasksCount {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for range tasksCount {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})
}

func TestRun_MZeroOrNegative_IgnoreErrors(t *testing.T) {
	tasksCount := 5
	calls := int32(0)
	// все задачи сразу возвращают ошибку, но мы ожидаем, что их всё равно запустят все
	tasks := make([]Task, 0, tasksCount)
	for range tasksCount {
		tasks = append(tasks, func() error {
			atomic.AddInt32(&calls, 1)
			return errors.New("fail")
		})
	}

	// m = 0 — по логике реализации должен игнорироваться
	err := Run(tasks, 3, 0)
	require.NoError(t, err)
	require.Equal(t, int32(tasksCount), atomic.LoadInt32(&calls),
		"должны выполнить все задачи, несмотря на ошибки")
}

func TestRun_ConcurrencyWithoutSleep(t *testing.T) {
	var started int32
	const n = 5
	const tasksCount = 5

	// barrier — чтобы задачи не завершились мгновенно
	startCh := make(chan struct{})
	tasks := make([]Task, 0, tasksCount)
	for range tasksCount {
		tasks = append(tasks, func() error {
			atomic.AddInt32(&started, 1)
			<-startCh // блокируемся, чтобы не закрыться сразу
			return nil
		})
	}

	// запускаем Run в фоне — он тут же должен стартовать n воркеров
	go func() {
		// разблокирует все задачи через секунду, чтобы Run завершился
		defer close(startCh)
		_ = Run(tasks, n, 1)
	}()

	// ждём, что в течение 100мс будет запущено ровно n воркеров
	require.Eventually(t, func() bool {
		return atomic.LoadInt32(&started) == n
	}, 100*time.Millisecond, 10*time.Millisecond,
		"должны одновременно стартовать все %d воркера", n)
}
