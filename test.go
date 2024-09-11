package main

import (
	"cmsApp/pkg/DES"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unsafe"
)

func replaceImageSrc(str, newSrc string) string {
	// 正则表达式匹配自定义的 Image 标签和 src 属性
	imageRegex := regexp.MustCompile(`(Image src=")([^"]+)(")`)
	return imageRegex.ReplaceAllString(str, fmt.Sprintf(`$1%s$3`, newSrc))
}
func appendParamToImageSrc(str string, paramName string, paramValue string) string {
	// 正则表达式匹配 Image 标签的 src 属性
	imageRegex := regexp.MustCompile(`(<Image [^>]*src=")([^"]+)(" [^>]*>)`)
	// 替换函数，用于在 src 属性值后追加参数
	replaceFunc := func(match string) string {
		parts := imageRegex.FindStringSubmatch(match)
		if len(parts) < 4 {
			return match // 如果没有匹配到src属性值，返回原始字符串
		}
		src := fmt.Sprintf("%s?%s=%s", parts[2], paramName, paramValue)
		return parts[1] + src + parts[3]
	}
	// 使用正则表达式替换所有匹配的 Image 标签的 src 属性
	return imageRegex.ReplaceAllStringFunc(str, replaceFunc)
}

func string2bytes(s string) []byte {
	if len(s) == 0 {
		return []byte("")
	}
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

func bytes2string(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	return *(*string)(unsafe.Pointer(&b))
}

type Encoder struct {
	src   []byte
	dst   []byte
	Error error
}

func newEncoder() Encoder {
	return Encoder{}
}

func (e Encoder) FromString(s string) Encoder {
	e.src = string2bytes(s)
	return e
}
func (e Encoder) ByHex() Encoder {
	if len(e.src) == 0 || e.Error != nil {
		return e
	}
	buf := make([]byte, hex.EncodedLen(len(e.src)))
	hex.Encode(buf, e.src)
	e.dst = buf
	return e
}
func (e Encoder) ToString() string {
	return bytes2string(e.dst)
}

func modifyImgSrc(content string, param string) string {
	// 使用正则表达式匹配 img 标签
	imgTagRegex := regexp.MustCompile(`<img\s+([^>]*)src\s*=\s*['"]([^'"]*)['"][^>]*>`)

	// 使用替换函数来修改匹配到的内容
	replacedContent := imgTagRegex.ReplaceAllStringFunc(content, func(match string) string {
		// 提取 src 属性的值
		submatches := imgTagRegex.FindStringSubmatch(match)
		if len(submatches) > 2 {
			attributes := submatches[1] // 所有的属性
			srcValue := submatches[2]   // src 的值
			// 在 src 值后面添加参数 t=xxx
			newSrc := srcValue + "?" + param
			// 构造新的 img 标签
			fmt.Println(attributes, "--attributes--")
			return fmt.Sprintf("<img %s src=\"%s\">", attributes, newSrc)
		}
		return match // 如果没有匹配到 src，则返回原字符串
	})

	return replacedContent
}

func main() {
	// 假设的加密字符串，需要替换为实际的加密字符串
	//encryptedStr := "cGeh YneCLLY32i9c71OWa36NzjAdZhk4xnUP/qjqhIIMOl8dNJ ElI6pjInzKLtk4vRLk4InglasLwwSXpy k61kck1bfra1awJHfnFCei 1kQyJBtw5Tgf7U KkiNwr9dMEufpsOBs16dhtPXFHwo7JetnRh62zXqjuoUc n83ssTG1upzWxNfvNesV/Mz7OeEhNIG byA62TF10dzOfPELOCvR8hgRmDgjK8KLBI4RKnHbOT6dPjfdl6Ay/ Rnz8dz5rZ9zNUW/wt0grjA4pYs1AMXFRViYGOIlqlgQE="
	//replacedStr := strings.Replace(encryptedStr, " ", "+", -1)
	//fmt.Println(AES.DecryptJsStr(replacedStr, "0123456789abcdef", "0123456789abcdef"))
	//keyStr := "0123456789abcdef" // 确保密钥长度为16, 24或32
	//ivStr := "0123456789abcdef"  // 确保 IV 长度为 16
	//
	//decryptedStr, err := aesDecryptStr(replacedStr, keyStr, ivStr)
	//if err != nil {
	//	fmt.Println("Error decrypting string:", err)
	//	return
	//}
	//
	//fmt.Println("Decrypted string:", decryptedStr)
	//html := `<img src="http://localhost:3015/api/image/static/1832334041920049152.png">爱的火花`
	//start := time.Now() // 开始时间
	//ssss := time.Now().Format("2006-01-02 15:04:05")
	//content, _ := DES.DesCbcEncryptBase64([]byte(ssss), []byte("upload89"), nil)
	//res, _ := parser.AddImgSrcT(html, content)
	//fmt.Println(res)
	//fmt.Println(strings.ContainsImgTag(html))
	//res, _ := DES.DesCbcDecryptByBase64(content, []byte("upload89"), nil)
	//fmt.Println(content, stringx.Bytes2string(res))
	//elapsed := time.Since(start) // 计算运行时间
	//fmt.Printf("Function took222 %s\n", elapsed)

	//fmt.Println(HEXX.EncodeHex("f"))
	//fmt.Println(HEXX.DecodeHex("30313233343536373839616263646566"))
	//fmt.Println(HTMLX.AppendParamToImageSrc(html, "t", "xxxx"))
	content := "http://localhost:3015/api/image/static/p/1830142143692279808.png"
	filename := filepath.Base(content)
	filenameWithoutExt := strings.TrimSuffix(filename, filepath.Ext(filename))
	//lastFourChars := filenameWithoutExt[len(filenameWithoutExt)-4:]
	encryptStr, _ := DES.DesCbcEncryptBase64([]byte(filenameWithoutExt+time.Now().Format("2006-01-02")), []byte("upload89"), nil)
	//desTime, err := DES.DesCbcDecryptByBase64("EzBD8o2zZJI0ZKS/PZ7xpysdGVdJ+ohJ", string2bytes("upload89"), nil)
	fmt.Println(filenameWithoutExt, encryptStr)

	timeTemplate := "2006-01-02"
	lastStamp, err := time.ParseInLocation(timeTemplate, "2024-09-07", time.Local)
	fmt.Println(lastStamp, err)

	str := "Hello World! This is a test string."
	// 使用Split函数分割字符串，这将分割所有空格
	parts := strings.SplitN(str, " ", 2)
	fmt.Println("第一部分:", parts[1])
}
