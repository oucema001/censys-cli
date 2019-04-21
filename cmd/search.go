package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/mohae/struct2csv"
	"github.com/oucema001/censys-cli/util"
	"github.com/oucema001/censys-go/censys"
	"github.com/spf13/cobra"
)

var typeSearch string
var query string
var output string
var csv bool

func init() {
	searchCmd.Flags().StringVarP(&typeSearch, "type", "t", "website", "type of search can be ipv4 certificate or website")
	searchCmd.Flags().StringVarP(&query, "query", "q", "", "query to search for")
	searchCmd.Flags().StringVarP(&output, "output", "o", "search.csv", "file location to output file to")
	searchCmd.Flags().BoolVarP(&csv, "csv", "c", true, "csv")
	rootCmd.AddCommand(searchCmd)
}

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for a query",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := util.NewClient()
		util.Panic(err)
		querySearch := &censys.SearchQuery{
			Query:   query,
			Page:    2,
			Fields:  []string{},
			Flatten: true,
		}
		var ty = censys.WEBSITESSEARCH
		switch typeSearch {
		case "website":
			ty = censys.WEBSITESSEARCH
		case "ipv4":
			ty = censys.IPV4SEARCH
		case "certificate":
			ty = censys.CERTIFICATESSEARCH
		}

		if typeSearch == "" {
			fmt.Println("You have to chose a search type")
		}
		search, err := c.Search(context.Background(), querySearch, ty)
		util.Panic(err)
		s := util.PrettyPrint(search)
		fmt.Println(s)
		if csv {
			var search censys.Search
			err := json.Unmarshal([]byte(s), &search)
			util.Panic(err)
			exportCsv(search, output)
		}
	},
}

func exportCsv(c censys.Search, output string) {
	results := c.Results
	csvDataFile, err := os.Create(output)
	util.Panic(err)
	w := struct2csv.NewWriter(csvDataFile)

	w.UseCRLF()
	err = w.WriteStructs(results)
	if err != nil {
		os.Exit(1)
	}
	w.Flush()
}
