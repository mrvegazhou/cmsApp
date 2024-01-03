package slideCaptcha

import (
	"cmsApp/pkg/AES"
	"cmsApp/pkg/utils/random"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/image/colornames"
	"image/color"
	"log"
	"math"
	"strconv"
	"strings"
)

type BlockPuzzleCaptchaService struct {
	point   SlideCaptchaPoint
	factory *CaptchaServiceFactory
}

func NewBlockPuzzleCaptchaService(factory *CaptchaServiceFactory) *BlockPuzzleCaptchaService {
	// 初始化静态资源
	ImgSetUp()

	return &BlockPuzzleCaptchaService{
		factory: factory,
	}
}

func (b *BlockPuzzleCaptchaService) Check(token string, pointJson string) error {
	cache := b.factory.GetCache()

	codeKey := fmt.Sprintf(CodeKeyPrefix, token)

	cachePointInfo := cache.Get(codeKey)

	if cachePointInfo == "" {
		return errors.New("验证码已失效")
	}

	// 解析结构体
	cachePoint := &SlideCaptchaPoint{}
	userPoint := &SlideCaptchaPoint{}
	err := json.Unmarshal([]byte(cachePointInfo), cachePoint)

	if err != nil {
		return err
	}

	// 解密前端传递过来的数据
	userPointJson := AES.Decrypt(pointJson, cachePoint.SecretKey)

	err = json.Unmarshal([]byte(userPointJson), userPoint)

	if err != nil {
		return err
	}

	// 校验两个点是否符合
	if math.Abs(float64(cachePoint.X-userPoint.X)) <= float64(b.factory.config.BlockPuzzle.Offset) && cachePoint.Y == userPoint.Y {
		return nil
	}

	return errors.New("验证失败")
}

func (b *BlockPuzzleCaptchaService) Verification(token string, pointJson string) error {
	err := b.Check(token, pointJson)
	if err != nil {
		return err
	}
	codeKey := fmt.Sprintf(CodeKeyPrefix, token)
	b.factory.GetCache().Delete(codeKey)
	return nil
}

// Get 获取验证码图片信息
func (b *BlockPuzzleCaptchaService) Get() (map[string]interface{}, error) {
	// 初始化背景图片
	backgroundImage := GetBackgroundImage()

	// 为背景图片设置水印
	colorStr := strings.Split(b.factory.config.Watermark.Color, ",")
	colorR, _ := strconv.ParseUint(colorStr[0], 10, 8)
	colorG, _ := strconv.ParseUint(colorStr[1], 10, 8)
	colorB, _ := strconv.ParseUint(colorStr[2], 10, 8)
	colorA, _ := strconv.ParseUint(colorStr[3], 10, 8)
	colorRGBA := color.RGBA{
		R: uint8(colorR),
		G: uint8(colorG),
		B: uint8(colorB),
		A: uint8(colorA),
	}
	backgroundImage.SetText(b.factory.config.Watermark.Text, b.factory.config.Watermark.FontSize, colorRGBA)

	// 初始化模板图片
	templateImage := GetTemplateImage()

	// 构造前端所需图片
	b.pictureTemplatesCut(backgroundImage, templateImage)

	originalImageBase64, err := backgroundImage.Base64()
	jigsawImageBase64, err := templateImage.Base64()

	if err != nil {
		return nil, err
	}

	data := make(map[string]interface{})
	data["originalImageBase64"] = originalImageBase64
	data["jigsawImageBase64"] = jigsawImageBase64
	data["secretKey"] = b.point.SecretKey
	data["token"] = random.GetUuid()

	codeKey := fmt.Sprintf(CodeKeyPrefix, data["token"])
	jsonPoint, err := json.Marshal(b.point)
	if err != nil {
		log.Printf("point json Marshal err: %v", err)
		return nil, err
	}

	b.factory.GetCache().Set(codeKey, string(jsonPoint), b.factory.config.CacheExpireSec)

	return data, nil
}

func (b *BlockPuzzleCaptchaService) pictureTemplatesCut(backgroundImage *ImageUtil, templateImage *ImageUtil) {
	// 生成拼图坐标点
	b.generateJigsawPoint(backgroundImage, templateImage)
	// 裁剪模板图
	b.cutByTemplate(backgroundImage, templateImage, b.point.X, 0)

	// 插入干扰图
	for {
		newTemplateImage := GetTemplateImage()
		if newTemplateImage.Src != templateImage.Src {
			offsetX := random.RandomInt(0, backgroundImage.Width-newTemplateImage.Width-5)
			if math.Abs(float64(newTemplateImage.Width-offsetX)) > float64(newTemplateImage.Width/2) {
				b.interferenceByTemplate(backgroundImage, newTemplateImage, offsetX, b.point.Y)
				break
			}
		}
	}
}

// 插入干扰图
func (b *BlockPuzzleCaptchaService) interferenceByTemplate(backgroundImage *ImageUtil, templateImage *ImageUtil, x1 int, y1 int) {
	xLength := templateImage.Width
	yLength := templateImage.Height

	for x := 0; x < xLength; x++ {
		for y := 0; y < yLength; y++ {
			// 如果模板图像当前像素点不是透明色 copy源文件信息到目标图片中
			isOpacity := templateImage.IsOpacity(x, y)

			// 当前模板像素在背景图中的位置
			backgroundX := x + x1
			backgroundY := y + y1

			// 当不为透明时
			if !isOpacity {
				// 背景图区域模糊
				backgroundImage.VagueImage(backgroundX, backgroundY)
			}

			//防止数组越界判断
			if x == (xLength-1) || y == (yLength-1) {
				continue
			}

			rightOpacity := templateImage.IsOpacity(x+1, y)
			downOpacity := templateImage.IsOpacity(x, y+1)

			//描边处理，,取带像素和无像素的界点，判断该点是不是临界轮廓点,如果是设置该坐标像素是白色
			if (isOpacity && !rightOpacity) || (!isOpacity && rightOpacity) || (isOpacity && !downOpacity) || (!isOpacity && downOpacity) {
				backgroundImage.RgbaImage.SetRGBA(backgroundX, backgroundY, colornames.White)
			}
		}
	}
}

func (b *BlockPuzzleCaptchaService) cutByTemplate(backgroundImage *ImageUtil, templateImage *ImageUtil, x1, y1 int) {
	xLength := templateImage.Width
	yLength := templateImage.Height

	for x := 0; x < xLength; x++ {
		for y := 0; y < yLength; y++ {
			// 如果模板图像当前像素点不是透明色 copy源文件信息到目标图片中
			isOpacity := templateImage.IsOpacity(x, y)

			// 当前模板像素在背景图中的位置
			backgroundX := x + x1
			backgroundY := y + y1

			// 当不为透明时
			if !isOpacity {
				// 获取原图像素
				backgroundRgba := backgroundImage.RgbaImage.RGBAAt(backgroundX, backgroundY)
				// 将原图的像素扣到模板图上
				templateImage.SetPixel(backgroundRgba, x, y)
				// 背景图区域模糊
				backgroundImage.VagueImage(backgroundX, backgroundY)
			}

			//防止数组越界判断
			if x == (xLength-1) || y == (yLength-1) {
				continue
			}

			rightOpacity := templateImage.IsOpacity(x+1, y)
			downOpacity := templateImage.IsOpacity(x, y+1)

			//描边处理，,取带像素和无像素的界点，判断该点是不是临界轮廓点,如果是设置该坐标像素是白色
			if (isOpacity && !rightOpacity) || (!isOpacity && rightOpacity) || (isOpacity && !downOpacity) || (!isOpacity && downOpacity) {
				templateImage.RgbaImage.SetRGBA(x, y, colornames.White)
				backgroundImage.RgbaImage.SetRGBA(backgroundX, backgroundY, colornames.White)
			}
		}
	}
}

// 生成模板图在背景图中的随机坐标点
func (b *BlockPuzzleCaptchaService) generateJigsawPoint(backgroundImage *ImageUtil, templateImage *ImageUtil) {
	widthDifference := backgroundImage.Width - templateImage.Width
	heightDifference := backgroundImage.Height - templateImage.Height

	x, y := 0, 0

	if widthDifference <= 0 {
		x = 5
	} else {
		x = random.RandomInt(100, widthDifference-100)
	}
	if heightDifference <= 0 {
		y = 5
	} else {
		y = random.RandomInt(5, heightDifference)
	}
	point := SlideCaptchaPoint{X: x, Y: y}
	point.SetSecretKey(random.RandString(16))
	b.point = point
}
