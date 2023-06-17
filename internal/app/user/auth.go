package user

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

type WXLoginResp struct {
	DataList []struct {
		Data struct {
			PhoneNumber string `json:"phoneNumber"`
		} `json:"data"`
	} `json:"data_list"`
}

func (s *Service) WXLogin(openid string, cloudID string) (bool, error) {
	// 合成url, 这里的appId和secret是在微信公众平台上获取的
	url := fmt.Sprintf("http://api.weixin.qq.com/wxa/getopendata?openid=%s", openid)
	// set body
	body, err := json.Marshal(map[string]interface{}{
		"cloudid_list": []string{cloudID},
	})
	// 创建http post请求
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return false, err
	}
	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	// 发送请求
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	// 解析http请求中body 数据到我们定义的结构体中
	wxResp := WXLoginResp{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&wxResp); err != nil {
		return false, err
	}
	return true, nil
}

// 将一个字符串进行MD5加密后返回加密后的字符串
func GetMD5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
