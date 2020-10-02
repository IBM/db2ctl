package cmd

import (
	"fmt"

	"github.com/IBM/db2ctl/internal/command"
	"github.com/spf13/cobra"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "defines the upstream proxy server to synchronize the contents to the local control node",
	Long: `defines the upstream proxy server to synchronize the contents to the local control node. For example:

	db2ctl client`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("client called")

		err := command.New(cmd.Flags()).
			// Client().
			Error

		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clientCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clientCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
