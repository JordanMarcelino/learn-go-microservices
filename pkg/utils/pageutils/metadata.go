package pageutils

import (
	"math"

	"github.com/jordanmarcelino/learn-go-microservices/pkg/dto"
)

func NewMetadata(count, page, limit int64) *dto.PageMetaData {
	totalItems := count
	totalPage := int64(math.Ceil(float64(totalItems) / float64(limit)))

	return &dto.PageMetaData{
		Page:      page,
		Size:      limit,
		TotalItem: totalItems,
		TotalPage: totalPage,
	}
}
