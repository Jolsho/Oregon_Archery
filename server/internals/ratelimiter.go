package internals

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Rate struct {
	Count int
	Expires  time.Time
}

type RateLimiter struct {
	TimeOut					map[string]*Rate 	// IP -> timeout + rate
	TimeoutMux				*sync.RWMutex

	BehaviourTracker 		map[string]*Rate 	 // IP -> bad_behaviour_rate
	BehaviourMux			*sync.RWMutex

	SessionCreationRates	map[string]*Rate	 // IP -> rate
	SessionRateMux			*sync.RWMutex

	DataRates 				map[string]*Rate 	 // Session Cookie -> rate
	DataRateMux				*sync.RWMutex
}

func New_Rate_Limiter() *RateLimiter {
	return &RateLimiter{
		TimeoutMux: &sync.RWMutex{},
		TimeOut: make(map[string]*Rate, 32),

		BehaviourMux: &sync.RWMutex{},
		BehaviourTracker: make(map[string]*Rate, 32),

		SessionRateMux: &sync.RWMutex{},
		SessionCreationRates: make(map[string]*Rate, 32),

		DataRates: make(map[string]*Rate, 32),
		DataRateMux: &sync.RWMutex{},
	};
}

func (limiter *RateLimiter) handle_timeout(ip string, timeout time.Duration) {
	if rate, exists := limiter.TimeOut[ip]; exists {
		rate.Expires = time.Now().Add(timeout);
		// TODO -- handle count here for repeated timeouts?
		// WOULD NEED TO CHANGE CLEANING LOGIC
		// NOT NECCESSARY, BAD ACTORS WONT USE SAME IP ANYWAY
		// COULD WRITE SOMETHING FOR THIS ON THE SYSTEM
		// BY SCANNING LOGS OR SOMETHING AND NOTIFYING A FIREWALL
	} else {
		limiter.TimeOut[ip] = &Rate{
			Expires: time.Now().Add(timeout),
			Count: 1,
		}
	}
}

func (limiter *RateLimiter)is_timed_out(ip string) (expires string, is_timedout bool) {
	limiter.TimeoutMux.Lock();
	defer limiter.TimeoutMux.Unlock();
	if rate, exists := limiter.TimeOut[ip]; exists {
		if rate.Expires.Before(time.Now()) {
			delete(limiter.TimeOut, ip);
		} else {
			error_str := fmt.Sprintf("%s", 
				rate.Expires.UTC().Format(time.RFC1123),
			);
			return error_str, true;
		}
	}
	return "", false;
}

const INFRACTION_RATE_LIMIT_EXCEEDED = 1;
const INFRACTION_INVALID_COOKIE = 2;

func bad_behaviour_score(infraction int) int {
	switch infraction {
	case INFRACTION_RATE_LIMIT_EXCEEDED: return 30;
	case INFRACTION_INVALID_COOKIE: return 10;
	}
	return 0;
}


const RATE_SESSION_CREATION = 1;
const RATE_DATA = 2;
const RATE_BEHAVIOUR = 3;


const SESSION_CREATION_THRESHOLD = 10;
const DATA_THRESHOLD = 25;
const RATE_LIFESPAN = 60 * time.Second;
const RATE_TIMEOUT = 120 * time.Second;


const BEHAVIOUR_THRESHOLD = 30;
const BEHAVIOUR_RATE_LIFESPAN = 1 * time.Hour;
const BEHAVIOUR_TIMEOUT = 12 * time.Hour;


func (limiter *RateLimiter)handle_rates(
	rate int, key string, score ...int,
) (time.Duration, error, int) {

	var rate_map *map[string]*Rate;
	var threshold, points int;
	var lifespan time.Duration;
	var ban_time time.Duration;
	var mux *sync.RWMutex;

	switch rate {

	case RATE_SESSION_CREATION: {
		rate_map = &limiter.SessionCreationRates;
		threshold = SESSION_CREATION_THRESHOLD;
		points = 1;
		lifespan = RATE_LIFESPAN;
		ban_time = RATE_TIMEOUT;
		mux = limiter.SessionRateMux;
	}
		
	case RATE_DATA: {
		rate_map = &limiter.DataRates;
		threshold = DATA_THRESHOLD;
		points = 1;
		lifespan = RATE_LIFESPAN;
		ban_time = RATE_TIMEOUT;
		mux = limiter.DataRateMux;
	}

	case RATE_BEHAVIOUR: {
		rate_map = &limiter.BehaviourTracker;
		threshold = BEHAVIOUR_THRESHOLD;
		points = score[0];
		lifespan = BEHAVIOUR_RATE_LIFESPAN;
		ban_time = BEHAVIOUR_TIMEOUT;
		mux = limiter.BehaviourMux;
	}

	default:
		return 0, errors.New("INVALID RATE LIMITER"), http.StatusInternalServerError;
	}

	mux.Lock();
	defer mux.Unlock();

	if rate, exists := (*rate_map)[key]; exists {

		if rate.Expires.Before(time.Now()) {
			rate.Expires = time.Now().Add(lifespan);
			rate.Count = 0;
		}
		rate.Count += points;

		if rate.Count > threshold {
			return ban_time, errors.New("EXCEEDED RATE"), http.StatusTooManyRequests;
		}
	} else {
		(*rate_map)[key] = &Rate{
			Expires: time.Now().Add(RATE_LIFESPAN),
			Count: 1,
		};
	}

	return 0, nil, http.StatusOK;
}



/////////////////////////////////////////////
/////////// CLEANING LOGIC /////////////////
///////////////////////////////////////////


const CLEANING_INTERVAL = 24 * time.Hour;
const RETRY_CLEAN_TIMEOUT = 1 * time.Hour;

type CleaningTask struct {
	Mapping *map[string]*Rate
	Mux		*sync.RWMutex
	When	time.Time
}

func (task *CleaningTask) get_when() time.Time { 
	return task.When; 
}

func (task *CleaningTask) clean() {
	now := time.Now()
	if task.Mux.TryLock() {
		for key, rate := range *task.Mapping {
			if rate.Expires.Before(now) {
				delete(*task.Mapping, key)
			}
		}
		task.Mux.Unlock()
		task.When = time.Now().Add(CLEANING_INTERVAL)
	} else {
		task.When = time.Now().Add(RETRY_CLEAN_TIMEOUT)
	}
}

func (limiter *RateLimiter) start_cleaner(done <-chan struct{}, confirm chan struct{}) {
	initial := time.Now().Add(CLEANING_INTERVAL)

	tasks := []*CleaningTask{
		{Mapping: &limiter.TimeOut, Mux: limiter.TimeoutMux, When: initial},
		{Mapping: &limiter.BehaviourTracker, Mux: limiter.BehaviourMux, When: initial},
		{Mapping: &limiter.DataRates, Mux: limiter.DataRateMux, When: initial},
	}

	taskMux := &sync.Mutex{}
	needsCleaning := make(chan struct{}, len(tasks))

	// SCHEDULER
	go Task_Scheduler(&tasks, taskMux, CLEANING_INTERVAL, needsCleaning, done)

	for {
		select {
		case <-needsCleaning:
			taskMux.Lock()

			// pop earliest
			task := tasks[0]
			copy(tasks, tasks[1:])
			tasks = tasks[:len(tasks)-1]

			// update schedule
			task.clean()

			// reinsert sorted by When
			i := 0
			for i < len(tasks) && tasks[i].When.Before(task.When) {
				i++
			}

			tasks = append(tasks, &CleaningTask{}) // grow
			copy(tasks[i+1:], tasks[i:])
			tasks[i] = task

			taskMux.Unlock()

		case <-done:
			confirm <- struct{}{};
			return;
		}
	}
}


