package config

import (
	"testing"

	"github.com/bmizerany/assert"
)

func TestConfig(t *testing.T) {
	configs, _ := NewConfig("testdata/endpoints.yaml")
	assert.Equal(t, 2, len(configs), "Expecting one valid endpoint")
}
