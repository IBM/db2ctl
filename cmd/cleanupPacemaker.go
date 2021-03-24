package cmd

import (
	"github.com/IBM/db2ctl/internal/command"
	"github.com/spf13/cobra"
)

// cleanupPacemakerCmd represents the cleanupPacemaker command
var cleanupPacemakerCmd = &cobra.Command{
	Use:   "pacemaker",
	Short: "Cleanup pacemaker",
	Long: `Cleanup pacemaker. For example:

	db2ctl cleanup pacemaker -c db2ctl.yaml`,
	RunE: func(cmd *cobra.Command, args []string) error {

		err := command.New(cmd.Flags()).
			ParseYaml(confFile).
			GenerateConfigFilesFromDir(pacemakerCleanupDir).
			RunBashScripts().
			DeleteStateForDir(pacemakerInstallDir).
			TimeTaken().
			Error

		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	cleanupCmd.AddCommand(cleanupPacemakerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cleanupPacemakerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cleanupPacemakerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
