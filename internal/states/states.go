package states

import (
    "io"
    "log"
    "net/http"
    "os"

    "terraform-http-backend/internal/utils"
)

func HandleStates(w http.ResponseWriter, r *http.Request, dataDir string) {
    statefilePath, dir := utils.GetFilePaths(r.URL.Path, dataDir)
    switch r.Method {
    case http.MethodGet:
        readState(w, r, statefilePath)
    case http.MethodPost, http.MethodPut:
        writeState(w, r, dir, statefilePath)
    case http.MethodDelete:
        deleteState(w, r, statefilePath)
    default:
        utils.MethodNotAllowed(w, r)
    }
}

func readState(w http.ResponseWriter, r *http.Request, statefilePath string) {
    data, err := os.ReadFile(statefilePath)
    if err != nil {
        utils.HandleFileError(w, r, statefilePath, err)
        return
    }
    w.WriteHeader(http.StatusOK)
    w.Write(data)
}

func writeState(w http.ResponseWriter, r *http.Request, dir, statefilePath string) {
    if err := os.MkdirAll(dir, 0755); err != nil {
        utils.HTTPError(w, "Error creating directory", err)
        return
    }
    file, err := os.Create(statefilePath)
    if err != nil {
        utils.HTTPError(w, "Error creating file", err)
        return
    }
    defer file.Close()
    if _, err := io.Copy(file, r.Body); err != nil {
        utils.HTTPError(w, "Error writing to file", err)
        return
    }
    w.WriteHeader(http.StatusOK)
    log.Printf("Updated state %s", statefilePath)
}

func deleteState(w http.ResponseWriter, r *http.Request, statefilePath string) {
    if err := os.Remove(statefilePath); err != nil {
        utils.HandleFileError(w, r, statefilePath, err)
        return
    }
    w.WriteHeader(http.StatusOK)
    log.Printf("Deleted state %s", statefilePath)
}