package api

import (
	"cmsApp/configs"
	"cmsApp/internal/constant"
	"cmsApp/internal/dao"
	"cmsApp/internal/models"
	"cmsApp/pkg/loggers"
	"cmsApp/pkg/utils/filesystem"
	HTMLX "cmsApp/pkg/utils/htmlx"
	"cmsApp/pkg/utils/snowflake"
	"cmsApp/pkg/utils/stringx"
	"context"
	"errors"
	"github.com/jinzhu/copier"
	"github.com/spf13/cast"
	_ "image/jpeg"
	"net/url"
	"path"
	"path/filepath"
	str "strings"
	"sync"
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

// 生成上传的文件夹
func (ser *apiImgsService) GetUploadDirs(name string) (string, error) {
	if name == "" {
		return "", errors.New(constant.UPLOAD_DIR_ERR)
	}
	var dir1 string
	var dir2 string
	newName := str.Split(name, ".")[0]
	if len(newName) > 10 {
		name1 := stringx.Reverse(newName[0:10])
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

// todo: 此处可加redis缓存
func (ser *apiImgsService) CheckImageIsYours(userId uint64, name string) (bool, error) {
	condition := map[string]interface{}{
		"name": name,
	}
	imgInfo, err := ser.Dao.GetImgInfo(condition)
	if err != nil {
		return false, errors.New(constant.IMAGE_NOT_EXIST_ERR)
	}
	if imgInfo.UserId != userId {
		return false, nil
	}
	return true, nil
}

// 生成新的文件名 存储地址 完整地址
func (ser *apiImgsService) Gen4Upload(fileName string) (basePath, imageName, dstPath, pathStr string, err error) {
	basePath = configs.App.Upload.BasePath
	// 重新生成图片名称
	imageName = snowflake.GenIDString()
	ext := path.Ext(fileName)
	imageName = imageName + ext
	dstPath, err = ser.GetUploadDirs(imageName)
	// 上传的完整路径
	pathStr, _ = url.JoinPath(basePath, dstPath)
	return
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

func (ser *apiImgsService) UploadImage(req models.AppImgTempUploadReq, userId uint64) (imgId, resourceId uint64, fullPath string, imgName string, fileName string, err error) {
	fileName = req.File.Filename
	resourceId = cast.ToUint64(req.ResourceId)

	imgId, _, imgName, fullPath, err = NewApiImgsTempService().SaveImage(req, userId)
	return
}

// 替换评论里的图片，并转移图片路径到正式图片库
func (ser *apiImgsService) HandleImgs(htmlContent string) (string, error) {
	imgUrls := HTMLX.GetImgSrcs(htmlContent)
	newContent := HTMLX.AppendParamToImageSrc(htmlContent)
	if len(imgUrls) == 0 {
		return newContent, nil
	}

	go func(imgUrls []string) {
		var names []string
		for _, img := range imgUrls {
			filename := filepath.Base(img)
			if filename != "" {
				names = append(names, filename)
			}
		}
		var imgsTemps []models.ImgsTemp
		imgsTemps, err := NewApiImgsTempService().GetImages(names)
		if err != nil {
			loggers.LogError(context.Background(), "service", "获取临时图片失败", map[string]string{"error": err.Error()})
			return
		}
		imgModels := []models.Imgs{}
		for _, img := range imgsTemps {
			imgModel := models.Imgs{}
			img.Id = 0
			copier.CopyWithOption(&imgModel, img, copier.Option{IgnoreEmpty: true, DeepCopy: true})
			imgModels = append(imgModels, imgModel)
		}
		// 保存到正式表
		err = ser.Dao.CreateImages(imgModels)
		if err != nil {
			loggers.LogError(context.Background(), "service", "创建图片信息失败", map[string]string{"error": err.Error()})
			return
		}
		// 删除临时表
		err = NewApiImgsTempService().DeleteImages(names)
		if err != nil {
			loggers.LogError(context.Background(), "service", "删除临时表失败", map[string]string{"error": err.Error()})
			return
		}
	}(imgUrls)

	return newContent, nil
}
