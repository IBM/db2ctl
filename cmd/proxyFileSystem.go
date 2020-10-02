package cmd

import (
	"fmt"

	"github.com/IBM/db2ctl/internal/command"
	"github.com/spf13/cobra"
)

// proxyFileSystemCmd represents the proxyFileSystem command
var proxyFileSystemCmd = &cobra.Command{
	Use:   "proxyFileSystem",
	Short: "download the contents to a file system in which there is no access to the machine in DMZ",
	Long: `download the contents to a file system in which there is no access to the machine inDMZ. For example:

	db2ctl proxy file system`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("proxyFileSystem called")

		err := command.New(cmd.Flags()).
			// ProxyFileSystem().
			Error

		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(proxyFileSystemCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// proxyFileSystemCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// proxyFileSystemCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
