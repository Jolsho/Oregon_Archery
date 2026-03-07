package state

import (
	"errors"
	"slices"
	"sync"
	"time"
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
	Name      string `json:"name"`
	Threshold int    `json:"threshold"`
}

type Event struct {
	Title         string          `json:"title"`
	IsOwn         bool            `json:"is_own"`
	IsPersisted   bool            `json:"is_persisted"`
	Leaders       map[string]int  `json:"leaders"`
	Divisions     []Division      `json:"divisions"`
	Teams         map[string]Team `json:"teams"`
	Kind          string          `json:"kind"`
	ScoresPerTeam int             `json:"scores_per_team"`
	Expires       time.Time       `json:"expires"`
	CreatedAt     time.Time       `json:"created_at"`
	Secret        string          `json:"-"`
}
const EVENT_LIFESPAN = 24 * time.Hour

func DefaultEvent() Event {
	return Event{
		Divisions: DIVISIONS,
		Kind:      "OUTDOOR",
		IsPersisted: true,
		CreatedAt: time.Now(),
		Expires: time.Now().Add(EVENT_LIFESPAN),
	}
}

func (ev *Event) Sanitize() error {
	if ev.Title == "" {
		return errors.New("NO TITLE")
	}

	ev.IsPersisted = true;

	if len(ev.Divisions) > 10 {
		return errors.New("NO MANY DIVISIONS")
	}

	if (time.Since(ev.CreatedAt) < 0) {
		ev.CreatedAt = time.Now();
		ev.Expires = ev.CreatedAt.Add(EVENT_LIFESPAN);
	}

	if !slices.Contains(kinds, ev.Kind) {
		return errors.New("INVALID KIND")
	}
	return nil
}

const CLEANING_INTERVAL = 10 * time.Minute

type EventCleaningTask struct {
	Array *[]Event
	Mux   *sync.RWMutex
	When  time.Time
}

func (task *EventCleaningTask) Get_when() time.Time { return task.When }
func (task *EventCleaningTask) Clean() {
	now := time.Now()

	task.Mux.Lock()
	defer task.Mux.Unlock()

	events := *task.Array
	kept := events[:0] // reuse underlying array
	next := now.Add(CLEANING_INTERVAL)

	for _, event := range events {
		if event.Expires.After(now) {
			kept = append(kept, event)

			if event.Expires.Before(next) {
				next = event.Expires
			}
		}
	}

	*task.Array = kept
	task.When = next
}
