package state

import (
	"server/utils"
	"sync"
	"time"
)

type State struct {
	Events       []Event
	EventMux     *sync.RWMutex
	workers		 *utils.WorkGroup
}

func New_State() *State {

	state := State{
		Events:       Events,
		EventMux:     &sync.RWMutex{},
		workers: 	  utils.NewWorkGroup(),
	}

	// EVENT CLEANER
	tasks, taskMux, needsCleaning := utils.NewCleaningTasks(
		&EventCleaningTask{
			Array: &state.Events,
			Mux:   state.EventMux,
			When:  time.Now().Add(CLEANING_INTERVAL),
		},
	)
	go utils.Task_Scheduler(&tasks, taskMux, CLEANING_INTERVAL, needsCleaning, state.workers)
	go utils.RunCleaner(needsCleaning, taskMux, tasks, state.workers)

	return &state
}
func (state *State) Event_exists(title string) *Event {
	for i := range state.Events {
		if state.Events[i].Title == title {
			return &state.Events[i]
		}
	}
	return nil
}

func (state *State) Shutdown() {
	state.workers.Cancel();
	state.workers.WG.Wait();
}
