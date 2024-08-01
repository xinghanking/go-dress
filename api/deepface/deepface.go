package deepface

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"go-dress/config"
	"go-dress/database"
	"go-dress/models"
	"go-dress/models/utils"
	"image"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/xinghanking/session"
)

type Params struct {
	Category uint     `form:"category" default:"1"`
	Action   []string `form:"action" default:"0,1"`
}

func Exec(context *gin.Context) {
	if context.Request.Host != "ai.han-dress.cn" {
		context.JSON(http.StatusForbidden, gin.H{})
		return
	}
	var params Params
	err := context.ShouldBind(&params)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		panic(err)
	}
	userId := session.Get("uid")
	if userId == nil {
		//context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"code": 4})
		//return
	}
	uid, err := strconv.Atoi(userId.(string))
	if err != nil {
		panic(err)
	}
	database.GetDb()
	var response map[string]interface{}
	file, header, err := context.Request.FormFile("image")
	if err != nil {
		context.JSON(http.StatusBadRequest, err)
		panic(err)
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)
	avatarData, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	postData := map[string]interface{}{"avatar_data": "data:image/png;base64," + base64.StdEncoding.EncodeToString(avatarData), "actions": []string{"age"}}
	response, err = utils.ReqPost(&config.RemoteDeepfaceApis, postData, &config.InvalidDeepfaceApis)
	if err != nil {
		config.RemoteDeepfaceApis = config.InvalidDeepfaceApis
		config.InvalidDeepfaceApis = []string{}
		response, err = utils.ReqPost(&config.RemoteDeepfaceApis, postData, &config.InvalidDeepfaceApis)
		if err != nil {
			fmt.Println(err)
			context.JSON(http.StatusInternalServerError, err)
			panic(err)
		}
	}
	//response, err = utils.PostJson("http://127.0.0.1:9001/analyze", postData)
	img, err := imaging.Decode(bytes.NewReader(avatarData))
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusInternalServerError, err)
		panic(err)
	}
	width := float64(img.Bounds().Max.X - img.Bounds().Min.X)
	height := float64(img.Bounds().Max.Y - img.Bounds().Min.Y)
	fileSize := header.Size
	fileName := header.Filename
	imageInfo := models.Image{
		ImageName:   fileName,
		Size:        uint(fileSize),
		ExtName:     "png",
		Category:    params.Category,
		OperatorUid: uint(uid),
	}
	imageInfo.OwnerUid = imageInfo.OperatorUid
	imageId := models.CreateImageId(imageInfo)
	dir, f := getUploadDir(imageId)
	src := dir + f + ".png"
	avatarUrl := config.BASE_URL_IMAGE + src
	dir = config.IMAGE_DIR + dir
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			context.JSON(http.StatusInternalServerError, err)
			panic(err)
		}
	}
	if response["code"].(float64) == 1 && response["data"] != nil {
		data := response["data"].(map[string]any)
		region := data["region"].(map[string]any)
		if region != nil {
			reg := map[string]float64{
				"x":      region["x"].(float64),
				"y":      region["y"].(float64),
				"w":      region["w"].(float64),
				"h":      region["h"].(float64),
				"width":  width,
				"height": height,
			}
			rec := utils.ChangeCronRang(reg, 2)
			rect := image.Rect(rec["x"], rec["y"], rec["x"]+rec["w"], rec["y"]+rec["h"])
			err = utils.ImageCrop(img, rect, config.IMAGE_DIR+src)
			if err == nil {
				session.Set("selfie", avatarUrl)
				session.Set("selfie_image_id", imageId)
				context.JSON(http.StatusOK, gin.H{"code": config.RESPONSE_SUCCESS, "data": gin.H{"age": data["age"], "avatar": avatarUrl, "gender": data["gender"]}})
				return
			}
		}
	}
	context.JSON(http.StatusOK, response)
}

func getUploadDir(id uint) (string, string) {
	path := strconv.Itoa(int(id))
	dir := ""
	n := len(path)
	if n < 3 {
		return dir, path
	}
	var f string
	if n%2 == 0 {
		f = path[n-2:]
		n = n - 2
	} else {
		f = path[n-1:]
		n = n - 1
	}
	for i := 0; i < n; i += 2 {
		dir += path[i:i+2] + "/"
	}
	return dir, f
}
