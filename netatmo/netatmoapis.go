// Package netatmo (callapis.go) :
package netatmo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/urfave/cli"
)

const (
	netatmoApi = "https://api.netatmo.com/"
)

// chkResErr : Check response error.
func chkResErr(r []byte) bool {
	var rs map[string]interface{}
	json.Unmarshal(r, &rs)
	if _, ok := rs["error"]; ok {
		return true
	}
	return false
}

// GetTokens : Retrieve tokens.
func GetTokens(val url.Values) ([]byte, error) {
	var err error
	r := &RequestParams{
		Method:      "POST",
		APIURL:      netatmoApi + "oauth2/token",
		Data:        strings.NewReader(val.Encode()),
		Contenttype: "application/x-www-form-urlencoded",
		Dtime:       30,
	}
	res, err := r.fetch()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%v\n%v", err, res))
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%v\n%v", err, string(body)))
	}
	if chkResErr(body) {
		return nil, errors.New(fmt.Sprintf("Error: Refresh token may be broken or revoked. Please retrive it again. Please run with the parameters of client id, client secret, mail address and password for Netatmo, again."))
	}
	defer res.Body.Close()
	return body, nil
}

// getNetatmoValues : Retrieve values from Netatmo APIs.
func (r *RequestParams) getNetatmoValues() ([]byte, error) {
	var err error
	res, err := r.fetch()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%v\n%v", err, res))
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%v\n%v", err, string(body)))
	}
	if chkResErr(body) {
		return nil, errors.New(fmt.Sprintf("%v", string(body)))
	}
	return body, nil
}

// callNetatmoApis : Call Netatmo's APIs.
func callNetatmoApis(url string) ([]byte, error) {
	r := &RequestParams{
		Method:      "GET",
		APIURL:      url,
		Data:        nil,
		Contenttype: "application/x-www-form-urlencoded",
		Dtime:       60,
	}
	body, err := r.getNetatmoValues()
	return body, err
}

// Getpublicdata : https://dev.netatmo.com/en-US/resources/technical/reference/weatherapi/getpublicdata
func Getpublicdata(c *cli.Context, accesstoken string, coordinates []float64) ([]byte, error) {
	tokenparams := url.Values{}
	tokenparams.Set("access_token", accesstoken)
	tokenparams.Set("lat_ne", strconv.FormatFloat(coordinates[0], 'f', 15, 64))
	tokenparams.Set("lon_ne", strconv.FormatFloat(coordinates[1], 'f', 15, 64))
	tokenparams.Set("lat_sw", strconv.FormatFloat(coordinates[2], 'f', 15, 64))
	tokenparams.Set("lon_sw", strconv.FormatFloat(coordinates[3], 'f', 15, 64))
	url := netatmoApi + "api/getpublicdata?" + tokenparams.Encode()
	return callNetatmoApis(url)
}

// Getmeasure : https://dev.netatmo.com/en-US/resources/technical/reference/common/getmeasure
func Getmeasure(c *cli.Context, accesstoken string) ([]byte, error) {
	datebegin, err := time.Parse(time.RFC3339Nano, c.String("datebegin"))
	if err != nil {
		return nil, err
	}
	dateend, err := time.Parse(time.RFC3339Nano, c.String("dateend"))
	if err != nil {
		return nil, err
	}
	tokenparams := url.Values{}
	tokenparams.Set("access_token", accesstoken)
	tokenparams.Set("device_id", c.String("deviceid"))
	tokenparams.Set("module_id", c.String("moduleid"))
	tokenparams.Set("type", c.String("type"))
	tokenparams.Set("scale", c.String("scale"))
	tokenparams.Set("date_begin", strconv.FormatInt(datebegin.Unix(), 10))
	tokenparams.Set("date_end", strconv.FormatInt(dateend.Unix(), 10))
	tokenparams.Set("limit", c.String("limit"))
	tokenparams.Set("real_time", func(scale string) string {
		if scale != "max" {
			return "true"
		} else {
			return "false"
		}
	}(c.String("scale")))
	url := netatmoApi + "api/getmeasure?" + tokenparams.Encode()
	return callNetatmoApis(url)
}

// GetStationsData : https://dev.netatmo.com/resources/technical/reference/weatherstation/getstationsdata
func GetStationsData(accesstoken string) ([]byte, error) {
	url := netatmoApi + "api/getstationsdata?access_token=" + accesstoken
	return callNetatmoApis(url)
}
