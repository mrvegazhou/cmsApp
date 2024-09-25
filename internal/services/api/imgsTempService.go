package api

import (
	"cmsApp/internal/constant"
	"cmsApp/internal/dao"
	"cmsApp/internal/models"
	"cmsApp/pkg/uploader"
	"cmsApp/pkg/utils/filesystem"
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/spf13/cast"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type apiImgsTempService struct {
	Dao    *dao.ImgsTempDao
	ImgDao *dao.ImgsDao
}

var (
	instanceApiImgsTempService *apiImgsTempService
	onceApiImgsTempService     sync.Once
)

func NewApiImgsTempService() *apiImgsTempService {
	onceApiImgsTempService.Do(func() {
		instanceApiImgsTempService = &apiImgsTempService{
			Dao:    dao.NewImgsTempDao(),
			ImgDao: dao.NewImgsDao(),
		}
	})
	return instanceApiImgsTempService
}

// 转移封面图片到正式图片表
func (ser *apiImgsTempService) move2ArticleCoverImg(imgName string) error {
	condition := map[string]interface{}{
		"name": imgName,
		"type": 2,
	}
	imgInfo, err := ser.Dao.GetImgInfo(condition)
	if err != nil {
		filePath, err := NewApiImgsService().GetImageDirs(imgName)
		if err == nil {
			flag, _ := filesystem.FileExists(filePath)
			if flag {
				os.Remove(filePath)
				// 判断文件夹是否为空
				dirPath := filepath.Dir(filePath)
				notNullPath, _ := filesystem.IsDirEmpty(dirPath)
				if notNullPath {
					os.Remove(dirPath)
				}
			}
		}
		return errors.New(constant.ARTICLE_SAVE_COVER_IMG_ERR)
	}
	// 保存进图片表
	imgModel := models.Imgs{}
	copier.CopyWithOption(&imgModel, imgInfo, copier.Option{IgnoreEmpty: true, DeepCopy: true})
	_, err = ser.ImgDao.CreateImage(imgModel)
	if err != nil {
		return errors.New(constant.ARTICLE_SAVE_COVER_IMG_ERR)
	}
	// 删除tmp中图片数据
	ser.DeleteImage(imgInfo.Id)
	return nil
}

// 从临时图片表转移到文章图片表
func (ser *apiImgsTempService) move2ArticleImgs(imgNames []string, resourceId uint64) error {
	conditions := map[string][]interface{}{
		"resource_id": {"= ?", resourceId},
		"type":        {"= ?", 1}, // 文章类型
	}
	imgs, err := ser.Dao.GetImgs(conditions)
	if err == nil {
		foundImgs := []models.ImgsTempFields{}
		imgsPaths := make([]string, 0, len(imgs))
		for _, img := range imgs {
			flag := false
			for _, imgName := range imgNames {
				if imgName == img.Name {
					temp := models.ImgsTempFields{}
					err := copier.Copy(&temp, &img)
					if err == nil {
						foundImgs = append(foundImgs, temp)
					}
					// 说明图片在html中
					flag = true
				}
			}
			// 说明临时图片不在html，需要删除
			if flag == false {
				imgsPaths = append(imgsPaths, filepath.Join(img.Path, img.Name))
			}
		}

		if len(foundImgs) == 0 {
			return nil
		}
		imgModel := []models.Imgs{}
		copier.CopyWithOption(&imgModel, foundImgs, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		err = ser.ImgDao.CreateImages(imgModel)

		if err != nil {
			return errors.New(constant.IMAGE_UPLOAD_ERR)
		} else {
			go func(imgsPaths []string) {
				// 删除临时图片
				for _, filePath := range imgsPaths {
					// 判断文件是否存在
					flag, _ := filesystem.FileExists(filePath)
					if flag {
						os.Remove(filePath)
						// 判断文件夹是否为空
						dirPath := filepath.Dir(filePath)
						notNullPath, _ := filesystem.IsDirEmpty(dirPath)
						if notNullPath {
							os.Remove(dirPath)
						}
					}

				}
				// 删除tmp中图片数据
				ser.DeleteImage(resourceId)
			}(imgsPaths)
		}
	}
	return nil
}

func (ser *apiImgsTempService) DeleteImage(resourceId uint64) (err error) {
	conds := make(map[string][]interface{})
	conds["resource_id"] = []interface{}{"=?", resourceId}
	_, err = ser.Dao.DeleteImage(conds)
	return
}

func (ser *apiImgsTempService) DeleteImages(names []string) (err error) {
	if len(names) == 0 {
		return
	}
	conds := make(map[string][]interface{})
	conds["name"] = []interface{}{"in (?)", names}
	_, err = ser.Dao.DeleteImage(conds)
	return
}

func (ser *apiImgsTempService) GetImages(names []string) (imgs []models.ImgsTemp, err error) {
	if len(names) == 0 {
		return imgs, err
	}
	conds := make(map[string][]interface{})
	conds["name"] = []interface{}{"in (?)", names}
	imgs, err = ser.Dao.GetImgs(conds)
	return
}

func (ser *apiImgsTempService) GetImageInfo(name string) (models.ImgsTemp, error) {
	info := models.ImgsTemp{}
	if &name == nil {
		return info, errors.New(constant.IMAGE_NOT_EXIST_ERR)
	}
	condition := map[string]interface{}{
		"name": name,
	}
	info, err := ser.Dao.GetImgInfo(condition)
	if err != nil {
		return info, err
	}
	return info, nil
}

func (ser *apiImgsTempService) SaveImage(imgReq models.AppImgTempUploadReq, userId uint64) (uint64, string, string, string, error) {
	oriFileName := imgReq.File.Filename
	basePath, imageName, dstPath, pathStr, err := NewApiImgsService().Gen4Upload(oriFileName)
	if err != nil {
		return 0, "", imageName, "", errors.New(constant.UPLOAD_DIR_ERR)
	}

	stor := uploader.LocalStorage{}
	filePath, width, height, err := stor.Save(imgReq.File, pathStr, imageName)
	if err != nil {
		stor.RomovePath(basePath, dstPath)
		return 0, pathStr, imageName, "", errors.New(constant.IMAGE_UPLOAD_ERR)
	}
	// 存储到临时图片表
	img := models.ImgsTemp{}
	img.CreateTime = time.Now()
	img.UpdateTime = time.Now()
	img.Path = pathStr
	img.Name = imageName
	img.Tags = oriFileName
	img.Type = cast.ToUint(imgReq.Type)
	img.Width = width
	img.Height = height
	img.Tags = imgReq.Tags
	img.ResourceId = cast.ToUint64(imgReq.ResourceId)
	img.UserId = userId
	imgId, err := ser.Dao.CreateImgsTemp(img)
	if err != nil {
		stor.RomoveFile(filePath)
		stor.RomovePath(basePath, dstPath)
		return 0, pathStr, imageName, "", errors.New(constant.IMAGE_UPLOAD_ERR)
	}
	fullPath, _ := url.JoinPath(pathStr, imageName)
	return imgId, pathStr, imageName, fullPath, nil
}

func (ser *apiImgsTempService) DeleteImageById(uid, imgId uint64, imgName string) error {

	conds := make(map[string][]interface{})
	conds["id"] = []interface{}{"=?", imgId}
	conds["user_id"] = []interface{}{"=?", uid}
	conds["name"] = []interface{}{"=?", imgName}
	rowsAffected, err := ser.Dao.DeleteImage(conds)
	if err != nil {
		return errors.New(constant.DEL_IMAGE_ERR)
	}
	// 删除文件
	filePath, err := NewApiImgsService().GetImageDirs(imgName)
	fmt.Println(rowsAffected, filePath, uid, imgId, imgName, "----filePath img----")
	if err != nil {
		return errors.New(constant.DEL_IMAGE_ERR)
	}
	stor := uploader.LocalStorage{}
	stor.RomoveFile(filePath)
	return nil
}
