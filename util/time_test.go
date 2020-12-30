package util

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNow(t *testing.T) {
	c1 := time.Now()
	c2 := Now()
	c3 := RoundNow(5)

	assert.NotEqual(t, c1, c2)
	assert.Equal(t, c3.Nanosecond(), 0)
	assert.Equal(t, c3.Second()%5, 0)
	fmt.Println(c1)
	fmt.Println(c2)
	fmt.Println(c3)
}

func TestTimeCanonical(t *testing.T) {
	t1, _ := time.Parse(time.RFC3339Nano, "2006-01-02T15:04:04.999999999Z")
	t2, _ := time.Parse(time.RFC3339Nano, "2006-01-02T15:04:05.111111111Z")
	t3, _ := time.Parse(time.RFC3339Nano, "2006-01-02T15:04:01.999999999Z")
	t4, _ := time.Parse(time.RFC3339Nano, "2006-01-02T15:04:09.999999999Z")
	c1 := RoundTime(t1, 10)
	c2 := RoundTime(t2, 10)
	c3 := RoundTime(t3, 10)
	c4 := RoundTime(t4, 10)
	fmt.Println(c1)
	fmt.Println(c2)
	fmt.Println(c3)
	fmt.Println(c4)

	assert.Equal(t, c1.Nanosecond(), 0)
	assert.Equal(t, c2.Nanosecond(), 0)
	assert.Equal(t, c3.Nanosecond(), 0)
	assert.Equal(t, c4.Nanosecond(), 0)

	assert.Equal(t, c1.Second(), 0)
	assert.Equal(t, c2.Second(), 10)
	assert.Equal(t, c3.Second(), 0)
	assert.Equal(t, c4.Second(), 10)

}
