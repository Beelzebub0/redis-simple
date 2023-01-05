package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer client.Close()

	p, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(p)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		switch r.Method {
		case "GET":
			key := r.URL.Query().Get("key")
			val, err := client.Get(ctx, key).Result()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			fmt.Fprint(w, val)
		case "POST":
			key := r.URL.Query().Get("key")
			val := r.URL.Query().Get("val")
			err := client.Set(ctx, key, val, 0).Err()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			fmt.Fprint(w, "OK")
		}
	})

	http.ListenAndServe(":8080", nil)
}
