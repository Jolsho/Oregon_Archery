package internals

import (
	"sync"
	"errors"
	"slices"
)

type Member struct {
	Name     string `json:"name"`
	Division int    `json:"division"`
	Score    int    `json:"score"`
	XCount   int    `json:"x_count"`
}

type Team struct {
	Name    string   `json:"name"`
	Score   int      `json:"score"`
	XCount  int      `json:"x_count"`
	Members []Member `json:"members"`
}

type Division struct {
	Name    string  `json:"name"`
	Threshold int  	`json:"threshold"`
}

type Event struct {
	Title     string          `json:"title"`
	IsOwn	  bool			  `json:"is_own"`
	Leaders   map[string]int  `json:"leaders"`
	Divisions []Division      `json:"divisions"`
	Teams     map[string]Team `json:"teams"`
	Kind      string          `json:"kind"`
	Secret 	  string 		  `json:"-"`
	ScoresPerTeam int         `json:"scores_per_team"`
}

func DefaultEvent() Event {
	return Event{
		Divisions: DIVISIONS,
		Kind: "OUTDOOR",
	};
}

func (ev *Event) Sanitize() error {
	if (ev.Title == "") {
		return errors.New("NO TITLE");
	}

	ev.Divisions = DIVISIONS;

	if (!slices.Contains(kinds, ev.Kind)) {
		return errors.New("INVALID KIND");
	}


	return nil;
}



// TODO -- NEED EVENT CLEANING LOGIC

type State struct {
	Events []Event
	EventMux    *sync.RWMutex
}
func New_State() *State {
	return &State{
		Events: Events,
		EventMux: &sync.RWMutex{},

	};
}
func (state *State) Event_exists(title string) *Event {
	for i := range state.Events {
		if state.Events[i].Title == title {
			return &state.Events[i]
		}
	}
	return nil
}
