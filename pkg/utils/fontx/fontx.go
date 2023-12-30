package fontx

import (
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"log"
	"os"
	"unicode"
)

// GetFont 获取一个字体对象
func GetFont(src string) *truetype.Font {
	fontSourceBytes, err := os.ReadFile(src)
	if err != nil {
		log.Println("读取字体失败:", err)
	}

	trueTypeFont, err := freetype.ParseFont(fontSourceBytes)

	if err != nil {
		log.Println("解析字体失败:", err)
	}

	return trueTypeFont
}

func GetEnOrChLength(text string) int {
	enCount, zhCount := 0, 0

	for _, t := range text {
		if unicode.Is(unicode.Han, t) {
			zhCount++
		} else {
			enCount++
		}
	}

	chOffset := (25/2)*zhCount + 5
	enOffset := enCount * 8

	return chOffset + enOffset
}