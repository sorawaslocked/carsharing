package dto

import (
	"fmt"

	sharedmodel "carsharing/shared/model"
)

type scanner interface {
	Scan(dest ...any) error
}

func column(tableAlias, name string) string {
	if tableAlias == "" {
		return name
	}
	return fmt.Sprintf("%s.%s", tableAlias, name)
}

func ImageKeysToImages(keys []string) []sharedmodel.Image {
	images := make([]sharedmodel.Image, len(keys))
	for i, k := range keys {
		images[i] = sharedmodel.Image{Key: k}
	}
	return images
}

func ImagesToKeys(images []sharedmodel.Image) []string {
	keys := make([]string, len(images))
	for i, img := range images {
		keys[i] = img.Key
	}
	return keys
}
