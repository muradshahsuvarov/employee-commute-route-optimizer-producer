package RouteFinder

import (
	"encoding/json"
	"fmt"
	"main/src/Consumer"
	"main/src/KafkaResponse"
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

	fmt.Println("SURKAY 1")
	var modeStr string = strings.Join(mode, ";")
	var _url string = fmt.Sprintf("https://routematching.hereapi.com/v8/match/routelinks?apiKey=%s&mode=%s&waypoint0=%s&waypoint1=%s&routeMatch=%d", apiKey, modeStr, waypoint0, waypoint1, routeMatch)
	var _id string = Random.GenerateRandomID(10, "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+_)(*&^%$#@!)")
	var _type string = "HERE_API"
	var _server string = "localhost:9092"
	var _producer_topic string = "ecro_req_topic"
	var _consumer_topic string = "ecro_res_topic"
	var _producerPropertiesFile = "C:\\kafka\\config\\producer.properties"
	var _consumerPropertiesFile = "C:\\kafka\\config\\consumer.properties"

	fmt.Println("SURKHAY 2")

	var routeResponse RouteResponse = RouteResponse{}

	fmt.Println("SURKHAY 3")

	var kResChan <-chan KafkaResponse.KafkaResponse = make(chan KafkaResponse.KafkaResponse)
	var responseChannel <-chan Response.Response = make(chan Response.Response)

	fmt.Printf("SURKHAY 4: %s %s %s %s %d %s \n", _id, _type, _server, _producer_topic, 0, _consumerPropertiesFile)
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		fmt.Println("SURKHAY 5 START")
		kResChan = Consumer.ConsumeMessages(_id, _type, _server, _consumer_topic, 0, _consumerPropertiesFile)
		for kafkaResponse := range kResChan {
			fmt.Println("SURKHAY 5: ", kafkaResponse)
			fmt.Println("SURKHAY 5: ", kafkaResponse.Message)
			var byteMessage []byte = []byte(kafkaResponse.Message)
			json.Unmarshal(byteMessage, &routeResponse)
			wg.Done()
			fmt.Println("SURKHAY 6: ", routeResponse)
			return
		}
	}()

	wg.Add(1)
	go func() {
		responseChannel = Producer.ProduceMessage(_id, _server, _producer_topic, _type, _url, _producerPropertiesFile)
		select {
		case res := <-responseChannel:
			if res.Id == _id {
				fmt.Println("SURKHAY 7: ", res)
				wg.Done()
				return // Exit the goroutine after sending the response
			}
		}
	}()

	fmt.Println("SURKHAY 8")
	wg.Wait()
	fmt.Println("--------------------DONE---------------------")
	return routeResponse
}
