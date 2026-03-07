package network

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
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

	Logger *Logger

	PrivKey	*ecdsa.PrivateKey
	PubKey 	*ecdsa.PublicKey

	workers *utils.WorkGroup
}

func New_Networker(key_path, logger_path string) *Networker {
	priv, err := utils.LoadPrivateKey(key_path);
	if err != nil { 
		panic("CANT LOAD PRIVATE KEY"); 
	}

	workers := utils.NewWorkGroup()

	// LOGGER
	logger := New_logger(logger_path);
	go logger.start_writer(workers);
	logger.Log(DEBUG_LEVEL, "LOGGER STARTED");

	return &Networker{
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
