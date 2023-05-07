package RouteFinder

import (
	"fmt"
	"main/src/Producer"
	"main/src/Random"
	"main/src/Response"
	"strings"
)

type RouteResponse struct {
	Response struct {
		Route []struct {
			Mode struct {
				Type           string   `json:"type"`
				TransportModes []string `json:"transportModes"`
				TrafficMode    string   `json:"trafficMode"`
			} `json:"mode"`
			BoatFerry bool `json:"boatFerry"`
			RailFerry bool `json:"railFerry"`
			Waypoint  []struct {
				LinkID         string `json:"linkId"`
				MappedPosition struct {
					Latitude  float64 `json:"latitude"`
					Longitude float64 `json:"longitude"`
				} `json:"mappedPosition"`
				OriginalPosition struct {
					Latitude  float64 `json:"latitude"`
					Longitude float64 `json:"longitude"`
				} `json:"originalPosition"`
				Spot                        float64 `json:"spot"`
				ConfidenceValue             float64 `json:"confidenceValue"`
				Elevation                   float64 `json:"elevation"`
				HeadingDegreeNorthClockwise float64 `json:"headingDegreeNorthClockwise"`
				HeadingMatched              float64 `json:"headingMatched"`
				MatchDistance               float64 `json:"matchDistance"`
				MinError                    float64 `json:"minError"`
				RouteLinkSeqNrMatched       int     `json:"routeLinkSeqNrMatched"`
				SpeedMps                    float64 `json:"speedMps"`
				Timestamp                   int     `json:"timestamp"`
			} `json:"waypoint"`
			Leg []struct {
				Length     int `json:"length"`
				TravelTime int `json:"travelTime"`
				Link       []struct {
					LinkID          string    `json:"linkId"`
					Length          float64   `json:"length"`
					RemainDistance  int       `json:"remainDistance"`
					RemainTime      int       `json:"remainTime"`
					Shape           []float64 `json:"shape"`
					FunctionalClass int       `json:"functionalClass"`
					Confidence      float64   `json:"confidence"`
					SegmentRef      string    `json:"segmentRef"`
				} `json:"link"`
				TrafficTime     int `json:"trafficTime"`
				BaseTime        int `json:"baseTime"`
				RefReplacements struct {
					Zero string `json:"0"`
					One  string `json:"1"`
				} `json:"refReplacements"`
			} `json:"leg"`
			Summary struct {
				TravelTime  int           `json:"travelTime"`
				Distance    int           `json:"distance"`
				BaseTime    int           `json:"baseTime"`
				TrafficTime int           `json:"trafficTime"`
				Flags       []interface{} `json:"flags"`
			} `json:"summary"`
		} `json:"route"`
		Warnings []interface{} `json:"warnings"`
		Language string        `json:"language"`
	} `json:"response"`
}

// Return type must be RouteResponse
func (rr *RouteResponse) GetRouteFromAtoB(apiKey string, mode []string, waypoint0 string, waypoint1 string, routeMatch int32) Response.Response {
	var modeStr string = strings.Join(mode, ";")
	var url string = fmt.Sprintf("https://routematching.hereapi.com/v8/match/routelinks?apiKey=%s&mode=%s&waypoint0=%s&waypoint1=%s&routeMatch=%d", apiKey, modeStr, waypoint0, waypoint1, routeMatch)
	var id string = Random.GenerateRandomID(10, "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+_)(*&^%$#@!)")

	responseChannel := Producer.ProduceMessage(id, "localhost:9092", "ecro_req_topic", "HERE_API", url, "C:\\kafka\\config\\producer.properties")
	responseReceived := make(chan Response.Response) // Channel to receive the response

	go func() {
		for {
			select {
			case res := <-responseChannel:
				if res.Id == id {
					responseReceived <- res // Send the response through the channel
					return                  // Exit the goroutine after sending the response
				}
			}
		}
	}()

	output := <-responseReceived // Wait for the response to be received from the channel

	return output
}

// func (rr *RouteResponse) GetTheShortestLocation(apiKey string, mode []string, waitPoint string, waitPoints []string, routeMatch int) RouteResponse {

// 	var routeResponses []RouteResponse = make([]RouteResponse, 0)

// 	// Unbuffered channel
// 	ch := make(chan RouteResponse)

// 	for _, value := range waitPoints {
// 		go func(val string) {
// 			ch <- rr.GetRouteFromAtoB(apiKey, mode, waitPoint, val, int32(routeMatch))
// 		}(value)
// 	}

// 	// Blocking operation on the unbuffered channel ch
// 	for i := 0; i < len(waitPoints); i++ {
// 		routeResponses = append(routeResponses, <-ch)
// 	}

// 	// Sorting operation. Sorting is in ascending order
// 	sort.Slice(routeResponses, func(i, j int) bool {
// 		return routeResponses[i].Response.Route[0].Summary.TravelTime < routeResponses[j].Response.Route[0].Summary.TravelTime
// 	})

// 	return routeResponses[len(routeResponses)-1]
// }
