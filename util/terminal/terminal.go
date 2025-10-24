//nolint:forbidigo // enable printing function for terminal package
package terminal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
)

var terminalSupported = false

func init() {
	terminalSupported = CheckTerminalSupported()
}

// CheckTerminalSupported returns true if the current terminal supports
// ANSI escape sequences and advanced terminal features.
func CheckTerminalSupported() bool {
	bad := map[string]bool{"": true, "dumb": true, "cons25": true}

	return !bad[strings.ToLower(os.Getenv("TERM"))]
}

func FatalErrorCheck(err error) {
	if err != nil {
		if terminalSupported {
			fmt.Printf("\033[31m%s\033[0m\n", err.Error())
		} else {
			fmt.Printf("%s\n", err.Error())
		}

		os.Exit(1)
	}
}

func PrintErrorMsgf(format string, args ...any) {
	format = "[ERROR] " + format
	if terminalSupported {
		// Print error msg with red color
		format = fmt.Sprintf("\033[31m%s\033[0m", format)
	}
	fmt.Printf(format+"\n", args...)
}

func PrintSuccessMsgf(format string, a ...any) {
	if terminalSupported {
		// Print successful msg with green color
		format = fmt.Sprintf("\033[32m%s\033[0m", format)
	}
	fmt.Printf(format+"\n", a...)
}

func PrintWarnMsgf(format string, a ...any) {
	if terminalSupported {
		// Print warning msg with yellow color
		format = fmt.Sprintf("\033[33m%s\033[0m", format)
	}
	fmt.Printf(format+"\n", a...)
}

func PrintInfoMsgf(format string, a ...any) {
	fmt.Printf(format+"\n", a...)
}

func PrintInfoMsgBoldf(format string, a ...any) {
	if terminalSupported {
		format = fmt.Sprintf("\033[1m%s\033[0m", format)
	}
	fmt.Printf(format+"\n", a...)
}

func PrintLine() {
	fmt.Println()
}

func PrintJSONData(data []byte) {
	var out bytes.Buffer
	err := json.Indent(&out, data, "", "   ")
	FatalErrorCheck(err)

	PrintInfoMsgf(out.String())
}

func PrintJSONObject(obj any) {
	data, err := json.Marshal(obj)
	FatalErrorCheck(err)

	PrintJSONData(data)
}

func ProgressBar(totalSize int64, barWidth int) *progressbar.ProgressBar {
	if barWidth < 15 {
		barWidth = 15
	}

	opts := []progressbar.Option{
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetWidth(barWidth),
		progressbar.OptionSetElapsedTime(false),
		progressbar.OptionSetPredictTime(false),
		progressbar.OptionShowDescriptionAtLineEnd(),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	}

	return progressbar.NewOptions64(totalSize, opts...)
}
