package helper

import (
	"slices"
	"testing"
)

func TestGetMedia(t *testing.T) {
	t.Setenv("TEST_MEDIA", "12,323,4")

	result := GetMedia("TEST_MEDIA")

	if !slices.Equal(result, []string{"12", "323", "4"}) {
		t.Fatalf("'%v' does not equal the expected values", result)
	}
}

func TestGetEnvVar(t *testing.T) {
	t.Setenv("TEST_VAR", "my-value")

	result := GetEnvVar("TEST_VAR")

	if result != "my-value" {
		t.Fatalf("'%v' does not equal the expected value", result)
	}
}
