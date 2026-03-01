package network

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"log"
	"server/utils"
	"sync"

	"github.com/gorilla/websocket"
)

type WsMsg struct {
	Msg 	string 			`json:"msg"`
	Payload json.RawMessage `json:"payload"`
}

type WSConn struct {
	Nonce string
	Ip	  string
	Conn *websocket.Conn
	Ctx context.Context
	Cancel context.CancelFunc
	Out chan WsMsg
}

type Networker struct {
	Conns	map[string]*WSConn
	ConnsMux *sync.RWMutex
	Upgrader *websocket.Upgrader

	RateLimiter *RateLimiter
	Logger *Logger

	PrivKey	*ecdsa.PrivateKey
	PubKey 	*ecdsa.PublicKey

	workers *utils.WorkGroup
}

func New_Networker() *Networker {
	priv, err := utils.LoadPrivateKey("KEYS.txt");
	if err != nil { 
		panic("CANT LOAD PRIVATE KEY"); 
	}

	workers := utils.NewWorkGroup()

	// LOGGER
	logger := New_logger();
	go logger.start_writer(workers);
	logger.Log(DEBUG_LEVEL, "LOGGER STARTED");

	// RATE LIMITER
	limiter := New_Rate_Limiter();
	go limiter.start_cleaner(workers)
	logger.Log(DEBUG_LEVEL, "LIMITER STARTED CLEANING");


	return &Networker{
		RateLimiter: limiter,
		Logger: logger,

		Upgrader: New_Upgrader(),
		Conns: make(map[string]*WSConn, 32),
		ConnsMux: &sync.RWMutex{},

		PrivKey: priv, 
		PubKey: &priv.PublicKey,

		workers: workers,
	};
}

func (net *Networker) Shutdown() {
	err := utils.SavePrivateKey("/var/ohsal/KEYS.txt", net.PrivKey)
	if err != nil { 
		log.Fatal("STATE_PERSIST::SAVE_KEY::",err.Error()) 
	}

	net.ConnsMux.Lock();
	for _, conn := range net.Conns { conn.Cancel(); }
	net.ConnsMux.Unlock();

	net.workers.Cancel();
	net.workers.WG.Wait();
}


func (net *Networker) Bad_Behaviour(infraction int, ip string) bool {
	timeout, err, status := net.RateLimiter.Handle_behaviour(ip, infraction);
	if err != nil { 

		log := fmt.Sprintf("TIMEDOUT %s because %s", ip, err.Error());
		net.Logger.Log(INFO_LEVEL, log);

		net.RateLimiter.Handle_timeout(ip, timeout, status)

		return false; 
	};
	return true;
}
