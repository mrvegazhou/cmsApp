package random

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
	mathRand "math/rand"
	"time"
)

// GetUuid 获取UUID
func GetUuid() string {
	b := make([]byte, 16)
	io.ReadFull(rand.Reader, b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

/*
*
生成随机字符串
*/
func RandString(codeLen int) string {
	startTime := time.Now()
	// 1. 定义原始字符串
	rawStr := "jkwangagDGFHGSERKILMJHSNOPQR546413890_"
	// 2. 定义一个buf，并且将buf交给bytes往buf中写数据
	buf := make([]byte, 0, codeLen)
	b := bytes.NewBuffer(buf)
	// 随机从中获取
	ran := mathRand.New(mathRand.NewSource(time.Now().UnixNano()))
	for rawStrLen := len(rawStr); codeLen > 0; codeLen-- {
		randNum := ran.Intn(rawStrLen)
		b.WriteByte(rawStr[randNum])
	}
	//time.Since(startTime)
	//cost := int(time.Since(startTime) / time.Second)
	fmt.Println(time.Since(startTime))
	return b.String()
}

func RandomNumber(len int) string {
	var numbers = []byte{0, 1, 2, 3, 4, 5, 7, 8, 9}
	var container string
	length := bytes.NewReader(numbers).Len()

	for i := 1; i <= len; i++ {
		random, err := rand.Int(rand.Reader, big.NewInt(int64(length)))
		if err != nil {

		}
		container += fmt.Sprintf("%d", numbers[random.Int64()])
	}
	return container
}

func RandomInt(min, max int) int {
	if min >= max || max == 0 {
		return max
	}
	ran := mathRand.New(mathRand.NewSource(time.Now().UnixNano()))
	return ran.Intn(max-min) + min
}
