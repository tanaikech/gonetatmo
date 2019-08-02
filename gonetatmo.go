// Package main (gonetatmo.go) :
package main

import (
	"os"

	"github.com/urfave/cli"
)

const (
	appname = "gonetatmo"
)

// createHelp : Create help document.
func createHelp() *cli.App {
	a := cli.NewApp()
	a.Name = appname
	a.Author = "tanaike [ https://github.com/tanaikech/" + appname + " ] "
	a.Email = "tanaike@hotmail.com"
	a.Usage = "Retrieve values from own Netatmo."
	a.Version = "1.0.1"
	a.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "raw",
			Usage: "Display raw data which retrieved from Netatmo. At default, simple data is displayed.",
		},
		cli.BoolFlag{
			Name:  "json, j",
			Usage: "Output data which retrieved from Netatmo as json. At default, simple data is displayed.",
		},
		cli.StringFlag{
			Name:  "clientid",
			Usage: "Client ID for accessing to Netatmo. You can retrieve this by registered your application at https://dev.netatmo.com/",
		},
		cli.StringFlag{
			Name:  "clientsecret",
			Usage: "Client secret for accessing to Netatmo. You can retrieve this by registered your application at https://dev.netatmo.com/",
		},
		cli.StringFlag{
			Name:  "email",
			Usage: "E-mail that you use when you login to Netatmo. This is not saved to the config file.",
		},
		cli.StringFlag{
			Name:  "password",
			Usage: "Password that you use when you login to Netatmo. This is not saved to the config file.",
		},
		cli.StringFlag{
			Name:  "googleapikey, key",
			Usage: "API key for using Google Map.",
		},
	}
	a.Commands = []cli.Command{
		{
			Name:        "getmeasure",
			Aliases:     []string{"m"},
			Usage:       "-di \"12:34:56:78:90:12\" -b 2018-01-23T12:00:00+09:00 -e 2018-01-23T13:00:00+09:00",
			Description: "Retrieve values from device ID you have. Please read 'https://dev.netatmo.com/en-US/resources/technical/reference/common/getmeasure'.",
			Action:      handler,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "deviceid, di",
					Usage: "Mac address of the device.",
				},
				cli.StringFlag{
					Name:  "moduleid, mi",
					Usage: "Mac address of the module.",
				},
				cli.StringFlag{
					Name:  "scale, sc",
					Usage: "Timelapse between two measurements. You can select from 30min, 1hour, 3hours, 1day, 1week and 1month. max means all values.",
					Value: "max",
				},
				cli.StringFlag{
					Name:  "type, ty",
					Usage: "Category of data you want. About the detail, please the URL of description.",
					Value: "Temperature,Humidity",
				},
				cli.StringFlag{
					Name:  "datebegin, b",
					Usage: "Timestamp (ISO8601) of the first measure to retrieve. For example, 'YYYY-MM-DDThh:mm:dd+00:00'.",
				},
				cli.StringFlag{
					Name:  "dateend, e",
					Usage: "Timestamp (ISO8601) of the last measure to retrieve. For example, 'YYYY-MM-DDThh:mm:dd+00:00'.",
				},
				cli.StringFlag{
					Name:  "limit, l",
					Usage: "Maximum number of measurements (default and max are 1024)",
					Value: "1024",
				},
			},
		},
		{
			Name:        "getpublicdata",
			Aliases:     []string{"p"},
			Usage:       "-a \"tokyo station\"",
			Description: "Retrieve all netatmo's values and average values for area from area information. Please read 'https://dev.netatmo.com/en-US/resources/technical/reference/weatherapi/getpublicdata'.",
			Action:      handler,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "address, a",
					Usage: "Input place name, postal code and address of spot you want to retrieve. e.g. `-a tokyo station`",
				},
				cli.Float64Flag{
					Name:  "latitude, lat",
					Usage: "Input center latitude of spot you want to retrieve.",
				},
				cli.Float64Flag{
					Name:  "longitude, lon",
					Usage: "Input center longitude of spot you want to retrieve.",
				},
				cli.Float64Flag{
					Name:  "range, r",
					Usage: "Input range of area. Unit is kilometers. Default is a square area 10 kilometers on a side.",
					Value: 10,
				},
				cli.StringFlag{
					Name:  "requireddata, re",
					Usage: "To filter stations based on relevant measurements you want (e.g. rain will only return stations with rain gauges). Default is no filter.",
				},
				cli.StringFlag{
					Name:  "filter, f",
					Usage: "True to exclude station with abnormal temperature measures.",
					Value: "false",
				},
				cli.StringFlag{
					Name:  "language, lng",
					Usage: "Language for Google MAP. (ISO 639-1)",
					Value: "en",
				},
				cli.StringFlag{
					Name:  "type, t",
					Usage: "Data you want to display. If you want temperature and pressure, please input temperature and pressure. This cannot be used for the option 'raw'.",
					Value: "temperature,pressure,humidity,rain,wind",
				},
				cli.BoolFlag{
					Name:  "raw",
					Usage: "Display raw data which retrieved from Netatmo. At default, simple data is displayed.",
				},
				cli.BoolFlag{
					Name:  "json, j",
					Usage: "Output as json data. Default is data for displaying to terminal.",
				},
			},
		},
	}
	return a
}

// main : Main of this script
func main() {
	a := createHelp()
	a.Action = handler
	a.Run(os.Args)
}
