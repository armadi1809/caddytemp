package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "caddytemp",
	Short: "caddytemp is a CLI tool to generate common Caddyfile configurations",
	Long: `A CLI application helping with the generation of common Caddyfile configurations. A list of all available configs
		can be found at https://caddyserver.com/docs/caddyfile/patterns`,
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
