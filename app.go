package main

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/blinkbean/dingtalk"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

const (
	// 考生号
	studentNo string = ``
	// 身份证号
	idNo string = ``
	// 姓名
	name string = ``
	// 钉钉群 Webhook 链接内的 token
	dingToken string = ""
	// 阿里云市场里 3023data 识别验证码 API 的凭证
	appCode string = ""
	// 3023data 识别验证码 API 的凭证
	key string = ""
)

func main() {
	args := []string{studentNo, idNo, name, dingToken}

	for _, arg := range args {
		if arg == "" {
			log.Fatal("需要先补全配置文件才能运行")
		}
	}

	message := "没有查询到你的信息"
	times := 0
	messageTime, err := time.Parse(time.RFC3339, "2021-07-19T12:57:00.000+08:00")
	if err != nil {
		log.Fatal(err)
	}

	// 每隔 15 分钟查询一次, 如果结果没变化在一小时内不会发送重复消息
	for {
		scrape(&message, &messageTime, &times)
		time.Sleep(15 * time.Minute)
	}
}

//实现:
//1. 请求 http://zsb.jlu.edu.cn/Index/gkluqucx.html
//2. 获取一个一次性的 hash, 通过 `meta name="hash"`
//3. 获取验证码图片 base64 之后请求 API 识别, 获取一串数字结果
//4. 请求 http://zsb.jlu.edu.cn/index/gkcxresult.html , 提取响应结果 `.system-message`>
//5. 将结果发送到钉钉的 Webhook.
func scrape(prevResult *string, lastMessageTime *time.Time, duplicatedTimes *int) {
	client := getClient()

	hash := getHash(client)
	fmt.Printf("hash: %s\n", hash)

	verificationCode := getVerificationCode(client)
	fmt.Printf("verificationCode: %s\n", verificationCode)

	//5. 请求 http://zsb.jlu.edu.cn/index/gkcxresult.html , 提取响应结果 `.system-message`> DONE
	result := queryResult(client, verificationCode, hash)
	fmt.Println(fmt.Sprintf("time: %s, result: %s", time.Now().Format(time.RFC3339), result))

	if time.Now().Sub(*lastMessageTime).Seconds() < (1 * time.Hour).Seconds() {
		if *prevResult == result {
			*duplicatedTimes++

			fmt.Printf("和上次查询结果相同, 不发送消息, 已重复 %d 次\n", *duplicatedTimes)
			return
		}

		if strings.Contains(result, "验证码输入错误") {
			*duplicatedTimes++

			fmt.Println("验证码错误")
			return
		}
	}

	// 沉默了一小时, 或者有新的结果
	*prevResult = result
	*lastMessageTime = time.Now()
	*duplicatedTimes = 0

	send(result)
}

func getClient() *http.Client {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}

	return &http.Client{
		Transport: &http.Transport{
			// This is insecure; use only in dev environments.
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Jar: jar,
	}
}

func getHash(client *http.Client) string {
	//1. 请求 http://zsb.jlu.edu.cn/Index/gkluqucx.html DONE
	indexHtml := getIndexHtml(client)

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(indexHtml))
	if err != nil {
		log.Fatal(err)
	}

	//2. 获取一个一次性的 hash, 通过 `meta name="hash"`
	hash := ""
	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		if name, _ := s.Attr("name"); name == "hash" {
			hash, _ = s.Attr("content")
		}
	})
	return hash
}

//3. 获取验证码图片
//4. 验证码 base64 之后请求 API 识别, 获取一串数字结果
func getVerificationCode(client *http.Client) string {
	base64Str, err := getCodeImage(client)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(base64Str)

	verificationCode, err := identifyCode(client, base64Str)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return verificationCode
}

func getCodeImage(client *http.Client) (base64Str string, err error) {
	imageUrl := "http://zsb.jlu.edu.cn/index/verify/tp/1.html"

	//Get the response bytes from the url
	response, err := client.Get(imageUrl)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	return base64.StdEncoding.EncodeToString(bodyBytes), nil
}

func getIndexHtml(client *http.Client) string {
	// Request the HTML page.
	res, err := client.Get("http://zsb.jlu.edu.cn/Index/gkluqucx.html")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		// handle err
		log.Fatal(err)
	}
	bodyStr := string(bodyBytes)

	return bodyStr
}

func queryResult(client *http.Client, verificationCode string, hash string) string {
	params := url.Values{}
	params.Add("ksh", studentNo)
	params.Add("sfzh", idNo)
	params.Add("xm", name)
	params.Add("yzm", verificationCode)
	params.Add("hash", hash)
	str := params.Encode()
	body := strings.NewReader(str)

	req, err := http.NewRequest("POST", "http://zsb.jlu.edu.cn/index/gkcxresult.html", body)
	if err != nil {
		log.Fatal(err)
	}

	req.Host = "zsb.jlu.edu.cn"
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Add("Accept-Charset", "utf-8")
	// 不支持 gzip
	req.Header.Set("Accept-Encoding", "")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Origin", "http://zsb.jlu.edu.cn")
	req.Header.Set("Referer", "http://zsb.jlu.edu.cn/Index/gkluqucx.html")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyStr := string(bodyBytes)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(bodyStr))
	if err != nil {
		log.Fatal(err)
	}

	// 未查询到结果
	errorMessage := doc.Find(".system-errorMessage").Text()
	if errorMessage != "" {
		return mapped(errorMessage)
	}

	// 成功
	ems := "尚未寄出"
	doc.Find("td").Each(func(i int, s *goquery.Selection) {
		if left := s.Text(); strings.Contains(left, "EMS编号：") {
			value := s.Next().Text()
			if value != "" {
				ems = value
			}
		}
	})

	return fmt.Sprintf("EMS编号: %s", ems)
}

func identifyCode(client *http.Client, base64Str string) (string, error) {
	params := url.Values{}
	params.Add("image", base64Str)
	params.Add("maxlength", `8`)
	params.Add("minlength", `1`)
	params.Add("type", `5005`)
	str := params.Encode()
	body := strings.NewReader(str)

	req, err := http.NewRequest("POST", "http://api.3023data.com/ocr/captcha", body)
	if err != nil {
		log.Fatal(err)
	}

	// 有 appCode 时代表是在阿里云上测试, 30 次之后需要使用 3023data 的 key 替换 appCode
	if appCode != "" {
		req, err = http.NewRequest("POST", "http://302307.market.alicloudapi.com/ocr/captcha", body)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("Authorization", "APPCODE "+appCode)
	}

	req.Header.Set("key", key)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		// handle err
		log.Fatal(err)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	bodyStr := string(bodyBytes)

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("识别验证码失败, 返回状态码: %d, 响应: %s\n", resp.StatusCode, bodyStr)
	}

	r := CodeResponse{}
	err = json.Unmarshal(bodyBytes, &r)
	if err != nil {
		log.Fatal(err)
	}

	if r.Code != 0 {
		return "", errors.New(fmt.Sprintf("识别验证码失败, 响应: %s\n", bodyStr))
	}

	return r.Data.Captcha, nil
}

// 6. 将结果发送到钉钉的 Webhook.
func send(result string) {
	robot := dingtalk.InitDingTalk([]string{dingToken}, "")

	err := robot.SendMarkDownMessage("查询结果", fmt.Sprintf(
		`
%s, 查询结果: 

%s

> 每隔 15 分钟查询一次
> 
> 如果结果没变化, 一小时内不会发送重复消息
`, time.Now().Format(time.Kitchen), result))

	if err != nil {
		log.Fatal(err)
	}
}

func mapped(result string) string {
	if strings.Contains(result, "没有查询到你的信息") {
		return "没有查询到你的信息"
	}

	if strings.Contains(result, "验证码输入错误") {
		return "验证码输入错误"
	}

	return result
}

type CodeResponse struct {
	Code int `json:"code"`
	Data struct {
		Captcha string `json:"captcha"`
		Type    int    `json:"type"`
		Length  int    `json:"length"`
		ID      string `json:"id"`
	} `json:"data"`
}
