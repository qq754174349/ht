package wechat

import (
	"context"
	"encoding/json"
	"github.com/qq754174349/ht/ht-frame/logger"
	"github.com/qq754174349/ht/ht-frame/redis"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	appid     string
	appSecret string
}

type accessTokenResp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type Code2SessionResp struct {
	SessionKey string `json:"session_key"`
	UnionId    string `json:"unionid"`
	ErrMsg     string `json:"errmsg"`
	OpenId     string `json:"openid"`
	ErrCode    int    `json:"errcode"`
}

func NewClient() *Client {
	return &Client{
		appid:     "wx963823904279a30e",
		appSecret: "e55a06f73e1fa8f41d0ff806f56ca1f3",
	}
}

func (c *Client) getAccToken() string {
	redisDb, _ := redis.Get()
	key := "crm:wechat:accessToken"
	accessToken, err := redisDb.Get(context.Background(), key).Result()
	if err != nil {
		logger.Error(err.Error())
	}
	if accessToken == "" {
		getTokenUrl := "https://api.weixin.qq.com/cgi-bin/token"
		param := url.Values{}
		param.Set("appid", c.appid)
		param.Set("secret", c.appSecret)
		param.Set("grant_type", "client_credential")

		req, _ := http.NewRequest(http.MethodGet, getTokenUrl+"?"+param.Encode(), nil)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			logger.Info("获取微信AppToken失败", err.Error())
		}
		all, err := io.ReadAll(resp.Body)
		if err != nil {
			return ""
		}
		accessToken := &accessTokenResp{}
		err = json.Unmarshal(all, accessToken)
		if err != nil {
			logger.Error(err.Error())
		}
		redisDb, _ := redis.Get()
		redisDb.Set(context.Background(), key, accessToken.AccessToken, time.Duration(accessToken.ExpiresIn-60)*time.Second)
	}

	return accessToken
}

func (c *Client) Code2Session(code string) (*Code2SessionResp, error) {
	reqUrl := "https://api.weixin.qq.com/sns/jscode2session"
	param := url.Values{}
	param.Set("appid", c.appid)
	param.Set("secret", c.appSecret)
	param.Set("js_code", code)
	param.Set("grant_type", "authorization_code")
	req, _ := http.NewRequest(http.MethodGet, reqUrl+"?"+param.Encode(), nil)
	client := &http.Client{}
	do, err := client.Do(req)
	if err != nil {
		logger.Error("微信小程序登录错误", err.Error())
		return nil, err
	}
	all, err := io.ReadAll(do.Body)
	if err != nil {
		return nil, err
	}
	resp := Code2SessionResp{}
	json.Unmarshal(all, &resp)
	return &resp, nil
}
