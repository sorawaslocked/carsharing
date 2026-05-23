package model

import sharedmodel "carsharing/shared/model"

type LocationFilter struct {
	Location sharedmodel.Location
	RadiusKM float64
}
