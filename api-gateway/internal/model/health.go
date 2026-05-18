package model

import "time"

type ServiceHealth struct {
	Name    string
	Status  string
	Version string

	Timestamp     time.Time
	UptimeSeconds uint64

	Dependencies []DependencyHealth
}

type DependencyHealth struct {
	Name   string
	Status string

	LatencyMS *uint32

	Error *string
}
