package utils

import (
    "errors"
    "net/http"
    "net/http/httptest"
    "os"
    "path/filepath"
    "testing"
)

func TestGetFilePaths(t *testing.T) {
    dataDir := "/testdata"
    tests := []struct {
        inputPath       string
        expectedFull    string
        expectedDir     string
    }{
        {
            inputPath:    "/statefile.tfstate",
            expectedFull: filepath.Join("/testdata", "statefile.tfstate"),
            expectedDir:  filepath.Dir(filepath.Join("/testdata", "statefile.tfstate")),
        },
        {
            inputPath:    "nested/dir/statefile.tfstate",
            expectedFull: filepath.Join("/testdata", "nested/dir/statefile.tfstate"),
            expectedDir:  filepath.Dir(filepath.Join("/testdata", "nested/dir/statefile.tfstate")),
        },
    }

    for _, test := range tests {
        fullPath, dir := GetFilePaths(test.inputPath, dataDir)
        if fullPath != test.expectedFull {
            t.Errorf("GetFilePaths(%q, %q) fullPath = %q; want %q", test.inputPath, dataDir, fullPath, test.expectedFull)
        }
        if dir != test.expectedDir {
            t.Errorf("GetFilePaths(%q, %q) dir = %q; want %q", test.inputPath, dataDir, dir, test.expectedDir)
        }
    }
}

func TestMethodNotAllowed(t *testing.T) {
    rr := httptest.NewRecorder()
    req := httptest.NewRequest("TRACE", "/", nil)

    MethodNotAllowed(rr, req)

    if status := rr.Code; status != http.StatusMethodNotAllowed {
        t.Errorf("MethodNotAllowed returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
    }

    expectedBody := "Method not allowed\n"
    if rr.Body.String() != expectedBody {
        t.Errorf("MethodNotAllowed returned wrong body: got %q want %q", rr.Body.String(), expectedBody)
    }
}

func TestHandleFileError(t *testing.T) {
    tests := []struct {
        err             error
        expectedStatus  int
        expectedBody    string
    }{
        {
            err:            os.ErrNotExist,
            expectedStatus: http.StatusNotFound,
            expectedBody:   "404 page not found\n",
        },
        {
            err:            errors.New("some other error"),
            expectedStatus: http.StatusInternalServerError,
            expectedBody:   "some other error\n",
        },
    }

    for _, test := range tests {
        rr := httptest.NewRecorder()
        req := httptest.NewRequest("GET", "/", nil)
        filePath := "/fake/path"

        HandleFileError(rr, req, filePath, test.err)

        if status := rr.Code; status != test.expectedStatus {
            t.Errorf("HandleFileError returned wrong status code for error %v: got %v want %v",
                test.err, status, test.expectedStatus)
        }

        if rr.Body.String() != test.expectedBody {
            t.Errorf("HandleFileError returned wrong body for error %v: got %q want %q",
                test.err, rr.Body.String(), test.expectedBody)
        }
    }
}

func TestHTTPError(t *testing.T) {
    rr := httptest.NewRecorder()
    err := errors.New("internal server error")

    HTTPError(rr, "Test error message", err)

    if status := rr.Code; status != http.StatusInternalServerError {
        t.Errorf("HTTPError returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
    }

    expectedBody := "internal server error\n"
    if rr.Body.String() != expectedBody {
        t.Errorf("HTTPError returned wrong body: got %q want %q", rr.Body.String(), expectedBody)
    }
}