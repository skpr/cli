package version

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrint(t *testing.T) {
	var buffer bytes.Buffer

	params := PrintParams{}

	err := Print(&buffer, params)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "version not found")

	params.ClientVersion = "0.0.1"
	err = Print(&buffer, params)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "build date not found")

	params.ClientBuildDate = "2020-08-05-T16:41:12+1000"
	err = Print(&buffer, params)
	assert.NoError(t, err)

	assert.True(t, strings.Contains(buffer.String(), params.ClientVersion))
	assert.True(t, strings.Contains(buffer.String(), params.ClientBuildDate))
}
