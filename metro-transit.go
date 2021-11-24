package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const api_url string = "https://svc.metrotransit.org/NexTrip/"


// helper function to get url and return json matching go interface/struct

func getJson(url string, target interface{}) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return json.NewDecoder(res.Body).Decode(target)
}


type ProviderMap map[int]string

func GetProviders() ProviderMap {
	var providers []struct {
		Name string `json:"Text"`
		Value int `json:"Value,string"`
	}
	err := getJson(api_url+"Providers?format=json", &providers)

	if err != nil {
		fmt.Print("got an error: ", err)
	}
	
	var m ProviderMap = make(ProviderMap)

	for _, p := range providers {
		m[p.Value] = p.Name
	}

	return m
}

type RouteID int
// Routes endpoint
type Route struct {
	Description string `json:"Description"`
	Provider int `json:"ProviderID,string"`
	Id RouteID `json:"Route,string"`
}

func GetRoutes() []Route {
	var routes []Route
	err := getJson(api_url + "Routes?format=json", &routes)
	if err != nil {
		fmt.Print("got an error: ", err)
	}
	return routes
}

type Stop struct {
	Name string `json:"Text"`
	Id string `json:"Value"`
}

func GetStops(route RouteID, direction int) []Stop {
	var stops []Stop
	err := getJson(api_url + fmt.Sprintf("Stops/%d/%d?format=json", route, direction), &stops)
	if err != nil {
		fmt.Print("got an error: ", err)
	}
	return stops
}

// for depatures, we want to parse the strange string they used and turn it into a time.
// we use a custom type

type CustomTime struct {
	time.Time
}

func (ct *CustomTime) UnmarshalJSON(b []byte) (err error) {
	str := strings.Trim(string(b), "\"\\/Date()") // now its 12354678-0500
	str = strings.Split(str, "-")[0]
	if str == "null" {
		ct.Time = time.Time{}
		return
	}
	millis, err := strconv.Atoi(str)
	ct.Time = time.UnixMilli(int64(millis))
	return
}

type Departure struct {
	Actual bool
	DepartureTime CustomTime
}

func GetDepartures(route RouteID, direction int, stopID string) []Departure {
	var deps []Departure
	err := getJson(api_url + fmt.Sprintf("%d/%d/%s?format=json", route, direction, stopID), &deps)

	if err != nil {
		fmt.Print("got an error: ", err)
	}
	return deps
}

func main() {
	providers := GetProviders()

	for id, name := range(providers) {
		fmt.Printf("Provider Name: %s\tProvider Value: %d\n",name, id)
	}

	routes := GetRoutes()

	fmt.Printf("%#v\n", routes[0])

	fmt.Printf("%#v", GetDepartures(902, 1, "EABK"))
}
