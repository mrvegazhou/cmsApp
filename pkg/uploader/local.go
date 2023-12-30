package uploader

import (
	"io"
	"mime/multipart"
	"os"
	"strings"
)

type LocalStorage struct {
}

func (stor LocalStorage) Save(file *multipart.FileHeader, dst string, fileName string) (string, error) {

	var (
		filePath string
	)
	name := file.Filename
	if fileName != "" {
		name = fileName
	}
	filePath = dst + string(os.PathSeparator) + name

	src, err := file.Open()
	if err != nil {
		return filePath, err
	}
	defer src.Close()

	err = os.MkdirAll(dst, os.ModePerm)
	if err != nil {
		return filePath, err
	}

	out, err := os.Create(filePath)
	if err != nil {
		return filePath, err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return filePath, err
}

// 检查指定路径是否为文件夹
func (stor LocalStorage) IsDir(name string) bool {
	if info, err := os.Stat(name); err == nil {
		return info.IsDir()
	}
	return false
}

func (stor LocalStorage) RomovePath(basePath string, dst string) {
	if basePath == "" {
		return
	}
	if stor.IsDir(dst) {
		pathArr := strings.Split(dst, string(os.PathSeparator))
		tempPath := basePath
		for _, value := range pathArr {
			tempPath = tempPath + string(os.PathSeparator) + value
			dir, _ := os.ReadDir(tempPath)
			if len(dir) == 0 {
				os.Remove(dst)
			}
		}
	}
}

func (stor LocalStorage) RomoveFile(filePath string) {
	if filePath == "" {
		return
	}
	os.Remove(filePath)
}

func main() {
	//fmt.Println("begin:")
	//stor := LocalStorage{}
	//tmp := "/Users/vega/workspace/codes/golang_space/gopath/src/work/cmsApp/uploadfile"
	//stor.RomovePath(tmp)
}
