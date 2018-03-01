// Package netatmo (googleapis.go) :
package netatmo

import (
	"errors"
	"fmt"
	"io/ioutil"
)

const (
	geocodingApi = "https://maps.googleapis.com/maps/api/geocode/json?key="
)

// getGoogleValues : Retrieve values from Google APIs.
func (r *RequestParams) getGoogleValues() ([]byte, error) {
	var err error
	res, err := r.fetch()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%v\n%v", err, res))
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%v\n%v", err, string(body)))
	}
	return body, nil
}

// callGoogleApis : Call Google's APIs.
func callGoogleApis(url string) ([]byte, error) {
	r := &RequestParams{
		Method:      "GET",
		APIURL:      url,
		Data:        nil,
		Contenttype: "application/x-www-form-urlencoded",
		Dtime:       30,
	}
	body, err := r.getNetatmoValues()
	return body, err
}

// Geocoding : https://developers.google.com/maps/documentation/geocoding/intro?hl=en
func Geocoding(key, address, lng string) ([]byte, error) {
	url := geocodingApi + key + "&address=" + address + "&language=" + lng
	return callGoogleApis(url)
}
