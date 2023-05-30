package main

import (
	"encoding/json"
	"strconv"
	"time"
)

type Metadata struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Latitude  float64 `json:"-"`
	Longitude float64 `json:"-"`
}

type rawMetadata struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Latitude  string `json:"lat"`
	Longitude string `json:"lon"`
}

type WindData struct {
	Time    time.Time `json:"-"`
	Speed   float64   `json:"-"`
	Degrees float64   `json:"-"`
	Ordinal string    `json:"dr"`
	Gust    float64   `json:"-"`
	F       string    `json:"f"`
}

type rawWindData struct {
	Time    string `json:"t"`
	Speed   string `json:"s"`
	Degrees string `json:"d"`
	Ordinal string `json:"dr"`
	Gust    string `json:"g"`
	F       string `json:"f"`
}

type WeatherData struct {
	Metadata Metadata   `json:"metadata"`
	Data     []WindData `json:"data"`
}

func (m *Metadata) UnmarshalJSON(data []byte) error {
	var raw rawMetadata
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	var err error
	if m.Latitude, err = strconv.ParseFloat(raw.Latitude, 64); err != nil {
		return err
	}
	if m.Longitude, err = strconv.ParseFloat(raw.Longitude, 64); err != nil {
		return err
	}

	m.ID = raw.ID
	m.Name = raw.Name

	return nil
}

func (w *WindData) UnmarshalJSON(data []byte) error {
	var raw rawWindData
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	t, err := time.Parse("2006-01-02 15:04", raw.Time)
	if err != nil {
		return err
	}
	w.Time = t

	if w.Speed, err = strconv.ParseFloat(raw.Speed, 64); err != nil {
		return err
	}
	if w.Degrees, err = strconv.ParseFloat(raw.Degrees, 64); err != nil {
		return err
	}
	if w.Gust, err = strconv.ParseFloat(raw.Gust, 64); err != nil {
		return err
	}

	w.Ordinal = raw.Ordinal
	w.F = raw.F

	return nil
}
