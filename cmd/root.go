package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "go-paperless-cli",
	Short: "A go based paperless-ng API client",
	Long:  `go-paperless-cli is a go based paperless-ng API client`,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
}
