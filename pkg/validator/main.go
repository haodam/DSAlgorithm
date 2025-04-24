package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"log"
)

type people struct {
	Name  string `validate:"required,min=3,max=100"`
	Email string `validate:"required,email"`
	Age   int    `validate:"gte=21"`
}

func main() {
	u := people{
		Name:  "Jack",
		Email: "jack@gmail.com",
		Age:   17,
	}
	validate := validator.New()
	err := validate.Struct(u)

	if err != nil {
		log.Println("Validation Failed")
		for _, e := range err.(validator.ValidationErrors) {
			fmt.Printf("Field %s, Error: %s\n", e.Field(), e.Tag())
		}
	} else {
		log.Println("Validation Success")
	}
}
