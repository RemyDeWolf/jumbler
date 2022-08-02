package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitDefault(t *testing.T) {
	config, err := GetDefault()

	assert.NoError(t, err)
	assert.Equal(t, "", config.Ext)
}
