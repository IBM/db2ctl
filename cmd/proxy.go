package cmd

import (
	"fmt"

	"github.com/IBM/db2ctl/internal/command"
	"github.com/spf13/cobra"
)

// proxyCmd represents the proxy command
var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "Download at an address to the airgap environment",
	Long: `Download at an address to the airgap environment. For example:

	db2ctl proxy -c s3config.yaml`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Proxy Command to download file from IBM S3 COS to airgap environment")

		err := command.New(cmd.Flags()).
			Error

		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(proxyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// proxyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// proxyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
