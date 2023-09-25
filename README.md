# weather-api

## Overview
API to take a lat/long from the request and check with third party weather API for provided location.

 - Checks for highest temperature over the next 5 days, selecting the lowest humidity day if multiple have the same hihighest temperature, or the first date if humidity is the same.

 - Limits to 6 decimal places for the lat/long values

 - Only accepts UK locations.

## Usage
Update main.go to replace `apiKey` with a valid API key in the line:
```go
p := weatherprovider.NewOpenWeatherProvider("http://api.openweathermap.org", "apiKey")
```

To run the API, from the root folder of the project, execute:
```
go run .
```

You can then use your preferred method of accessing (browser, postman, curl...)

Example 1 (Brighton, UK, should return a valid response):

```
localhost:8080/weather?lat=50.821462&long=-0.140056
```

Example 2 (Paris, FR, should return an invalid response):

```
localhost:8080/weather?lat=48.8535&long=2.3484
```

Example 2 (Too many decimal places, should return an invalid response):

```
localhost:8080/weather?lat=48.8535123456&long=2.3484
```

## Status
The functionality of the calling and returning the warmest day from the API is largely complete.

The second requirement for storing the call history and allowing querying of that history has only started to be stubbed out with TODOs.