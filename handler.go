// Package main (handler.go) :
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/tanaikech/gonetatmo/netatmo"
	"github.com/urfave/cli"
)

// locationData : Structure for retrieved location data.
type locationData struct {
	Results []struct {
		AddressComponents []struct {
			LongName  string   `json:"long_name"`
			ShortName string   `json:"short_name"`
			Types     []string `json:"types"`
		} `json:"address_components"`
		FormattedAddress string `json:"formatted_address"`
		Geometry         struct {
			Bounds struct {
				Northeast struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"northeast"`
				Southwest struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"southwest"`
			} `json:"bounds"`
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
			LocationType string `json:"location_type"`
			Viewport     struct {
				Northeast struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"northeast"`
				Southwest struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"southwest"`
			} `json:"viewport"`
		} `json:"geometry"`
		PlaceID string   `json:"place_id"`
		Types   []string `json:"types"`
	} `json:"results"`
	Status string `json:"status"`
}

// dispTable : Display results using tablewriter.
func dispTable(header []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.SetBorder(false)
	table.SetAlignment(tablewriter.ALIGN_RIGHT)
	table.AppendBulk(data)
	table.Render()
}

// dispGetpublicdata : Display data retrieved by getpublicdata.
func dispGetpublicdata(inputtedTypes string, allData []byte, outRaw, outJson bool) {
	if outRaw {
		fmt.Println(string(allData))
		return
	}
	types := strings.Split(inputtedTypes, ",")
	for i, e := range types {
		types[i] = strings.TrimSpace(e)
	}
	pubdat := parsePublicdata(types, allData)
	if outJson {
		outjson, err := json.Marshal(pubdat)
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
		fmt.Println(string(outjson))
	} else {
		cA := calcAverage(setSearchValues(types), pubdat)
		var data [][]string
		var header []string
		header, data = createOutputFormatForgetPublicData(cA, data)
		dispTable(header, data)
	}
	return
}

// getpublicdata : https://dev.netatmo.com/en-US/resources/technical/reference/weatherapi/getpublicdata
func (m *materials) getpublicdata(c *cli.Context) {
	if m.configFile.GoogleApiKey == "" {
		fmt.Printf("Error: Please input your API key for using Google MAP API.\n\n $ gonetatmo --key ###\n\n")
		os.Exit(1)
	}
	if c.String("address") != "" && c.Float64("latitude") == 0 && c.Float64("longitude") == 0 {
		res, err := netatmo.Geocoding(
			m.configFile.GoogleApiKey,
			strings.Replace(c.String("address"), " ", "+", -1),
			c.String("language"),
		)
		if err != nil {
			fmt.Printf("%v, %v\n", err, res)
			os.Exit(1)
		}
		l := &locationData{}
		json.Unmarshal(res, &l)
		for _, e := range l.Results {
			coordinates, err := netatmo.GetCoordinates(c.Float64("range"), e.Geometry.Location.Lat, e.Geometry.Location.Lng, 10)
			if err != nil {
				fmt.Printf("%v, %v\n", err, coordinates)
				os.Exit(1)
			}
			allData, err := netatmo.Getpublicdata(c, m.configFile.tokens.Accesstoken, coordinates)
			if err != nil {
				fmt.Printf("%v, %v\n", err, coordinates)
				os.Exit(1)
			}
			if !c.Bool("raw") && !c.Bool("json") {
				h := []string{"Properties", "Values"}
				o := [][]string{
					[]string{"Time", m.para.pstart.In(time.Local).Format("20060102 15:04:05 MST")},
					[]string{"Formatted address", e.FormattedAddress},
					[]string{"Center(Latitude)", strconv.FormatFloat(e.Geometry.Location.Lat, 'f', 10, 64)},
					[]string{"Center(Longitude)", strconv.FormatFloat(e.Geometry.Location.Lng, 'f', 10, 64)},
					[]string{"North east corner(Latitude)", strconv.FormatFloat(coordinates[0], 'f', 10, 64)},
					[]string{"North east corner(Longitude)", strconv.FormatFloat(coordinates[1], 'f', 10, 64)},
					[]string{"South west corner(Latitude)", strconv.FormatFloat(coordinates[2], 'f', 10, 64)},
					[]string{"South west corner(Longitude)", strconv.FormatFloat(coordinates[3], 'f', 10, 64)},
				}
				dispTable(h, o)
				fmt.Printf("\n")
			}
			dispGetpublicdata(c.String("type"), allData, c.Bool("raw"), c.Bool("json"))
		}
	}
	if c.String("address") == "" && (c.Float64("latitude") != 0 || c.Float64("longitude") != 0) {
		coordinates, err := netatmo.GetCoordinates(c.Float64("range"), c.Float64("latitude"), c.Float64("longitude"), 10)
		if err != nil {
			fmt.Printf("%v, %v\n", err, coordinates)
			os.Exit(1)
		}
		allData, err := netatmo.Getpublicdata(c, m.configFile.tokens.Accesstoken, coordinates)
		if err != nil {
			fmt.Printf("%v, %v\n", err, allData)
			os.Exit(1)
		}
		dispGetpublicdata(c.String("type"), allData, c.Bool("raw"), c.Bool("json"))
	}
	return
}

// getmeasure : https://dev.netatmo.com/resources/technical/reference/common/getmeasure
func (m *materials) getmeasure(c *cli.Context) {
	allData, err := netatmo.Getmeasure(c, m.configFile.tokens.Accesstoken)
	if err != nil {
		fmt.Printf("%v\n%v\n", err, string(allData))
		os.Exit(1)
	}
	fmt.Println(string(allData))
	return
}

// getStationsData : https://dev.netatmo.com/resources/technical/reference/weatherstation/getstationsdata
func (m *materials) getStationsData(c *cli.Context) {
	allData, err := netatmo.GetStationsData(m.configFile.tokens.Accesstoken)
	if err != nil {
		fmt.Printf("%v\n%v\n", err, string(allData))
		os.Exit(1)
	}
	if c.Bool("raw") {
		fmt.Println(string(allData))
		return
	}
	sData := parseStationsData(allData)
	if c.Bool("json") {
		fmt.Println(string(sData))
	} else {
		od := &stations{}
		json.Unmarshal(sData, &od)
		var data [][]string
		var header []string
		for i, e := range od.Stations {
			header, data = od.createOutputFormatForgetStationsData(i, e, data)
		}
		dispTable(header, data)
	}
}

// handler : Initialize of "para".
func handler(c *cli.Context) {
	m := initParams()
	if err := m.chkCfg(c); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	switch c.Command.Names()[0] {
	case "getmeasure":
		m.getmeasure(c)
	case "getpublicdata":
		m.getpublicdata(c)
	default:
		m.getStationsData(c)
	}
	return
}
