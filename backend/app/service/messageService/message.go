package messageService

import (
	"bytes"
	"fmt"
	"github.com/valyala/fasthttp"
	"message/app/model"
	"message/config/database"
	"time"
)

func CreateMessage(message *model.Message) error {
	return database.DB.Create(message).Error
}

func ListMessage(userId, fromId, num, MessageType int) ([]*model.Message, error) {
	ansList := make([]*model.Message, 0)
	database.DB.Model(&model.Message{}).Where(&model.Message{
		SendType: MessageType,
		SendTo:   userId,
		From:     fromId,
	}).Order("ctime DESC").Limit(num).Find(&ansList)
	return ansList, nil
}

func GetMessageList(userId int, out time.Time) ([]*model.Message, error) {
	ansList := make([]*model.Message, 0)
	database.DB.Model(&model.Message{}).Where("send_to = ? AND ctime > ?", userId, out).Order("ctime DESC").Find(&ansList)
	return ansList, nil
}

func Translate(word string) (string, error) {
	url := "https://fanyi.so.com/index/search?eng=1&validate=&ignore_trans=0&query=hello"
	headers := map[string]string{
		"Accept":             "application/json, text/plain, */*",
		"Accept-Encoding":    "gzip, deflate, br",
		"Accept-Language":    "zh-CN,zh;q=0.9",
		"Content-Length":     "0",
		"Cookie":             "QiHooGUID=F02A63E0BCB72DB4A01C21FA023475E1.1703769301607; Q_UDID=00b0237e-501b-1360-b2eb-96b79d1ac5ec; __guid=144965027.253643186935022000.1703769305042.223; count=2",
		"Origin":             "https://fanyi.so.com",
		"Pro":                "fanyi",
		"Referer":            "https://fanyi.so.com/",
		"Sec-Ch-Ua":          `"Not_A Brand";v="8", "Chromium";v="120", "Google Chrome";v="120"`,
		"Sec-Ch-Ua-Mobile":   "?0",
		"Sec-Ch-Ua-Platform": `"Windows"`,
		"Sec-Fetch-Dest":     "empty",
		"Sec-Fetch-Mode":     "cors",
		"Sec-Fetch-Site":     "same-origin",
		"User-Agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	}
	form := fmt.Sprintf("eng=%s&ignore_trans=0&query=%s", 1, word)
	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()

	req.SetRequestURI(url)
	req.Header.SetMethod("POST")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	req.SetBody([]byte(form))
	req.Header.SetContentType("application/x-www-form-urlencoded")

	err := fasthttp.Do(req, res)
	if err != nil {
		fmt.Println("请求出错:", err)
	}

	body := res.Body()
	start := bytes.Index(body, []byte(`"fanyi":"`))
	if start == -1 {
		fmt.Println("翻译结果未找到")
	}
	start += len(`"fanyi":"`)
	end := bytes.IndexByte(body[start:], '"')
	if end == -1 {
		fmt.Println("翻译结果解析失败")
	}
	translated := body[start : start+end]
	return string(translated), nil
}
