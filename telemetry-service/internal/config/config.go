package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	GRPCPort             string  `yaml:"grpc_port"              env:"GRPC_PORT"               env-default:":50051"`
	DBUrl                string  `yaml:"db_url"                 env:"DB_URL"                  env-default:"postgres://postgres:postgres@localhost:5432/car_rental?sslmode=disable"`
	OSRMUrl              string  `yaml:"osrm_url"               env:"OSRM_URL"                env-default:"http://router.project-osrm.org"`
	OSRMProfile          string  `yaml:"osrm_profile"           env:"OSRM_PROFILE"            env-default:"car"`
	TelemetryInterval    string  `yaml:"telemetry_interval"     env:"TELEMETRY_INTERVAL"      env-default:"15s"`
	SpeedKmh             float64 `yaml:"speed_kmh"              env:"SPEED_KMH"               env-default:"80"`
	FuelConsumPerKm      float32 `yaml:"fuel_consum_per_km"     env:"FUEL_CONSUM_PER_KM"      env-default:"0.2"`
	BatteryConsumPerKm   float32 `yaml:"battery_consum_per_km"  env:"BATTERY_CONSUM_PER_KM"   env-default:"0.25"`
	NATSUrl              string  `yaml:"nats_url"               env:"NATS_URL"                env-default:"nats://localhost:4222"`
	TripStartedSubject   string  `yaml:"trip_started_subject"    env:"TRIP_STARTED_SUBJECT"    env-default:"trip.started"`
	TripEndedSubject     string  `yaml:"trip_ended_subject"      env:"TRIP_ENDED_SUBJECT"      env-default:"trip.ended"`
	TripCancelledSubject string  `yaml:"trip_cancelled_subject"  env:"TRIP_CANCELLED_SUBJECT"  env-default:"trip.cancelled"`
}

func Load(path string) (*Config, error) {
	cfg := &Config{}
	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
