package downloader

import "fmt"

type chunk struct {
	start, end int64
}

func createChunks(contentLength, totalChunks int64) []*chunk {
	chunks := make([]*chunk, 0, totalChunks)
	chunkSize := contentLength / totalChunks
	for i := int64(0); i < totalChunks; i++ {
		start := i * chunkSize
		end := start + chunkSize - 1
		// adjust the end for the last chunk
		if i == totalChunks-1 {
			end = contentLength - 1
		}
		chunks = append(chunks, &chunk{start: start, end: end})
	}

	return chunks
}

func (c *chunk) rangeHeader() string {
	return fmt.Sprintf("bytes=%d-%d", c.start, c.end)
}

func (c *chunk) size() int64 {
	return (c.end + 1) - c.start
}
