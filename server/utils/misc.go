package utils

import (
	"sync"
	"time"
)


type Task interface {
	Get_when() time.Time
}

func NewCleaningTasks(
	raw_tasks ...CleaningTask,
) (tasks []CleaningTask, taskMux *sync.Mutex, needsCleaning chan struct{}) {
	taskMux = &sync.Mutex{};
	needsCleaning = make(chan struct{}, len(tasks));
	tasks = raw_tasks;
	return;
}

func Task_Scheduler[T Task](
	tasks *[]T,
	mux *sync.Mutex, 
	default_interval time.Duration, 
	notifier chan struct{},
	group *WorkGroup,
) {
	group.WG.Add(1);
	var timer *time.Timer

	for {
		mux.Lock()
		if len(*tasks) == 0 {
			mux.Unlock()
			select {
			case <-time.After(default_interval):
				continue
			case <-group.Ctx.Done():
				group.WG.Done();
				return
			}
		}

		when := (*tasks)[0].Get_when()
		mux.Unlock()

		wait := max(time.Until(when), 0)

		if timer == nil {
			timer = time.NewTimer(wait)
		} else {
			timer.Reset(wait)
		}

		select {
		case <-timer.C:
			select {
			case notifier <- struct{}{}:
			case <-group.Ctx.Done():
				group.WG.Done();
				return
			}
		case <-group.Ctx.Done():
			group.WG.Done();
			return
		}
	}
}

type CleaningTask interface {
	Clean()
	Get_when() time.Time
}


func RunCleaner(
	needsCleaning <- chan struct{}, 
	taskMux *sync.Mutex, 
	tasks *[]CleaningTask,
	group *WorkGroup,
) {
	group.WG.Add(1);
	for {
		select {
		case <-needsCleaning:
			taskMux.Lock()

			// pop earliest
			task := (*tasks)[0]
			copy((*tasks), (*tasks)[1:])
			(*tasks) = (*tasks)[:len(*tasks)-1]

			// update schedule
			task.Clean()

			// reinsert sorted by When
			i := 0
			when := task.Get_when();
			for i < len(*tasks) && (*tasks)[i].Get_when().Before(when) {
				i++
			}

			*tasks = append(*tasks, nil) // grow
			copy((*tasks)[i+1:], (*tasks)[i:])
			(*tasks)[i] = task

			taskMux.Unlock()

		case <-group.Ctx.Done():
			group.WG.Done();
			return;
		}
	}
}
