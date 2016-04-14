package api

type ImageType struct {
	ID           int32  `json:"id"`
	Name         string `json:"name" binding:"required" gorm:"size:32;not null"`
	Path         string `json:"path" binding:"required" gorm:"size:255;not null"`
	Thumb_width  uint16 `json:"thumb_width" binding:"required" gorm:"not null"`
	Thumb_height uint16 `json:"thumb_height" binding:"required" gorm:"not null"`
	Created_at   int32  `json:"created_at" binding:"required" gorm:"not null"`
	Updated_at   int32  `json:"updated_at" binding:"required" gorm:"not null"`
}
