package api

import (
	"cmsApp/configs"
	"cmsApp/internal/constant"
	"cmsApp/internal/dao"
	"cmsApp/internal/models"
	"cmsApp/pkg/uploader"
	"cmsApp/pkg/utils/filesystem"
	"cmsApp/pkg/utils/snowflake"
	"cmsApp/pkg/utils/strings"
	"errors"
	_ "image/jpeg"
	"net/url"
	"path"
	str "strings"
	"sync"
	"time"
)

type apiImgsService struct {
	Dao    *dao.ImgsDao
	TmpDao *dao.ImgsTempDao
}

var (
	instanceApiImgsService *apiImgsService
	onceApiImgsService     sync.Once
)

func NewApiImgsService() *apiImgsService {
	onceApiImgsService.Do(func() {
		instanceApiImgsService = &apiImgsService{
			Dao:    dao.NewImgsDao(),
			TmpDao: dao.NewImgsTempDao(),
		}
	})
	return instanceApiImgsService
}

func (ser *apiImgsService) GetUploadDirs(name string) (string, error) {
	if name == "" {
		return "", errors.New(constant.UPLOAD_DIR_ERR)
	}
	var dir1 string
	var dir2 string
	newName := str.Split(name, ".")[0]
	if len(newName) > 10 {
		name1 := strings.Reverse(newName[0:10])
		name2 := newName[10:]
		dir1, dir2 = filesystem.GetUploadDirs2(name1, name2)
	} else {
		dir1, dir2 = filesystem.GetUploadDirs(newName)
	}
	if dir1 == "" || dir2 == "" {
		return "", errors.New(constant.UPLOAD_DIR_ERR)
	}
	pathStr, _ := url.JoinPath(dir1, dir2)
	return pathStr, nil
}

func (ser *apiImgsService) GetImageDirs(name string) (string, error) {
	basePath := configs.App.Upload.BasePath
	if name == "" {
		return basePath + "404.png", nil
	}
	path, err := ser.GetUploadDirs(name)
	if err != nil {
		return "", err
	}
	pathStr, _ := url.JoinPath(basePath, path, name)
	return pathStr, nil
}

func (ser *apiImgsService) SaveImage(articleImg models.AppArticleUploadImage, userId uint64) (uint64, string, string, error) {
	basePath := configs.App.Upload.BasePath
	// 重新生成图片名称
	imageName := snowflake.GenIDString()
	oriFileName := articleImg.File.Filename
	ext := path.Ext(oriFileName)
	imageName = imageName + ext
	dstPath, err := ser.GetUploadDirs(imageName)
	if err != nil {
		return 0, "", imageName, errors.New(constant.UPLOAD_DIR_ERR)
	}
	stor := uploader.LocalStorage{}
	// 上传的完整路径
	pathStr, _ := url.JoinPath(basePath, dstPath)
	filePath, width, height, err := stor.Save(articleImg.File, pathStr, imageName)
	if err != nil {
		stor.RomovePath(basePath, dstPath)
		return 0, pathStr, imageName, errors.New(constant.IMAGE_UPLOAD_ERR)
	}
	// 存储到临时图片表
	img := models.ImgsTemp{}
	img.CreateTime = time.Now()
	img.UpdateTime = time.Now()
	img.Path = pathStr
	img.Name = imageName
	img.Tags = oriFileName
	img.Type = articleImg.Type
	img.Width = width
	img.Height = height
	img.Tags = articleImg.Tags
	img.ResourceId = articleImg.ArticleId
	img.UserId = userId
	imgId, err := ser.TmpDao.CreateImgsTemp(img)
	if err != nil {
		stor.RomoveFile(filePath)
		stor.RomovePath(basePath, dstPath)
		return 0, pathStr, imageName, errors.New(constant.IMAGE_UPLOAD_ERR)
	}
	return imgId, pathStr, imageName, nil
}

func (ser *apiImgsService) GetImagesByUserId(userId uint64, page int, pageSize int) ([]models.ImgsListResp, int, int, error) {
	conds := make(map[string][]interface{})
	exp := []interface{}{"= ?", userId}
	conds["user_id"] = exp
	resList, page, totalPage, err := ser.Dao.GetImgs(conds, page, pageSize)
	imgList := []models.ImgsListResp{}
	for i := 0; i < len(resList); i++ {
		obj := models.ImgsListResp{}
		obj.Id = resList[i].Id
		obj.Name = resList[i].Name
		obj.Tags = resList[i].Tags
		obj.Width = resList[i].Width
		obj.Height = resList[i].Height
		imgList = append(imgList, obj)
	}
	return imgList, page, totalPage, err
}

func (ser *apiImgsService) DeleteImage(name string) (err error) {
	conds := make(map[string][]interface{})
	exp := []interface{}{"=", name}
	conds["name"] = exp
	err = ser.Dao.DeleteImage(conds)
	return
}
