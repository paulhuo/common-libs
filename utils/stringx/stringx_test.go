package stringx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandom(t *testing.T) {
	s := Random(10)
	assert.Equal(t, 10, len(s))
}

func TestToSnake(t *testing.T) {
	assert.Equal(t, "hello_world", ToSnake("HelloWorld"))
	assert.Equal(t, "foo_bar", ToSnake("foo bar"))
	assert.Equal(t, "foo_bar", ToSnake("foo-bar"))
}

func TestToCamel(t *testing.T) {
	assert.Equal(t, "HelloWorld", ToCamel("hello_world"))
	assert.Equal(t, "FooBar", ToCamel("foo-bar"))
	assert.Equal(t, "Hello", ToCamel("hello"))
}

func TestTruncate(t *testing.T) {
	assert.Equal(t, "hello", Truncate("hello", 10))
	assert.Equal(t, "he...", Truncate("hello world", 5))
	assert.Equal(t, "he", Truncate("hello", 2))
}
