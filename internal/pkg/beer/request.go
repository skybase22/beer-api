package beer

import (
	"beer-api/internal/models"
	"mime/multipart"
)

type createRequest struct {
	Name        string                `json:"name" form:"name"`
	Type        string                `json:"type" form:"type"`
	Description string                `json:"description" form:"description"`
	File        *multipart.FileHeader `form:"file"`
}

type updateRequest struct {
	ID          uint                  `json:"id" path:"id" validate:"required"`
	Name        *string               `json:"name" form:"name"`
	Type        *string               `json:"type" form:"type"`
	Description *string               `json:"description" form:"description"`
	File        *multipart.FileHeader `form:"file"`
}

type getOneRequest struct {
	ID int64 `json:"id" query:"id" path:"id" validate:"required"`
}

type GetAllRequest struct {
	models.PageForm
}
