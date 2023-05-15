package RouteFinder

import (
	"encoding/json"
	"fmt"
	"main/src/Consumer"
	"main/src/Producer"
	"main/src/Random"
	"main/src/Response"
	"strings"
	"sync"
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

func (rr *RouteResponse) GetRouteFromAtoB(apiKey string, mode []string, waypoint0 string, waypoint1 string, routeMatch int32) RouteResponse {

	var modeStr string = strings.Join(mode, ";")
	var _url string = fmt.Sprintf("https://routematching.hereapi.com/v8/match/routelinks?apiKey=%s&mode=%s&waypoint0=%s&waypoint1=%s&routeMatch=%d", apiKey, modeStr, waypoint0, waypoint1, routeMatch)
	var _id string = Random.GenerateRandomID(10, "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+_)(*&^%$#@!)")
	var _type string = "HERE_API"
	var _server string = "localhost:9092"
	var _topic string = "ecro_req_topic"
	var _propertiesFile = "C:\\kafka\\config\\producer.properties"

	responseChannel := Producer.ProduceMessage(_id, _server, _topic, _type, _url, _propertiesFile)
	responseReceived := make(chan Response.Response) // Channel to receive the response

	wg := sync.WaitGroup{}

	go func() {
		for {
			select {
			case res := <-responseChannel:
				if res.Id == _id {
					responseReceived <- res // Send the response through the channel
					close(responseReceived)
					return // Exit the goroutine after sending the response
				}
			}
		}
	}()

	var routeResponse RouteResponse = RouteResponse{}

	output := <-responseReceived // Wait for the response to be received from the channel

	kResChan := Consumer.ConsumeMessages(_id, _type, _server, _topic, output.Partition, _propertiesFile)

	wg.Add(1)
	// Use a separate goroutine to consume messages from kResChan to avoid the infinite blocking possibilty
	go func() {
		for kafkaResponse := range kResChan {
			var byteMessage []byte = []byte(kafkaResponse.Message)
			json.Unmarshal(byteMessage, &routeResponse)
			wg.Done()
			return
		}
	}()
	wg.Wait()

	return routeResponse
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
