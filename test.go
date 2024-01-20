/*
一个简单认证模块, 防止接口完全公开被无脑调用

适用于内网, 低延时,可信度较高环境中的同一套服务的内部使用
仅需要使用统一的 KEY, 做简单的类似防盗链的加密认证方式

Decode 系列函数返回两个结果, 是否验证通过和错误码,
当通过当时候, 错误码为 nil, 不通过原因可通过错误码获得
其中 NoErr 系列只返回是否验证通过

Encode 系列函数需要注意, 当使用返回两个结果的函数时候,
其返回的密钥串中不包含时间戳信息
*/

package main

import (
	"github.com/disintegration/imaging"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/pkg/errors"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	"reflect"
)

func AddWatermarkForImage(oriImage image.Image, uid string) (*image.RGBA, error) {
	watermarkedImage := image.NewRGBA(oriImage.Bounds())
	draw.Draw(watermarkedImage, oriImage.Bounds(), oriImage, image.Point{}, draw.Src)

	// 生成水印的图片
	watermark, err := MakeImageByText(uid, color.Transparent)
	if err != nil {
		return nil, err
	}
	rotatedWatermark := imaging.Rotate(watermark, 30, color.Transparent)

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
	return watermarkedImage, nil
}

// MakeImageByText 根据文本内容制作一个仅包含该文本内容的图片
func MakeImageByText(text string, bgColor color.Color) (image.Image, error) {
	fontSize := float64(15)
	freetypeCtx := MakeFreetypeCtx(fontSize)

	width, height := int(fontSize)*len(text), int(fontSize)*2
	rgbaRect := image.NewRGBA(image.Rect(0, 0, width, height))

	// 仅当非透明时才做一次额外的渲染
	if bgColor != color.Transparent {
		bgUniform := image.NewUniform(bgColor)
		draw.Draw(rgbaRect, rgbaRect.Bounds(), bgUniform, image.Pt(0, 0), draw.Src)
	}

	freetypeCtx.SetClip(rgbaRect.Rect)
	freetypeCtx.SetDst(rgbaRect)
	pt := freetype.Pt(0, int(freetypeCtx.PointToFixed(fontSize)>>6))
	_, err := freetypeCtx.DrawString(text, pt)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return rgbaRect, nil
}

// MustParseFont 通过单测来保证该方法必不会 panic
func MustParseFont() *truetype.Font {
	ft, err := freetype.ParseFont(gomono.TTF)
	if err != nil {
		panic(err)
	}
	return ft
}

func MakeFreetypeCtx(fontSize float64) *freetype.Context {
	fontColor := color.RGBA{R: 0, G: 0, B: 0, A: 50}
	fontColorUniform := image.NewUniform(fontColor)

	freetypeCtx := freetype.NewContext()
	freetypeCtx.SetDPI(100)
	freetypeCtx.SetFont(MustParseFont())
	freetypeCtx.SetFontSize(fontSize)
	freetypeCtx.SetSrc(fontColorUniform)
	freetypeCtx.SetHinting(font.HintingNone)
	return freetypeCtx
}

func main() {
	//fileBytes, _ := ioutil.ReadFile("/Users/vega/workspace/codes/golang_space/gopath/src/work/cmsApp/uploadfile/404.png")
	//reader := bytes.NewBuffer(fileBytes)
	//w, err := watermark.New(reader, ".png", 2, watermark.Center)
	//fmt.Println(err, "---err--")
	//err = w.MarkFile("/Users/vega/workspace/codes/golang_space/gopath/src/work/cmsApp/uploadfile/")

	//watermark, _ := MakeImageByText("fuck", color.RGBA{R: 10, G: 20, B: 0, A: 0})
	//f, _ := os.Open("/Users/vega/workspace/codes/golang_space/gopath/src/work/cmsApp/uploadfile/404.png")
	//gopherImg, _, err := image.Decode(f)
	//watermarkedImg, _ := AddWatermarkForImage(gopherImg, "ffuck")
	//
	//outFile, err := os.Create("gopher2.png")
	//defer outFile.Close()
	//if err != nil {
	//	panic(err)
	//}
	//b := bufio.NewWriter(outFile)
	//err = png.Encode(b, watermarkedImg)
	//if err != nil {
	//	panic(err)
	//}
	//err = b.Flush()
	//if err != nil {
	//	panic(err)
	//}

	//fmt.Println(err)
	//t1 := time.Now().Unix()     //获取本地现在时间
	//time.Sleep(time.Second * 2) //延时2秒
	//t2 := time.Now().Unix()
	//fmt.Println(t2 - t1)
	type A struct {
		Field1 string
		Field2 int
	}

	type B struct {
		Field1 string
		Field2 int
		Field3 bool
	}

	reflect.StructField{
		A,
		B,
	}
}
