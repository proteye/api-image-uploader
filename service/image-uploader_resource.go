package service

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/nfnt/resize"
	"github.com/proteye/api-image-uploader/api"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const IMAGE_POST_FIELD = "image"
const META_COUNT_FIELD = "image_count"

const THUMB_SUFFIX = "_thumb"
const THUMB_DIR = "thumbs/"

const ORDER_IMAGE_TYPE = "order"
const USER_IMAGE_TYPE = "user"

type ImageUploaderResource struct {
	db     gorm.DB
	config ServiceConfig
}

func (ir *ImageUploaderResource) UploadOrderImage(c *gin.Context) {
	response, apiError := SaveImage(ir, c, ORDER_IMAGE_TYPE)
	if apiError != nil {
		log.Print(apiError.Message)
		c.JSON(500, apiError)
	} else {
		log.Print(*response)
		c.JSON(201, response)
	}
}

func (ir *ImageUploaderResource) UploadUserImage(c *gin.Context) {
	response, apiError := SaveImage(ir, c, USER_IMAGE_TYPE)
	if apiError != nil {
		log.Print(apiError.Message)
		c.JSON(500, apiError)
	} else {
		log.Print(*response)
		c.JSON(201, response)
	}
}

func (ir *ImageUploaderResource) GetAllImages(c *gin.Context) {
	var images []api.Image

	ir.db.Order("created_at desc").Find(&images)

	c.JSON(200, images)
}

func (ir *ImageUploaderResource) GetImage(c *gin.Context) {
	id, err := ir.getId(c)
	if err != nil {
		c.JSON(400, api.NewError("Problem decoding ID sent", 0))
		return
	}

	var image api.Image

	if ir.db.First(&image, id).RecordNotFound() {
		c.JSON(404, api.NewError("Image is not found", 1))
	} else {
		c.JSON(200, image)
	}
}

func (ir *ImageUploaderResource) UpdateImage(c *gin.Context) {
	id, err := ir.getId(c)
	if err != nil {
		c.JSON(400, api.NewError("Problem decoding ID sent", 0))
		return
	}

	var image api.Image

	if c.Bind(&image) != nil {
		c.JSON(400, api.NewError("Problem decoding body", 1))
		return
	}

	image.ID = int32(id)

	var existing api.Image

	if ir.db.First(&existing, id).RecordNotFound() {
		c.JSON(404, api.NewError("Image is not found", 2))
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
		c.JSON(400, api.NewError("Problem decoding ID sent", 0))
		return
	}

	var image api.Image

	if ir.db.First(&image, id).RecordNotFound() {
		c.JSON(404, api.NewError("Image is not found", 1))
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

func UniqFilename(origFilename string, thumbSuffix string) (string, string, error) {
	var err error = nil
	var filename, thumbname string

	if origFilename == "" {
		err = errors.New("Filename is empty")
		return filename, thumbname, err
	}

	h := sha1.New()
	time_now := int(time.Now().UnixNano())
	filename_byte := []byte(origFilename + strconv.Itoa(time_now))
	h.Write(filename_byte)
	uniqid := hex.EncodeToString(h.Sum(nil))
	file_ext := filepath.Ext(origFilename)
	filename = uniqid + file_ext
	thumbname = uniqid + thumbSuffix + file_ext

	return filename, thumbname, err
}

func SaveImage(ir *ImageUploaderResource, c *gin.Context, imageTypeName string) (*api.Response, *api.Error) {
	var response *api.Response
	var api_error *api.Error
	var apiImage api.Image
	var apiImageType api.ImageType
	var apiMeta api.Meta

	if imageTypeName == "" {
		api_error = api.NewError("Function argument imageTypeName is empty", 0)
		return response, api_error
	}

	if ir.db.Where("name = ?", imageTypeName).First(&apiImageType).RecordNotFound() {
		api_error = api.NewError("ImageType not found", 1)
		return response, api_error
	}

	if ir.db.Where("name = ?", META_COUNT_FIELD).First(&apiMeta).RecordNotFound() {
		api_error = api.NewError("Meta is not found", 2)
		return response, api_error
	}

	file, header, err := c.Request.FormFile(IMAGE_POST_FIELD)
	if err != nil {
		api_error = api.NewError(err.Error(), 3)
		return response, api_error
	}

	filename, thumbname, err := UniqFilename(header.Filename, ir.config.Thumb_suffix)
	if err != nil {
		api_error = api.NewError(err.Error(), 4)
		return response, api_error
	}
	file_type := header.Header["Content-Type"][0]
	file_path := ir.config.Upload_dir + apiImageType.Path + "/" + filename
	thumb_path := ir.config.Upload_dir + apiImageType.Path + ir.config.Thumb_dir + "/" + thumbname
	file_url := ir.config.Address + ir.config.Upload_url + apiImageType.Path + "/" + filename
	thumb_url := ir.config.Address + ir.config.Upload_url + apiImageType.Path + ir.config.Thumb_dir + "/" + thumbname

	log.Print(filename)
	log.Print(thumbname)
	log.Print(header.Filename)

	/* Save large image */
	out, err := os.Create(file_path)
	if err != nil {
		api_error = api.NewError(err.Error(), 5)
		return response, api_error
	}

	defer out.Close()

	_, err = io.Copy(out, file)

	if err != nil {
		api_error = api.NewError(err.Error(), 6)
		return response, api_error
	}

	/* Save thumbnail */
	var thumb_img image.Image

	out, err = os.Open(file_path)

	if err != nil {
		api_error = api.NewError(err.Error(), 7)
		return response, api_error
	}

	defer out.Close()

	if file_type == api.JpegType {
		thumb_img, err = jpeg.Decode(out)
		if err != nil {
			api_error = api.NewError(err.Error(), 8)
			return response, api_error
		}
	} else if file_type == api.PngType {
		thumb_img, err = png.Decode(out)
		if err != nil {
			api_error = api.NewError(err.Error(), 9)
			return response, api_error
		}
	} else {
		api_error = api.NewError("Invalid image file type! Only jpg, jpeg, png available", 10)
		return response, api_error
	}

	thumb := resize.Thumbnail(uint(apiImageType.Thumb_width), uint(apiImageType.Thumb_height), thumb_img, resize.Lanczos3)

	out, err = os.Create(thumb_path)

	if err != nil {
		api_error = api.NewError(err.Error(), 11)
		return response, api_error
	}

	defer out.Close()

	if file_type == api.JpegType {
		jpeg.Encode(out, thumb, nil)
	} else if file_type == api.PngType {
		png.Encode(out, thumb)
	}

	apiImage.ImageTypeID = apiImageType.ID
	apiImage.Filename = filename
	apiImage.Original = header.Filename
	apiImage.Mime_type = file_type
	apiImage.Created_at = int32(time.Now().Unix())
	apiImage.Updated_at = int32(time.Now().Unix())

	ir.db.Save(&apiImage)

	if ir.db.Error != nil {
		api_error = api.NewError(err.Error(), 12)
		return response, api_error
	}

	responseImage := api.ResponseImage{ID: apiImage.ID, Thumb: thumb_url, Large: file_url}
	response = &api.Response{Image: responseImage, Count: apiMeta.Value_int + 1}

	return response, api_error
}
