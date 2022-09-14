package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Flight struct {
	FlightPath     [][2]string `json:"flightPath"`
}


func apiResponse (w http.ResponseWriter, r *http.Request) {
	// Set the return Content-Type as JSON like before
	w.Header().Set( "Content-Type" , "application/json" )

	// Change the response depending on the method being requested
	switch r.Method {
	case "GET" :
		w.WriteHeader(http.StatusOK)
		w.Write([] byte ( `{"message": "Please use Post method"}` ))
	case "POST" :
		w.WriteHeader(http.StatusOK)
		response:=parseRequest(r)
		str, err := json.Marshal(response)
		if err != nil {
			fmt.Println(err)
		}
		w.Write(str)
	default :
		w.WriteHeader(http.StatusNotFound)
		w.Write([] byte ( `{"message": "Can't find method requested"}` ))
	}
}
func parseRequest (r *http.Request) map[string][2]string {
	defer r.Body.Close()
	 response:=make (map[string][2]string,1)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("read err")
	}
	var flight Flight
	if err = json.Unmarshal(body, &flight); err != nil {
		fmt.Printf("Unmarshal err, %v", err)
		return nil
	}
	if len(flight.FlightPath)>=1{
		response["flight"]=sortFlightPath(flight)
		return response
	}else{
		fmt.Println("flightPath is nil")
		return nil
	}
}

func sortFlightPath (flight Flight) [2]string{
	start := make(map[string]int)
	end := make(map[string]int)
	var output [2]string
	for _, flightPath := range flight.FlightPath {
		if _, ok := start[flightPath[0]]; !ok {
			start[flightPath[0]] = 1
		}
		if _, ok := end[flightPath[1]]; !ok {
			end[flightPath[1]] = 1

		}
	}
	for _, flightPath := range flight.FlightPath {
		if _, ok := start[flightPath[1]]; !ok {
			output[1] = 	flightPath[1]
		}
		if _, ok := end[flightPath[0]]; !ok {
			output[0] = 	flightPath[0]
		}
	}
	return output
	}

func main () {
	http.HandleFunc( "/" ,apiResponse)
	log.Fatal(http.ListenAndServe( ":8080" , nil ))
}
