package binancef

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/gorilla/websocket"
	"github.com/valyala/fastjson"
	"reflect"
	"testing"
)

func TestBinanceF_Receive(t *testing.T) {
	type fields struct {
		ws      *websocket.Conn
		symbols map[string]*actor.PID
		c       *actor.Context
	}
	type args struct {
		c *actor.Context
	}
	var tests []struct {
		name   string
		fields fields
		args   args
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BinanceF{
				ws:      tt.fields.ws,
				symbols: tt.fields.symbols,
				c:       tt.fields.c,
			}
			b.Receive(tt.args.c)
		})
	}
}

func TestBinanceF_handleAggTrade(t *testing.T) {
	type fields struct {
		ws      *websocket.Conn
		symbols map[string]*actor.PID
		c       *actor.Context
	}
	type args struct {
		symbol string
		data   *fastjson.Value
	}
	var tests []struct {
		name   string
		fields fields
		args   args
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BinanceF{
				ws:      tt.fields.ws,
				symbols: tt.fields.symbols,
				c:       tt.fields.c,
			}
			b.handleAggTrade(tt.args.symbol, tt.args.data)
		})
	}
}

func TestBinanceF_start(t *testing.T) {
	type fields struct {
		ws      *websocket.Conn
		symbols map[string]*actor.PID
		c       *actor.Context
	}
	type args struct {
		c *actor.Context
	}
	var tests []struct {
		name   string
		fields fields
		args   args
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BinanceF{
				ws:      tt.fields.ws,
				symbols: tt.fields.symbols,
				c:       tt.fields.c,
			}
			b.start(tt.args.c)
		})
	}
}

func TestBinanceF_wsLoop(t *testing.T) {
	type fields struct {
		ws      *websocket.Conn
		symbols map[string]*actor.PID
		c       *actor.Context
	}
	var tests []struct {
		name   string
		fields fields
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BinanceF{
				ws:      tt.fields.ws,
				symbols: tt.fields.symbols,
				c:       tt.fields.c,
			}
			b.wsLoop()
		})
	}
}

func TestNew(t *testing.T) {
	var tests []struct {
		name string
		want actor.Producer
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createWsEndpoint(t *testing.T) {
	var tests []struct {
		name string
		want string
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createWsEndpoint(); got != tt.want {
				t.Errorf("createWsEndpoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_splitStream(t *testing.T) {
	type args struct {
		stream string
	}
	var tests []struct {
		name  string
		args  args
		want  string
		want1 string
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := splitStream(tt.args.stream)
			if got != tt.want {
				t.Errorf("splitStream() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("splitStream() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestBinanceF_Receive1(t *testing.T) {
	type fields struct {
		ws      *websocket.Conn
		symbols map[string]*actor.PID
		c       *actor.Context
	}
	type args struct {
		c *actor.Context
	}
	var tests []struct {
		name   string
		fields fields
		args   args
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BinanceF{
				ws:      tt.fields.ws,
				symbols: tt.fields.symbols,
				c:       tt.fields.c,
			}
			b.Receive(tt.args.c)
		})
	}
}

func TestBinanceF_handleAggTrade1(t *testing.T) {
	type fields struct {
		ws      *websocket.Conn
		symbols map[string]*actor.PID
		c       *actor.Context
	}
	type args struct {
		symbol string
		data   *fastjson.Value
	}
	var tests []struct {
		name   string
		fields fields
		args   args
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BinanceF{
				ws:      tt.fields.ws,
				symbols: tt.fields.symbols,
				c:       tt.fields.c,
			}
			b.handleAggTrade(tt.args.symbol, tt.args.data)
		})
	}
}

func TestBinanceF_start1(t *testing.T) {
	type fields struct {
		ws      *websocket.Conn
		symbols map[string]*actor.PID
		c       *actor.Context
	}
	type args struct {
		c *actor.Context
	}
	var tests []struct {
		name   string
		fields fields
		args   args
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BinanceF{
				ws:      tt.fields.ws,
				symbols: tt.fields.symbols,
				c:       tt.fields.c,
			}
			b.start(tt.args.c)
		})
	}
}

func TestBinanceF_wsLoop1(t *testing.T) {
	type fields struct {
		ws      *websocket.Conn
		symbols map[string]*actor.PID
		c       *actor.Context
	}
	var tests []struct {
		name   string
		fields fields
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BinanceF{
				ws:      tt.fields.ws,
				symbols: tt.fields.symbols,
				c:       tt.fields.c,
			}
			b.wsLoop()
		})
	}
}

func TestNew1(t *testing.T) {
	var tests []struct {
		name string
		want actor.Producer
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createWsEndpoint1(t *testing.T) {
	var tests []struct {
		name string
		want string
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createWsEndpoint(); got != tt.want {
				t.Errorf("createWsEndpoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_splitStream1(t *testing.T) {
	type args struct {
		stream string
	}
	var tests []struct {
		name  string
		args  args
		want  string
		want1 string
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := splitStream(tt.args.stream)
			if got != tt.want {
				t.Errorf("splitStream() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("splitStream() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
