package model

import "errors"

var ErrInvalidImageType = errors.New("must be a valid image type")

type ImageType string

const (
	ImageTypeIDFront             ImageType = "id_front"
	ImageTypeIDBack              ImageType = "id_back"
	ImageTypeDrivingLicenseFront ImageType = "driving_license_front"
	ImageTypeDrivingLicenseBack  ImageType = "driving_license_back"
)

var validImageTypes = map[ImageType]struct{}{
	ImageTypeIDFront:             {},
	ImageTypeIDBack:              {},
	ImageTypeDrivingLicenseFront: {},
	ImageTypeDrivingLicenseBack:  {},
}

func AllImageTypes() []ImageType {
	return []ImageType{
		ImageTypeIDFront,
		ImageTypeIDBack,
		ImageTypeDrivingLicenseFront,
		ImageTypeDrivingLicenseBack,
	}
}

func ImageTypeFromString(s string) (ImageType, error) {
	it := ImageType(s)
	if _, ok := validImageTypes[it]; ok {
		return it, nil
	}
	return "", ErrInvalidImageType
}

func (t ImageType) String() string {
	return string(t)
}
