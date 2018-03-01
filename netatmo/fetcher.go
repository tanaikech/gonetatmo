// Package netatmo (fetcher.go) :
package netatmo

import (
	"io"
	"net/http"
	"time"
)

// requestParams : Parameters for FetchAPI
type RequestParams struct {
	Method      string
	APIURL      string
	Data        io.Reader
	Contenttype string
	Accesstoken string
	Dtime       int64
}

// fetch : Fetch data from Google Drive
func (r *RequestParams) fetch() (*http.Response, error) {
	req, err := http.NewRequest(r.Method, r.APIURL, r.Data)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", r.Contenttype)
	client := &http.Client{Timeout: time.Duration(r.Dtime) * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
