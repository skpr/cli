package list

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/skpr/api/pb"
)

func TestPrint(t *testing.T) {
	list := []*pb.Config{
		{
			Key:   "foo",
			Value: "xxxxx",
			Type:  pb.ConfigType_User,
		},
		{
			Key:   "bar",
			Value: "yyyyyy",
			Type:  pb.ConfigType_System,
		},
	}

	var b bytes.Buffer

	err := Print(&b, list, "dev", 5, false)
	assert.NoError(t, err)

	assert.Contains(t, b.String(), "yyyyy...")
	assert.Contains(t, b.String(), "xxxxx")
	assert.Contains(t, b.String(), "Values have been trimmed.")
}
