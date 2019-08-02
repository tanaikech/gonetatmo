// Package main (parsedisp.go) :
package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	mestimeThreshold = 3600 // [second] Threshold for time of measure cycle.
)

// getstationsdataSt : Structure of getstationsdata.
type getstationsdataStForParse struct {
	Body struct {
		Devices []interface{} `json:"devices"`
		User    struct {
			Mail string `json:"mail"`
		} `json:"user"`
	} `json:"body"`
	Status string `json:"status"`
}

// stations : For detail version.
type stations struct {
	Stations []stationsdataForOutput `json:"stations,omitempty"`
}

// stationsdataForOutput : For detail version.
type stationsdataForOutput struct {
	Inside  []insideData  `json:"insideData,omitempty"`
	Outside []outsideData `json:"outsideData,omitempty"`
}

// insideData : Structure for data of inside device.
type insideData struct {
	Id               string  `json:"id,omitempty"`
	StationName      string  `json:"station_name,omitempty"`
	TimeUtc          int64   `json:"time_utc,omitempty"`
	MesTime          string  `json:"Measurement_time,omitempty"`
	AbsolutePressure float64 `json:"AbsolutePressure,omitempty"`
	Noise            int     `json:"Noise,omitempty"`
	Temperature      float64 `json:"Temperature,omitempty"`
	TempTrend        string  `json:"temp_trend,omitempty"`
	Humidity         float64 `json:"Humidity,omitempty"`
	Pressure         float64 `json:"Pressure,omitempty"`
	PressureTrend    string  `json:"pressure_trend,omitempty"`
	CO2              int     `json:"CO2,omitempty"`
	DateMaxTemp      int64   `json:"date_max_temp,omitempty"`
	DateMinTemp      int64   `json:"date_min_temp,omitempty"`
	MinTemp          float64 `json:"min_temp,omitempty"`
	MaxTemp          float64 `json:"max_temp,omitempty"`
	WifiStatus       int     `json:"wifi_status,omitempty"`
	FirmWare         int     `json:"firmware,omitempty"`
}

// outsideData : Structure for data of outside device.
type outsideData struct {
	Id             string  `json:"id,omitempty"`
	ModuleName     string  `json:"module_name,omitempty"`
	TimeUtc        int64   `json:"time_utc,omitempty"`
	MesTime        string  `json:"Measurement_time,omitempty"`
	Temperature    float64 `json:"Temperature,omitempty"`
	TempTrend      string  `json:"temp_trend,omitempty"`
	Humidity       float64 `json:"Humidity,omitempty"`
	DateMaxTemp    int64   `json:"date_max_temp,omitempty"`
	DateMinTemp    int64   `json:"date_min_temp,omitempty"`
	MinTemp        float64 `json:"min_temp,omitempty"`
	MaxTemp        float64 `json:"max_temp,omitempty"`
	BatteryVp      int     `json:"battery_vp,omitempty"`
	BatteryPercent int     `json:"battery_percent,omitempty"`
	RfStatus       int     `json:"rf_status,omitempty"`
	FirmWare       int     `json:"firmware,omitempty"`
}

// publicData : Structure for public data.
type publicData struct {
	Body []struct {
		ID    string `json:"_id"`
		Place struct {
			Location []float64 `json:"location"` // [0]longitude, [1]latitude
			Altitude int       `json:"altitude"`
			Timezone string    `json:"timezone"`
		} `json:"place"`
		Mark        int                    `json:"mark"`
		Measures    map[string]interface{} `json:"measures"`
		Modules     []string               `json:"modules"`
		ModuleTypes map[string]interface{} `json:"module_types"`
	} `json:"body"`
	Status     string  `json:"status"`
	TimeExec   float64 `json:"time_exec"`
	TimeServer int     `json:"time_server"`
}

// createOutputFormatForgetPublicData : Create output format from results for getPublicData.
func createOutputFormatForgetPublicData(rrr []map[string]interface{}, data [][]string) ([]string, [][]string) {
	header := []string{"", "average", "number"}
	for _, e := range rrr {
		temp := make([]string, 3)
		for j, f := range e {
			if strings.Index(j, "_c") == -1 {
				if val, ok := f.(float64); ok {
					strconv.FormatFloat(val, 'f', 4, 64)
					temp[0] = j
					temp[1] = strconv.FormatFloat(val, 'f', 2, 64)
				}
			} else {
				if val, ok := f.(float64); ok {
					temp[2] = strconv.FormatInt(int64(val), 10)
				}
			}
		}
		data = append(data, temp)
	}
	return header, data
}

// calcAverage : Calculate average values from retrieved public data.
func calcAverage(sv []string, res []map[string]interface{}) []map[string]interface{} {
	rrr := []map[string]interface{}{}
	chk := 0
	for _, s := range sv {
		rr := map[string]interface{}{}
		var total float64
		var cn float64
		for _, e := range res {
			if tt, ok := e[s].(float64); ok {
				total += tt
				cn += 1
			}
		}
		rr[s] = func(a float64, b int) float64 {
			s := math.Pow(10, float64(b))
			r := math.Floor(a*s+.5) / s
			if math.IsNaN(r) {
				chk += 1
			}
			return r
		}(total/cn, 2)
		rr[s+"_c"] = cn
		rrr = append(rrr, rr)
	}
	if chk == len(sv) {
		fmt.Printf("## Data was not returned from Netatmo. Please try again.\n")
	}
	return rrr
}

// setSearchValues : Set search values.
func setSearchValues(search []string) []string {
	var sv []string
	for _, s := range search {
		if s == "rain" {
			sv = append(sv, []string{"rain_60min", "rain_24h"}...)
		} else if s == "wind" {
			sv = append(sv, []string{"wind_strength", "gust_strength"}...)
		} else {
			sv = append(sv, s)
		}
	}
	return sv
}

// parsePublicdata : Parse retrieved public data.
func parsePublicdata(search []string, data []byte) []map[string]interface{} {
	nt := time.Now().Unix()
	pb := &publicData{}
	json.Unmarshal(data, &pb)
	res := []map[string]interface{}{}
	for _, e := range pb.Body {
		t1 := map[string]interface{}{}
		for _, f := range e.Measures {
			if v, ok := f.(map[string]interface{})["res"]; ok {
				for k, g := range v.(map[string]interface{}) {
					i64, err := strconv.ParseInt(k, 10, 64)
					if err != nil {
						fmt.Fprintf(os.Stderr, "Error: %v\n", err)
						os.Exit(1)
					}
					if nt-i64 < mestimeThreshold {
						for l, h := range g.([]interface{}) {
							for _, s := range search {
								if f.(map[string]interface{})["type"].([]interface{})[l].(string) == s {
									t1[s] = h
								}
							}
						}
					}
				}
			}
			for _, s := range search {
				if s == "rain" || s == "wind" {
					if v, ok := f.(map[string]interface{})[s+"_timeutc"]; ok {
						if nt-int64(v.(float64)) < mestimeThreshold {
							for k, g := range f.(map[string]interface{}) {
								if k == s+"_timeutc" {
									continue
								}
								t1[k] = g
							}
						}
					}
				}
			}
		}
		res = append(res, t1)
	}
	return res
}

// transpose : Transpose slice. Slice of (n x m) to (m x n).
func transpose(slice [][]string) [][]string {
	xl := len(slice[0])
	yl := len(slice)
	result := make([][]string, xl)
	for i := range result {
		result[i] = make([]string, yl)
	}
	for i := 0; i < xl; i++ {
		for j := 0; j < yl; j++ {
			result[i][j] = slice[j][i]
		}
	}
	return result
}

// createOutputFormatForgetStationsData : Create output format from results for getStationsData.
func (sim *stations) createOutputFormatForgetStationsData(i int, e stationsdataForOutput, data [][]string) ([]string, [][]string) {
	header := []string{""}
	col1 := []string{
		"ID",
		"Status",
		"Measurement time",
		"Temperature [C]",
		"Temperature trend",
		"Humidity [%]",
		"Pressure [hPa]",
		"Pressure trend",
		"CO2 [ppm]",
		"Noise [dB]",
		"WifiStatus",
		"Battery [%]",
		"Firmware",
	}
	data = append(data, col1)
	for j, f := range e.Inside {
		header = append(header, "in")
		date := time.Unix(f.TimeUtc, 0)
		out := date.In(time.Local).Format("20060102 15:04:05 MST")
		sim.Stations[i].Inside[j].MesTime = out
		sim.Stations[i].Inside[j].TimeUtc = 0
		status := func(t int64) string {
			if time.Now().Unix()-t > mestimeThreshold {
				return "Not working!"
			}
			return "Working."
		}(f.TimeUtc)
		temp := []string{
			f.Id,
			status,
			out,
			strconv.FormatFloat(f.Temperature, 'f', 1, 64),
			f.TempTrend,
			strconv.FormatFloat(f.Humidity, 'f', 1, 64),
			strconv.FormatFloat(f.Pressure, 'f', 1, 64),
			f.PressureTrend,
			strconv.Itoa(f.CO2),
			strconv.Itoa(f.Noise),
			strconv.Itoa(f.WifiStatus),
			"",
			strconv.Itoa(f.FirmWare),
		}
		data = append(data, temp)
	}
	for j, f := range e.Outside {
		header = append(header, "out")
		date := time.Unix(f.TimeUtc, 0)
		out := date.In(time.Local).Format("20060102 15:04:05 MST")
		sim.Stations[i].Outside[j].MesTime = out
		sim.Stations[i].Outside[j].TimeUtc = 0
		status := func(t int64) string {
			if time.Now().Unix()-t > mestimeThreshold {
				return "Not working!"
			}
			return "Working."
		}(f.TimeUtc)
		temp := []string{
			f.Id,
			status,
			out,
			strconv.FormatFloat(f.Temperature, 'f', 1, 64),
			f.TempTrend,
			strconv.FormatFloat(f.Humidity, 'f', 1, 64),
			"",
			"",
			"",
			"",
			strconv.Itoa(f.RfStatus),
			strconv.Itoa(f.BatteryPercent),
			strconv.Itoa(f.FirmWare),
		}
		data = append(data, temp)
	}
	return header, transpose(data)
}

// getInsideData : Retrieve data from inside devices.
func (so *stationsdataForOutput) getInsideData(e interface{}) {
	inDat := &insideData{}
	s := reflect.ValueOf(inDat).Elem()
	typeOfT := s.Type()
	inDat.Id = e.(map[string]interface{})["_id"].(string)
	inDat.WifiStatus = int(e.(map[string]interface{})["wifi_status"].(float64))
	inDat.FirmWare = int(e.(map[string]interface{})["firmware"].(float64))
	for i := 0; i < s.NumField(); i++ {
		for j, f := range e.(map[string]interface{})["dashboard_data"].(map[string]interface{}) {
			if strings.Split(typeOfT.Field(i).Tag.Get("json"), ",")[0] == j {
				fl := s.FieldByName(typeOfT.Field(i).Name)
				switch fl.Kind() {
				case reflect.Int, reflect.Int64:
					c, _ := f.(float64)
					fl.SetInt(int64(c))
				case reflect.Float64:
					f64, _ := f.(float64)
					fl.SetFloat(f64)
				case reflect.String:
					fl.SetString(f.(string))
				}
			}
			if inDat.TimeUtc > 0 && inDat.MesTime == "" {
				date := time.Unix(inDat.TimeUtc, 0)
				inDat.MesTime = date.In(time.Local).Format("20060102_15:04:05_MST")
			}
		}
	}
	so.Inside = append(so.Inside, *inDat)
}

// getOutsideData : Retrieve data from outside devices.
func (so *stationsdataForOutput) getOutsideData(e interface{}) {
	for _, f := range e.(map[string]interface{})["modules"].([]interface{}) {
		otDat := &outsideData{}
		s := reflect.ValueOf(otDat).Elem()
		typeOfT := s.Type()
		otDat.Id = f.(map[string]interface{})["_id"].(string)
		otDat.RfStatus = int(f.(map[string]interface{})["rf_status"].(float64))
		otDat.FirmWare = int(f.(map[string]interface{})["firmware"].(float64))
		otDat.BatteryPercent = int(f.(map[string]interface{})["battery_percent"].(float64))
		otDat.BatteryVp = int(f.(map[string]interface{})["battery_vp"].(float64))
		if f.(map[string]interface{})["reachable"].(bool) {
			for i := 0; i < s.NumField(); i++ {
				for j, f := range f.(map[string]interface{})["dashboard_data"].(map[string]interface{}) {
					if strings.Split(typeOfT.Field(i).Tag.Get("json"), ",")[0] == j {
						fl := s.FieldByName(typeOfT.Field(i).Name)
						switch fl.Kind() {
						case reflect.Int, reflect.Int64:
							c, _ := f.(float64)
							fl.SetInt(int64(c))
						case reflect.Float64:
							f64, _ := f.(float64)
							fl.SetFloat(f64)
						case reflect.String:
							fl.SetString(f.(string))
						}
					}
					if otDat.TimeUtc > 0 && otDat.MesTime == "" {
						date := time.Unix(otDat.TimeUtc, 0)
						otDat.MesTime = date.In(time.Local).Format("20060102_15:04:05_MST")
					}
				}
			}
		}
		so.Outside = append(so.Outside, *otDat)
	}
}

// parseStationsData : Parse stations data
func parseStationsData(res []byte) []byte {
	s := &stations{}
	rs := &getstationsdataStForParse{}
	json.Unmarshal(res, &rs)
	for _, e := range rs.Body.Devices {
		so := &stationsdataForOutput{}
		so.getInsideData(e)
		so.getOutsideData(e)
		s.Stations = append(s.Stations, *so)
	}
	si, err := json.Marshal(s)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	return si
}
