package alias

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpandTemplate(t *testing.T) {
	alias := "mysql image pull dev"
	args := []string{}

	result := ExpandTemplate(alias, args)
	assert.Equal(t, "mysql image pull dev", result)

	alias = "exec $1 -- drush uli"
	args = []string{"dev"}

	result = ExpandTemplate(alias, args)
	assert.Equal(t, "exec dev -- drush uli", result)

	alias = "exec $1 -- drush status $1"
	args = []string{"dev"}

	result = ExpandTemplate(alias, args)
	assert.Equal(t, "exec dev -- drush status dev", result)

	alias = "exec $1 -- drush php:exec $2"
	args = []string{"dev", "phpinfo()"}

	result = ExpandTemplate(alias, args)
	assert.Equal(t, "exec dev -- drush php:exec phpinfo()", result)
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
