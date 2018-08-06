package catalog

import (
	"testing"
)

func TestRandom(t *testing.T) {
	result := random()
	t.Log(result)
}
