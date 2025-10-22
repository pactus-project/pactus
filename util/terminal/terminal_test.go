package terminal

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// captureOutput is a helper function to capture the printed output of a function.
func captureOutput(fun func()) string {
	// Redirect stdout to a buffer
	oldStdout := os.Stdout
	reader, writer, _ := os.Pipe()
	os.Stdout = writer

	// Capture the printed output
	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, reader)
		outC <- buf.String()
	}()

	// Execute the function
	fun()

	// Reset stdout
	_ = writer.Close()
	os.Stdout = oldStdout
	out := <-outC

	return out
}

func TestPrintNotSupported(t *testing.T) {
	terminalSupported = false
	output := captureOutput(func() {
		PrintJSONObject([]int{1, 2, 3})
		PrintLine()
		PrintInfoMsgBoldf("This is PrintInfoMsgBoldf: %s", "msg")
		PrintInfoMsgf("This is PrintInfoMsgf: %s", "msg")
		PrintSuccessMsgf("This is PrintSuccessMsgf: %s", "msg")
		PrintWarnMsgf("This is PrintWarnMsgf: %s", "msg")
		PrintErrorMsgf("This is PrintErrorMsgf: %s", "msg")
	})

	expected := "[\n   1,\n   2,\n   3\n]\n" +
		"\n" +
		"This is PrintInfoMsgBoldf: msg\n" +
		"This is PrintInfoMsgf: msg\n" +
		"This is PrintSuccessMsgf: msg\n" +
		"This is PrintWarnMsgf: msg\n" +
		"[ERROR] This is PrintErrorMsgf: msg\n"

	assert.Equal(t, expected, output)
}

func TestPrintSupported(t *testing.T) {
	terminalSupported = true
	output := captureOutput(func() {
		PrintJSONObject([]int{1, 2, 3})
		PrintLine()
		PrintInfoMsgBoldf("This is PrintInfoMsgBoldf: %s", "msg")
		PrintInfoMsgf("This is PrintInfoMsgf: %s", "msg")
		PrintSuccessMsgf("This is PrintSuccessMsgf: %s", "msg")
		PrintWarnMsgf("This is PrintWarnMsgf: %s", "msg")
		PrintErrorMsgf("This is PrintErrorMsgf: %s", "msg")
	})

	expected := "[\n   1,\n   2,\n   3\n]\n" +
		"\n" +
		"\x1b[1mThis is PrintInfoMsgBoldf: msg\x1b[0m\n" +
		"This is PrintInfoMsgf: msg\n" +
		"\x1b[32mThis is PrintSuccessMsgf: msg\x1b[0m\n" +
		"\x1b[33mThis is PrintWarnMsgf: msg\x1b[0m\n" +
		"\x1b[31m[ERROR] This is PrintErrorMsgf: msg\x1b[0m\n"

	assert.Equal(t, expected, output)
}
