/*
Package main (doc.go) :
This is a CLI tool to retrieve data from Netatmo.

# Features of this CLI tool is as follows.

1. Retrieves data from Netatmo.

2. Use Getstationsdata : https://dev.netatmo.com/resources/technical/reference/weatherstation/getstationsdata

3. Use Getmeasure : https://dev.netatmo.com/resources/technical/reference/common/getmeasure

4. Use Getpublicdata : https://dev.netatmo.com/en-US/resources/technical/reference/weatherapi/getpublicdata

5. For using Getpublicdata, you can retrieve data by inputting the address using Google Maps Geocoding API.

---------------------------------------------------------------

# Usage
Help

$ gonetatmo --help

Retrieving access token

$ gonetatmo --clientid ### --clientsecret ### --email ### --password ###

Input Google API key for using Google Maps Geocoding API

$ gonetatmo --key ###


Using Getstationsdata

$ gonetatmo

Use Getmeasure

$ gonetatmo m -di "12:34:56:78:90:12" -b "2018-01-23T12:00:00+09:00" -e "2018-01-23T13:00:00+09:00"

Use Getpublicdata

$ gonetatmo p -a "tokyo station"

$ gonetatmo p -lat 35.681167 -lon 139.767052

You can give the range of area using --range

$ gonetatmo p -a "tokyo station" --range 50

This means that data is retrieved from a square area 50 kilometers on a side with the center of "tokyo station".


You can see the detail information at https://github.com/tanaikech/gonetatmo

---------------------------------------------------------------
*/
package main
