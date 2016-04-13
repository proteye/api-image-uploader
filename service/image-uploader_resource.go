package service

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/proteye/api-image-uploader/api"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const UPLOAD_DIR = "/var/www/levsha/web/uploads/"
const THUMB_DIR = "thumbs/"

type ImageUploaderResource struct {
	db gorm.DB
}

func (ir *ImageUploaderResource) UploadOrderImage(c *gin.Context) {
	var image api.Image
	var imageType api.ImageType
	var meta api.Meta

	if ir.db.Where("name = ?", "order").First(&imageType).RecordNotFound() {
		log.Print("ImageType not found")
		c.JSON(500, gin.H{"error": "ImageType not found"})
	}

	file, header, err := c.Request.FormFile("image")
	if err != nil {
		log.Fatal(err)
	}

	filename, err := UniqFilename(header.Filename)
	if err != nil {
		log.Print(err)
		c.JSON(422, gin.H{"error": err.Error()})
	}
	filetype := header.Header["Content-Type"][0]
	filepath := UPLOAD_DIR + imageType.Path + "/" + filename

	log.Print(filename)
	log.Print(header.Filename)

	out, err := os.Create(filepath)

	if err != nil {
		log.Fatal(err)
	}

	defer out.Close()

	_, err = io.Copy(out, file)

	if err != nil {
		log.Fatal(err)
	}

	image.ImageTypeID = imageType.Id
	image.Filename = filename
	image.Original = header.Filename
	image.Mime_type = filetype
	image.Created_at = int32(time.Now().Unix())
	image.Updated_at = int32(time.Now().Unix())

	ir.db.Save(&image)

	if ir.db.Where("name = ?", "image_count").First(&meta).RecordNotFound() {
		log.Print("Meta not found")
		meta.Value_int = 0
	}

	response := gin.H{"path": filepath, "count": meta.Value_int}
	c.JSON(201, response)
}

func (ir *ImageUploaderResource) UploadUserImage(c *gin.Context) {
	var image api.Image
	var imageType api.ImageType

	if ir.db.Where("name = ?", "user").First(&imageType).RecordNotFound() {
		log.Print("ImageType not found")
		c.JSON(500, gin.H{"error": "ImageType not found"})
	}

	file, header, err := c.Request.FormFile("image")
	if err != nil {
		log.Fatal(err)
	}

	filename, err := UniqFilename(header.Filename)
	if err != nil {
		log.Print(err)
		c.JSON(422, gin.H{"error": err.Error()})
	}
	filetype := header.Header["Content-Type"][0]
	filepath := UPLOAD_DIR + imageType.Path + "/" + filename

	log.Print(filename)
	log.Print(header.Filename)

	out, err := os.Create(filepath)

	if err != nil {
		log.Fatal(err)
	}

	defer out.Close()

	_, err = io.Copy(out, file)

	if err != nil {
		log.Fatal(err)
	}

	image.ImageTypeID = imageType.Id
	image.Filename = filename
	image.Original = header.Filename
	image.Mime_type = filetype
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
	image.Mime_type = api.JpegType
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

func UniqFilename(filename string) (string, error) {
	var err error = nil

	if filename == "" {
		err = errors.New("Filename is empty")
		return filename, err
	}

	h := sha1.New()
	time_now := int(time.Now().UnixNano())
	filename_byte := []byte(filename + strconv.Itoa(time_now))
	h.Write(filename_byte)
	new_filename := hex.EncodeToString(h.Sum(nil)) + filepath.Ext(filename)
	return new_filename, err
}
