package setup

import (
	"fmt"

	"github.com/rocajuanma/anvil/pkg/constants"
	"github.com/spf13/cobra"
)

var SetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Batch installation of applications using brew",
	Long:  constants.SETUP_COMMAND_LONG_DESCRIPTION,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("setup called")
	},
}

func init() {}
