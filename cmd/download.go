package cmd

import (
	"fmt"

	"github.com/IBM/db2ctl/internal/command"
	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download from IBM S3",
	Long: `Download from IBM S3. For example:

	db2ctl download -c s3config.yaml`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Download command to download file from IBM S3 COS")

		err := command.New(cmd.Flags()).
			ParseS3Config(confFile).
			DownloadFromS3().
			Error

		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downloadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
