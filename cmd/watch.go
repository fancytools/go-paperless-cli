package cmd

import (
	"github.com/fancytools/go-paperless-cli/pkg/fileUpload"
	"github.com/spf13/cobra"
)

var endpoint, token, watchedDir string
var deleteFile bool

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watches on file writes in a directory and uploads them",
	Long:  `"Watches on file writes in a directory and uploads them`,
	Args:  cobra.NoArgs,

	Run: func(cmd *cobra.Command, args []string) {
		fileUpload.StartWatching(watchedDir, endpoint, token, deleteFile)
	},
}

func init() {
	rootCmd.AddCommand(watchCmd)

	watchCmd.Flags().StringVarP(&endpoint, "endpoint", "e", "http://paperlesss.local:8080", "Address of paperless-ng instance")
	watchCmd.Flags().StringVarP(&token, "token", "t", "", "Token of your paperless-ng instance")
	watchCmd.Flags().StringVarP(&watchedDir, "watchedDir", "d", "./", "Path to the watched dir")
	watchCmd.Flags().BoolVar(&deleteFile, "deleteFile", false, "Delete successful transferred files")
}
