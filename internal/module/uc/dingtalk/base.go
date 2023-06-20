package dingtalk

import (
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/module/uc/constant"
	"github.com/yockii/celestial/internal/module/uc/model"
	"github.com/yockii/ruomu-core/cache"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/tidwall/gjson"
)

var clientMap = make(map[uint64]*client)

type client struct {
	sourceId   uint64
	client     *http.Client
	baseUrl    string
	newBaseUrl string

	appKey    string
	appSecret string

	agentId string

	accessToken string
	deadline    time.Time
}

func checkClient(source *model.ThirdSource) *client {
	c, ok := clientMap[source.ID]
	if !ok {
		c = &client{}
		c.init(source)
		clientMap[source.ID] = c
	}
	return c
}

func (c *client) init(source *model.ThirdSource) {
	c.sourceId = source.ID
	conf := gjson.Parse(source.Configuration)
	c.appKey = conf.Get("appKey").String()
	c.appSecret = conf.Get("appSecret").String()
	c.client = &http.Client{}
	c.baseUrl = "https://oapi.dingtalk.com"
	c.newBaseUrl = "https://api.dingtalk.com"

	c.agentId = conf.Get("agentId").String()
}

func (c *client) doGetAccessToken() (string, int64, error) {
	q := url.Values{}
	q.Set("appkey", c.appKey)
	q.Set("appsecret", c.appSecret)

	resp, err := c.client.Get(c.baseUrl + "/gettoken?" + q.Encode())
	if err != nil {
		return "", 0, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	res, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Errorln(err)
		return "", 0, err
	}
	resJson := gjson.ParseBytes(res)
	// 错误处理
	return resJson.Get("access_token").String(), resJson.Get("expires_in").Int(), nil
}

func (c *client) getAccessToken() error {
	redisConn := cache.Get()
	defer func(redisConn redis.Conn) {
		_ = redisConn.Close()
	}(redisConn)
	key := constant.RedisKeyDingtalkAccessToken + ":" + strconv.FormatUint(c.sourceId, 10)
	token, _ := redis.String(redisConn.Do("GET", key))
	if token != "" {
		ttl, _ := redis.Int(redisConn.Do("TTL", key))
		logger.Debug("获取存储的access token", token, "剩余有效时间(s)", ttl)
		if ttl > 60 {
			c.accessToken = token
			c.deadline = time.Now().Local().Add(time.Duration(ttl-10) * time.Second) // 减10秒
			return nil
		}
	}

	token, expired, err := c.doGetAccessToken()
	if err != nil {
		return err
	}
	c.accessToken = token
	c.deadline = time.Now().Local().Add(time.Duration(expired) * time.Second)
	_, _ = redisConn.Do("SETEX", key, expired, token)
	return nil
}

func (c *client) checkAccessToken() error {
	if c.accessToken == "" || time.Now().Local().After(c.deadline) {
		return c.getAccessToken()
	}
	return nil
}
