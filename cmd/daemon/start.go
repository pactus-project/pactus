package main

import (
	"bytes"
	"fmt"
	"net/http"
	_ "net/http/pprof" // #nosec
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/wallet"
	"github.com/spf13/cobra"
)

// buildStartCmd builds a sub-command to starts the Pactus blockchain node.
func buildStartCmd(parentCmd *cobra.Command) {
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start the Pactus blockchain node",
	}

	parentCmd.AddCommand(startCmd)

	workingDirOpt := startCmd.Flags().StringP("working-dir", "w",
		cmd.PactusHomeDir(), "A path to the working directory to read the wallet and node files")

	passwordOpt := startCmd.Flags().StringP("password", "p", "", "The wallet password")

	pprofOpt := startCmd.Flags().String("pprof", "", "debug pprof server address(not recommended to expose to internet)")

	startCmd.Run = func(_ *cobra.Command, _ []string) {
		workingDir, _ := filepath.Abs(*workingDirOpt)
		// change working directory
		err := os.Chdir(workingDir)
		cmd.FatalErrorCheck(err)

		pidFile := filepath.Join(workingDir, ".pactus.pid")
		// Check for already running instance
		if isAlreadyRunning(pidFile) {
			fmt.Println("An instance of Pactus is already running. Please stop it before starting a new one.")
			return
		}

		if *pprofOpt != "" {
			cmd.PrintWarnMsgf("Starting Debug pprof server on: http://%s/debug/pprof/", *pprofOpt)
			server := &http.Server{
				Addr:              *pprofOpt,
				ReadHeaderTimeout: 3 * time.Second,
			}
			go func() {
				err := server.ListenAndServe()
				cmd.FatalErrorCheck(err)
			}()
		}

		passwordFetcher := func(wlt *wallet.Wallet) (string, bool) {
			if !wlt.IsEncrypted() {
				return "", true
			}

			var password string
			if *passwordOpt != "" {
				password = *passwordOpt
			} else {
				password = cmd.PromptPassword("Wallet password", false)
			}
			return password, true
		}
		node, _, err := cmd.StartNode(
			workingDir, passwordFetcher)
		cmd.FatalErrorCheck(err)

		// Write current PID to the file
		err = os.WriteFile(pidFile, []byte(fmt.Sprintf("%d", os.Getpid())), 0o600)
		cmd.FatalErrorCheck(err)

		cmd.TrapSignal(func() {
			os.Remove(pidFile)
			node.Stop()
			cmd.PrintInfoMsgf("Exiting ...")
		})

		// run forever (the node will not be returned)
		select {}
	}
}

// isAlreadyRunning checks if an instance of the application is already running.
func isAlreadyRunning(pidFile string) bool {
	if data, err := os.ReadFile(pidFile); err == nil {
		pid, err := strconv.Atoi(string(data))
		if err == nil && pidExists(pid) {
			return true // PID found and process is running
		}
	}
	return false
}

// pidExists checks if a given PID is currently active.
func pidExists(pid int) bool {
	if pid < 0 {
		return false
	}

	if runtime.GOOS == "windows" {
		pidStr := strconv.Itoa(pid)
		windowsCmd := exec.Command("tasklist", "/FI", "PID eq "+pidStr)
		var out bytes.Buffer
		windowsCmd.Stdout = &out
		err := windowsCmd.Run()
		if err != nil {
			return false
		}
		return strings.Contains(out.String(), strconv.Itoa(pid))
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	// On Unix systems, FindProcess always succeeds and the call to Signal does not kill the process
	return process.Signal(syscall.Signal(0)) == nil
}
