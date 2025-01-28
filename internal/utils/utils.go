package utils

import (
    "log"
    "net/http"
    "os"
    "path/filepath"
)

func GetFilePaths(path, dataDir string) (string, string) {
    filePath := filepath.Clean(path)
    fullPath := filepath.Join(dataDir, filePath)
    return fullPath, filepath.Dir(fullPath)
}

func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
    log.Printf("Method not allowed: %s", r.Method)
    http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func HandleFileError(w http.ResponseWriter, r *http.Request, filePath string, err error) {
    if os.IsNotExist(err) {
        // log.Printf("File not found: %s", filePath)
        http.NotFound(w, r)
    } else {
        HTTPError(w, "Error accessing file", err)
    }
}

func HTTPError(w http.ResponseWriter, message string, err error) {
    log.Printf("%s: %v", message, err)
    http.Error(w, err.Error(), http.StatusInternalServerError)
}