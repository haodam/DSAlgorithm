package symbol

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/haodam/DSAlgorithm/crypto/event"
	"log"
)

type Symbol struct {
	pair event.Pair
}

func New(pair event.Pair) actor.Producer {
	return func() actor.Receiver {
		//log.Printf("Creating new Symbol actor for pair: %+v", pair)
		return &Symbol{
			pair: pair,
		}
	}
}

func (s *Symbol) Receive(c *actor.Context) {
	//log.Printf("Symbol actor received message: %T", c.Message())
	switch v := c.Message().(type) {
	case actor.Started:
		log.Printf("Symbol actor started for pair: %v", s.pair)
	case event.Trader:
		log.Printf("Received trader event: %+v", v)
	default:
		log.Printf("Unhandled message type: %T", v)
	}
}
