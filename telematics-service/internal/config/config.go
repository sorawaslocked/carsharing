package config

import "os"

type Config struct {
	GRPCPort           string
	NATSUrl            string
	OSRMUrl            string
	OSRMProfile        string
	TripStartedSubject string
	TripEndedSubject   string
}

func Load() *Config {
	return &Config{
		GRPCPort:           getEnv("GRPC_PORT", ":50051"),
		NATSUrl:            getEnv("NATS_URL", "nats://localhost:4222"),
		OSRMUrl:            getEnv("OSRM_URL", "http://router.project-osrm.org"),
		OSRMProfile:        getEnv("OSRM_PROFILE", "car"),
		TripStartedSubject: getEnv("TRIP_STARTED_SUBJECT", "events.trip.started"),
		TripEndedSubject:   getEnv("TRIP_ENDED_SUBJECT", "events.trip.ended"),
	}
}

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}
