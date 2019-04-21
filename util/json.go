package util

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"

	"github.com/mohae/struct2csv"
	"github.com/oucema001/censys-go/censys"
)

//Keys struct to hold the keys
type Keys struct {
	ID     string `json:"app_id"`
	Secret string `json:"app_secret"`
}

//EncodeKeystoFile encodes app secrets to json file
func EncodeKeystoFile(appID string, appSecret string) error {
	var err error
	keys := Keys{
		ID:     appID,
		Secret: appSecret,
	}
	a, err := json.Marshal(keys)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(".keys.json", a, 0644)
	if err != nil {
		return err
	}
	return err
}

func decodeKeysFile(fileName string) (*Keys, error) {
	var keys Keys
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	dec.Decode(&keys)
	return &keys, nil
}

//NewClient returns new censys Client
func NewClient() (*censys.Client, error) {
	keys, err := decodeKeysFile(".keys.json")
	if err != nil {
		return nil, err
	}
	c := censys.NewClient(nil, keys.ID, keys.Secret)
	return c, nil
}

//Panic prints the error and exits
func Panic(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

//PrettyPrint returns a pretty pring of a struct
func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

type results struct {
	//IPV4 results
	IP                   string   `json:"ip,omitempty,omitempty"`
	Protocols            []string `json:"protocols,omitempty"`
	Country              string   `json:"location.country,omitempty"`
	RegisteredCountry    string   `json:"location.registered_country,omitempty"`
	Longitude            float64  `json:"location.longitude,omitempty"`
	Latitude             float64  `json:"location.latitude,omitempty"`
	City                 string   `json:"location.city,omitempty"`
	RegisteredCountyCode string   `json:"location.registered_country_code,omitempty"`
	CountryCode          string   `json:"location.country_code,omitempty"`
	Province             string   `json:"location.province,omitempty"`
	PostalCode           string   `json:"location.postal_code,omitempty"`
	TimeZone             string   `json:"location.timezone,omitempty"`
	Continent            string   `json:"location.continent,omitempty"`
	//Certificate Results
	FingerprintSha256 string `json:"parsed.fingerprint_sha256,omitempty"`
	SubjectDn         string `json:"parsed.subject_dn,omitempty"`
	IssuerDn          string `json:"parsed.issuer_dn,omitempty"`
	//Website Results
	Domain    string `json:"domain,omitempty"`
	AlexaRank int    `json:"alexa_rank,omitempty"`
}

//ExportSearchCsv export search results to csv file
func ExportSearchCsv(s string, output string) {
	var j censys.Search
	err := json.Unmarshal([]byte(s), &j)
	Panic(err)

	csvDataFile, err := os.Create(output)
	Panic(err)
	b := j.Results
	defer csvDataFile.Close()
	writer := csv.NewWriter(csvDataFile)
	var v []string
	writeSearchHeaders(writer, v)
	for _, a := range b {
		var s []string
		s = append(s, strconv.Itoa(a.AlexaRank))
		s = append(s, a.City)
		s = append(s, a.Continent)
		s = append(s, a.Country)
		s = append(s, a.CountryCode)
		s = append(s, a.Domain)
		s = append(s, a.FingerprintSha256)
		s = append(s, a.IP)
		s = append(s, a.PostalCode)
		s = append(s, a.Province)
		s = append(s, a.TimeZone)
		s = append(s, a.SubjectDn)
		for _, p := range a.Protocols {
			s = append(s, p)
		}
		f := strconv.FormatFloat(a.Latitude, 'f', 6, 64)
		s = append(s, f)

		//s = append(s)
		writer.Write(s)
	}
	writer.Flush()
}

func writeSearchHeaders(w *csv.Writer, v []string) {
	v = append(v, "AlexaRank")
	v = append(v, "City")
	v = append(v, "Continent")
	v = append(v, "Country")
	v = append(v, "CountryCode")
	v = append(v, "Domain")
	v = append(v, "FingerprintSha256")
	v = append(v, "IP")
	v = append(v, "PostalCode")
	v = append(v, "Province")
	v = append(v, "TimeZone")
	v = append(v, "SubjectDn")
	v = append(v, "Protocols")
	v = append(v, "Latitude")
	w.Write(v)
}

func export(i interface{}, writer *csv.Writer) {
	fields := reflect.TypeOf(i)
	fmt.Println(fields)
	//value := reflect.New(fields).Interface()
	values := reflect.ValueOf(i)
	fmt.Println(values)

	//num := fields.NumField()
	for i := 0; i < fields.Len(); i++ {
		field := fields.Field(i)
		value := values.Field(i)
		fmt.Print("Type:", field.Type, ",", field.Name, "=", value, "\n")
		var s []string
		switch value.Kind() {
		case reflect.String:
			{
				v := value.String()
				s = append(s, v)
			}
		case reflect.Int:
			{
				v := strconv.FormatInt(value.Int(), 10)
				s = append(s, v)
			}
		case reflect.Float64:
			{
				v := fmt.Sprintf("%f", value.Float())
				s = append(s, v)
			}
		}
		writer.Write(s)
	}
}

var j censys.Search

func ExportCsv(in interface{}) {
	//var j censys.Search
	//err := json.Unmarshal([]byte(s), &j)
	//Panic(err)
	//data := j.Results
	buff := &bytes.Buffer{}

	w := struct2csv.NewWriter(buff)
	w.UseCRLF()
	err := w.WriteStructs(in)
	if err != nil {
		// handle error
	}
	w.Flush()
	fmt.Println(buff)
}
