package tests

import (
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"go_final_project/pkg/db"
	"go_final_project/pkg/server"
)

func TestMain(m *testing.M) {
	if os.Getenv("TODO_PORT") == "" {
		os.Setenv("TODO_PORT", "7540")
	}

	if os.Getenv("TODO_DBFILE") == "" {
		os.Setenv("TODO_DBFILE", "../scheduler.db")
	}

	if err := db.Init(os.Getenv("TODO_DBFILE")); err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	go server.Start()

	url := "http://localhost:" + os.Getenv("TODO_PORT")
	deadline := time.Now().Add(5 * time.Second)
	for {
		if time.Now().After(deadline) {
			log.Fatalf("server did not start within timeout")
		}
		resp, err := http.Get(url)
		if err == nil {
			resp.Body.Close()
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	os.Exit(m.Run())
}
