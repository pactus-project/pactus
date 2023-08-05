package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "pactus-daemon",
		Short: "Pactus daemon",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("use --help")
		},
	}

	buildVersionCmd(rootCmd)
	buildInitCmd(rootCmd)
	buildStartCmd(rootCmd)
	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
