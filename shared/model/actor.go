package model

type ActorType string

const (
	ActorTypeSystem    ActorType = "system"
	ActorTypeUser      ActorType = "user"
	ActorTypeTelemetry ActorType = "telemetry"
)
