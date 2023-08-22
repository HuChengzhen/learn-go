package cache

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestErr(t *testing.T) {
	assert.True(t, errors.Is(fmt.Errorf("%w, key: %s", ErrKeyNotFound, "a"), ErrKeyNotFound))
}
