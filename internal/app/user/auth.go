package user

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"wxcloudrun-golang/internal/pkg/model"
)

type Service struct {
	UserDao *model.User
}

func NewService() *Service {
	return &Service{}
}

type WXLoginResp struct {
	DataList []struct {
		Json struct {
			Data struct {
				PhoneNumber string `json:"phoneNumber"`
			} `json:"data"`
		} `json:"json"`
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
	defer resp.Body.Close()
	bodys, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return false, err
	}
	fmt.Println(string(bodys))
	// print response
	fmt.Println(resp)
	// 解析http请求中body 数据到我们定义的结构体中
	wxResp := WXLoginResp{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&wxResp); err != nil {
		return false, err
	}
	// print data
	fmt.Println(wxResp)
	if wxResp.DataList != nil && len(wxResp.DataList) > 0 && wxResp.DataList[0].Json.Data.PhoneNumber == "" {
		_, err = s.UserDao.Create(&model.User{
			OpenID:      openid,
			Phone:       wxResp.DataList[0].Json.Data.PhoneNumber,
			CreatedTime: time.Now(),
			UpdatedTime: time.Now(),
		})
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

// 将一个字符串进行MD5加密后返回加密后的字符串
func GetMD5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
