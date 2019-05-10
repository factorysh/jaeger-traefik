package tiny

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplit(t *testing.T) {
	slugs := split("beuha:aussi", ':')
	assert.Equal(t, []string{"beuha", "aussi"}, slugs)
	slugs = split("a:b:c:", ':')
	assert.Equal(t, []string{"a", "b", "c", ""}, slugs)
}
