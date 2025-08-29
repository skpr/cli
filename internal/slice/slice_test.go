package slice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	list := []string{"foo", "bar"}
	assert.True(t, Contains(list, "foo"))
	assert.False(t, Contains(list, "baz"))
}

func TestRemove(t *testing.T) {
	list := []string{"foo", "bar"}
	list = Remove(list, "foo")
	assert.Equal(t, []string{"bar"}, list)
}

func TestEqual(t *testing.T) {
	assert.True(t, Equal([]string{"foo", "bar"}, []string{"foo", "bar"}))
	assert.True(t, Equal([]string{"foo", "bar"}, []string{"bar", "foo"}))
	assert.False(t, Equal([]string{"foo"}, []string{"foo", "bar"}))
	assert.False(t, Equal([]string{"foo", "bar"}, []string{"foo"}))
}

func TestAppendSlice(t *testing.T) {
	a := []string{
		"foo",
		"bar",
	}

	b := []string{
		"foo",
		"baz",
	}

	expect := []string{
		"foo",
		"bar",
		"baz",
	}

	assert.ElementsMatch(t, expect, AppendSlice(a, b))
}
