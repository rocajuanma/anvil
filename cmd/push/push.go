package push

import (
	"fmt"

	"github.com/rocajuanma/anvil/pkg/constants"
	"github.com/spf13/cobra"
)

var PushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push assets to Github",
	Long:  constants.PUSH_COMMAND_LONG_DESCRIPTION,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("push called")
	},
}

func init() {
}
