package cmd

import (
	"github.com/IBM/db2ctl/internal/command"
	"github.com/spf13/cobra"
)

// runLinbitCmd represents the linbit command
var runLinbitCmd = &cobra.Command{
	Use:   "linbit",
	Short: "Install linbit",
	Long: `Install linbit. For example:

	db2ctl install linbit -c db2ctl-sample.yaml`,
	RunE: func(cmd *cobra.Command, args []string) error {

		err := command.New(cmd.Flags()).
			ParseYaml(confFile).
			GenerateBinPackingCSV().
			ReadBinPackingCSV().
			GenerateConfigFilesFromDir(linbitInstallDir).
			RunBashScripts().
			DeleteStateForDir(linbitCleanupDir).
			TimeTaken().
			Error

		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	installCmd.AddCommand(runLinbitCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runLinbitCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runLinbitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
