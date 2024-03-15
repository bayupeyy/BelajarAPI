package handler

type ActivityRequest struct {
	Kegiatan string `json:"kegiatan" form:"kegiatan" validate:"required"`
}
