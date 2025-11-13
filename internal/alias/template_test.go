package alias

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpandTemplate(t *testing.T) {
	alias := "exec $1 -- drush uli"
	args := []string{"dev"}

	result := ExpandTemplate(alias, args)
	assert.Equal(t, "exec dev -- drush uli", result)
}

func TestCountTemplateArgs(t *testing.T) {
	result := CountTemplateArgs("mysql image pull dev")
	assert.Equal(t, 0, result)

	result = CountTemplateArgs("exec $1 -- drush uli")
	assert.Equal(t, 1, result)

	result = CountTemplateArgs("exec $1 -- drush sql-dump --structure-tables-key=common --result-file=/mnt/private/$1-db.sql")
	assert.Equal(t, 1, result)

	result = CountTemplateArgs("exec $1 -- drush php:cli $2")
	assert.Equal(t, 2, result)
}
