package main

import (
	"github.com/anthdm/hollywood/actor"
	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/haodam/DSAlgorithm/crypto/actor/consumer/binancef"
	"log"
)

type Panel struct {
	position rl.Vector2
	width    float32
	height   float32
	title    string
}

func NewPanel(pos rl.Vector2, width float32, height float32) *Panel {
	return &Panel{
		position: pos,
		width:    width,
		height:   height,
	}
}

func (p *Panel) update() {}

func (p *Panel) getDrawPos(x, y float32) rl.Vector2 {
	return rl.NewVector2(p.position.X+x, p.position.Y+24+y+p.height)
}

func (p *Panel) render() {
	gui.Panel(rl.NewRectangle(p.position.X, p.position.Y, p.width, p.height), p.title)
	pos := p.getDrawPos(10, 10)
	rl.DrawCircle(int32(pos.X), int32(pos.Y), 2, rl.Red)
}

func main() {

	//log.Println("Initializing actor engine...")
	e, err := actor.NewEngine(actor.NewEngineConfig())
	if err != nil {
		log.Fatal("Failed to initialize actor engine:", err)
	}

	//log.Println("Spawning BinanceF actor...")
	e.Spawn(binancef.New(), "binancef")
	//log.Println("Actor system is running. Waiting for messages...")
	select {}
	return

	//rl.InitWindow(1200, 600, "Market Crypto")
	//defer rl.CloseWindow()
	//
	//rl.SetTargetFPS(60)
	//
	//gui.SetStyle(0, gui.BACKGROUND_COLOR, 0x2d2d2dff)
	//
	//panel := NewPanel(rl.NewVector2(100, 100), 800, 600)
	//panel.title = "Heatmap - Binance Future"
	//
	//for !rl.WindowShouldClose() {
	//	panel.update()
	//
	//	rl.BeginDrawing()
	//	rl.ClearBackground(rl.Black)
	//
	//	panel.render()
	//
	//	rl.EndDrawing()
	//}
}
