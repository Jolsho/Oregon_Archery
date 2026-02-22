package data

import (
	"math/rand"
	"sync"
	crypt_rand "crypto/rand" 
	"encoding/base64"
)

type State struct {
	Events []Event
	Mux    *sync.RWMutex
}
func (state *State) Event_exists(title string) *Event {
	for _, event := range state.Events {
		if event.Title == title {
			return &event;
		}
	}
	return nil;
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

type Event struct {
	Title     string          `json:"title"`
	IsOwn	  bool			  `json:"is_own"`
	Leaders   map[string]int  `json:"leaders"`
	Divisions []string        `json:"divisions"`
	Teams     map[string]Team `json:"teams"`
	Secret 	  string 		   `json:"-"`
}

// assumed to exist elsewhere
var DIVISIONS = []string{"OPEN", "MODERN", "OLYMPIC", "TRADITIONAL"}
func randDivision() int {
	return rand.Intn(len(DIVISIONS));
}


func Random_secret() (string, error) {
	rnd_secret := make([]byte, 8);
	n, err := crypt_rand.Read(rnd_secret);
	if n < 8 || err != nil { return "", err }
	return base64.RawURLEncoding.EncodeToString(rnd_secret), nil
}

func rand_secret_ignored() string {
	sec, _ := Random_secret();
	return sec;
}

var Events = []Event{
	{
		Title:     "Thurston vs Jesuit",
		Leaders:   map[string]int{},
		Divisions: DIVISIONS,
		Secret: rand_secret_ignored(),
		Teams: map[string]Team{
			"Thurston": {
				Name:   "Thurston",
				Score:  0,
				XCount: 0,
				Members: []Member{
					{Name: "Alex Morgan", Division: randDivision(), Score: 412, XCount: 6},
					{Name: "Chris Nolan", Division: randDivision(), Score: 398, XCount: 3},
					{Name: "Taylor Reed", Division: randDivision(), Score: 430, XCount: 11},
					{Name: "Jordan Kim", Division: randDivision(), Score: 445, XCount: 9},
					{Name: "Sam Ortega", Division: randDivision(), Score: 376, XCount: 2},
					{Name: "Riley Chen", Division: randDivision(), Score: 459, XCount: 14},
				},
			},
			"Jesuit": {
				Name:   "Jesuit",
				Score:  0,
				XCount: 0,
				Members: []Member{
					{Name: "Evan Brooks", Division: randDivision(), Score: 438, XCount: 10},
					{Name: "Noah Patel", Division: randDivision(), Score: 471, XCount: 18},
					{Name: "Liam O'Connor", Division: randDivision(), Score: 401, XCount: 5},
					{Name: "Diego Ramirez", Division: randDivision(), Score: 455, XCount: 12},
					{Name: "Mason Wright", Division: randDivision(), Score: 429, XCount: 7},
				},
			},
		},
	},

	{
		Title:     "Oregon Outdoor State Championship",
		Leaders:   map[string]int{},
		Divisions: DIVISIONS,
		Secret: rand_secret_ignored(),
		Teams: map[string]Team{
			"Central Catholic": {
				Name:   "Central Catholic",
				Score:  0,
				XCount: 0,
				Members: []Member{
					{Name: "Sophia Lin", Division: randDivision(), Score: 472, XCount: 19},
					{Name: "Marcus Hill", Division: randDivision(), Score: 488, XCount: 22},
					{Name: "Tyler Bennett", Division: randDivision(), Score: 405, XCount: 6},
					{Name: "Avery Collins", Division: randDivision(), Score: 451, XCount: 13},
					{Name: "Julian Perez", Division: randDivision(), Score: 479, XCount: 20},
					{Name: "Harper Nguyen", Division: randDivision(), Score: 392, XCount: 3},
					{Name: "Kevin Foster", Division: randDivision(), Score: 436, XCount: 9},
				},
			},
			"Jesuit": {
				Name:   "Jesuit",
				Score:  0,
				XCount: 0,
				Members: []Member{
					{Name: "Isabella Cruz", Division: randDivision(), Score: 469, XCount: 18},
					{Name: "Andrew Park", Division: randDivision(), Score: 492, XCount: 24},
					{Name: "Sean Murphy", Division: randDivision(), Score: 410, XCount: 7},
					{Name: "Daniel Kim", Division: randDivision(), Score: 447, XCount: 11},
					{Name: "Victor Gomez", Division: randDivision(), Score: 475, XCount: 19},
				},
			},
			"West Linn": {
				Name:   "West Linn",
				Score:  0,
				XCount: 0,
				Members: []Member{
					{Name: "Olivia Hart", Division: randDivision(), Score: 455, XCount: 14},
					{Name: "Ryan Cooper", Division: randDivision(), Score: 399, XCount: 4},
					{Name: "Ethan Zhou", Division: randDivision(), Score: 483, XCount: 21},
					{Name: "Lucas Grant", Division: randDivision(), Score: 442, XCount: 10},
					{Name: "Nina Petrova", Division: randDivision(), Score: 468, XCount: 17},
					{Name: "Cole Harrison", Division: randDivision(), Score: 387, XCount: 2},
				},
			},
		},
	},

	{
		Title:     "OHSAL Public Open",
		Leaders:   map[string]int{},
		Divisions: DIVISIONS,
		Secret: rand_secret_ignored(),
		Teams: map[string]Team{
			"Springfield": {
				Name:   "Springfield",
				Score:  0,
				XCount: 0,
				Members: []Member{
					{Name: "Parker Young", Division: randDivision(), Score: 432, XCount: 8},
					{Name: "Logan Reed", Division: randDivision(), Score: 467, XCount: 16},
					{Name: "Dylan Moore", Division: randDivision(), Score: 390, XCount: 3},
					{Name: "Elliot Baker", Division: randDivision(), Score: 446, XCount: 12},
					{Name: "Miguel Santos", Division: randDivision(), Score: 458, XCount: 14},
				},
			},
			"North Salem": {
				Name:   "North Salem",
				Score:  0,
				XCount: 0,
				Members: []Member{
					{Name: "Adrian Fox", Division: randDivision(), Score: 419, XCount: 7},
					{Name: "Ben Wallace", Division: randDivision(), Score: 401, XCount: 5},
					{Name: "Ivan Petrov", Division: randDivision(), Score: 471, XCount: 18},
					{Name: "Theo Ramirez", Division: randDivision(), Score: 437, XCount: 10},
					{Name: "Kira Volkov", Division: randDivision(), Score: 460, XCount: 15},
					{Name: "Jesse Long", Division: randDivision(), Score: 382, XCount: 2},
				},
			},
			"McKenzie": {
				Name:   "McKenzie",
				Score:  0,
				XCount: 0,
				Members: []Member{
					{Name: "Carter Mills", Division: randDivision(), Score: 428, XCount: 9},
					{Name: "Wyatt Dean", Division: randDivision(), Score: 395, XCount: 4},
					{Name: "Hugo Laurent", Division: randDivision(), Score: 463, XCount: 17},
					{Name: "Tariq Hassan", Division: randDivision(), Score: 441, XCount: 11},
					{Name: "Miles Carter", Division: randDivision(), Score: 452, XCount: 13},
				},
			},
			"Ashland": {
				Name:   "Ashland",
				Score:  0,
				XCount: 0,
				Members: []Member{
					{Name: "Noelle Wright", Division: randDivision(), Score: 435, XCount: 9},
					{Name: "Grant Olson", Division: randDivision(), Score: 388, XCount: 3},
					{Name: "Soren Dahl", Division: randDivision(), Score: 469, XCount: 18},
					{Name: "Priya Nair", Division: randDivision(), Score: 449, XCount: 12},
					{Name: "Leo Martinez", Division: randDivision(), Score: 457, XCount: 14},
				},
			},
		},
	},

	{
		Title:     "4A State Qualifiers",
		Leaders:   map[string]int{},
		Divisions: DIVISIONS,
		Secret: rand_secret_ignored(),
		Teams: map[string]Team{
			"Roseburg": {
				Name:   "Roseburg",
				Score:  0,
				XCount: 0,
				Members: []Member{
					{Name: "Hunter Blake", Division: randDivision(), Score: 423, XCount: 7},
					{Name: "Spencer King", Division: randDivision(), Score: 392, XCount: 4},
					{Name: "Dominic Ruiz", Division: randDivision(), Score: 468, XCount: 17},
					{Name: "Evan Price", Division: randDivision(), Score: 441, XCount: 10},
					{Name: "Jonah Klein", Division: randDivision(), Score: 459, XCount: 15},
				},
			},
			"Silverton": {
				Name:   "Silverton",
				Score:  0,
				XCount: 0,
				Members: []Member{
					{Name: "Aiden Ross", Division: randDivision(), Score: 436, XCount: 11},
					{Name: "Calvin Brooks", Division: randDivision(), Score: 401, XCount: 6},
					{Name: "Mateo Alvarez", Division: randDivision(), Score: 472, XCount: 19},
					{Name: "Toby Chen", Division: randDivision(), Score: 448, XCount: 12},
					{Name: "Roman Novak", Division: randDivision(), Score: 461, XCount: 16},
					{Name: "Finn Carter", Division: randDivision(), Score: 387, XCount: 2},
				},
			},
			"Stayton": {
				Name:   "Stayton",
				Score:  0,
				XCount: 0,
				Members: []Member{
					{Name: "Levi Turner", Division: randDivision(), Score: 429, XCount: 8},
					{Name: "Brady Owens", Division: randDivision(), Score: 396, XCount: 4},
					{Name: "Rafael Costa", Division: randDivision(), Score: 465, XCount: 17},
					{Name: "Nolan Pierce", Division: randDivision(), Score: 444, XCount: 11},
					{Name: "Emil Johansson", Division: randDivision(), Score: 458, XCount: 14},
				},
			},
		},
	},
}
