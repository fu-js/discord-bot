package main

import (
	"fmt"
	"github.com/fu-js/discord-bot/cmd/viblo"
	"github.com/spf13/cobra"
	"os"
)

var root = &cobra.Command{
	Use:   "discord",
	Short: "discord bot",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func main() {
	root.AddCommand(viblo.Cmd)
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
