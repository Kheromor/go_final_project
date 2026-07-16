package server

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"go_final_project/pkg/api"
)

const defaultPort = "7540"
const defaultWebDir = "./web"

func Start() {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = defaultPort
	}

	webDir := getWebDir()
	api.RegisterHandlers()
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	addr := ":" + port
	log.Printf("starting web server on %s, serving %s", addr, webDir)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

func getWebDir() string {
	if dir := os.Getenv("TODO_WEB_DIR"); dir != "" {
		return dir
	}

	wd, err := os.Getwd()
	if err == nil {
		candidates := []string{
			filepath.Join(wd, defaultWebDir),
			filepath.Join(wd, "..", "web"),
			filepath.Join(wd, "..", "..", "web"),
		}
		for _, candidate := range candidates {
			if info, err := os.Stat(candidate); err == nil && info.IsDir() {
				return candidate
			}
		}
	}

	return defaultWebDir
}
