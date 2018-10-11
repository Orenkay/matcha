package localisation

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	googleAPIKey  = "AIzaSyBciG2mapnXT-z59x40gmE_cT_7W61Mb8M"
	ipstackAPIKey = "dba51c193cd0388325c3de2d35c6caa1"
)

type IPStackData struct {
	City    string  `json:"city"`
	Country string  `json:"country_name"`
	Lat     float64 `json:"latitude"`
	Lng     float64 `json:"longitude"`
}

type PlaceDataLocation struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type PlaceDataGeometry struct {
	Location PlaceDataLocation `json:"location"`
}

type PlaceDataCandidates struct {
	Address  string            `json:"formatted_address"`
	Geometry PlaceDataGeometry `json:"geometry"`
}

type PlaceData struct {
	Candidates []PlaceDataCandidates `json:"candidates"`
}

type PlaceByIDData struct {
	Result PlaceDataCandidates `json:"result"`
}

type Place struct {
	Lat     float64
	Lng     float64
	Address string
}

func PlaceByIP(ip string) (*Place, error) {
	ip = "92.169.74.153" // TODO: un-hardcode me please :'(
	resp, err := http.Get("http://api.ipstack.com/" + ip + "?access_key=" + ipstackAPIKey)
	if err != nil {
		return nil, err
	}

	d := &IPStackData{}
	err = json.NewDecoder(resp.Body).Decode(d)
	if err != nil {
		return nil, err
	}

	return &Place{
		Lat:     d.Lat,
		Lng:     d.Lng,
		Address: fmt.Sprintf("%s, %s", d.City, d.Country),
	}, nil
}

func PlaceByID(placeID string) (*Place, error) {
	resp, err := http.Get("https://maps.googleapis.com/maps/api/place/details/json?placeid=" + placeID + "&fields=formatted_address,geometry&key=" + googleAPIKey)
	if err != nil {
		return nil, err
	}

	d := &PlaceByIDData{}
	{
		err = json.NewDecoder(resp.Body).Decode(&d)
		if err != nil {
			return nil, err
		}
	}

	return &Place{
		Lat:     d.Result.Geometry.Location.Lat,
		Lng:     d.Result.Geometry.Location.Lng,
		Address: d.Result.Address,
	}, nil
}
