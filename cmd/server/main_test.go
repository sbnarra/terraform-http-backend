package main

import (
    "bytes"
    "encoding/base64"
    "encoding/json"
    "io/ioutil"
    "net/http"
    "os"
    "path/filepath"
    "testing"
    "time"
)

// TestIntegrationServer tests the full server integration
func TestIntegrationServer(t *testing.T) {
    // Set up environment variables
    os.Setenv("AUTH_USERNAME", "testuser")
    os.Setenv("AUTH_PASSWORD", "testpass")
    os.Setenv("DATA_DIR", "./testdata")
    os.Setenv("PORT", "8081") // Use a non-standard port for testing

    defer func() {
        // Clean up environment variables and test data after test
        os.Unsetenv("AUTH_USERNAME")
        os.Unsetenv("AUTH_PASSWORD")
        os.Unsetenv("DATA_DIR")
        os.Unsetenv("PORT")
        os.RemoveAll("./testdata")
    }()

    // Start the server in a separate goroutine
    go func() {
			main()
    }()

    // Wait briefly to ensure the server has started
    time.Sleep(500 * time.Millisecond)

    // Base URL for requests
    baseURL := "http://localhost:8081"

    // Prepare Authorization header
    authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte("testuser:testpass"))

    // Test state storage endpoints
    stateFilePath := "/states/test.tfstate"
    stateData := []byte(`{"version": 1}`)
    testStateEndpoints(t, baseURL, authHeader, stateFilePath, stateData)

    // Test lock endpoints
    lockFilePath := "/locks/test.tfstate"
    lockInfo := LockInfo{
        ID:        "test-lock-id",
        Operation: "OperationType",
        Info:      "Integration test lock",
        Who:       "tester",
        Version:   "0.12",
        Created:   time.Now().Format(time.RFC3339),
        Path:      stateFilePath,
    }
    testLockEndpoints(t, baseURL, authHeader, lockFilePath, lockInfo)
}

// testStateEndpoints performs integration tests on the state endpoints
func testStateEndpoints(t *testing.T, baseURL, authHeader, stateFilePath string, stateData []byte) {
    postReq, err := http.NewRequest(http.MethodPost, baseURL+stateFilePath, bytes.NewReader(stateData))
    if err != nil {
        t.Fatalf("Failed to create POST request: %v", err)
    }
    postReq.Header.Set("Authorization", authHeader)

    postResp, err := http.DefaultClient.Do(postReq)
    if err != nil {
        t.Fatalf("POST request failed: %v", err)
    }
    defer postResp.Body.Close()

    if postResp.StatusCode != http.StatusOK {
        t.Errorf("POST /states/ returned status %v; want %v", postResp.StatusCode, http.StatusOK)
    }

    // GET state data
    getReq, err := http.NewRequest(http.MethodGet, baseURL+stateFilePath, nil)
    if err != nil {
        t.Fatalf("Failed to create GET request: %v", err)
    }
    getReq.Header.Set("Authorization", authHeader)

    getResp, err := http.DefaultClient.Do(getReq)
    if err != nil {
        t.Fatalf("GET request failed: %v", err)
    }
    defer getResp.Body.Close()

    if getResp.StatusCode != http.StatusOK {
        t.Errorf("GET /states/ returned status %v; want %v", getResp.StatusCode, http.StatusOK)
    }

    responseData, err := ioutil.ReadAll(getResp.Body)
    if err != nil {
        t.Fatalf("Failed to read GET response body: %v", err)
    }

    if !bytes.Equal(responseData, stateData) {
        t.Errorf("GET /states/ returned data %v; want %v", string(responseData), string(stateData))
    }

    deleteReq, err := http.NewRequest(http.MethodDelete, baseURL+stateFilePath, nil)
    if err != nil {
        t.Fatalf("Failed to create DELETE request: %v", err)
    }
    deleteReq.Header.Set("Authorization", authHeader)

    deleteResp, err := http.DefaultClient.Do(deleteReq)
    if err != nil {
        t.Fatalf("DELETE request failed: %v", err)
    }
    defer deleteResp.Body.Close()

    if deleteResp.StatusCode != http.StatusOK {
        t.Errorf("DELETE /states/ returned status %v; want %v", deleteResp.StatusCode, http.StatusOK)
    }

    // Verify the state file has been deleted
    _, err = os.Stat(filepath.Join("./testdata/states", "test.tfstate"))
    if !os.IsNotExist(err) {
        t.Errorf("State file was not deleted")
    }
}

// testLockEndpoints performs integration tests on the lock endpoints
func testLockEndpoints(t *testing.T, baseURL, authHeader, lockFilePath string, lockInfo LockInfo) {
    // Acquire lock
    lockData, err := json.Marshal(lockInfo)
    if err != nil {
        t.Fatalf("Failed to marshal lock info: %v", err)
    }

    postReq, err := http.NewRequest("LOCK", baseURL+lockFilePath, bytes.NewReader(lockData))
    if err != nil {
        t.Fatalf("Failed to create LOCK request: %v", err)
    }
		
    postReq.Method = "LOCK"
    postReq.Header.Set("Authorization", authHeader)
    postReq.Header.Set("Content-Type", "application/json")

    postResp, err := http.DefaultClient.Do(postReq)
    if err != nil {
        t.Fatalf("LOCK request failed: %v", err)
    }
    defer postResp.Body.Close()

    if postResp.StatusCode != http.StatusOK {
        t.Errorf("LOCK /locks/ returned status %v; want %v", postResp.StatusCode, http.StatusOK)
    }

    // Attempt to acquire lock again, expect conflict
    conflictReq, err := http.NewRequest("LOCK", baseURL+lockFilePath, bytes.NewReader(lockData))
    if err != nil {
        t.Fatalf("Failed to create second LOCK request: %v", err)
    }
    conflictReq.Header.Set("Authorization", authHeader)
    conflictReq.Header.Set("Content-Type", "application/json")

    conflictResp, err := http.DefaultClient.Do(conflictReq)
    if err != nil {
        t.Fatalf("Second LOCK request failed: %v", err)
    }
    defer conflictResp.Body.Close()

    if conflictResp.StatusCode != http.StatusLocked {
        t.Errorf("Second LOCK /locks/ returned status %v; want %v", conflictResp.StatusCode, http.StatusLocked)
    }

    // Release lock
    deleteReq, err := http.NewRequest("UNLOCK", baseURL+lockFilePath, bytes.NewReader(lockData))
    if err != nil {
        t.Fatalf("Failed to create UNLOCK request: %v", err)
    }
    deleteReq.Header.Set("Authorization", authHeader)
    deleteReq.Header.Set("Content-Type", "application/json")

    deleteResp, err := http.DefaultClient.Do(deleteReq)
    if err != nil {
        t.Fatalf("UNLOCK request failed: %v", err)
    }
    defer deleteResp.Body.Close()

    if deleteResp.StatusCode != http.StatusOK {
        t.Errorf("UNLOCK /locks/ returned status %v; want %v", deleteResp.StatusCode, http.StatusOK)
    }

    // Verify the lock file has been deleted
    _, err = os.Stat(filepath.Join("./testdata/locks", "test.tfstate"))
    if !os.IsNotExist(err) {
        t.Errorf("Lock file was not deleted")
    }
}

// LockInfo represents the structure of a lock file (duplicate for test)
type LockInfo struct {
    ID        string `json:"ID"`
    Operation string `json:"Operation"`
    Info      string `json:"Info"`
    Who       string `json:"Who"`
    Version   string `json:"Version"`
    Created   string `json:"Created"`
    Path      string `json:"Path"`
}