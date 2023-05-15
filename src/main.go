package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"main/src/Config"
	"main/src/Response"
	"main/src/RouteFinder"
	"net/http"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		fmt.Fprintf(w, "Method must be GET, not %s", r.Method)
		return
	}

	var resp Response.Response = Response.Response{}
	output, err := json.Marshal(resp)

	if err != nil {
		log.Println("Couldn't marshal the output")
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
}

func getRouteFromAtoBHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		log.Fatal("Method call has to be POST")
		return
	}

	r_body, err_0 := ioutil.ReadAll(r.Body)

	if err_0 != nil {
		log.Fatalf("Coudln't read the request body. Error: %s", err_0.Error())
		return
	}

	type body_struct struct {
		Mode       []string `json:"mode"`
		Waypoint1  string   `json:"waypoint1"`
		Waypoint2  string   `json:"waypoint2"`
		RouteMatch int32    `json:"routematch"`
	}

	var bs body_struct = body_struct{}

	err_1 := json.Unmarshal(r_body, &bs)

	if err_1 != nil {
		log.Fatalf("Coudln't unmarshall the request body. Error: %s", err_1.Error())
		return
	}

	var apiKey string = Config.Config{}.LoadConfig().HEREAPIKey[1]

	var resp RouteFinder.RouteResponse = (&RouteFinder.RouteResponse{}).GetRouteFromAtoB(apiKey, bs.Mode,
		bs.Waypoint1, bs.Waypoint2, bs.RouteMatch)

	data, _ := json.Marshal(resp)

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func main() {

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/getRouteFromAtoBHandler", getRouteFromAtoBHandler)

	fmt.Println("Listening at port 8080...")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Server couldn't start. Error:", err.Error())
		return
	} else {
		fmt.Println("Listening at port 8080...")
	}
}
