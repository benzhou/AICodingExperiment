package utils

import (
	"os"
	"testing"
)

func TestGetEnvOrDefault(t *testing.T) {
	// Test with existing env variable
	os.Setenv("TEST_VAR", "test_value")
	if got := GetEnvOrDefault("TEST_VAR", "default"); got != "test_value" {
		t.Errorf("GetEnvOrDefault() = %v, want %v", got, "test_value")
	}

	// Test with non-existing env variable
	if got := GetEnvOrDefault("NON_EXISTENT", "default"); got != "default" {
		t.Errorf("GetEnvOrDefault() = %v, want %v", got, "default")
	}
}

func TestGetEnvIntOrDefault(t *testing.T) {
	// Test with valid integer
	os.Setenv("TEST_INT", "123")
	if got := GetEnvIntOrDefault("TEST_INT", 0); got != 123 {
		t.Errorf("GetEnvIntOrDefault() = %v, want %v", got, 123)
	}

	// Test with invalid integer
	os.Setenv("TEST_INVALID_INT", "not_a_number")
	if got := GetEnvIntOrDefault("TEST_INVALID_INT", 456); got != 456 {
		t.Errorf("GetEnvIntOrDefault() = %v, want %v", got, 456)
	}
}

func TestGetEnvBoolOrDefault(t *testing.T) {
	// Test with valid boolean
	os.Setenv("TEST_BOOL", "true")
	if got := GetEnvBoolOrDefault("TEST_BOOL", false); !got {
		t.Errorf("GetEnvBoolOrDefault() = %v, want %v", got, true)
	}

	// Test with invalid boolean
	os.Setenv("TEST_INVALID_BOOL", "not_a_bool")
	if got := GetEnvBoolOrDefault("TEST_INVALID_BOOL", true); !got {
		t.Errorf("GetEnvBoolOrDefault() = %v, want %v", got, true)
	}
}
