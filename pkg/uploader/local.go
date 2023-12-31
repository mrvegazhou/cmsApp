package uploader

import (
	"image"
	"io"
	"mime/multipart"
	"os"
	"strings"
)

type LocalStorage struct {
}

func (stor LocalStorage) Save(file *multipart.FileHeader, dst string, fileName string) (string, int, int, error) {

	var (
		filePath string
	)
	width := 0
	height := 0

	name := file.Filename
	if fileName != "" {
		name = fileName
	}
	filePath = dst + string(os.PathSeparator) + name

	src, err := file.Open()
	if err != nil {
		return filePath, width, height, err
	}
	defer src.Close()

	im, _, err := image.DecodeConfig(src)
	if err == nil {
		width = im.Width
		height = im.Height
	}

	err = os.MkdirAll(dst, os.ModePerm)
	if err != nil {
		return filePath, width, height, err
	}

	out, err := os.Create(filePath)
	if err != nil {
		return filePath, width, height, err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return filePath, width, height, err
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
