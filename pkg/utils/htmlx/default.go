package HTMLX

import (
	"cmsApp/configs"
	"cmsApp/pkg/DES"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func GetImgSrcs(content string) []string {
	// 编译正则表达式，匹配img标签的src属性
	re := regexp.MustCompile(`<img src="([^"]+)"`)
	// 使用正则表达式查找所有匹配项
	matches := re.FindAllStringSubmatch(content, -1)
	// 遍历匹配项
	var srcValues []string
	for _, match := range matches {
		// match[0] 是整个匹配的字符串，match[1] 是第一个捕获组，即src的值
		srcValues = append(srcValues, match[1])
	}
	return srcValues
}

func AppendParamToImageSrc(str string) string {
	// 正则表达式匹配img标签
	imgTagRegex := regexp.MustCompile(`(<img [^>]*?src=")([^"]+)(".*?>)`)
	// 替换src属性，添加参数t=xxx
	newContent := imgTagRegex.ReplaceAllStringFunc(str, func(match string) string {
		// 分割匹配的字符串
		parts := imgTagRegex.FindStringSubmatch(match)
		// 添加参数t=xxx到src属性值
		filename := filepath.Base(parts[2])
		filenameWithoutExt := strings.TrimSuffix(filename, filepath.Ext(filename))
		needEncryptStr := fmt.Sprintf("%s %s", filenameWithoutExt, time.Now().Format("2006-01-02"))
		encryptStr, _ := DES.DesCbcEncryptBase64([]byte(needEncryptStr), []byte(configs.App.Upload.Key), nil)
		newSrc := fmt.Sprintf("%s?t=%s", parts[2], encryptStr)
		// 重新组合img标签
		return parts[1] + newSrc + parts[3]
	})
	return newContent
}
