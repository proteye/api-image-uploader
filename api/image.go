package api

type Image struct {
	Id          int32 `json:"id"`
	ImageType   ImageType
	ImageTypeID int32  `json:"image_type_id" gorm:"not null"`
	Filename    string `json:"filename" binding:"required" gorm:"size:255;not null"`
	Original    string `json:"original" binding:"required" gorm:"size:255;not null"`
	Mime_type   string `json:"mime_type" binding:"required" gorm:"size:16;not null"`
	Created_at  int32  `json:"created_at" binding:"required" gorm:"not null"`
	Updated_at  int32  `json:"updated_at" binding:"required" gorm:"not null"`
}

type ImageType struct {
	Id         int32  `json:"id"`
	Name       string `json:"filename" binding:"required" gorm:"size:32;not null"`
	Path       string `json:"original" binding:"required" gorm:"size:255;not null"`
	Created_at int32  `json:"created_at" binding:"required" gorm:"not null"`
	Updated_at int32  `json:"updated_at" binding:"required" gorm:"not null"`
}

type Meta struct {
	Name       string `json:"filename" binding:"required" gorm:"primary_key;size:32"`
	Value_str  string `json:"original" binding:"required" gorm:"size:255;default null"`
	Value_int  int32  `json:"value_int" binding:"required" gorm:"default null"`
	Created_at int32  `json:"created_at" binding:"required" gorm:"not null"`
	Updated_at int32  `json:"updated_at" binding:"required" gorm:"not null"`
}

const (
	PngType  string = "image/png"
	JpegType string = "image/jpeg"
	GifType  string = "image/gif"
)
