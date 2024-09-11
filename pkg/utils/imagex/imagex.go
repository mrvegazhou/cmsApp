package imagex

import (
	"bytes"
	"cmsApp/configs"
	"cmsApp/internal/constant"
	"cmsApp/pkg/utils/filesystem"
	"errors"
	"golang.org/x/image/webp"
	"image"
	"image/draw"
	_ "image/gif" // 导入需要支持的图片格式
	_ "image/jpeg"
	"image/png"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// OpenPngImage 打开png图片
func OpenImage(src string) image.Image {
	ff, err := os.Open(src)
	if err != nil {
		log.Printf("打开 %s 图片失败: %v", src, err)
	}

	img, err := png.Decode(ff)

	if err != nil {
		log.Printf("png %s decode  失败: %v", src, err)
	}

	return img
}

// ImageToRGBA 图片转rgba
func ImageToRGBA(img image.Image) *image.RGBA {
	// No conversion needed if image is an *image.RGBA.
	if dst, ok := img.(*image.RGBA); ok {
		return dst
	}

	// Use the image/draw package to convert to *image.RGBA.
	b := img.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(dst, dst.Bounds(), img, b.Min, draw.Src)
	return dst
}

// 检查图片后缀
func CheckImageExt(fileName string) bool {
	ext := filesystem.GetExt(fileName)
	for _, allowExt := range configs.App.Upload.ImageAllowExts {
		if strings.ToUpper(allowExt) == strings.ToUpper(ext) {
			return true
		}
	}
	return false
}

// 检查图片大小
func CheckImageSize(f multipart.File) bool {
	size, err := filesystem.GetSize(f)
	if err != nil {
		return false
	}
	return size <= configs.App.Upload.ImageMaxSize
}

func CheckImage(src string) error {
	isExists, _ := filesystem.FileExists(src)
	if isExists == false {
		return errors.New(constant.FILE_NOT_EXIST_ERR)
	}
	perm := filesystem.CheckPermission(src)
	if perm == true {
		return errors.New(constant.FILE_PERMISSION_ERR)
	}
	flag := IsImage(src)
	if flag == false {
		return errors.New(constant.DECODE_IMG_ERR)
	}
	return nil
}

func IsImage(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer file.Close()
	_, typ, err := image.Decode(file)
	if err != nil {
		return false
	} else {
		// 获取图片的类型
		switch typ {
		case `jpeg`:
		case `png`:
		case `gif`:
		case `bmp`:
		default:
			data, _ := os.ReadFile(filePath)
			// 尝试以 webp 进行解码
			_, err := webp.Decode(bytes.NewReader(data))
			if err != nil {
				return false
			}
		}
	}
	return true
}

func ExtractImageSrcs(htmls string) []string {
	var imgRE = regexp.MustCompile(`<img[^>]+\bsrc=["']([^"']+)["']`)
	imgs := imgRE.FindAllStringSubmatch(htmls, -1)
	out := make([]string, 0, len(imgs))
	for i := range imgs {
		out = append(out, filepath.Base(imgs[i][1]))
	}
	return out
}
