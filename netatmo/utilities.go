// Package netatmo (utilities.go) :
package netatmo

import (
	"errors"
	"fmt"
	"math"
)

// hBase : This is for GetCoordinates.
func hBase(aLatY1, aLonX1, bLatY2, bLonX2 float64) float64 {
	a1LatY1 := aLatY1 * math.Pi / 180
	b1LatY2 := bLatY2 * math.Pi / 180
	a := 6378137.000000
	e := math.Sqrt((math.Pow(a, 2) - math.Pow(6356752.314245, 2)) / math.Pow(a, 2))
	muY := (a1LatY1 + b1LatY2) / 2
	W := math.Sqrt(1 - math.Pow(e, 2)*math.Pow(math.Sin(muY), 2))
	return math.Sqrt(math.Pow((a1LatY1-b1LatY2)*(a*(1-math.Pow(e, 2))/math.Pow(W, 3)), 2) + math.Pow(((aLonX1*math.Pi/180)-(bLonX2*math.Pi/180))*(a/W)*math.Cos(muY), 2))
}

// GetCoordinates : Caluculate coordinates from center latitude and longitude.
func GetCoordinates(oneSide, cLatY, cLonX float64, n int) ([]float64, error) {
	if cLatY < -90 || cLatY > 90 {
		return nil, errors.New(fmt.Sprintf("Error: Wrong latitude."))
	}
	if cLonX < -180 || cLonX > 180 {
		return nil, errors.New(fmt.Sprintf("Error: Wrong longitude."))
	}
	oneSide *= 500
	p1LatY := func(bq, cLatY1, cLonX1 float64, n int) float64 {
		cLatY2 := 0.0
		nStep := 0.01
		initY := cLatY1
		for rn := 1; rn <= n; rn++ {
			for cLatY2 = initY; bq-hBase(cLatY1, cLonX1, cLatY2, cLonX1) > 0; cLatY2 += nStep {
			}
			initY = cLatY2 - nStep
			nStep *= 0.1
		}
		return cLatY2
	}(oneSide, cLatY, cLonX, n)
	p1LonX := func(bq, cLatY1, cLonX1 float64, n int) float64 {
		cLonX2 := 0.0
		nStep := 0.01
		initX := cLonX1
		for rn := 1; rn <= n; rn++ {
			for cLonX2 = initX; bq-hBase(cLatY1, cLonX1, cLatY1, cLonX2) > 0; cLonX2 += nStep {
			}
			initX = cLonX2 - nStep
			nStep *= 0.1
		}
		return cLonX2
	}(oneSide, cLatY, cLonX, n)
	return []float64{p1LatY, p1LonX, cLatY - (p1LatY - cLatY), cLonX - (p1LonX - cLonX)}, nil
}
