package model

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
	if _, ok := validImageTypes[it]; !ok {
		return "", ErrInvalidImageType
	}
	return it, nil
}

func (t ImageType) String() string {
	return string(t)
}
