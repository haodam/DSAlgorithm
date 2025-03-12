package vm

import "fmt"

type VendingMachine struct {
}

func NewVendingMachine() *VendingMachine {
	return &VendingMachine{}
}

func (vm *VendingMachine) Execute(money int64, brand string) string {
	return fmt.Sprintf("super hot %s", brand)
}
