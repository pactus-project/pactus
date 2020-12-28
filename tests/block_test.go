package tests

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func getBlockAt(t *testing.T, height int) string {
	url := fmt.Sprintf("http://%s/block/height/%d", tCurlAddress, height)
	for i := 0; i < 50; i++ {
		res, err := http.Get(url)
		if err == nil {
			if res.StatusCode == 200 {
				buf := new(bytes.Buffer)
				_, err := buf.ReadFrom(res.Body)
				assert.NoError(t, err)
				return buf.String()
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
	assert.NoError(t, fmt.Errorf("timeout"))
	return ""
}

func TestGeneratingBlocks(t *testing.T) {

	res := getBlockAt(t, 1)
	assert.Contains(t, res, "0000000000000000000000000000000000000000000000000000000000000000")
	fmt.Println(res)
}
