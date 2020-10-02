package cmd

import (
	"github.com/IBM/db2ctl/internal/command"
	"github.com/spf13/cobra"
)

// linbitCmd represents the linbit command
var linbitCleanupCmd = &cobra.Command{
	Use:   "linbit",
	Short: "Cleanup linbit",
	Long: `Cleanup linbit. For example:

	db2ctl cleanup linbit -c db2ctl-sample.yaml`,
	RunE: func(cmd *cobra.Command, args []string) error {

		err := command.New(cmd.Flags()).
			ParseYaml(confFile).
			GenerateConfigFilesFromDir(linbitCleanupDir).
			SetAutoYesEnabled(true).
			RunBashScripts().
			DeleteStateForDir(linbitInstallDir).
			TimeTaken().
			Error

		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	cleanupCmd.AddCommand(linbitCleanupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// linbitCleanupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// linbitCleanupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
