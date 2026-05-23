package validation

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidLatitudeRange  = errors.New("latitude must be between -90 and 90")
	ErrInvalidLongitudeRange = errors.New("longitude must be between -180 and 180")
	ErrInvalidRadiusRange    = fmt.Errorf("radius must be between %g and %g km", MinRadiusKM, MaxRadiusKM)
)
