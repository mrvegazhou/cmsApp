package watermark

import (
	"errors"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// 水印的位置
const (
	TopLeft Pos = iota
	TopRight
	BottomLeft
	BottomRight
	Center
	Tile
)

const (
	fontSize = 20
	padding  = 2
	angle    = 30
	//fontFile = "/Users/vega/workspace/codes/golang_space/gopath/src/work/cmsApp/uploadfile/SimSun.ttf"
)

// 允许做水印的图片类型
var allowExts = []string{
	".gif", ".jpg", ".jpeg", ".png",
}

// Pos 表示水印的位置
type Pos int

type Watermark struct {
	image          image.Image // 需要被水印的图片
	gifImg         *gif.GIF
	waterImagePath string // 水印图片地址
	//padding        int         `default:2` // 水印留的边白
	pos  Pos    `default:3` // 水印的位置
	text string // 水印内容
	//fontSize       int         `default:12`    // 水印字体大小
	//rotate int  `default:30`    // 如果是平铺水印 角度设置
	//isResize       bool        `default:true`  // 是否自动改变水印图片大小
	bgColor color.Color // 水印图片背景色
	ext     string
}

func IsAllowExt(ext string) bool {
	if ext == "" {
		panic("参数 ext 不能为空")
	}

	ext = strings.ToLower(ext)

	for _, e := range allowExts {
		if e == ext {
			return true
		}
	}
	return false
}

// r 为需要被水印的图片
func New(r io.Reader, ext string, pos Pos, text string, waterImagePath string) (w *Watermark, err error) {
	if !IsAllowExt(ext) {
		panic("无效的图片")
	}

	if pos != Tile {
		if pos < TopLeft || pos > Center {
			panic("无效的 pos 值")
		}
	}

	var gifImg *gif.GIF
	// 图片
	var img image.Image
	switch strings.ToLower(ext) {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(r)
	case ".png":
		img, err = png.Decode(r)
	case ".gif":
		gifImg, err = gif.DecodeAll(r)
		img = gifImg.Image[0]
	default:
		return nil, fmt.Errorf("格式错误")
	}

	return &Watermark{
		image:          img,
		gifImg:         gifImg,
		pos:            pos,
		text:           text,
		waterImagePath: waterImagePath,
		ext:            ext,
	}, nil
}

func (w *Watermark) WaterMakeDone(path string) (err error) {
	// 水印
	var waterImg image.Image
	var gifWaterImg *gif.GIF

	if w.text != "" {
		// 生成文字图片
		waterImg, err = w.makeImageByText(w.text)
	} else {
		file, err := os.Open(w.waterImagePath)
		if err != nil {
			return err
		}
		extension := filepath.Ext(w.waterImagePath)
		if extension == "" {
			return fmt.Errorf("%s has to be of type png, jpeg or gif", w.waterImagePath)
		}
		switch extension {
		case ".jpg", ".jpeg":
			waterImg, err = jpeg.Decode(file)
			break
		case ".png":
			waterImg, err = png.Decode(file)
			break
		case ".gif":
			gifWaterImg, err = gif.DecodeAll(file)
			// 如果原始图片不是gif 水印是gif 则把gif的第一帧水印到原始图片上
			waterImg = gifWaterImg.Image[0]
			break
		default:
			return fmt.Errorf("%s has to be of type png, jpeg or gif", w.waterImagePath)
		}
		if err != nil {
			return err
		}
	}
	// 判断尺寸
	bound := w.image.Bounds()
	var point image.Point
	width := bound.Dx()
	height := bound.Dy()

	if w.pos != Tile {
		// 通过布局方式判断坐标是否超过被水印的图片
		point = w.getPoint(width, height, waterImg)
		if err = w.checkTooLarge(point, bound); err != nil {
			return err
		}
	}
	return w.waterMarkDoing(path, waterImg, bound, point, gifWaterImg)
}

// 平铺水印图片
func (w *Watermark) waterMakeTile(path string, waterImg image.Image, bound image.Rectangle) error {
	watermarkedImage := image.NewRGBA(bound)
	draw.Draw(watermarkedImage, bound, w.image, image.Point{}, draw.Src)
	rotatedWatermark := imaging.Rotate(waterImg, angle, color.Transparent)
	x, y := 0, 0
	for y <= watermarkedImage.Bounds().Max.Y {
		for x <= watermarkedImage.Bounds().Max.X {
			offset := image.Pt(x, y)
			draw.Draw(watermarkedImage, rotatedWatermark.Bounds().Add(offset), rotatedWatermark, image.Point{}, draw.Over)
			// 稀疏一点, 稍微提升点速度
			x += rotatedWatermark.Bounds().Dx() * 2
		}
		y += rotatedWatermark.Bounds().Dy()
		x = 0
	}
	err := w.saveImage(watermarkedImage, path)
	return err
}

func (w *Watermark) waterMakePoint(path string, waterImg image.Image, bound image.Rectangle, point image.Point) error {
	watermarkedImage := image.NewRGBA(bound)
	draw.Draw(watermarkedImage, bound, w.image, image.Point{}, draw.Src)
	draw.Draw(watermarkedImage, bound, waterImg, point, draw.Over)
	err := w.saveImage(watermarkedImage, path)
	return err
}

// 给gif打水印
func (w *Watermark) waterMakePoint4Gif(path string, waterImg image.Image, point image.Point) error {
	for index, img := range w.gifImg.Image {
		watermarkedImage := image.NewPaletted(img.Bounds(), img.Palette)
		draw.Draw(watermarkedImage, watermarkedImage.Bounds(), img, image.Point{}, draw.Src)

		width := watermarkedImage.Bounds().Dx()
		height := watermarkedImage.Bounds().Dy()
		if waterImg.Bounds().Dx() > width || waterImg.Bounds().Dy() > height {
			waterImg, _ = w.resizeImage(waterImg, height, width)
		}
		draw.Draw(watermarkedImage, watermarkedImage.Bounds(), waterImg, point, draw.Over)
		w.gifImg.Image[index] = watermarkedImage
	}
	err := w.saveGif(w.gifImg, path)
	return err
}

// Mark 将水印写入 src 中，由 ext 确定当前图片的类型。
// bound 被水印图片的大小
func (w *Watermark) waterMarkDoing(path string, waterImg image.Image, bound image.Rectangle, point image.Point, gifWaterImg *gif.GIF) error {
	var err error
	ext := strings.ToLower(w.ext)
	// 原始图片是gif
	// gif没有平铺设置
	if ext == ".gif" {
		if w.pos == Tile {
			return fmt.Errorf("gif图片不允许平铺设置")
		}
		// 水印不是gif
		if gifWaterImg == nil {
			bound := w.gifImg.Image[0].Bounds()
			point := w.getPoint(bound.Dx(), bound.Dy(), waterImg)
			if err = w.checkTooLarge(point, bound); err != nil {
				return err
			}
			err = w.waterMakePoint4Gif(path, waterImg, point)
		} else { // 水印也是 GIF
			windex := 0
			wmax := len(w.gifImg.Image)
			for index, img := range w.gifImg.Image {
				dstImg := image.NewPaletted(img.Bounds(), img.Palette)
				draw.Draw(dstImg, dstImg.Bounds(), img, image.Point{}, draw.Src)

				// 获取对应帧数的水印图片
				if windex >= wmax {
					windex = 0
				}
				draw.Draw(dstImg, dstImg.Bounds(), gifWaterImg.Image[windex], point, draw.Over)
				windex++
				w.gifImg.Image[index] = dstImg
			}
			err = w.saveGif(w.gifImg, path)
		}
	} else {
		// 调整水印图片大小
		width := bound.Dx()
		height := bound.Dy()
		if waterImg.Bounds().Dx() > width || waterImg.Bounds().Dy() > height {
			waterImg, err = w.resizeImage(waterImg, height, width)
			if err != nil {
				return err
			}
		}
		if w.pos == Tile {
			err = w.waterMakeTile(path, waterImg, bound)
		} else {
			err = w.waterMakePoint(path, waterImg, bound, point)
		}
	}
	return err
}

func (w *Watermark) saveGif(gifImg *gif.GIF, path string) error {
	file, err := os.Create(path)
	err = gif.EncodeAll(file, gifImg)
	return err
}

// SaveImage Saves an image file into the secondary storage
func (w *Watermark) saveImage(img image.Image, path string) error {
	var extension string

	// read raw file
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	// Parse file extension
	extension = filepath.Ext(path)
	if extension == "" {
		return fmt.Errorf("%s has to be of type png, jpeg or gif", path)
	}
	switch extension {
	case ".jpg", ".jpeg":
		err = jpeg.Encode(file, img, nil)
		break
	case ".png":
		err = png.Encode(file, img)
		break
	default:
		return fmt.Errorf("%s has to be of type png, jpeg or gif", path)
	}
	return err
}

// 改变原始图片大小
func (w *Watermark) resizeImage(img image.Image, height, width int) (image.Image, error) {
	if img == nil {
		return nil, errors.New("water make image is nil")
	}
	currentBounds := img.Bounds()
	newBounds := image.Rect(0, 0, width, height)
	newImage := image.NewNRGBA(newBounds)
	for i := 0; i < newBounds.Dx(); i++ {
		for j := 0; j < newBounds.Dy(); j++ {
			atX := int(float64(i) * float64(currentBounds.Dx()) / float64(newBounds.Dx()))
			atY := int(float64(j) * float64(currentBounds.Dy()) / float64(newBounds.Dy()))
			colorAt := img.At(atX, atY)
			R, G, B, A := colorAt.RGBA()
			colorAtRGBA := color.NRGBA{R: uint8(R), G: uint8(G), B: uint8(B), A: uint8(A)}
			newImage.SetNRGBA(i, j, colorAtRGBA)
		}
	}
	return newImage, nil
}

func (w *Watermark) checkTooLarge(start image.Point, dst image.Rectangle) error {
	// 允许的最大高宽
	width := dst.Dx() - start.X - padding
	height := dst.Dy() - start.Y - padding
	if width < w.image.Bounds().Dx() || height < w.image.Bounds().Dy() {
		return fmt.Errorf("大小不合适")
	}
	return nil
}

// 根据位置获取
func (w *Watermark) getPoint(width, height int, waterImg image.Image) image.Point {
	var point image.Point

	switch w.pos {
	case TopLeft:
		point = image.Point{X: -padding, Y: -padding}
	case TopRight:
		point = image.Point{
			X: -(width - padding - waterImg.Bounds().Dx()),
			Y: -padding,
		}
	case BottomLeft:
		point = image.Point{
			X: -padding,
			Y: -(height - padding - waterImg.Bounds().Dy()),
		}
	case BottomRight:
		point = image.Point{
			X: -(width - padding - waterImg.Bounds().Dx()),
			Y: -(height - padding - waterImg.Bounds().Dy()),
		}
	case Center:
		point = image.Point{
			X: -(width - padding - waterImg.Bounds().Dx()) / 2,
			Y: -(height - padding - waterImg.Bounds().Dy()) / 2,
		}
	default:
		panic("无效的 pos 值")
	}

	return point
}

// 制作文字水印图片
func (w *Watermark) makeImageByText(text string) (image.Image, error) {
	if w.bgColor == nil {
		w.bgColor = color.RGBA{R: 0, G: 0, B: 0, A: 0}
	}
	fontSize := float64(fontSize)
	freetypeCtx := w.makeFreetypeCtx(fontSize)

	width, height := int(fontSize)*len(text), int(fontSize)*2
	rgbaRect := image.NewRGBA(image.Rect(0, 0, width, height))

	// 仅当非透明时才做一次额外的渲染
	if w.bgColor != color.Transparent {
		bgUniform := image.NewUniform(w.bgColor)
		draw.Draw(rgbaRect, rgbaRect.Bounds(), bgUniform, image.Pt(0, 0), draw.Src)
	}

	freetypeCtx.SetClip(rgbaRect.Rect)
	freetypeCtx.SetDst(rgbaRect)
	pt := freetype.Pt(0, int(freetypeCtx.PointToFixed(fontSize)>>6))
	_, err := freetypeCtx.DrawString(text, pt)
	if err != nil {
		return nil, err
	}
	return rgbaRect, nil
}

func (w *Watermark) makeFreetypeCtx(fontSize float64) *freetype.Context {

	fontColor := color.RGBA{R: 0, G: 0, B: 0, A: 50}
	fontColorUniform := image.NewUniform(fontColor)

	freetypeCtx := freetype.NewContext()
	// 设置像素密度
	freetypeCtx.SetDPI(100)
	// 指定字体
	freetypeCtx.SetFont(w.mustParseFont())
	//freetypeCtx.SetFont(font)
	freetypeCtx.SetFontSize(fontSize)
	// 指定文字颜色
	freetypeCtx.SetSrc(fontColorUniform)
	freetypeCtx.SetHinting(font.HintingNone)
	return freetypeCtx
}

// MustParseFont 通过单测来保证该方法必不会 panic
func (w *Watermark) mustParseFont() *truetype.Font {
	ft, err := freetype.ParseFont(gomono.TTF)
	if err != nil {
		panic(err)
	}
	return ft
}

//fileBytes, _ := ioutil.ReadFile("/Users/vega/workspace/codes/golang_space/gopath/src/work/cmsApp/uploadfile/404.png")
//reader := bytes.NewBuffer(fileBytes)
//img, _ := watermark.New(reader, ".png", watermark.BottomRight, "avc", "")
//err := img.WaterMakeDone("/Users/vega/workspace/codes/golang_space/gopath/src/work/cmsApp/uploadfile/202.png")
