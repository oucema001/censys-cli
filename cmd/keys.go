package cmd

import (
	"fmt"

	"github.com/oucema001/censys-cli/util"

	"github.com/spf13/cobra"
)

var appID string
var secret string

func init() {
	keysCmd.Flags().StringVarP(&appID, "appID", "a", "", "Application ID")
	keysCmd.Flags().StringVarP(&secret, "Secret", "s", "", "Application Secret")
	rootCmd.AddCommand(keysCmd)
}

var keysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Add Application Id and application secret to be used by the cli App",
	Run: func(cmd *cobra.Command, args []string) {
		err := util.EncodeKeystoFile(appID, secret)
		if err != nil {
			fmt.Println(err)
		}
	},
}
