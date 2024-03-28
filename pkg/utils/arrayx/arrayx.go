package arrayx

import (
	"fmt"
	"github.com/spf13/cast"
	"strings"
)

type MyType interface {
	uint | uint64 | int | int64 | string
}

func IsContain[T MyType](items []T, item T) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

type IntType interface {
	uint | uint64 | int | int64
}

// 整型数组转逗号分隔的字符串
func JoinIntArr2Str[T IntType](numArr []T) string {
	var strArr = make([]string, len(numArr))
	for k, v := range numArr {
		strArr[k] = fmt.Sprintf("%v", v)
	}
	return strings.Join(strArr, ",")
}

func String2Int(strArr []string) []int {
	res := make([]int, len(strArr))
	for index, val := range strArr {
		res[index] = cast.ToInt(val)
	}
	return res
}

func String2Uint64(strArr []string) []uint64 {
	res := make([]uint64, len(strArr))
	for index, val := range strArr {
		res[index] = cast.ToUint64(val)
	}
	return res
}
