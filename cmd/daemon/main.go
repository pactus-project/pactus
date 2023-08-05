package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// startCmd = &cobra.Command{
	// 	Use: "start",
	// 	Short: "Start the Pactus blockchain",
	// 	Run: Start(),
	// }

)

func main() {
	var rootCmd = &cobra.Command{
		Use: "pactus-daemon",
		Short: "Pactus daemon",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("use --help")
		},
	}

	var initCmd = &cobra.Command{
		Use: "init",
		Short: "Initialize the Pactus blockchain",
		Run: Init(),
	}

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.Execute()
}
