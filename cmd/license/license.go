package license

import (
	"github.com/spf13/cobra"
)

var matchLabels = map[string]string{"app": "discovery-engine"}
var port int64 = 9089

var LicenseCmd = &cobra.Command{
	Use:   "license",
	Short: "license is a pallete which contains license commands",
	Long:  `license is a pallete which contains license commands`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		return err
	},
}

func init() {

}
