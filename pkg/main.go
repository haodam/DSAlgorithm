package main

import (
	"fmt"
	"sync"
)

//type Store struct {
//	data  map[int]string
//	cache redis.Cacher
//}
//
//func NewStore(c redis.Cacher) *Store {
//	data := map[int]string{
//		1: "Elon Musk is the new owner of Twitter",
//		2: "Foo is not bar and bar is not baz",
//		3: "Musk watch Anthony GG",
//	}
//	return &Store{
//		data:  data,
//		cache: c,
//	}
//}
//
//func (s *Store) Get(key int) (string, error) {
//	val, ok := s.cache.Get(key)
//	if ok {
//		// busting the cache
//		if err := s.cache.Remove(key); err != nil {
//			fmt.Println(err)
//		}
//		fmt.Println("return the value from the cache")
//		return val, nil
//	}
//	val, ok = s.data[key]
//	if !ok {
//		return "", fmt.Errorf("key not found: %d", key)
//	}
//
//	if err := s.cache.Set(key, val); err != nil {
//		return "", err
//	}
//	fmt.Println("returning key from internal storage")
//	return val, nil
//}

func main() {

	//ctx := context.Background()
	//s := handling.Service{}
	//
	//err := s.Signup(ctx, "damanhhaogmail.com", "hao123")
	//if err != nil {
	//	if errors.Is(err, handling.ErrBadRequest) {
	//		fmt.Errorf("invalid email 401")
	//	}
	//}
	//rdb := redis.NewClient(&redis.Options{
	//	Addr:     "localhost:6379",
	//	Password: "", // no password set
	//	DB:       0,  // use default DB
	//})
	//ttl := time.Second * 5
	//s := NewStore(redis.NewRedisCache(rdb, ttl))
	//for i := 0; i < 2; i++ {
	//	val, err := s.Get(3)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	fmt.Println(val)
	//
	//	time.Sleep(5 * time.Second)
	//}

	var mu sync.Mutex
	mu.Lock()
	go func() {
		fmt.Println("Hello")
		mu.Unlock()
	}()
	mu.Lock()

}
