package cmd

import (
	"github.com/IBM/db2ctl/internal/command"
	"github.com/spf13/cobra"
)

// cleanupDB2Cmd represents the cleanupDB2 command
var cleanupDB2Cmd = &cobra.Command{
	Use:   "db2",
	Short: "Cleanup db2",
	Long: `Cleanup db2. For example:

	db2ctl cleanup db2 -c db2ctl-sample.yaml`,
	RunE: func(cmd *cobra.Command, args []string) error {

		err := command.New(cmd.Flags()).
			ParseYaml(confFile).
			GenerateConfigFilesFromDir(db2CleanupDir).
			RunBashScripts().
			DeleteStateForDir(db2InstallDir).
			TimeTaken().
			Error

		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	cleanupCmd.AddCommand(cleanupDB2Cmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cleanupDB2Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cleanupDB2Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
