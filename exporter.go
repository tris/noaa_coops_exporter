package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	labelNames = []string{
		"station_id",
		"station_name",
		"latitude",
		"longitude",
	}
)

func scrapeHandler(w http.ResponseWriter, r *http.Request) {
	stationID := r.URL.Query().Get("station_id")
	if stationID == "" {
		http.Error(w, "Missing required parameter: station_id", http.StatusBadRequest)
		return
	}

	// Fetch latest reading from NOAA
	resp, err := http.Get("https://api.tidesandcurrents.noaa.gov/api/prod/datagetter?product=wind&date=latest&time_zone=gmt&units=english&format=json&station=" + stationID)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var weatherData WeatherData
	if err := json.Unmarshal(body, &weatherData); err != nil {
		log.Println(err)
		return
	}

	if len(weatherData.Data) == 0 {
		log.Println("No wind data available for station_id " + stationID)
		return
	}

	// Create a new registry for this scrape
	registry := prometheus.NewRegistry()

	// Create labels
	labels := prometheus.Labels{
		"station_id":   weatherData.Metadata.ID,
		"station_name": weatherData.Metadata.Name,
		"latitude":     strconv.FormatFloat(weatherData.Metadata.Latitude, 'f', 6, 64),
		"longitude":    strconv.FormatFloat(weatherData.Metadata.Longitude, 'f', 6, 64),
	}

	// Create gauges and set values
	collector := &WeatherCollector{
		metrics: []*weatherMetric{
			newWeatherMetric("noaa_wind_speed", "Wind speed in miles per hour", labels, weatherData.Data[0].Speed, weatherData.Data[0].Time),
			newWeatherMetric("noaa_wind_gust", "Wind gust in miles per hour", labels, weatherData.Data[0].Gust, weatherData.Data[0].Time),
			newWeatherMetric("noaa_wind_direction", "Wind direction in degrees", labels, weatherData.Data[0].Degrees, weatherData.Data[0].Time),
		},
	}

	// Register the gauges with the registry
	registry.MustRegister(collector)

	// Use a promhttp.HandlerFor with the new registry to serve the metrics
	promhttp.HandlerFor(registry, promhttp.HandlerOpts{}).ServeHTTP(w, r)
}
