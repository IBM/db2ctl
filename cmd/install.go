package cmd

import (
	"github.com/spf13/cobra"
)

// installCmd represents the run command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install commands",
	Long: `Install commands. For example:

	db2ctl install <command> -c db2ctl-sample.yaml`,
}

func init() {
	rootCmd.AddCommand(installCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
