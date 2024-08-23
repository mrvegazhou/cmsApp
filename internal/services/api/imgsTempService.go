package api

import (
	"cmsApp/internal/constant"
	"cmsApp/internal/dao"
	"cmsApp/internal/models"
	"cmsApp/pkg/utils/filesystem"
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	"os"
	"path/filepath"
	"sync"
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

// 从临时图片表转移到文章图片表
func (ser *apiImgsTempService) move2ArticleImgs(imgNames []string, resourceId uint64) error {
	conditions := map[string][]interface{}{
		"resource_id": {"= ?", resourceId},
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
					fmt.Println(err, img, temp, "--temp--")
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
		fmt.Println(imgsPaths, foundImgs, "---imgsPaths--")
		if len(foundImgs) == 0 {
			return nil
		}
		imgModel := []models.Imgs{}
		copier.CopyWithOption(&imgModel, foundImgs, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		err = ser.ImgDao.CreateImages(imgModel)

		if err != nil {
			return errors.New(constant.IMAGE_UPLOAD_ERR)
		} else {
			var wg sync.WaitGroup
			wg.Add(1)
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
				wg.Done()
			}(imgsPaths)
			wg.Wait()
		}
	}
	return nil
}

func (ser *apiImgsTempService) DeleteImage(resourceId uint64) (err error) {
	conds := make(map[string][]interface{})
	conds["resource_id"] = []interface{}{"=", resourceId}
	err = ser.Dao.DeleteImage(conds)
	return
}
