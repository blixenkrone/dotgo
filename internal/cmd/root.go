package cmd

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/spf13/cobra"
)

type CobraFunc func(cmd *cobra.Command, args []string)

var rootCmd = &cobra.Command{
	Use:   "dotgo",
	Short: "dotgo is a tool to take your dotfiles to the moon",
}

func Execute() {
	initSelectDotfilesCmd()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
		debug.PrintStack()
	}
}
