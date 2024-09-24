package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/redis/go-redis/v9"
)

type Task struct {
	Description string   `json:"description"`
	Priority    *string  `json:"priority,omitempty"`
	Project     *string  `json:"project,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

func main() {
	// TODO: https://github.com/redis/go-redis?tab=readme-ov-file#connecting-via-a-redis-url
	redisUrl := os.Getenv("REDIS_URL")
	if redisUrl == "" {
		redisUrl = "localhost:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisUrl,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	ctx := context.Background()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://github.com/bytes-zone/tw-stash", 307)
	})
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) { fmt.Fprintln(w, "OK") })

	http.HandleFunc("/stash", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		var task Task
		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			http.Error(w, "Could not decode task", http.StatusBadRequest)
			return
		}

		if task.Description == "" {
			http.Error(w, "Description is required", http.StatusBadRequest)
			return
		}

		redisErr := rdb.LPush(ctx, "tasks", task)
		if redisErr != nil {
			log.Println(redisErr)
			http.Error(w, "Could not add task", http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Task added")
	})

	fmt.Println("Serving on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
