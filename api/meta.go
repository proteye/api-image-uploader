package api

type Meta struct {
	Name       string `json:"filename" binding:"required" gorm:"primary_key;size:32"`
	Value_str  string `json:"original" binding:"required" gorm:"size:255;default null"`
	Value_int  int32  `json:"value_int" binding:"required" gorm:"default null"`
	Created_at int32  `json:"created_at" gorm:"not null"`
	Updated_at int32  `json:"updated_at" gorm:"not null"`
}
