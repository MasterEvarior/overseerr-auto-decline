package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMedia(t *testing.T) {
	t.Setenv("TEST_MEDIA", "12,323,4")

	result := GetMedia("TEST_MEDIA")

	assert.EqualValues(t, []string{"12", "323", "4"}, result)
}

func TestGetMedia_EmptyString(t *testing.T) {
	t.Setenv("TEST_MEDIA", "")

	result := GetMedia("TEST_MEDIA")

	assert.EqualValues(t, []string{""}, result)
}

func TestGetEnvVar(t *testing.T) {
	t.Setenv("TEST_VAR", "my-value")

	result := GetEnvVar("TEST_VAR")

	assert.Equal(t, "my-value", result)
}
