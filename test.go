package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"
)

func postRequest(url string, data []byte) error {
	// 创建一个 HTTP 客户端
	client := &http.Client{
		Timeout: time.Second * 10, // 设置超时时间
	}

	// 创建一个请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	// 设置请求头，例如设置 Content-Type
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// 打印响应体
	println(string(body))

	return nil
}

func main() {
	url := "http://localhost:3015/api/article/toolBarData"
	data := []byte(`{"articleId": 199}`) // 你的 JSON 数据

	// 循环发送 POST 请求
	for i := 0; i < 5000; i++ {
		err := postRequest(url, data)
		if err != nil {
			println("Error:", err)
			return
		}
		//time.Sleep(0.1 * time.Second) // 等待一秒
	}
}
