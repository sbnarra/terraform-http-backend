package states

import (
    "bytes"
    "io/ioutil"
    "net/http"
    "net/http/httptest"
    "os"
    "path/filepath"
    "testing"
)

func TestHandleStatesGet(t *testing.T) {
    tempDir, err := ioutil.TempDir("", "testdata")
    if err != nil {
        t.Fatalf("Failed to create temp dir: %v", err)
    }
    defer os.RemoveAll(tempDir)

    testFilePath := filepath.Join(tempDir, "statefile.tfstate")
    testData := []byte(`{"version": 1}`)
    err = ioutil.WriteFile(testFilePath, testData, 0644)
    if err != nil {
        t.Fatalf("Failed to write test file: %v", err)
    }

    req := httptest.NewRequest(http.MethodGet, "/statefile.tfstate", nil)
    rr := httptest.NewRecorder()

    HandleStates(rr, req, tempDir)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }
    if rr.Body.String() != string(testData) {
        t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), string(testData))
    }
}

func TestHandleStatesPut(t *testing.T) {
    tempDir, err := ioutil.TempDir("", "testdata")
    if err != nil {
        t.Fatalf("Failed to create temp dir: %v", err)
    }
    defer os.RemoveAll(tempDir)

    testData := []byte(`{"version": 2}`)
    req := httptest.NewRequest(http.MethodPost, "/statefile.tfstate", bytes.NewReader(testData))
    rr := httptest.NewRecorder()

    HandleStates(rr, req, tempDir)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    testFilePath := filepath.Join(tempDir, "statefile.tfstate")
    data, err := ioutil.ReadFile(testFilePath)
    if err != nil {
        t.Fatalf("Failed to read test file: %v", err)
    }
    if string(data) != string(testData) {
        t.Errorf("File content mismatch: got %v want %v", string(data), string(testData))
    }
}

func TestHandleStatesDelete(t *testing.T) {
    tempDir, err := ioutil.TempDir("", "testdata")
    if err != nil {
        t.Fatalf("Failed to create temp dir: %v", err)
    }
    defer os.RemoveAll(tempDir)

    testFilePath := filepath.Join(tempDir, "statefile.tfstate")
    testData := []byte(`{"version": 1}`)
    err = ioutil.WriteFile(testFilePath, testData, 0644)
    if err != nil {
        t.Fatalf("Failed to write test file: %v", err)
    }

    req := httptest.NewRequest(http.MethodDelete, "/statefile.tfstate", nil)
    rr := httptest.NewRecorder()

    HandleStates(rr, req, tempDir)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    if _, err := os.Stat(testFilePath); !os.IsNotExist(err) {
        t.Errorf("File was not deleted")
    }
}

func TestHandleStatesMethodNotAllowed(t *testing.T) {
    tempDir, err := ioutil.TempDir("", "testdata")
    if err != nil {
        t.Fatalf("Failed to create temp dir: %v", err)
    }
    defer os.RemoveAll(tempDir)

    req := httptest.NewRequest(http.MethodPatch, "/statefile.tfstate", nil)
    rr := httptest.NewRecorder()

    HandleStates(rr, req, tempDir)

    if status := rr.Code; status != http.StatusMethodNotAllowed {
        t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
    }
}