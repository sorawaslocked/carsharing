package model

type Location struct {
	Latitude  float64
	Longitude float64
}

type Pagination struct {
	Limit  int64
	Offset int64
}

type ActorType string

const (
	ActorTypeUser   ActorType = "user"
	ActorTypeSystem ActorType = "system"
)

func (a ActorType) String() string {
	return string(a)
}
