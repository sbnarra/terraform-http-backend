package auth

import (
    "encoding/base64"
    "log"
    "net/http"
    "os"
    "strings"
)

var authEnabled bool
var authUsername string
var authPassword string

// Initialize sets up authentication based on environment variables
func Initialize() {
    authUsername = os.Getenv("AUTH_USERNAME")
    authPassword = os.Getenv("AUTH_PASSWORD")
    if authUsername != "" && authPassword != "" {
        authEnabled = true
        log.Println("Basic authentication enabled")
    } else {
        authEnabled = false
        log.Println("Warning: Basic authentication is disabled")
    }
}

// WithAuth is a middleware that provides HTTP Basic Authentication
func WithAuth(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if authEnabled {
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" || !checkAuth(authHeader) {
                w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
								log.Println("Unauthorized: " + r.URL.Path)
                return
            }
        }
        next(w, r)
    }
}

func checkAuth(authHeader string) bool {
    const prefix = "Basic "
    if !strings.HasPrefix(authHeader, prefix) {
        return false
    }
    authEncoded := strings.TrimPrefix(authHeader, prefix)
    authDecodedBytes, err := base64.StdEncoding.DecodeString(authEncoded)
    if err != nil {
        return false
    }
    authDecoded := string(authDecodedBytes)
    authPair := strings.SplitN(authDecoded, ":", 2)
    if len(authPair) != 2 {
        return false
    }
    return authPair[0] == authUsername && authPair[1] == authPassword
}