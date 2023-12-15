package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/haodam/DSAlgorithm/pkg/handling"
)

func main() {

	ctx := context.Background()
	s := handling.Service{}

	err := s.Signup(ctx, "damanhhaogmail.com", "hao123")
	if err != nil {
		if errors.Is(err, handling.ErrBadRequest) {
			fmt.Errorf("invalid email 401")
		}
	}

}
