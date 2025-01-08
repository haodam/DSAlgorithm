package binancef

import (
	"errors"
	"fmt"
	"github.com/anthdm/hollywood/actor"
	"github.com/gorilla/websocket"
	"github.com/haodam/DSAlgorithm/crypto/actor/symbol"
	"github.com/haodam/DSAlgorithm/crypto/event"
	"github.com/valyala/fastjson"
	"log"
	"net"
	"strconv"
	"strings"
)

const wsEndpoint = "wss://fstream.binance.com/stream?streams="

var symbols = []string{
	"BTCUSDT",
}

type BinanceF struct {
	ws      *websocket.Conn
	symbols map[string]*actor.PID
	c       *actor.Context
}

func (b *BinanceF) Receive(c *actor.Context) {
	//log.Printf("Received message in BinanceF: %T - %+v", c.Message(), c.Message())

	switch c.Message().(type) {
	case actor.Started:
		//log.Println("Actor started message received")
		b.start(c)
		b.c = c
		//default:
		//	log.Printf("Unhandled message of type: %T - Content: %+v", msg, msg)
	}
}

func New() actor.Producer {
	return func() actor.Receiver {
		//log.Println("Creating new BinanceF actor")
		return &BinanceF{
			symbols: make(map[string]*actor.PID),
		}
	}
}

func (b *BinanceF) start(c *actor.Context) {
	//log.Println("Starting BinanceF actor")
	// Initialize all the symbol actors as children
	for _, sym := range symbols {
		pair := event.Pair{
			Exchange: "binancef",
			Symbol:   strings.ToLower(sym),
		}
		//log.Printf("Spawning child actor for symbol: %s", pair.Symbol)
		pid := c.SpawnChild(symbol.New(pair), fmt.Sprintf("symbol-%s", pair.Symbol))
		if pid == nil {
			//log.Printf("Failed to spawn child actor for symbol: %s", pair.Symbol)
			continue
		}
		b.symbols[pair.Symbol] = pid
	}

	wsEndpoint := createWsEndpoint()
	//log.Printf("Connecting to WebSocket endpoint: %s", wsEndpoint)
	ws, _, err := websocket.DefaultDialer.Dial(wsEndpoint, nil)
	if err != nil {
		log.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	//log.Println("WebSocket connected successfully")
	b.ws = ws
	go b.wsLoop()
}

func (b *BinanceF) wsLoop() {
	//log.Println("WebSocket loop started")
	for {
		_, msg, err := b.ws.ReadMessage()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				//log.Println("WebSocket connection closed")
				break
			}
			//log.Printf("Error reading from WebSocket: %v", err)
			fmt.Println("error reading from ws connection", err)
			continue
		}
		//log.Printf("Received WebSocket message: %s", string(msg))

		parser := fastjson.Parser{}
		v, err := parser.ParseBytes(msg)
		if err != nil {
			//log.Printf("Failed to parse WebSocket message: %v", err)
			fmt.Println("failed to parse msg", err)
			continue
		}

		stream := v.GetStringBytes("stream")
		s, kind := splitStream(string(stream))
		//log.Printf("Parsed stream: %s, kind: %s", s, kind)
		data := v.Get("data")

		if kind == "aggTrade" {
			//log.Printf("Handling aggTrade for symbol: %s", s)
			b.handleAggTrade(s, data)
		} else {
			//log.Printf("Unhandled stream kind: %s", kind)
		}
	}
}

func (b *BinanceF) handleAggTrade(symbol string, data *fastjson.Value) {
	price, err := strconv.ParseFloat(string(data.GetStringBytes("p")), 64)
	if err != nil {
		//log.Printf("Failed to parse price for symbol %s: %v", symbol, err)
		return
	}
	qty, err := strconv.ParseFloat(string(data.GetStringBytes("q")), 64)
	if err != nil {
		//log.Printf("Failed to parse quantity for symbol %s: %v", symbol, err)
		return
	}
	isBuy := data.GetBool("m")
	unix := data.GetInt64("T")

	trade := event.Trader{
		Price: price,
		Qty:   qty,
		IsBuy: isBuy,
		Unix:  unix,
		Pair: event.Pair{
			Exchange: "binancef",
			Symbol:   symbol,
		},
	}
	//log.Printf("Parsed aggTrade: %+v", trade)

	symbolPID, ok := b.symbols[symbol]
	if !ok {
		//log.Printf("Could not find PID for symbol: %s", symbol)
		fmt.Println("could not find PID for symbol", symbol)
		return
	}
	//log.Printf("Sending trade to actor: %v", symbolPID)
	b.c.Send(symbolPID, trade)
}

func createWsEndpoint() string {
	var results []string
	for _, sym := range symbols {
		results = append(results, fmt.Sprintf("%s@aggTrade", strings.ToLower(sym)))
	}
	endpoint := fmt.Sprintf("%s%s", wsEndpoint, strings.Join(results, "/"))
	//log.Printf("Generated WebSocket endpoint: %s", endpoint)
	return endpoint
}

func splitStream(stream string) (string, string) {
	parts := strings.Split(stream, "@")
	if len(parts) != 2 {
		//log.Printf("Invalid stream format: %s", stream)
		return "", ""
	}
	return parts[0], parts[1]
}
