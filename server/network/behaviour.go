package network

import (
	"errors"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

const INFRACTION_RATE_LIMIT_EXCEEDED = 1;
const INFRACTION_INVALID_COOKIE = 2;
const INFRACTION_MALFORMED_DATA = 3;
const INFRACTION_ATTEMPTED_ACCESS = 4;

func bad_behaviour_score(infraction int) int {
	switch infraction {
	case INFRACTION_RATE_LIMIT_EXCEEDED: return 30;
	case INFRACTION_ATTEMPTED_ACCESS: return 30;

	case INFRACTION_INVALID_COOKIE: return 10;

	case INFRACTION_MALFORMED_DATA: return 5;
	}
	return 0;
}

const BEHAVIOUR_THRESHOLD = 30;
const BEHAVIOUR_RATE_LIFESPAN = 1 * time.Hour;
const BEHAVIOUR_TIMEOUT = 12 * time.Hour;

func (limiter *RateLimiter)Handle_behaviour(key string, infraction int) (time.Duration, error, int) {
	score := bad_behaviour_score(infraction);
	limiter.IpRatesMux.Lock();
	defer limiter.IpRatesMux.Unlock();

	if r, exists := limiter.IpRates[key]; exists {

		if r.BehaviourExpire.Before(time.Now()) {
			r.BehaviourExpire = time.Now().Add(BEHAVIOUR_RATE_LIFESPAN);
			r.Score = 0;
		}
		r.Score += score;

		if r.Score > BEHAVIOUR_THRESHOLD {
			return BEHAVIOUR_TIMEOUT, errors.New("EXCEEDED BEHAVIOUR RATE"), http.StatusTooManyRequests;
		}
	} else {
		limiter.IpRates[key] = &Rates{
			BehaviourExpire: time.Now().Add(BEHAVIOUR_RATE_LIFESPAN),
			Score: score,

			Timeout: time.Now(),
			Status: http.StatusOK,

			Limiter: rate.NewLimiter(IP_RATE, IP_BURST),
		};
	}

	return 0, nil, http.StatusOK;
}


