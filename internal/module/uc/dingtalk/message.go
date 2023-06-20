package dingtalk

import (
	"encoding/json"
	"errors"
	"fmt"
	logger "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"io"
	"strings"
)

type Message struct {
	MsgType string `json:"msgtype"` // text/image/voice/file/link/markdown/action_card
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
	Image struct {
		MediaId string `json:"media_id"`
	} `json:"image"`
	Voice struct {
		MediaId  string `json:"media_id"`
		Duration string `json:"duration"`
	} `json:"voice"`
	File struct {
		MediaId string `json:"media_id"`
	} `json:"file"`
	Link struct {
		Title      string `json:"title"`
		Text       string `json:"text"`
		PicUrl     string `json:"picUrl"`
		MessageUrl string `json:"messageUrl"`
	} `json:"link"`
	Markdown struct {
		Title string `json:"title"`
		Text  string `json:"text"`
	} `json:"markdown"`
	ActionCard struct {
		Title       string `json:"title"`
		Markdown    string `json:"markdown"`
		SingleTitle string `json:"single_title"`
		SingleUrl   string `json:"single_url"`
	} `json:"action_card"`
}

func (c *client) SendMessage(dingStaffIdList []string, msg *Message) (taskId string, err error) {
	if err = c.checkAccessToken(); err != nil {
		return "", err
	}
	if msg == nil || msg.MsgType == "" {
		return "", errors.New("msg is nil or msgtype is empty")
	}
	msgJson, err := json.Marshal(msg)
	uidList := strings.Join(dingStaffIdList, ",")
	body := fmt.Sprintf(
		"{\"agent_id\":%d,\"userid_list\":\"%s\",\"msg\":%s}",
		c.agentId, uidList, string(msgJson),
	)

	resp, err := c.client.Post(
		c.baseUrl+"/topapi/message/corpconversation/asyncsend_v2?access_token="+c.accessToken,
		"application/json",
		strings.NewReader(body),
	)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	resJson := gjson.ParseBytes(res)
	// 错误处理
	if resJson.Get("errcode").Int() != 0 {
		err = errors.New(resJson.Get("errmsg").String())
		logger.Error(err)
		return "", err
	}
	return resJson.Get("task_id").String(), nil
}

func (c *client) GetMessageSendProgress(taskId string) (success bool, err error) {
	if err = c.checkAccessToken(); err != nil {
		return
	}
	if taskId == "" {
		return false, errors.New("task_id is empty")
	}
	body := fmt.Sprintf(
		"{\"agent_id\":%d,\"task_id\":\"%s\"}",
		c.agentId, taskId,
	)
	resp, err := c.client.Post(
		c.baseUrl+"/topapi/message/corpconversation/getsendprogress?access_token="+c.accessToken,
		"application/json",
		strings.NewReader(body),
	)
	if err != nil {
		logger.Errorln(err)
		return false, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	res, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Errorln(err)
		return false, err
	}
	resJson := gjson.ParseBytes(res)
	// 错误处理
	if resJson.Get("errcode").Int() != 0 {
		err = errors.New(resJson.Get("errmsg").String())
		logger.Error(err)
		return false, err
	}
	return resJson.Get("progress.status").Int() == 2, nil
}
