package RouteResponse

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
				HeadingDegreeNorthClockwise int     `json:"headingDegreeNorthClockwise"`
				HeadingMatched              int     `json:"headingMatched"`
				MatchDistance               float64 `json:"matchDistance"`
				MinError                    int     `json:"minError"`
				RouteLinkSeqNrMatched       int     `json:"routeLinkSeqNrMatched"`
				SpeedMps                    int     `json:"speedMps"`
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
			} `json:"leg"`
		} `json:"route"`
	} `json:"response"`
}
