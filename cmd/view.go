package cmd

import (
	"context"
	"encoding/json"
	"os"

	"github.com/landoop/tableprinter"

	"github.com/mohae/struct2csv"
	"github.com/oucema001/censys-cli/util"
	"github.com/oucema001/censys-go/censys"
	"github.com/spf13/cobra"
)

var typeView string
var queryView string
var outputView string
var csvView bool

func init() {
	viewCmd.Flags().StringVarP(&typeView, "type", "t", "website", "view type can be website ipv4 or certificate")
	viewCmd.Flags().StringVarP(&queryView, "query", "q", "", "query to view")
	viewCmd.Flags().StringVarP(&outputView, "output", "o", "view.csv", "file to output csv to")
	viewCmd.Flags().BoolVarP(&csvView, "csv", "c", true, "export view to csv")
	rootCmd.AddCommand(viewCmd)
}

var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "View returns a view of the result",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := util.NewClient()
		util.Panic(err)
		var ty = censys.WEBSITESVIEW
		switch typeView {
		case "website":
			ty = censys.WEBSITESVIEW
		case "ipv4":
			ty = censys.IPV4VIEW
		case "certificate":
			ty = censys.CERTIFICATESVIEW
		}
		v, err := c.GetView(context.Background(), ty, queryView)
		util.Panic(err)
		//fmt.Println(util.PrettyPrint(v))
		//tableprinter.Print(os.Stdout, v.Num80.HTTP.Get.Headers)
		s, _ := json.MarshalIndent(v.Num443.HTTPSWww.TLS.Certificate.Parsed.Subject, "", "\t")
		tableprinter.PrintJSON(os.Stdout, s)
		//	tableprinter.Print(os.Stdout, v.Num443.HTTPSWww.TLS.Certificate)
		exportViewCsv(v, outputView)
	},
}

func exportViewCsv(c *censys.View, output string) {
	csvDataFile, err := os.Create(output)
	util.Panic(err)
	w := struct2csv.NewWriter(csvDataFile)

	w.UseCRLF()
	err = w.WriteStructs(c)
	if err != nil {
		os.Exit(1)
	}
	w.Flush()
}
