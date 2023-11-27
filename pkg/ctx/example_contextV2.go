package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

func main() {
	start := time.Now()
	ctx := context.Background()
	userID := 10
	val, err := fetchUserDataV2(ctx, userID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("result: ", val)
	fmt.Println("took", time.Since(start))

}

type Response struct {
	value int
	err   error
}

func fetchUserDataV2(ctx context.Context, userID int) (int, error) {

	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*200)
	defer cancel()

	respch := make(chan Response)

	go func() {
		val, err := fetchThirdPartyStuffWhichCanBeSlowV2()
		respch <- Response{
			value: val,
			err:   err,
		}
	}()
	for {
		select {
		case <-ctx.Done():
			return 0, fmt.Errorf("fetching data from third party took to long")
		case resp := <-respch:
			return resp.value, resp.err
		}
	}
}

func fetchThirdPartyStuffWhichCanBeSlowV2() (int, error) {
	time.Sleep(time.Millisecond * 150)
	return 999, nil

}

/*
	output
	result:  999
	took 163.2167ms

	Thời gian ở ví dụ 1 là took 511.4386ms
*/
