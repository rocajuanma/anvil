package draw

import (
	"fmt"

	"github.com/rocajuanma/anvil/pkg/constants"
	"github.com/rocajuanma/anvil/pkg/figure"
	"github.com/spf13/cobra"
)

var DrawCmd = &cobra.Command{
	Use:   "draw",
	Short: "Uses go-figure to generate ASCII text",
	Long:  constants.DRAW_COMMAND_LONG_DESCRIPTION,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("draw called")
		figure.Draw("anvil", args[0])
	},
}

func init() {}
