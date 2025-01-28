package locks

import (
    "bytes"
    "encoding/json"
    "io/ioutil"
    "net/http"
    "net/http/httptest"
    "os"
    "path/filepath"
    "testing"
)

func TestHandleLocksAcquire(t *testing.T) {
    tempDir, err := ioutil.TempDir("", "locktest")
    if err != nil {
        t.Fatalf("Failed to create temp dir: %v", err)
    }
    defer os.RemoveAll(tempDir)

    lockInfo := LockInfo{
        ID:        "test-lock-id",
        Operation: "OperationType",
        Info:      "Locking for test",
        Who:       "tester",
        Version:   "1.0",
        Created:   "2023-10-10T00:00:00Z",
        Path:      "/test/path",
    }
    lockData, err := json.Marshal(lockInfo)
    if err != nil {
        t.Fatalf("Failed to marshal lock info: %v", err)
    }

    req := httptest.NewRequest("LOCK", "/test-lock", bytes.NewReader(lockData))
    rr := httptest.NewRecorder()

    HandleLocks(rr, req, tempDir)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    lockFilePath := filepath.Join(tempDir, "/test-lock")
    if _, err := os.Stat(lockFilePath); os.IsNotExist(err) {
        t.Errorf("Lock file was not created")
    }
}

func TestHandleLocksAcquireConflict(t *testing.T) {
    tempDir, err := ioutil.TempDir("", "locktest")
    if err != nil {
        t.Fatalf("Failed to create temp dir: %v", err)
    }
    defer os.RemoveAll(tempDir)

    existingLockInfo := LockInfo{
        ID:        "existing-lock-id",
        Operation: "OperationType",
        Info:      "Existing lock",
        Who:       "existing-user",
        Version:   "1.0",
        Created:   "2023-10-10T00:00:00Z",
        Path:      "/test/path",
    }
    lockData, err := json.Marshal(existingLockInfo)
    if err != nil {
        t.Fatalf("Failed to marshal existing lock info: %v", err)
    }

    lockFilePath := filepath.Join(tempDir, "/test-lock")
    err = os.MkdirAll(filepath.Dir(lockFilePath), 0755)
    if err != nil {
        t.Fatalf("Failed to create lock directory: %v", err)
    }
    err = ioutil.WriteFile(lockFilePath, lockData, 0644)
    if err != nil {
        t.Fatalf("Failed to write lock file: %v", err)
    }

    newLockInfo := LockInfo{
        ID:        "new-lock-id",
        Operation: "OperationType",
        Info:      "Trying to acquire existing lock",
        Who:       "new-user",
        Version:   "1.0",
        Created:   "2023-10-11T00:00:00Z",
        Path:      "/test/path",
    }
    newLockData, err := json.Marshal(newLockInfo)
    if err != nil {
        t.Fatalf("Failed to marshal new lock info: %v", err)
    }

    req := httptest.NewRequest("LOCK", "/test-lock", bytes.NewReader(newLockData))
    rr := httptest.NewRecorder()

    HandleLocks(rr, req, tempDir)

    if status := rr.Code; status != http.StatusLocked {
        t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusLocked)
    }

    data, err := ioutil.ReadFile(lockFilePath)
    if err != nil {
        t.Fatalf("Failed to read lock file: %v", err)
    }

    var storedLockInfo LockInfo
    err = json.Unmarshal(data, &storedLockInfo)
    if err != nil {
        t.Fatalf("Failed to unmarshal lock file: %v", err)
    }

    if storedLockInfo.ID != existingLockInfo.ID {
        t.Errorf("Lock file was overwritten: got ID %v want %v", storedLockInfo.ID, existingLockInfo.ID)
    }
}

func TestHandleLocksRelease(t *testing.T) {
    tempDir, err := ioutil.TempDir("", "locktest")
    if err != nil {
        t.Fatalf("Failed to create temp dir: %v", err)
    }
    defer os.RemoveAll(tempDir)

    lockInfo := LockInfo{
        ID:        "test-lock-id",
        Operation: "OperationType",
        Info:      "Locking for test",
        Who:       "tester",
        Version:   "1.0",
        Created:   "2023-10-10T00:00:00Z",
        Path:      "/test/path",
    }
    lockData, err := json.Marshal(lockInfo)
    if err != nil {
        t.Fatalf("Failed to marshal lock info: %v", err)
    }

    lockFilePath := filepath.Join(tempDir, "/test-lock")
    err = os.MkdirAll(filepath.Dir(lockFilePath), 0755)
    if err != nil {
        t.Fatalf("Failed to create lock directory: %v", err)
    }
    err = ioutil.WriteFile(lockFilePath, lockData, 0644)
    if err != nil {
        t.Fatalf("Failed to write lock file: %v", err)
    }

    req := httptest.NewRequest("UNLOCK", "/test-lock", bytes.NewReader(lockData))
    rr := httptest.NewRecorder()

    HandleLocks(rr, req, tempDir)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    if _, err := os.Stat(lockFilePath); !os.IsNotExist(err) {
        t.Errorf("Lock file was not deleted")
    }
}

func TestHandleLocksReleaseConflict(t *testing.T) {
    tempDir, err := ioutil.TempDir("", "locktest")
    if err != nil {
        t.Fatalf("Failed to create temp dir: %v", err)
    }
    defer os.RemoveAll(tempDir)

    existingLockInfo := LockInfo{
        ID:        "existing-lock-id",
        Operation: "OperationType",
        Info:      "Locking for test",
        Who:       "tester",
        Version:   "1.0",
        Created:   "2023-10-10T00:00:00Z",
        Path:      "/test/path",
    }
    lockData, err := json.Marshal(existingLockInfo)
    if err != nil {
        t.Fatalf("Failed to marshal existing lock info: %v", err)
    }

    lockFilePath := filepath.Join(tempDir, "/test-lock")
    err = os.MkdirAll(filepath.Dir(lockFilePath), 0755)
    if err != nil {
        t.Fatalf("Failed to create lock directory: %v", err)
    }
    err = ioutil.WriteFile(lockFilePath, lockData, 0644)
    if err != nil {
        t.Fatalf("Failed to write lock file: %v", err)
    }

    unlockInfo := LockInfo{
        ID:        "different-lock-id",
        Operation: "OperationType",
        Info:      "Attempting to unlock with wrong ID",
        Who:       "wrong-user",
        Version:   "1.0",
        Created:   "2023-10-11T00:00:00Z",
        Path:      "/test/path",
    }
    unlockData, err := json.Marshal(unlockInfo)
    if err != nil {
        t.Fatalf("Failed to marshal unlock info: %v", err)
    }

    req := httptest.NewRequest("UNLOCK", "/test-lock", bytes.NewReader(unlockData))
    rr := httptest.NewRecorder()

    HandleLocks(rr, req, tempDir)

    if status := rr.Code; status != http.StatusConflict {
        t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusConflict)
    }

    if _, err := os.Stat(lockFilePath); os.IsNotExist(err) {
        t.Errorf("Lock file was deleted but should not have been")
    }

    if rr.Body.String() != string(lockData) {
        t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), string(lockData))
    }
}

func TestHandleLocksMethodNotAllowed(t *testing.T) {
    tempDir, err := ioutil.TempDir("", "locktest")
    if err != nil {
        t.Fatalf("Failed to create temp dir: %v", err)
    }
    defer os.RemoveAll(tempDir)

    req := httptest.NewRequest(http.MethodGet, "/test-lock", nil)
    rr := httptest.NewRecorder()

    HandleLocks(rr, req, tempDir)

    if status := rr.Code; status != http.StatusMethodNotAllowed {
        t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
    }
}