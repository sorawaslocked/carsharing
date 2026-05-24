module carsharing/telematics-service

go 1.26

require (
	carsharing/protos v0.0.0
	github.com/gojuno/go.osrm v0.1.0
	github.com/ilyakaznacheev/cleanenv v1.5.0
	github.com/lib/pq v1.12.3
	github.com/nats-io/nats.go v1.52.0
	github.com/paulmach/go.geo v0.0.0-20180829195134-22b514266d33
	google.golang.org/grpc v1.81.1
	google.golang.org/protobuf v1.36.11
)

replace carsharing/protos v0.0.0 => ../protos

require (
	github.com/BurntSushi/toml v1.6.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/klauspost/compress v1.18.5 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/nats-io/nkeys v0.4.15 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/paulmach/go.geojson v1.4.0 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/stretchr/testify v1.11.1 // indirect
	golang.org/x/crypto v0.49.0 // indirect
	golang.org/x/net v0.51.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
	golang.org/x/text v0.35.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260226221140-a57be14db171 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
)
