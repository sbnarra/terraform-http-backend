package auth

import (
    "encoding/base64"
    "net/http"
    "net/http/httptest"
    "os"
    "testing"
)

func TestInitializeAuthEnabled(t *testing.T) {
    os.Setenv("AUTH_USERNAME", "testuser")
    os.Setenv("AUTH_PASSWORD", "testpass")
    defer func() {
        os.Unsetenv("AUTH_USERNAME")
        os.Unsetenv("AUTH_PASSWORD")
    }()

    Initialize()

    if !authEnabled {
        t.Errorf("Expected authEnabled to be true, got false")
    }
    if authUsername != "testuser" || authPassword != "testpass" {
        t.Errorf("Authentication credentials not set correctly")
    }
}

func TestInitializeAuthDisabled(t *testing.T) {
    os.Unsetenv("AUTH_USERNAME")
    os.Unsetenv("AUTH_PASSWORD")

    Initialize()

    if authEnabled {
        t.Errorf("Expected authEnabled to be false, got true")
    }
}

func TestWithAuthAuthorized(t *testing.T) {
    authEnabled = true
    authUsername = "testuser"
    authPassword = "testpass"

    handler := WithAuth(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    })

    req := httptest.NewRequest("GET", "/", nil)
    authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte("testuser:testpass"))
    req.Header.Set("Authorization", authHeader)

    rr := httptest.NewRecorder()

    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }
}

func TestWithAuthUnauthorized(t *testing.T) {
    authEnabled = true
    authUsername = "testuser"
    authPassword = "testpass"

    handler := WithAuth(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    })

    req := httptest.NewRequest("GET", "/", nil)

    rr := httptest.NewRecorder()

    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusUnauthorized {
        t.Errorf("Handler returned wrong status code without auth header: got %v want %v", status, http.StatusUnauthorized)
    }

    req = httptest.NewRequest("GET", "/", nil)
    req.Header.Set("Authorization", "Basic invalidencodedstring")

    rr = httptest.NewRecorder()

    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusUnauthorized {
        t.Errorf("Handler returned wrong status code with invalid auth header: got %v want %v", status, http.StatusUnauthorized)
    }

    req = httptest.NewRequest("GET", "/", nil)
    wrongAuthHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte("wronguser:wrongpass"))
    req.Header.Set("Authorization", wrongAuthHeader)

    rr = httptest.NewRecorder()

    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusUnauthorized {
        t.Errorf("Handler returned wrong status code with incorrect credentials: got %v want %v", status, http.StatusUnauthorized)
    }
}

func TestWithAuthDisabled(t *testing.T) {
    authEnabled = false

    handler := WithAuth(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    })

    req := httptest.NewRequest("GET", "/", nil)

    rr := httptest.NewRecorder()

    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("Handler returned wrong status code when auth is disabled: got %v want %v", status, http.StatusOK)
    }
}

func TestCheckAuth(t *testing.T) {
    authUsername = "testuser"
    authPassword = "testpass"

    validAuthHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte("testuser:testpass"))

    if !checkAuth(validAuthHeader) {
        t.Errorf("checkAuth failed with valid credentials")
    }

    invalidAuthHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte("wronguser:wrongpass"))

    if checkAuth(invalidAuthHeader) {
        t.Errorf("checkAuth passed with invalid credentials")
    }

    malformedAuthHeader := "Basic invalidbase64"

    if checkAuth(malformedAuthHeader) {
        t.Errorf("checkAuth passed with malformed Authorization header")
    }

    noPrefixAuthHeader := base64.StdEncoding.EncodeToString([]byte("testuser:testpass"))

    if checkAuth(noPrefixAuthHeader) {
        t.Errorf("checkAuth passed with missing 'Basic ' prefix")
    }
}