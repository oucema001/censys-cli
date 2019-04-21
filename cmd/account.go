package cmd

import (
	"context"
	"fmt"

	"github.com/oucema001/censys-cli/util"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(accountCmd)
}

var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "Prints the details about the account you are using to interact with censys",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := util.NewClient()
		util.Panic(err)
		profile, err := c.GetProfile(context.Background())
		util.Panic(err)
		pr := util.PrettyPrint(profile)
		fmt.Println(pr)
	},
}
