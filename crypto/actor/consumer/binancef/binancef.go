package binancef

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/gorilla/websocket"
	"log"
)

const wsEndpoint = "wss://stream.binance.com/stream?streams="

var symbols = []string{
	"ETHUSDT",
	"BTCUSDT",
}

type BinanceF struct {
	ws      *websocket.Conn
	symbols map[string]*actor.PID
	c       *actor.Context
}

func (b *BinanceF) Receive(c *actor.Context) {
	switch c.Message().(type) {
	case *actor.Started:
		b.start(c)
		b.c = c
	}
}

func New() actor.Producer {
	return func() actor.Receiver {
		return &BinanceF{
			symbols: make(map[string]*actor.PID),
		}
	}
}

func (b *BinanceF) start(c *actor.Context) {
	// Initialize all the symbol actors as children
	//for _, sym := range symbols {
	//	pair :=
	//}

	ws, _, err := websocket.DefaultDialer.Dial(wsEndpoint, nil)
	if err != nil {
		log.Fatal(err)
	}
	b.ws = ws

}
