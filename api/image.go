package api

type Image struct {
	Id         int32  `json:"id"`
	File       string `json:"file" binding:"required" gorm:"size:255;not null"`
	Type       string `json:"type" gorm:"size:16;not null"`
	Created_at int32  `json:"created_at" gorm:"not null"`
	Updated_at int32  `json:"updated_at" gorm:"not null"`
}

const (
	PngType  string = "png"
	JpgType  string = "jpg"
	JpegType string = "jpeg"
	GifType  string = "gif"
)
