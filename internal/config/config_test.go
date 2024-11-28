package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	path := "../../configs.yml"
	assert.FileExists(t, path, "File exists")
	assert.NotPanics(t, func() {
		New(path)
	})
	config := New(path)
	assert.NotNil(t, config)
	assert.NotNil(t, config.App.ImageFormats)
	assert.Contains(t, config.App.ImageFormats, ".jpg", "image formats", config.App.ImageFormats)
}
