package cmd

import (
	"github.com/IBM/db2ctl/internal/command"
	"github.com/spf13/cobra"
)

// installDB2Cmd represents the installDB2 command
var installDB2Cmd = &cobra.Command{
	Use:   "db2",
	Short: "Install db2",
	Long: `Install db2. For example:

	db2ctl install db2 -c db2ctl-sample.yaml`,
	RunE: func(cmd *cobra.Command, args []string) error {

		err := command.New(cmd.Flags()).
			ParseYaml(confFile).
			GenerateConfigFilesFromDir(db2InstallDir).
			RunBashScripts().
			DeleteStateForDir(db2CleanupDir).
			TimeTaken().
			Error

		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	installCmd.AddCommand(installDB2Cmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installDB2Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installDB2Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
