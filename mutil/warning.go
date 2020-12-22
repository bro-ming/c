package mutil

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/bro-ming/c/mlog"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	dingClient *dingWarning
	once       sync.Once
)

func BuildDingWaring() *dingWarning {
	once.Do(func() {
		dingClient = &dingWarning{}
	})
	return dingClient
}

type dingWarning struct {
	text string
	hook string
}

func (dw *dingWarning) SetText(text string) *dingWarning {
	dw.text = text
	return dw
}

func (dw *dingWarning) SetHook(hook string) *dingWarning {
	dw.hook = hook
	return dw
}

// WarningHandler 异常预警机制
func (dw *dingWarning) Send() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("预警机制内部发生错误,明细%s", err)
		}
	}()

	if err := Verification(dw, "Text", "Hook"); err != nil {
		mlog.Errorf("钉钉预警机制未提供必要字段数据，停止推送！")
		return
	}

	reqUrl := sign(dw.hook)
	content := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"title": "严重错误警告",
			"text":  dw.text,
		},
		"at": map[string]interface{}{
			"atMobiles": []string{"", ""},
			"isAtAll":   true,
		},
	}

	if err := dingTo(content, reqUrl); err != nil {
		mlog.Errorf("钉钉预警机制推送发生异常！%s", err)
	}
}

// ding
func hmacSha256(stringToSign string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// Sign 发送钉钉消息
func sign(hook string) string {
	secret := "hg4EaTjisavDn2jwtoAEPai4k0ZyUNw05A72AC3n8hR8euGGH7LjaMm3ShwCZlBl"

	// webhook 从dingding群Hook获取
	webhook := hook

	timestamp := time.Now().UnixNano() / 1e6
	stringToSign := fmt.Sprintf("%d\n%s", timestamp, secret)
	sign := hmacSha256(stringToSign, secret)
	url := fmt.Sprintf("%s&timestamp=%d&sign=%s", webhook, timestamp, sign)
	return url
}

func dingTo(s map[string]interface{}, url string) error {
	b, _ := json.Marshal(s)
	resp, err := http.Post(url, "application/json",
		bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	// 发送结果
	mlog.Debugf("Ding waring result:%s", string(body))
	return nil
}

// Demo  调用Demo
func Demo() {
	text := "## <font color='#FF0000' size=6 face='黑体'> EX_SERVICE已全部失联,已无法继续提供数据！</font> \n\n" +
		"**时间**: " + time.Now().Format("2006/01/02 15:04:05") + "\n\n"
	hook := "https://oapi.dingtalk.com/robot/send?access_token=a5685f08bf7c2211b4374db6f3e237d1b093709da986f85161dc1c3c8115bc2e"
	go BuildDingWaring().SetText(text).SetHook(hook).Send()
}
