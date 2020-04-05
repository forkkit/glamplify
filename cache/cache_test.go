package cache

import (
	"gotest.tools/assert"
	"testing"
	"time"
)

func Test_SetGet(t *testing.T) {

	c := New()
	assert.Assert(t, c != nil, c)

	val := "a value"
	c.Set("k1", val, 1*time.Second)

	x, found := c.Get("k1")
	assert.Assert(t, found, found)

	v := x.(string)
	assert.Assert(t, v == "a value", v)
}

func Test_SetGet_Expiry(t *testing.T) {

	c := New()
	assert.Assert(t, c != nil, c)

	val := "a value"
	c.Set("k1", val, 1*time.Second)

	time.Sleep(2*time.Second)

	_, found := c.Get("k1")
	assert.Assert(t, !found, found)
}
