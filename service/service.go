package service

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/proteye/api-image-uploader/api"
	"time"
)

type Config struct {
	SvcHost    string
	DbUser     string
	DbPassword string
	DbHost     string
	DbName     string
	Service    ServiceConfig
}

type ServiceConfig struct {
	Address      string
	Upload_dir   string
	Upload_url   string
	Thumb_dir    string
	Thumb_suffix string
}

type ImageUploaderService struct {
}

func (s *ImageUploaderService) getDb(cfg Config) (*gorm.DB, error) {
	connectionString := cfg.DbUser + ":" + cfg.DbPassword + "@tcp(" + cfg.DbHost + ":3306)/" + cfg.DbName + "?charset=utf8&parseTime=True"
	db, err := gorm.Open("mysql", connectionString)

	if err != nil {
		panic("failed to connect database")
	}

	return db, nil
}

func (s *ImageUploaderService) Migrate(cfg Config) error {
	db, err := s.getDb(cfg)
	if err != nil {
		return err
	}

	db.SingularTable(true)

	db.AutoMigrate(&api.Image{})
	db.AutoMigrate(&api.ImageType{})
	db.AutoMigrate(&api.Meta{})

	db.Model(&api.Image{}).AddForeignKey("image_type_id", "image_type(id)", "RESTRICT", "RESTRICT")

	meta := api.Meta{
		Name:       "image_count",
		Value_int:  0,
		Created_at: int32(time.Now().Unix()),
		Updated_at: int32(time.Now().Unix()),
	}
	db.Create(&meta)

	imageType := api.ImageType{
		Name:         "order",
		Path:         "/image",
		Thumb_width:  320,
		Thumb_height: 240,
		Created_at:   int32(time.Now().Unix()),
		Updated_at:   int32(time.Now().Unix()),
	}
	db.Create(&imageType)

	imageType = api.ImageType{
		Name:         "user",
		Path:         "/user",
		Thumb_width:  150,
		Thumb_height: 150,
		Created_at:   int32(time.Now().Unix()),
		Updated_at:   int32(time.Now().Unix()),
	}
	db.Create(&imageType)

	return nil
}
func (s *ImageUploaderService) Run(cfg Config) error {
	db, err := s.getDb(cfg)
	if err != nil {
		return err
	}
	db.SingularTable(true)

	imageUploaderResource := &ImageUploaderResource{db: *db, config: cfg.Service}

	r := gin.Default()

	r.GET("/images", imageUploaderResource.GetAllImages)
	r.GET("/images/:id", imageUploaderResource.GetImage)
	r.POST("/images/order-upload", imageUploaderResource.UploadOrderImage)
	r.POST("/images/user-upload", imageUploaderResource.UploadUserImage)
	r.PUT("/images/:id", imageUploaderResource.UpdateImage)
	r.DELETE("/images/:id", imageUploaderResource.DeleteImage)

	r.Run(cfg.SvcHost)

	return nil
}
