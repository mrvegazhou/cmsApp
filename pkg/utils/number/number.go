package number

import (
	"cmsApp/internal/constant"
	"github.com/speps/go-hashids/v2"
	"github.com/spf13/cast"
	"math"
)

func RemoveRepeatedInArr[T uint64 | int | float32 | float64](s []T) []T {
	result := make([]T, 0)
	m := make(map[T]bool) //map的值不重要
	for _, v := range s {
		if _, ok := m[v]; !ok {
			result = append(result, v)
			m[v] = true
		}
	}
	return result
}

func ConvertUint64ToInt64Array(num uint64) []int64 {
	const int64Max = int64(^uint(0) >> 1)
	if IsUint64GreaterThanInt64Max(num) == false {
		return []int64{cast.ToInt64(num)}
	}
	var result []int64
	for num > 0 {
		// 模运算得到当前最低位的值，确保其不大于int64的最大值
		part := int64(num % 1000000000)
		// 将当前部分添加到结果数组中
		result = append(result, part)
		// 将数字右移，准备处理下一部分
		num /= 1000000000
	}
	return result
}

func Int64ArrayToUint64(arr []int64) uint64 {
	var result string
	const int64Max = int64(^uint(0) >> 1) // int64的最大值
	for _, val := range arr {
		if val < 0 || val > int64Max {
			panic("Error: Array contains value out of int64 range")
		}
		// 乘以相应的基数，然后加上当前的int64值
		result += cast.ToString(val)
	}
	return cast.ToUint64(result)
}

func IsUint64GreaterThanInt64Max(num uint64) bool {
	// int64的最大值计算方式：数学上的最大值减1
	int64Max := uint64(math.MaxInt64) + 1
	return num > int64Max
}

func NumToHashId(num uint64) (string, error) {
	hd := hashids.NewData()
	hd.Salt = constant.ARTICLE_ID_SECRET
	hd.MinLength = 30
	h, err := hashids.NewWithData(hd)
	e, err := h.EncodeInt64(ConvertUint64ToInt64Array(num))
	return e, err
}

func HashIdToNum(id string) (uint64, error) {
	hd := hashids.NewData()
	hd.Salt = constant.ARTICLE_ID_SECRET
	hd.MinLength = 30
	h, _ := hashids.NewWithData(hd)
	numArr, err := h.DecodeInt64WithError(id)
	if err != nil {
		return 0, err
	}
	return Int64ArrayToUint64(numArr), nil
}
