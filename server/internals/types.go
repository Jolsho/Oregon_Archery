package internals

import (
	"crypto/ecdsa"
	"errors"
	"log"
	"slices"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type State struct {
	Events []Event
	Mux    *sync.RWMutex

	TimeOut					map[string]time.Time // IP -> timeout_expires
	BehaviourTracker 		map[string]*Rate 	 // IP -> bad_behaviour_rate
	SessionCreationRates	map[string]*Rate	 // IP -> rate
	DataRates 				map[string]*Rate 	 // Session Cookie -> rate

	Conns	map[string]*websocket.Conn
	Upgrader *websocket.Upgrader

	PrivKey	*ecdsa.PrivateKey
	PubKey 	*ecdsa.PublicKey
}

func New_State() *State {
	priv, err := LoadPrivateKey("KEYS.txt");
	if err != nil { panic("CANT LOAD PRIVATE KEY"); }
	return &State{
		Events: Events,
		Mux: &sync.RWMutex{},

		TimeOut: make(map[string]time.Time, 32),
		BehaviourTracker: make(map[string]*Rate, 32),
		SessionCreationRates: make(map[string]*Rate, 32),
		DataRates: make(map[string]*Rate, 32),

		Upgrader: New_Upgrader(),
		Conns: make(map[string]*websocket.Conn, 32),

		PrivKey: priv, PubKey: &priv.PublicKey,
	};
}

func (state *State) Persist() {
	err := SavePrivateKey("KEYS.txt", state.PrivKey)
	if err != nil { 
		log.Fatal("STATE_PERSIST::SAVE_KEY::",err.Error()) 
	}
}


func (state *State) Event_exists(title string) *Event {
	for i := range state.Events {
		if state.Events[i].Title == title {
			return &state.Events[i]
		}
	}
	return nil
}


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

