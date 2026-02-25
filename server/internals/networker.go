package internals

import (
	"crypto/ecdsa"
	"log"

	"github.com/gorilla/websocket"
)

type Networker struct {
	Conns	map[string]*websocket.Conn
	Upgrader *websocket.Upgrader

	RateLimiter *RateLimiter

	PrivKey	*ecdsa.PrivateKey
	PubKey 	*ecdsa.PublicKey

	dones []chan struct{}
	confirm_done chan struct{}
}

func New_Networker() *Networker {
	priv, err := LoadPrivateKey("KEYS.txt");
	if err != nil { 
		panic("CANT LOAD PRIVATE KEY"); 
	}

	const CHILD_ROUTINES_COUNT = 1;
	dones := make([]chan struct{}, CHILD_ROUTINES_COUNT);
	confirm_done := make(chan struct{}, CHILD_ROUTINES_COUNT);

	// LIMITER IS A CHILD ROUTINE OF STATE
	limiter := New_Rate_Limiter();
	dones[0] = make(chan struct{}, 1);
	go limiter.start_cleaner(dones[0], confirm_done)

	return &Networker{
		RateLimiter: limiter,

		Upgrader: New_Upgrader(),
		Conns: make(map[string]*websocket.Conn, 32),

		PrivKey: priv, PubKey: &priv.PublicKey,

		dones: dones,
		confirm_done: confirm_done,
	};
}

func (net *Networker) Shutdown() {
	err := SavePrivateKey("KEYS.txt", net.PrivKey)
	if err != nil { 
		log.Fatal("STATE_PERSIST::SAVE_KEY::",err.Error()) 
	}

	i := 0;
	for _, done := range net.dones {
		done <- struct{}{};
		i++;
	}

	for i > 0 {
		<-net.confirm_done
		i--
	}
}
