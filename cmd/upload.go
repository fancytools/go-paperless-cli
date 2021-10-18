package cmd

import (
	"io/fs"
	"log"
	"os"

	"github.com/fancytools/go-paperless-cli/pkg/fileUpload"
	"github.com/spf13/cobra"
)

var filePath string

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload all files in a directory or file",
	Long:  `Upload all files in a directory or file`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fileSystem := os.DirFS(filePath)
		files := make(map[string]struct{}, 100)

		err := fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
			files[path] = struct{}{}
			return nil
		})
		if err != nil {
			log.Println(err)
			return
		}

		for file := range files {
			fileUpload.UploadFile(endpoint, file, token, false)
		}
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)

	uploadCmd.Flags().StringVarP(&endpoint, "endpoint", "e", "http://paperlesss.local:8080", "Address of paperless-ng instance")
	uploadCmd.Flags().StringVarP(&token, "token", "t", "", "Token of your paperless-ng instance")
	uploadCmd.Flags().StringVarP(&filePath, "filePath", "d", "./", "Path to file or directory to upload")
	uploadCmd.Flags().BoolVar(&deleteFile, "deleteFile", false, "Delete successful transferred files")
}
