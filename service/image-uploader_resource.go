package service

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/proteye/api-image-uploader/api"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

type ImageUploaderResource struct {
	db gorm.DB
}

func (ir *ImageUploaderResource) UploadImage(c *gin.Context) {
	var image api.Image
	file, header, err := c.Request.FormFile("image")
	filename := header.Filename

	log.Print(header.Filename)

	out, err := os.Create("/tmp/" + filename)

	if err != nil {
		log.Fatal(err)
	}

	defer out.Close()

	_, err = io.Copy(out, file)

	if err != nil {
		log.Fatal(err)
	}

	image.File = filename
	image.Type = api.PngType
	image.Created_at = int32(time.Now().Unix())
	image.Updated_at = int32(time.Now().Unix())

	ir.db.Save(&image)

	c.JSON(201, image)
}

func (ir *ImageUploaderResource) CreateImage(c *gin.Context) {
	var image api.Image

	if c.Bind(&image) != nil {
		c.JSON(400, api.NewError("problem decoding body"))
		return
	}
	image.Type = api.JpgType
	image.Created_at = int32(time.Now().Unix())
	image.Updated_at = int32(time.Now().Unix())

	ir.db.Save(&image)

	c.JSON(201, image)
}

func (ir *ImageUploaderResource) GetAllImages(c *gin.Context) {
	var images []api.Image

	ir.db.Order("created_at desc").Find(&images)

	c.JSON(200, images)
}

func (ir *ImageUploaderResource) GetImage(c *gin.Context) {
	id, err := ir.getId(c)
	if err != nil {
		c.JSON(400, api.NewError("problem decoding id sent"))
		return
	}

	var image api.Image

	if ir.db.First(&image, id).RecordNotFound() {
		c.JSON(404, gin.H{"error": "not found"})
	} else {
		c.JSON(200, image)
	}
}

func (ir *ImageUploaderResource) UpdateImage(c *gin.Context) {
	id, err := ir.getId(c)
	if err != nil {
		c.JSON(400, api.NewError("problem decoding id sent"))
		return
	}

	var image api.Image

	if c.Bind(&image) != nil {
		c.JSON(400, api.NewError("problem decoding body"))
		return
	}
	image.Id = int32(id)

	var existing api.Image

	if ir.db.First(&existing, id).RecordNotFound() {
		c.JSON(404, api.NewError("not found"))
	} else {
		image.Created_at = existing.Created_at
		image.Updated_at = int32(time.Now().Unix())
		ir.db.Save(&image)
		c.JSON(200, image)
	}

}

func (ir *ImageUploaderResource) DeleteImage(c *gin.Context) {
	id, err := ir.getId(c)
	if err != nil {
		c.JSON(400, api.NewError("problem decoding id sent"))
		return
	}

	var image api.Image

	if ir.db.First(&image, id).RecordNotFound() {
		c.JSON(404, api.NewError("not found"))
	} else {
		ir.db.Delete(&image)
		c.Data(204, "application/json", make([]byte, 0))
	}
}

func (ir *ImageUploaderResource) getId(c *gin.Context) (int32, error) {
	idStr := c.Params.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Print(err)
		return 0, err
	}
	return int32(id), nil
}
