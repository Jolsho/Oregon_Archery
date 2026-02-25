package internals

import (
	"errors"
	"net/http"
	"sync"
	"time"
)

type Rate struct {
	Count int
	Expires  time.Time
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


func (state *State)handle_rates(
	rate int, key string, score ...int,
) (time.Duration, error, int) {

	var rate_map *map[string]*Rate;
	var threshold, points int;
	var lifespan time.Duration;
	var ban_time time.Duration;
	var mux *sync.RWMutex;

	switch rate {

	case RATE_SESSION_CREATION: {
		rate_map = &state.SessionCreationRates;
		threshold = SESSION_CREATION_THRESHOLD;
		points = 1;
		lifespan = RATE_LIFESPAN;
		ban_time = RATE_TIMEOUT;
		mux = state.SessionRateMux;
	}
		
	case RATE_DATA: {
		rate_map = &state.DataRates;
		threshold = DATA_THRESHOLD;
		points = 1;
		lifespan = RATE_LIFESPAN;
		ban_time = RATE_TIMEOUT;
		mux = state.DataRateMux;
	}

	case RATE_BEHAVIOUR: {
		rate_map = &state.BehaviourTracker;
		threshold = BEHAVIOUR_THRESHOLD;
		points = score[0];
		lifespan = BEHAVIOUR_RATE_LIFESPAN;
		ban_time = BEHAVIOUR_TIMEOUT;
		mux = state.BehaviourMux;
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
