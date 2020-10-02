package cmd

import (
	"github.com/IBM/db2ctl/internal/command"
	"github.com/spf13/cobra"
)

// cleanupAllCmd represents the cleanupAll command
var cleanupAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Cleanup all modules",
	Long: `Cleanup all modules. For example:

	db2ctl cleanup all -c db2ctl-sample.yaml`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := command.New(cmd.Flags()).
			ParseYaml(confFile).
			GenerateConfigFilesFromDir("cleanup").
			RunBashScripts().
			DeleteStateForDir("install").
			TimeTaken().
			Error

		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	cleanupCmd.AddCommand(cleanupAllCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cleanupAllCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cleanupAllCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
