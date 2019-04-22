package util

import (
	"github.com/oucema001/censys-go/censys"
	"github.com/olekukonko/tablewriter"
	"strconv"
	"os"
	"fmt"
)

//RenderSearchWebsite renders a table to view website results
func RenderSearchWebsite(res *censys.Search ,typesearch string){
	table := tablewriter.NewWriter(os.Stdout)
	switch typesearch {
	case "website":
		table.SetHeader([]string{"Domain", "Rank"})
		table.SetCaption(true, "Website Results.")
		for _, v := range res.Results {
		data:= []string{v.Domain,strconv.Itoa(v.AlexaRank)}
		table.Append(data)
		}
	case "certificate" :
		table.SetHeader([]string{"FingerPrint", "Issuer","Subject"})
		table.SetCaption(true, "Certificate Results.")
		for _, v := range res.Results {
			data:= []string{v.FingerprintSha256,v.IssuerDn,v.SubjectDn}
			table.Append(data)
		}
	case "ipv4":
		table.SetHeader([]string{"Ip", "Country","RegisteredCountry","Longitude","Latitude","City","RegisteredCountyCode","CountryCode",
		"Province","PostalCode","TimeZone","Continent"})
		table.SetCaption(true, "Ipv4 Results.")
		for _, v := range res.Results {
			data:= []string{v.IP,v.Country,v.RegisteredCountry,fmt.Sprintf("%f",v.Longitude),fmt.Sprintf("%f",v.Latitude),v.City,
			v.RegisteredCountyCode,
			v.CountryCode,
			v.Province,v.PostalCode,v.TimeZone,v.Continent}
			table.Append(data)
		}

	}
	table.Render() 
}
