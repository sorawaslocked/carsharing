package dto

import (
	"strconv"

	"carsharing/api-gateway/internal/model"
	sharedmodel "carsharing/shared/model"
	"github.com/gin-gonic/gin"
)

type location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Image struct {
	Key string `json:"key"`
	URL string `json:"url"`
}

func toImage(m sharedmodel.Image) Image {
	return Image{Key: m.Key, URL: m.URL}
}

func toImages(ms []sharedmodel.Image) []Image {
	if len(ms) == 0 {
		return nil
	}
	imgs := make([]Image, len(ms))
	for i, m := range ms {
		imgs[i] = toImage(m)
	}
	return imgs
}

func optToImage(m *sharedmodel.Image) *Image {
	if m == nil {
		return nil
	}
	img := toImage(*m)
	return &img
}

type ImageUploadData struct {
	PresignedPutURL string `json:"presignedPutURL"`
	ObjectKey       string `json:"objectKey"`
}

type ImageUploadResponse struct {
	UploadData ImageUploadData `json:"uploadData"`
}

func ToImageUploadDataResponse(m sharedmodel.ImageUploadData) ImageUploadData {
	return ImageUploadData{
		PresignedPutURL: m.PresignedPutURL,
		ObjectKey:       m.ObjectKey,
	}
}

func pagination(c *gin.Context) (*sharedmodel.Pagination, error) {
	var p sharedmodel.Pagination
	paginationEmpty := true

	if v := c.Query("limit"); v != "" {
		vInt, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, model.ErrInvalidQueryParam
		}

		p.Limit = vInt
		paginationEmpty = false
	}
	if v := c.Query("offset"); v != "" {
		vInt, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, model.ErrInvalidQueryParam
		}

		p.Offset = vInt
		paginationEmpty = false
	}

	if paginationEmpty {
		return nil, nil
	}

	return &p, nil
}

func IDParam(c *gin.Context) (string, error) {
	id := c.Param("id")

	if id == "" {
		return "", model.ErrEmptyIDParam
	}

	return id, nil
}
