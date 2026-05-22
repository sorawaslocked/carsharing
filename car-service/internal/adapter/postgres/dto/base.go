package dto

import (
	sharedmodel "carsharing/shared/model"
	"fmt"
)

type scanner interface {
	Scan(dest ...any) error
}

type ArgsBuilder struct {
	Args []any
}

func (b *ArgsBuilder) Add(arg any) string {
	b.Args = append(b.Args, arg)

	return fmt.Sprintf("$%d", len(b.Args))
}

func BuildPagination(b *ArgsBuilder, p *sharedmodel.Pagination) string {
	if p == nil {
		return ""
	}

	return " LIMIT " + b.Add(p.Limit) + " OFFSET " + b.Add(p.Offset)
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
