package catalog

import (
	"fmt"
	"testing"
)

func TestRandom(t *testing.T) {
	result := random()
	fmt.Println(result)
}
