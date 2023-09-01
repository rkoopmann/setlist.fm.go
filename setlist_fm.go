package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type ArtistObject struct {
	MbId           string `json:"mbid"`
	TmId           *int   `json:"tmid,omitempty"`
	Name           string `json:"name"`
	SortName       string `json:"sortName"`
	Disambiguation string `json:"disambiguation"`
	Url            string `json:"url"`
}
type CoordsObject struct {
	Lat  float32 `json:"lat"`
	Long float32 `json:"long"`
}
type CountryObject struct {
	Code string `json:"code"`
	Name string `json:"name"`
}
type CityObject struct {
	Id        string        `json:"id"`
	Name      string        `json:"name"`
	State     string        `json:"state"`
	StateCode string        `json:"stateCode"`
	Coords    CoordsObject  `json:"coords"`
	Country   CountryObject `json:"country"`
}
type VenueObject struct {
	Id   string     `json:"id"`
	Name string     `json:"name"`
	City CityObject `json:"city"`
	Url  string     `json:"url"`
}
type TourObject struct {
	Name string `json:"name"`
}
type SongObject struct {
	Name  string        `json:"name"`
	With  *ArtistObject `json:"with,omitempty"`
	Cover *ArtistObject `json:"artist,omitempty"`
	Info  *string       `json:"info,omitempty"`
	Tape  *bool         `json:"tape,omitempty"`
}
type SetObject struct {
	Name   string       `json:"name"`
	Encore int          `json:"encore"`
	Song   []SongObject `json:"song"`
}
type SetsList struct {
	Set []SetObject `json:"set"`
}
type SetlistObject struct {
	Id          string       `json:"id"`
	VersionId   string       `json:"versionId"`
	EventDate   string       `json:"eventDate"`
	LastUpdated string       `json:"lastUpdated"`
	Artist      ArtistObject `json:"artist"`
	Venue       VenueObject  `json:"venue"`
	Tour        *TourObject  `json:"tour,omitempty"`
	Sets        SetsList     `json:"sets"`
	Info        *string      `json:"info,omitempty"`
	Url         string       `json:"url"`
}
type ResponseSetlist struct {
	Type         string          `json:"type"`
	ItemsPerPage int             `json:"itemsPerPage"`
	Page         int             `json:"page"`
	Total        int             `json:"total"`
	Setlist      []SetlistObject `json:"setlist"`
}
type EventObject struct {
	Id          string
	Date        string
	Venue       string
	City        string
	State       string
	Artist      string
	Tour        string
	SongsPlayed int
	Link        string
}
type EventsList struct {
	Event []EventObject
}

func (eventlist *EventsList) AddEvent(event EventObject) []EventObject {
	eventlist.Event = append(eventlist.Event, event)
	return eventlist.Event
}
func check(err error, note string) {
	if err != nil {
		log.Println("Checking for error with " + note)
		log.Println(err)
	}
}
func getEventsAttendedByUser(user, apiKey string) (events ResponseSetlist) {
	endpoint := "attended"
	page := 1

	// build out the url to be requested
	apiUrl, err := url.Parse("https://api.setlist.fm")
	check(err, "url.Parse")
	apiUrl.Path = "rest/1.0/user/" + user + "/" + endpoint
	apiUrlQs := url.Values{}
	apiUrlQs.Set("p", strconv.Itoa(page))
	apiUrl.RawQuery = apiUrlQs.Encode()

	// build out the request
	req, err := http.NewRequest("GET", apiUrl.String(), nil)
	check(err, "http.NewRequest")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("x-api-key", apiKey)

	// execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	check(err, "client.Do")

	// grab the data
	fullBody, err := ioutil.ReadAll(resp.Body)
	check(err, "ioutil.ReadAll")

	// close the resp.Body
	err = resp.Body.Close()
	check(err, "Problems closing resp.Body")

	// parse pagination things
	var pagination = new(ResponseSetlist)
	err = json.Unmarshal(fullBody, &pagination)
	check(err, "Unmarshal pagination")
	var pages int
	if pagination.Total != pagination.ItemsPerPage*pagination.Total/pagination.ItemsPerPage {
		pages = pagination.Total / pagination.ItemsPerPage
	} else {
		pages = pagination.Total/pagination.ItemsPerPage + 1
	}
	log.Println(endpoint + " - page 1 of " + strconv.Itoa(pages))

	var MainData ResponseSetlist
	err = json.Unmarshal(fullBody, &MainData)
	check(err, "Unmarshal "+endpoint)

	for page := 2; page <= pages; page++ {
		log.Println(endpoint + " - page " + strconv.Itoa(page) + " of " + strconv.Itoa(pages))

		// build out the url to be requested
		apiUrl, err := url.Parse("https://api.setlist.fm")
		check(err, "url.Parse")
		apiUrl.Path = "rest/1.0/user/" + user + "/" + endpoint
		apiUrlQs := url.Values{}
		apiUrlQs.Set("p", strconv.Itoa(page))
		apiUrl.RawQuery = apiUrlQs.Encode()

		// build out the request
		req, err := http.NewRequest("GET", apiUrl.String(), nil)
		check(err, "http.NewRequest")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("x-api-key", apiKey)

		// execute the request
		client := &http.Client{}
		resp, err := client.Do(req)
		check(err, "client.Do")

		// grab the data
		fullBody, err := ioutil.ReadAll(resp.Body)
		check(err, "ioutil.ReadAll")

		// close the resp.Body
		err = resp.Body.Close()
		check(err, "Problems closing resp.Body")

		var tempData ResponseSetlist
		err = json.Unmarshal(fullBody, &tempData)
		check(err, "Unmarshal "+endpoint)

		MainData.Setlist = append(MainData.Setlist, tempData.Setlist...)

		time.Sleep(150 * time.Millisecond)
	}
	return MainData
}
func cleanString(s string) (cleaned string) {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "'", "")
	s = strings.ReplaceAll(s, "&", "+")
	s = strings.ReplaceAll(s, "”", "")
	s = strings.ReplaceAll(s, "“", "")
	s = strings.ReplaceAll(s, "\"", "")
	cleaned = s
	return cleaned
}
func writeJsonEventsList(events ResponseSetlist, path string) {
	var el = new(EventsList)
	for _, event := range events.Setlist {
		eventDate, _ := time.Parse("02-01-2006", event.EventDate)
		eventDateString := fmt.Sprint(eventDate.Format("2006-01-02"))
		//eventYearString := fmt.Sprint(eventDate.Format("2006"))
		//relativeLink := "/event/" + eventYearString + "/" + cleanString(event.Artist.Name) + "/"

		e := EventObject{
			Id:     event.Id,
			Date:   eventDateString,
			Venue:  event.Venue.Name,
			City:   event.Venue.City.Name,
			State:  event.Venue.City.State,
			Artist: event.Artist.Name,
			//Tour:   tour,
			//SongsPlayed: len(event.Sets[].Set),
			Link: event.Url}
		el.AddEvent(e)
	}
	fileout := path + "list.json"
	_ = os.Remove(fileout) // remove from prior run
	//check(err, "Cannot remove the file")
	f, err := os.OpenFile(fileout, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	check(err, "Cannot create the file")

	prettyjson, _ := json.MarshalIndent(el, "", "\t")
	_, err = f.Write([]byte(prettyjson))
	check(err, "Cannot write json")

	err = f.Close()
	check(err, "Problem closing the file")
}
func writeJsonEventFiles(events ResponseSetlist, path string) {
	for _, event := range events.Setlist {
		// write event-artist.json
		eventDate, _ := time.Parse("02-01-2006", event.EventDate)
		eventDateString := fmt.Sprint(eventDate.Format("2006-01-02"))
		fileout := path + eventDateString + "-" + cleanString(event.Artist.Name) + ".json"
		_ = os.Remove(fileout) // remove from prior run
		//check(err, "Cannot remove the file")
		f, err := os.OpenFile(fileout, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		check(err, "Cannot create the file")

		prettyjson, _ := json.MarshalIndent(event, "", "\t")
		_, err = f.Write([]byte(prettyjson))
		check(err, "Cannot write json")

		err = f.Close()
		check(err, "Problem closing the file")
	}
}

func main() {
	type Configuration struct {
		User       string `json:"user"`
		ApiKey     string `json:"apiKey"`
		OutputPath string `json:"outputPath"`
	}
	file, _ := os.Open("configuration.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}

	EventArtists := getEventsAttendedByUser(configuration.User, configuration.ApiKey)
	writeJsonEventsList(EventArtists, configuration.OutputPath)
	writeJsonEventFiles(EventArtists, configuration.OutputPath)
}
