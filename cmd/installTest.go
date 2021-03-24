package cmd

import (
	"github.com/IBM/db2ctl/internal/command"
	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run test script",
	Long: `Run test script
	
	db2ctl install test`,
	RunE: func(cmd *cobra.Command, args []string) error {

		const dirName = "test"

		err := command.New(cmd.Flags()).
			SetAutoYesEnabled(true).
			SetRunOnLocal(true).
			GenerateConfigFilesFromDir(dirName).
			RunBashScripts().
			TimeTaken().
			Error

		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	installCmd.AddCommand(testCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
