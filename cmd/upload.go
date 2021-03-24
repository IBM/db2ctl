package cmd

import (
	"fmt"

	"github.com/IBM/db2ctl/internal/command"
	"github.com/spf13/cobra"
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload file",
	Long: `Upload file to IBM S3 COS. For example:

	db2ctl upload -c s3config.yaml`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Upload command to upload a file to IBM S3 COS")

		err := command.New(cmd.Flags()).
			ParseS3Config(confFile).
			UploadToS3().
			Error

		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// uploadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// uploadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
