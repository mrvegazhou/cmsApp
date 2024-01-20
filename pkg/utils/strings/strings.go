package strings

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strings"
	"unicode"
)

/*
*
密码加密
*/
func Encryption(password string, salt string) string {
	str := fmt.Sprintf("%s%s", password, salt)
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

/**
*首字母大写
**/
func StrFirstToUpper(str string) (string, string, string) {

	var upperStr string
	var firstStr string
	var secondUp string
	temp := strings.Split(str, "_")

	for y := 0; y < len(temp); y++ {
		vv := []rune(temp[y])

		firstStr += string(vv[0])
		vv[0] -= 32
		upperStr += string(vv)
		if y == 0 {
			secondUp += temp[0]
		} else {
			secondUp += string(vv)
		}
	}
	return upperStr, firstStr, secondUp
}

/*
*比较第二个slice一第一个slice的区别
 */
func CompareSlice(first []string, second []string) (add []string, incre []string) {

	secondMap := make(map[string]struct{})

	for _, v := range second {
		secondMap[v] = struct{}{}
	}

	for _, v := range first {
		_, ok := secondMap[v]
		if !ok {
			incre = append(incre, v)
		} else {
			delete(secondMap, v)
		}
	}

	for k, _ := range secondMap {
		add = append(add, k)
	}

	return
}

/**
* 组装字符串
 */
func JoinStr(items ...interface{}) string {
	if len(items) == 0 {
		return ""
	}
	var builder strings.Builder
	for _, v := range items {
		builder.WriteString(v.(string))
	}
	return builder.String()
}

// EncodeMD5 生成 MD5
func EncodeMD5(value string) string {
	m := md5.New()
	m.Write([]byte(value))
	return hex.EncodeToString(m.Sum(nil))
}

var (
	// ErrInvalidStartPosition is an error that indicates the start position is invalid.
	ErrInvalidStartPosition = errors.New("start position is invalid")
	// ErrInvalidStopPosition is an error that indicates the stop position is invalid.
	ErrInvalidStopPosition = errors.New("stop position is invalid")
)

// Contains checks if str is in list.
func Contains(list []string, str string) bool {
	for _, each := range list {
		if each == str {
			return true
		}
	}

	return false
}

// Filter filters chars from s with given filter function.
func Filter(s string, filter func(r rune) bool) string {
	var n int
	chars := []rune(s)
	for i, x := range chars {
		if n < i {
			chars[n] = x
		}
		if !filter(x) {
			n++
		}
	}

	return string(chars[:n])
}

// FirstN returns first n runes from s.
func FirstN(s string, n int, ellipsis ...string) string {
	var i int

	for j := range s {
		if i == n {
			ret := s[:j]
			for _, each := range ellipsis {
				ret += each
			}
			return ret
		}
		i++
	}

	return s
}

// HasEmpty checks if there are empty strings in args.
func HasEmpty(args ...string) bool {
	for _, arg := range args {
		if len(arg) == 0 {
			return true
		}
	}

	return false
}

// Join joins any number of elements into a single string, separating them with given sep.
// Empty elements are ignored. However, if the argument list is empty or all its elements are empty,
// Join returns an empty string.
func Join(sep byte, elem ...string) string {
	var size int
	for _, e := range elem {
		size += len(e)
	}
	if size == 0 {
		return ""
	}

	buf := make([]byte, 0, size+len(elem)-1)
	for _, e := range elem {
		if len(e) == 0 {
			continue
		}

		if len(buf) > 0 {
			buf = append(buf, sep)
		}
		buf = append(buf, e...)
	}

	return string(buf)
}

// NotEmpty checks if all strings are not empty in args.
func NotEmpty(args ...string) bool {
	return !HasEmpty(args...)
}

// Remove removes given strs from strings.
func Remove(strings []string, strs ...string) []string {
	out := append([]string(nil), strings...)

	for _, str := range strs {
		var n int
		for _, v := range out {
			if v != str {
				out[n] = v
				n++
			}
		}
		out = out[:n]
	}

	return out
}

// Reverse reverses s.
func Reverse(s string) string {
	a := func(s string) *[]rune {
		var b []rune
		for _, k := range []rune(s) {
			defer func(v rune) {
				b = append(b, v)
			}(k)
		}
		return &b
	}(s)
	return string(*a)
}

// Substr returns runes between start and stop [start, stop)
// regardless of the chars are ascii or utf8.
func Substr(str string, start, stop int) (string, error) {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		return "", ErrInvalidStartPosition
	}

	if stop < 0 || stop > length {
		return "", ErrInvalidStopPosition
	}

	return string(rs[start:stop]), nil
}

// TakeOne returns valid string if not empty or later one.
func TakeOne(valid, or string) string {
	if len(valid) > 0 {
		return valid
	}

	return or
}

// TakeWithPriority returns the first not empty result from fns.
func TakeWithPriority(fns ...func() string) string {
	for _, fn := range fns {
		val := fn()
		if len(val) > 0 {
			return val
		}
	}

	return ""
}

// ToCamelCase returns the string that converts the first letter to lowercase.
func ToCamelCase(s string) string {
	for i, v := range s {
		return string(unicode.ToLower(v)) + s[i+1:]
	}

	return ""
}

var Placeholder PlaceholderType

type (
	// PlaceholderType represents a placeholder type.
	PlaceholderType = struct{}
)

// Union merges the strings in first and second.
func Union(first, second []string) []string {
	set := make(map[string]PlaceholderType)

	for _, each := range first {
		set[each] = Placeholder
	}
	for _, each := range second {
		set[each] = Placeholder
	}

	merged := make([]string, 0, len(set))
	for k := range set {
		merged = append(merged, k)
	}

	return merged
}

func convertNBytes(n, b float64) float64 {
	bits := b * 8
	x := math.Pow(2, bits)
	y := math.Pow(2, bits-1)
	return math.Mod(n+y, x) - y
}

func GetHashCode(s string) float64 {
	h := float64(0)
	n := float64(len(s))
	for idx, _ := range s {
		i := float64(idx)
		x := math.Pow(31, n-1-i)
		h = h + float64(s[idx])*x
	}
	return convertNBytes(h, 4)
}

func CheckMobile(phone string) bool {
	// 匹配规则
	// ^1第一位为一
	// [345789]{1} 后接一位345789 的数字
	// \\d \d的转义 表示数字 {9} 接9位
	// $ 结束符
	regRuler := "^1[345789]{1}\\d{9}$"

	// 正则调用规则
	reg := regexp.MustCompile(regRuler)

	// 返回 MatchString 是否匹配
	return reg.MatchString(phone)
}

// CheckIdCard 检验身份证
func CheckIdCard(card string) bool {
	//18位身份证 ^(\d{17})([0-9]|X)$
	// 匹配规则
	// (^\d{15}$) 15位身份证
	// (^\d{18}$) 18位身份证
	// (^\d{17}(\d|X|x)$) 18位身份证 最后一位为X的用户
	regRuler := "(^\\d{15}$)|(^\\d{18}$)|(^\\d{17}(\\d|X|x)$)"

	// 正则调用规则
	reg := regexp.MustCompile(regRuler)

	// 返回 MatchString 是否匹配
	return reg.MatchString(card)
}

// 识别电子邮箱
func CheckEmail(email string) bool {
	result, _ := regexp.MatchString(`^([\w\.\_\-]{2,10})@(\w{1,}).([a-z]{2,4})$`, email)
	if result {
		return true
	}
	return false
}
