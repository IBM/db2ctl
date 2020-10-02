package cmd

import (
	"github.com/IBM/db2ctl/internal/command"
	"github.com/spf13/cobra"
)

// generateDB2 represents the db2 command
var generateDB2 = &cobra.Command{
	Use:   "db2",
	Short: "Generate db2 install scripts",
	Long:  `Generate db2 install scripts`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := command.New(cmd.Flags()).
			ParseYaml(confFile).
			GenerateConfigFilesFromDir(db2InstallDir).
			Error

		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	generateCmd.AddCommand(generateDB2)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateDB2.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateDB2.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
