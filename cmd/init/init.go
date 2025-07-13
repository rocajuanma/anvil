package init

import (
	"fmt"

	"github.com/rocajuanma/anvil/pkg/constants"
	"github.com/spf13/cobra"
)

var PullCmd = &cobra.Command{
	Use:   "init",
	Short: "Initiatizes Anvil CLI tool",
	Long:  constants.INIT_LONG_DESCRIPTION,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("init called")
	},
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
