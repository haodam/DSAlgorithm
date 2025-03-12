package main

import (
	"fmt"
	"github.com/haodam/DSAlgorithm/pkg/abstraction/vm"
)

type vendingMachine interface {
	Execute(money int64, brand string) string
}

type Application struct {
	vm vendingMachine
}

func (ap Application) Run() {
	myDrink := ap.vm.Execute(100, "CocaCola")
	fmt.Println(myDrink)
}

func newApplication(vm vendingMachine) *Application {
	return &Application{vm}
}

func main() {
	vendingMachines := vm.NewVendingMachine()
	app := newApplication(vendingMachines)
	app.Run()
}
