package main

import (
    "log"
    "net/http"
    "os"

    "terraform-http-backend/internal/auth"
    "terraform-http-backend/internal/config"
    "terraform-http-backend/internal/locks"
    "terraform-http-backend/internal/states"
)

func main() {
    // Initialize authentication
    auth.Initialize()

    // Get data directory from environment or use default
    dataDir := config.GetEnv("DATA_DIR", "./data")
    log.Printf("Storing data in '%s'", dataDir)
    createDataDir(dataDir)

    // Set up HTTP handlers with authentication
    http.HandleFunc("/states/", auth.WithAuth(func(w http.ResponseWriter, r *http.Request) {
        states.HandleStates(w, r, dataDir)
    }))
    http.HandleFunc("/locks/", auth.WithAuth(func(w http.ResponseWriter, r *http.Request) {
        locks.HandleLocks(w, r, dataDir)
    }))

    // Start the server
    startServer()
}

func createDataDir(dataDir string) {
    if err := os.MkdirAll(dataDir, 0755); err != nil {
        log.Fatalf("Failed to create storage root directory: %v", err)
    }
}

func startServer() {
    port := config.GetEnv("PORT", "9944")
    log.Printf("Starting server on port %s", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}