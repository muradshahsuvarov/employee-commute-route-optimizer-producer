# Employee Commute Route Optimizer Producer

The Route Finder API is a simple HTTP server that provides an interface to retrieve routes between two waypoints using the HERE API.

## Installation

1. Clone the repository:

```shell
git clone https://github.com/muradshahsuvarov/employee-commute-route-optimizer-producer
```

2. Install dependencies:

```shell
go mod download
```

## Configuration

The application requires an API key from HERE to access the routing services. Follow the steps below to configure the API key:

1. Open the `config.json` file located in the `src/Config` directory.

2. Replace the placeholder value `<HERE-API-KEY>` with your HERE API key.

## Usage

1. Start the server:

```shell
go run src/main.go
```

2. The server will start listening on port 8080. You can access the following endpoints:

- **GET /:** Returns a default response in JSON format.

- **POST /getRouteFromAtoBHandler:** Accepts a JSON payload containing the mode of transportation, waypoints, and route matching flag. It returns the route information between the specified waypoints.

## Request Payload

The `getRouteFromAtoBHandler` endpoint expects a POST request with the following JSON payload:

```json
{
  "mode": ["car"],
  "waypoint0": "Start Point",
  "waypoint1": "End Point",
  "routematch": 0
}
```

- `mode` (array of strings): Specifies the mode of transportation. Currently, only "car" is supported.

- `waypoint0` (string): The starting point of the route.

- `waypoint1` (string): The destination of the route.

- `routematch` (integer): Flag indicating whether to perform route matching. Set to 0 for exact matching or 1 for route deviation tolerance.

## Response Format

The response from the `getRouteFromAtoBHandler` endpoint will be in the following format:

```json
{
  "response": {
    "route": [
      {
        "mode": {
          "type": "matched",
          "transportModes": ["car"],
          "trafficMode": "enabled"
        },
        "boatFerry": false,
        "railFerry": false,
        "waypoint": [
          {
            "linkId": "+127812120",
            "mappedPosition": {
              "latitude": 37.77391,
              "longitude": -122.43128
            },
            ...
          },
          ...
        ],
        "leg": [
          {
            "length": 3249,
            "travelTime": 433,
            "link": [
              {
                "linkId": "127812120",
                "length": 42.24,
                "remainDistance": 3249,
                ...
              },
              ...
            ],
            ...
          },
          ...
        ],
        "summary": {
          "travelTime": 1234,
          "distance": 5678,
          ...
        }
      }
    ],
    "warnings": [],
    "language": "en"
  }
}
```

The response contains detailed route information, including modes of transportation, waypoints, legs, and a summary of the route.

## Contributing

Murad Shahsuvarov
