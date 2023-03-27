package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func main() {
	// 构造客户端
	client := http.Client{}
	// 构造Get请求
	URL := "https://api.bilibili.com/x/web-interface/popular?ps=20&pn=1"
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		fmt.Println("req err", err)
	}
	// 添加请求头
	req.Header.Set("authority", "api.bilibili.com")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36")
	req.Header.Set("accept", "*/*")
	req.Header.Set("referer", "https://www.bilibili.com/v/popular/all/?spm_id_from=333.1007.0.0")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9")
	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("发送请求失败", err)
	}
	defer resp.Body.Close()
	// 读取响应json到内存bodyTest中
	bodyTest, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("io err", err)
	}
	var resultList Bilibilihot
	_ = json.Unmarshal(bodyTest, &resultList)
	for i, data := range resultList.Data.List {
		fmt.Printf("第 %d 个视频信息：", i)
		fmt.Println("bvid:", data.Bvid)
		fmt.Println("pic_url :", data.Pic)
		fmt.Println("title :", data.Title)
		fmt.Println("owner :", data.Owner.Name)
		fmt.Println("########################################################")
	}

}

type Bilibilihot struct {
	Data struct {
		List []struct {
			Bvid  string `json:"bvid"`
			Owner struct {
				Face string `json:"face"`
				Mid  int64  `json:"mid"`
				Name string `json:"name"`
			} `json:"owner"`
			Pic   string `json:"pic"`
			Title string `json:"title"`
			Tname string `json:"tname"`
		} `json:"list"`
		NoMore bool `json:"no_more"`
	} `json:"data"`
}
