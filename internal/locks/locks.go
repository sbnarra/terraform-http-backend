package locks

import (
    "encoding/json"
    "log"
    "net/http"
    "os"

    "terraform-http-backend/internal/utils"
)

// LockInfo represents the structure of a lock file
type LockInfo struct {
    ID        string `json:"ID"`
    Operation string `json:"Operation"`
    Info      string `json:"Info"`
    Who       string `json:"Who"`
    Version   string `json:"Version"`
    Created   string `json:"Created"`
    Path      string `json:"Path"`
}

// HandleLocks processes lock-related HTTP requests
func HandleLocks(w http.ResponseWriter, r *http.Request, dataDir string) {
    lockfilePath, lockDir := utils.GetFilePaths(r.URL.Path, dataDir)
    switch r.Method {
    case "LOCK", http.MethodPost, http.MethodPut:
        acquireLock(w, r, lockfilePath, lockDir)
    case "UNLOCK", http.MethodDelete:
        releaseLock(w, r, lockfilePath)
    default:
        utils.MethodNotAllowed(w, r)
    }
}

func acquireLock(w http.ResponseWriter, r *http.Request, lockfilePath, lockDir string) {
    if lockExists(w, lockfilePath) {
        return
    }
    lockInfo, err := decodeLockInfo(r)
    if err != nil {
        utils.HTTPError(w, "Error decoding lock info", err)
        return
    }
    writeLock(w, lockfilePath, lockDir, lockInfo)
}

func releaseLock(w http.ResponseWriter, r *http.Request, lockfilePath string) {
    lockData, err := os.ReadFile(lockfilePath)
    if err != nil {
        utils.HandleFileError(w, r, lockfilePath, err)
        return
    }
    existingLockInfo, err := parseLockData(lockData)
    if err != nil {
        utils.HTTPError(w, "Error unmarshaling lock data", err)
        return
    }
    unlockInfo, err := decodeLockInfo(r)
    if err != nil {
        utils.HTTPError(w, "Error decoding unlock info", err)
        return
    }
    if unlockInfo.ID != existingLockInfo.ID {
        httpConflict(w, lockData)
        return
    }
    removeLock(w, lockfilePath, unlockInfo)
}

func lockExists(w http.ResponseWriter, lockfilePath string) bool {
    if _, err := os.Stat(lockfilePath); err == nil {
        lockData, _ := os.ReadFile(lockfilePath)
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusLocked)
        w.Write(lockData)
        log.Printf("Lock already held for %s", lockfilePath)
        return true
    } else if !os.IsNotExist(err) {
        utils.HTTPError(w, "Error checking lock file", err)
        return true
    }
    return false
}

func decodeLockInfo(r *http.Request) (LockInfo, error) {
    var lockInfo LockInfo
    err := json.NewDecoder(r.Body).Decode(&lockInfo)
    return lockInfo, err
}

func writeLock(w http.ResponseWriter, lockfilePath, lockDir string, lockInfo LockInfo) {
    if err := os.MkdirAll(lockDir, 0755); err != nil {
        utils.HTTPError(w, "Error creating lock directory", err)
        return
    }
    lockData, err := json.Marshal(lockInfo)
    if err != nil {
        utils.HTTPError(w, "Error marshaling lock info", err)
        return
    }
    if err := os.WriteFile(lockfilePath, lockData, 0644); err != nil {
        utils.HTTPError(w, "Error writing lock file", err)
        return
    }
    w.WriteHeader(http.StatusOK)
    log.Printf("Lock acquired for %s by %s", lockfilePath, lockInfo.Who)
}

func removeLock(w http.ResponseWriter, lockfilePath string, unlockInfo LockInfo) {
    if err := os.Remove(lockfilePath); err != nil {
        utils.HTTPError(w, "Error removing lock file", err)
        return
    }
    w.WriteHeader(http.StatusOK)
    log.Printf("Lock released for %s by %s", lockfilePath, unlockInfo.Who)
}

func parseLockData(lockData []byte) (LockInfo, error) {
    var lockInfo LockInfo
    err := json.Unmarshal(lockData, &lockInfo)
    return lockInfo, err
}

func httpConflict(w http.ResponseWriter, lockData []byte) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusConflict)
    w.Write(lockData)
}