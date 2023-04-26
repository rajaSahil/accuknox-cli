package get

import (
	"github.com/spf13/cobra"
)

var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "get is a pallete which contains get commands",
	Long:  `get is a pallete which contains get commands`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		return err
	},
}

func init() {

}
