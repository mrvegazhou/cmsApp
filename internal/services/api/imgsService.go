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
	"net/url"
	"path"
	str "strings"
	"sync"
	"time"
)

type apiImgsService struct {
	Dao *dao.ImgsDao
}

var (
	instanceApiImgsService *apiImgsService
	onceApiImgsService     sync.Once
)

func NewApiImgsService() *apiImgsService {
	onceApiImgsService.Do(func() {
		instanceApiImgsService = &apiImgsService{
			Dao: dao.NewImgsDao(),
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

func (ser *apiImgsService) SaveImage(articleImg models.AppArticleUploadImage) (uint64, string, string, error) {
	basePath := configs.App.Upload.BasePath
	imageName := snowflake.GenIDString()
	ext := path.Ext(articleImg.File.Filename)
	imageName = imageName + ext
	dstPath, err := ser.GetUploadDirs(imageName)
	if err != nil {
		return 0, "", imageName, errors.New(constant.UPLOAD_DIR_ERR)
	}
	stor := uploader.LocalStorage{}
	// 完整路径
	pathStr, _ := url.JoinPath(basePath, dstPath)
	filePath, err := stor.Save(articleImg.File, pathStr, imageName)
	if err != nil {
		stor.RomovePath(basePath, dstPath)
		return 0, pathStr, imageName, errors.New(constant.IMAGE_UPLOAD_ERR)
	}
	img := models.Imgs{}
	img.CreateTime = time.Now()
	img.UpdateTime = time.Now()
	img.Path = pathStr
	img.Name = imageName
	img.Type = 1
	img.Tags = articleImg.Tags
	img.ResourceId = articleImg.ArticleId
	imgId, err := ser.Dao.CreateImage(img)
	if err != nil {
		stor.RomoveFile(filePath)
		stor.RomovePath(basePath, dstPath)
		return 0, pathStr, imageName, errors.New(constant.IMAGE_UPLOAD_ERR)
	}
	return imgId, pathStr, imageName, nil
}
