package service

import (
	"api-image-uploader/api"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type Config struct {
	SvcHost    string
	DbUser     string
	DbPassword string
	DbHost     string
	DbName     string
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
	return nil
}
func (s *ImageUploaderService) Run(cfg Config) error {
	db, err := s.getDb(cfg)
	if err != nil {
		return err
	}
	db.SingularTable(true)

	imageUploaderResource := &ImageUploaderResource{db: *db}

	r := gin.Default()

	r.GET("/images", imageUploaderResource.GetAllImages)
	r.GET("/images/:id", imageUploaderResource.GetImage)
	r.POST("/images", imageUploaderResource.CreateImage)
	r.POST("/images/upload", imageUploaderResource.UploadImage)
	r.PUT("/images/:id", imageUploaderResource.UpdateImage)
	r.DELETE("/images/:id", imageUploaderResource.DeleteImage)

	r.Run(cfg.SvcHost)

	return nil
}
