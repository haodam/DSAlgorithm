package cmd

import (
	"fmt"
	"github.com/anthdm/hollywood/actor"
)

type HelloActor struct{}

func (h *HelloActor) Receiver(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case string:
		fmt.Printf("Received message: %s\n", msg)
	case int:
		fmt.Printf("Received number: %d\n", msg)
	}
}

func main() {
	props := actor.New
}
