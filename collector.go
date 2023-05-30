package main

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"google.golang.org/protobuf/proto"
)

type weatherMetric struct {
	desc *prometheus.Desc
	value float64
	labels prometheus.Labels
	timestamp time.Time
}

type WeatherCollector struct {
	metrics []*weatherMetric
}

func (c *WeatherCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range c.metrics {
		ch <- m.Desc()
	}
}

func (c *WeatherCollector) Collect(ch chan<- prometheus.Metric) {
	for _, m := range c.metrics {
		ch <- m
	}
}

func (c *weatherMetric) Desc() *prometheus.Desc {
	return c.desc
}

func (c *weatherMetric) Write(m *dto.Metric) error {
	m.Label = []*dto.LabelPair{}
	for k, v := range c.labels {
		m.Label = append(m.Label, &dto.LabelPair{
			Name:  proto.String(k),
			Value: proto.String(v),
		})
	}
	m.Gauge = &dto.Gauge{Value: &c.value}
	m.TimestampMs = proto.Int64(c.timestamp.UnixNano() / int64(time.Millisecond))
	return nil
}

func newWeatherMetric(name string, help string, labels prometheus.Labels, value float64, timestamp time.Time) *weatherMetric {
	return &weatherMetric{
		desc: prometheus.NewDesc(name, help, nil, labels),
		value: value,
		labels: labels,
		timestamp: timestamp,
	}
}
