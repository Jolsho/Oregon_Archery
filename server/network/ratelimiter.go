package network

import (
	"errors"
	"fmt"
	"net/http"
	"server/utils"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

const IP_RATE = 5;
const IP_BURST = 10;
const LAST_USED_THRESHOLD = 10 * time.Minute;
const RATE_EXCEEDED_TIMEOUT = 30 * time.Second;

type Rates struct {
	Score 	int
	BehaviourExpire time.Time

	Timeout time.Time
	Status int

	Limiter *rate.Limiter
	LastSeen time.Time
}

type RateLimiter struct {
	IpRates 		map[string]*Rates 	 // IP -> Rates
	IpRatesMux		*sync.RWMutex
}

func New_Rate_Limiter() *RateLimiter {
	return &RateLimiter{
		IpRatesMux: &sync.RWMutex{},
		IpRates: make(map[string]*Rates, 32),
	};
}

func (limiter *RateLimiter) Handle_timeout(ip string, timeout time.Duration, status int) {
	limiter.IpRatesMux.Lock();
	defer limiter.IpRatesMux.Unlock();

	if r, exists := limiter.IpRates[ip]; exists {
		r.Timeout = time.Now().Add(timeout);
		r.Status = status;
	} else {
		limiter.IpRates[ip] = &Rates{
			Score: 0,
			BehaviourExpire: time.Now().Add(timeout),

			Timeout: time.Now(),
			Status: status,

			Limiter: rate.NewLimiter(IP_RATE, IP_BURST),
		}
	}
}

func (limiter *RateLimiter)Is_timed_out(ip string) (expires string, status int) {
	limiter.IpRatesMux.Lock();
	defer limiter.IpRatesMux.Unlock();
	if rate, exists := limiter.IpRates[ip]; exists {
		rate.LastSeen = time.Now();
		if rate.Status != http.StatusOK {
			if rate.Timeout.After(rate.LastSeen) {
				error_str := fmt.Sprintf("%s", 
					rate.Timeout.Local().Format("01/02/2006 03:04:05 PM"),
				);
				return error_str, rate.Status;
			}
			rate.Status = http.StatusOK;

		}
	}
	return "", http.StatusOK;
}

func (limiter *RateLimiter) Handle_rates(ip string) (time.Duration, error, int) {

	limiter.IpRatesMux.Lock();
	if r ,exists := limiter.IpRates[ip]; exists {
		if !r.Limiter.Allow() {
			limiter.IpRatesMux.Unlock();
			return RATE_EXCEEDED_TIMEOUT, errors.New("RATE LIMIT EXCEEDED"), http.StatusTooManyRequests
		}
	} else {
		limiter.IpRates[ip] = &Rates{
			Score: 0,
			BehaviourExpire: time.Now().Add(BEHAVIOUR_RATE_LIFESPAN),

			Timeout: time.Now(),
			Status: http.StatusOK,

			Limiter: rate.NewLimiter(IP_RATE, IP_BURST),
		}
	}
	limiter.IpRatesMux.Unlock();

	return 0, nil, http.StatusOK;
}


/////////////////////////////////////////////
/////////// CLEANING LOGIC /////////////////
///////////////////////////////////////////


const CLEANING_INTERVAL = 24 * time.Hour;
const RETRY_CLEAN_TIMEOUT = 1 * time.Hour;


type CleaningTask struct {
	Mapping *map[string]*Rates
	Mux		*sync.RWMutex
	When	time.Time
}

func (task *CleaningTask) Get_when() time.Time { 
	return task.When; 
}


func (task *CleaningTask) Clean() {
	task.Mux.Lock()
	for key, rate := range *task.Mapping {
		if rate.Status == http.StatusOK && 
		time.Since(rate.LastSeen) > LAST_USED_THRESHOLD {
			delete(*task.Mapping, key)
		}
	}
	task.Mux.Unlock()
	task.When = time.Now().Add(CLEANING_INTERVAL)
}

func (limiter *RateLimiter) start_cleaner(
	group *utils.WorkGroup,
) {
	tasks, taskMux, needsCleaning := utils.NewCleaningTasks(
		&CleaningTask{
			Mapping: &limiter.IpRates, 
			Mux: limiter.IpRatesMux, 
			When: time.Now().Add(CLEANING_INTERVAL),
		},
	);
	// SCHEDULER
	go utils.Task_Scheduler(&tasks, taskMux, CLEANING_INTERVAL, needsCleaning, group)
	utils.RunCleaner(needsCleaning, taskMux, tasks, group);
}
