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
	secret := os.Getenv("SECRET")
	if secret == "" {
		log.Fatal("SECRET is required")
	}

	redisUrl := os.Getenv("REDIS_URL")
	if redisUrl == "" {
		redisUrl = "redis://localhost:6379"
	}
	opts, err := redis.ParseURL(redisUrl)
	if err != nil {
		log.Fatal(err)
	}

	rdb := redis.NewClient(opts)

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

		r.Header.Get("Authorization")
		if r.Header.Get("Authorization") != fmt.Sprintf("Bearer %s", secret) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
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

		szd, err := json.Marshal(task)
		if err != nil {
			http.Error(w, "Could not marshal task", http.StatusInternalServerError)
			return
		}

		resp := rdb.LPush(ctx, "tasks", szd)
		if resp.Err() != nil {
			log.Println(resp.Err())
			http.Error(w, "Could not add task", http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Task added")
	})

	http.HandleFunc("/slurp", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		r.Header.Get("Authorization")
		if r.Header.Get("Authorization") != fmt.Sprintf("Bearer %s", secret) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		resp := rdb.LRange(ctx, "tasks", 0, -1)
		if resp.Err() != nil {
			log.Println(resp.Err())
			http.Error(w, "Could not fetch tasks", http.StatusInternalServerError)
			return
		}

		bytes := []byte{'['}
		for i, task := range resp.Val() {
			if i != 0 {
				bytes = append(bytes, ',')
			}

			bytes = append(bytes, task...)
		}
		bytes = append(bytes, ']')

		w.Header().Set("Content-Type", "application/json")

		_, err := w.Write(bytes)
		if err != nil {
			log.Println(err)
			http.Error(w, "Could not write tasks", http.StatusInternalServerError)
			return
		}

		delResp := rdb.Del(ctx, "tasks")
		if delResp.Err() != nil {
			log.Println(delResp.Err())
			http.Error(w, "Could not trim tasks", http.StatusInternalServerError)
			return
		}
	})

	fmt.Println("Serving on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
