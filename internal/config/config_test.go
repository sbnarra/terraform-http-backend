package config

import (
    "os"
    "testing"
)

func TestGetEnv(t *testing.T) {
    testCases := []struct {
        description string
        envKey      string
        envValue    string
        fallback    string
        expected    string
    }{
        {
            description: "Environment variable is set",
            envKey:      "EXISTING_KEY",
            envValue:    "value1",
            fallback:    "default_value",
            expected:    "value1",
        },
        {
            description: "Environment variable is empty",
            envKey:      "EMPTY_KEY",
            envValue:    "",
            fallback:    "default_value",
            expected:    "default_value",
        },
        {
            description: "Environment variable is not set",
            envKey:      "NON_EXISTING_KEY",
            envValue:    "",
            fallback:    "default_value",
            expected:    "default_value",
        },
    }

    for _, tc := range testCases {
        t.Run(tc.description, func(t *testing.T) {
            if tc.envValue != "" || tc.envKey == "EMPTY_KEY" {
                os.Setenv(tc.envKey, tc.envValue)
                defer os.Unsetenv(tc.envKey)
            } else {
                os.Unsetenv(tc.envKey)
            }

            result := GetEnv(tc.envKey, tc.fallback)

            if result != tc.expected {
                t.Errorf("GetEnv(%q, %q) = %q; want %q", tc.envKey, tc.fallback, result, tc.expected)
            }
        })
    }
}