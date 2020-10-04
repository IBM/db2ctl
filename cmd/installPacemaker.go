package cmd

import (
	"github.com/IBM/db2ctl/internal/command"
	"github.com/spf13/cobra"
)

// installPacemaker represents the runPacemaker command
var installPacemaker = &cobra.Command{
	Use:   "pacemaker",
	Short: "Install pacemaker",
	Long: `Install pacemaker. For example:

	db2ctl install pacemaker -c db2ctl.yaml`,
	RunE: func(cmd *cobra.Command, args []string) error {

		err := command.New(cmd.Flags()).
			ParseYaml(confFile).
			GenerateBinPackingCSV().
			ReadBinPackingCSV().
			GenerateMappingFile().
			ReadMappingCSV().
			GenerateConfigFilesFromDir(pacemakerInstallDir).
			RunBashScripts().
			DeleteStateForDir(pacemakerCleanupDir).
			TimeTaken().
			Error

		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	installCmd.AddCommand(installPacemaker)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installPacemaker.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installPacemaker.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
