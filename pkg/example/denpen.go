package main

import "fmt"

type UserAuthenicator interface {
	execute(phone, password string)
}

// ////////// VERSION 1 ////////////
type LoginUseCase struct {
}

func (l LoginUseCase) execute(phone, password string) {
	fmt.Println("execute v1")
}

func NewLoginUseCase() LoginUseCase {
	return LoginUseCase{}
}

//////////// END VERSION 1 ////////////

// ////////// VERSION 2 ////////////
type LoginUseCaseV2 struct {
}

func (l LoginUseCaseV2) execute(phone, password string) {
	fmt.Println("execute v2")
}

func NewLoginUseCaseV2() LoginUseCaseV2 {
	return LoginUseCaseV2{}
}

//////////// END VERSION 2 ////////////

type UserHandler struct {
	UseCase UserAuthenicator
}

func (u UserHandler) login(phone, password string) {
	u.UseCase.execute(phone, password)
}

func main() {
	var userHandler = UserHandler{
		UseCase: NewLoginUseCaseV2(),
	}

	userHandler.login("0973901734", "123456")
}
