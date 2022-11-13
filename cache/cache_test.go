package cache

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	c *Cache
)

func init() {
	c = New()
}

// A B C D E F G H I J K L M N O P Q R S T
func setDataA2T() {
	c.Clean()
	for i := 0; i < 20; i++ {
		t := string(rune(i + 65))
		c.Set(t, t)
	}
}

func PrintKeys(t *testing.T, c *Cache) {
	for k := c.keys.Front(); k != nil; k = k.Next() {
		fmt.Print(k.Value, " ")
	}
	fmt.Println()
}

func TestResetValue(t *testing.T) {
	c.Clean()

	c.Set("A", "aaa")
	assert.Equal(t, "aaa", c.Get("A"))

	c.Set("A", map[string]interface{}{"reset_a": "aaaa"})
	assert.Equal(t, map[string]interface{}{"reset_a": "aaaa"}, c.Get("A"))
}

func TestSetValueAfterMaxLen(t *testing.T) {
	setDataA2T()

	assert.Equal(t, "A", c.keys.Front().Value)
	c.Set("Z", "Z")
	assert.Equal(t, "B", c.keys.Front().Value)
	assert.Equal(t, "Z", c.keys.Back().Value)

	assert.Nil(t, c.Get("A"))
}
