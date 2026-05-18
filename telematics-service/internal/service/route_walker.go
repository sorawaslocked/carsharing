package service

import (
	"context"
	"log/slog"
	"math"
	"math/rand/v2"

	osrm "github.com/gojuno/go.osrm"
	geo "github.com/paulmach/go.geo"
)

const (
	fallbackCircleRadiusDeg = 0.003 // ~330 m radius
	fallbackCirclePoints    = 60    // 60 points → one every 6°
)

type latLng struct {
	lat, lng float64
}

// routeWalker advances a position along a polyline of waypoints at a
// configurable speed, returning updated coordinates on each call to advance.
type routeWalker struct {
	waypoints []latLng
	segIdx    int
	segFrac   float64 // fractional progress within current segment [0, 1)
}

// done reports whether the walker has reached the end of its waypoint list.
func (w *routeWalker) done() bool {
	return w.segIdx >= len(w.waypoints)-1
}

// advance moves the walker forward by metersToAdvance along the polyline and
// returns the interpolated position. If the end of the route is reached, the
// last waypoint coordinates are returned.
func (w *routeWalker) advance(metersToAdvance float64) (lat, lng float64) {
	remaining := metersToAdvance

	for w.segIdx < len(w.waypoints)-1 {
		from := w.waypoints[w.segIdx]
		to := w.waypoints[w.segIdx+1]

		segLen := haversineMeters(from.lat, from.lng, to.lat, to.lng)
		if segLen < 0.1 { // skip degenerate segments
			w.segIdx++
			w.segFrac = 0
			continue
		}

		distInSeg := segLen * (1.0 - w.segFrac)
		if remaining <= distInSeg {
			w.segFrac += remaining / segLen
			lat = from.lat + w.segFrac*(to.lat-from.lat)
			lng = from.lng + w.segFrac*(to.lng-from.lng)
			return
		}

		remaining -= distInSeg
		w.segIdx++
		w.segFrac = 0
	}

	last := w.waypoints[len(w.waypoints)-1]
	return last.lat, last.lng
}

// fetchRoute requests a route from OSRM starting at (lat, lng) toward a
// randomly chosen destination 2–5 km away. Falls back to a circular path if
// OSRM is unavailable or returns an empty route.
func (s *SimulationService) fetchRoute(ctx context.Context, lat, lng float64) *routeWalker {
	destLat, destLng := randomDestination(lat, lng)

	// geo.Point is [2]float64{longitude, latitude}
	coords := osrm.NewGeometryFromPointSet(geo.PointSet{
		{lng, lat},
		{destLng, destLat},
	})

	resp, err := s.osrmClient.Route(ctx, osrm.RouteRequest{
		Profile:     s.osrmProfile,
		Coordinates: coords,
		Steps:       osrm.StepsTrue,
		Overview:    osrm.OverviewFull,
	})
	if err != nil || len(resp.Routes) == 0 {
		slog.Warn("OSRM unavailable, using fallback circle", "error", err)
		return fallbackCircleWalker(lat, lng)
	}

	var waypoints []latLng
	for _, leg := range resp.Routes[0].Legs {
		for _, step := range leg.Steps {
			// step.Geometry embeds geo.Path; Points() returns []geo.Point where [0]=lng, [1]=lat
			for _, pt := range step.Geometry.Points() {
				waypoints = append(waypoints, latLng{lat: pt[1], lng: pt[0]})
			}
		}
	}

	if len(waypoints) < 2 {
		slog.Warn("OSRM returned too few waypoints, using fallback circle")
		return fallbackCircleWalker(lat, lng)
	}

	return &routeWalker{waypoints: waypoints}
}

// randomDestination picks a random point 2–5 km from (lat, lng).
func randomDestination(lat, lng float64) (float64, float64) {
	distDeg := 0.018 + rand.Float64()*0.027 // ~2–5 km expressed in degrees
	angle := rand.Float64() * 2 * math.Pi
	cosLat := math.Cos(lat * math.Pi / 180)
	if cosLat < 0.001 {
		cosLat = 0.001
	}
	return lat + distDeg*math.Cos(angle), lng + distDeg*math.Sin(angle)/cosLat
}

// fallbackCircleWalker generates a closed circular route of ~330 m radius
// around the given centre for use when OSRM is unavailable.
func fallbackCircleWalker(lat, lng float64) *routeWalker {
	cosLat := math.Cos(lat * math.Pi / 180)
	if cosLat < 0.001 {
		cosLat = 0.001
	}
	pts := make([]latLng, fallbackCirclePoints+1)
	for i := 0; i < fallbackCirclePoints; i++ {
		angle := 2 * math.Pi * float64(i) / float64(fallbackCirclePoints)
		pts[i] = latLng{
			lat: lat + fallbackCircleRadiusDeg*math.Cos(angle),
			lng: lng + fallbackCircleRadiusDeg*math.Sin(angle)/cosLat,
		}
	}
	pts[fallbackCirclePoints] = pts[0] // close the loop
	return &routeWalker{waypoints: pts}
}

// haversineMeters returns the great-circle distance in metres between two
// latitude/longitude coordinates.
func haversineMeters(lat1, lng1, lat2, lng2 float64) float64 {
	const R = 6_371_000.0
	φ1 := lat1 * math.Pi / 180
	φ2 := lat2 * math.Pi / 180
	Δφ := (lat2 - lat1) * math.Pi / 180
	Δλ := (lng2 - lng1) * math.Pi / 180
	a := math.Sin(Δφ/2)*math.Sin(Δφ/2) +
		math.Cos(φ1)*math.Cos(φ2)*math.Sin(Δλ/2)*math.Sin(Δλ/2)
	return R * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}
