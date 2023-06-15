package dingtalk

import (
	"context"
	"encoding/json"
	"fmt"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/module/uc/model"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

// GetUserOutsideDingtalk 根据code获取用户信息（浙政钉外三方应用扫码登录)
func GetUserOutsideDingtalk(source *model.ThirdSource, code string) (staffId, unionId string, err error) {
	c := checkClient(source)
	if c == nil {
		return "", "", nil
	}
	var accessToken string
	accessToken, err = c.GetUserAccessToken(code)
	if err != nil {
		logger.Errorln(err)
		return
	}
	var userInfo string
	userInfo, err = c.GetUserInfoByUserAccessToken(accessToken)
	if err != nil {
		logger.Errorln(err)
		return
	}
	uj := gjson.Parse(userInfo)
	if uj.Get("errorcode").Int() != 0 {
		err = fmt.Errorf("获取用户信息失败，错误码：%d，错误信息：%s", uj.Get("errorcode").Int(), uj.Get("errmsg").String())
		logger.Errorln(err)
		return
	}
	if uj.Get("code").Exists() && uj.Get("message").Exists() {
		err = fmt.Errorf("获取用户信息失败，错误码：%s，错误信息：%s", uj.Get("code").String(), uj.Get("message").String())
		logger.Errorln(err)
		return
	}
	unionId = uj.Get("unionId").String()
	if unionId == "" {
		err = fmt.Errorf("unionId为空")
		logger.Errorln(err)
		return
	}
	staffId, err = c.GetStaffIdByUnionId(unionId)
	if err != nil {
		logger.Errorln(err)
		return
	}
	return
}

func (c *client) GetStaffIdByUnionId(unionId string) (staffId string, err error) {
	err = c.checkAccessToken()
	if err != nil {
		return
	}

	body := map[string]string{
		"unionid": unionId,
	}
	bodyJson, err := json.Marshal(body)
	if err != nil {
		logger.Errorln(err)
		return
	}
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	request, err := http.NewRequestWithContext(timeoutCtx, "POST", c.baseUrl+"/topapi/user/getbyunionid?access_token="+c.accessToken, strings.NewReader(string(bodyJson)))
	if err != nil {
		logger.Errorln(err)
		return
	}
	request.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		logger.Errorln(err)
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	res, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Errorln(err)
		return
	}
	resJson := gjson.ParseBytes(res)
	if resJson.Get("errcode").Int() != 0 {
		err = fmt.Errorf(resJson.Get("errmsg").String())
		logger.Errorln(err)
		return
	}
	staffId = resJson.Get("result.userid").String()
	return
}

func (c *client) GetUserAccessToken(code string) (accessToken string, err error) {
	body := map[string]string{
		"clientId":     c.appKey,
		"clientSecret": c.appSecret,
		"code":         code,
		"grantType":    "authorization_code",
	}
	bodyJson, err := json.Marshal(body)
	if err != nil {
		logger.Errorln(err)
		return
	}
	resp, err := http.Post(c.newBaseUrl+"/v1.0/oauth2/userAccessToken", "application/json", strings.NewReader(string(bodyJson)))
	if err != nil {
		logger.Errorln(err)
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	resJson := gjson.ParseBytes(res)
	// TODO 错误处理
	accessToken = resJson.Get("accessToken").String()
	return
}

// GetUserInfoByUserAccessToken 根据用户accessToken获取用户信息
//
//	{
//	 "nick" : "zhangsan",
//	 "avatarUrl" : "https://xxx",
//	 "mobile" : "150xxxx9144",
//	 "openId" : "123",
//	 "unionId" : "z21HjQliSzpw0Yxxxx",
//	 "email" : "zhangsan@alibaba-inc.com",
//	 "stateCode" : "86"
//	}
func (c *client) GetUserInfoByUserAccessToken(userAccessToken string) (userJson string, err error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	request, err := http.NewRequestWithContext(timeoutCtx, http.MethodGet, c.newBaseUrl+"/v1.0/contact/users/me", nil)
	request.Header.Add("x-acs-dingtalk-access-token", userAccessToken)

	// 进行请求
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		logger.Errorln(err)
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("获取用户信息失败，错误码：%d, 错误信息: %s", resp.StatusCode, string(res))
		logger.Errorln(err)
		return
	}
	userJson = string(res)
	return
}

func (c *client) getUserInfoByCode(code string) (staffId string, err error) {
	if err = c.checkAccessToken(); err != nil {
		return "", err
	}
	q := url.Values{}
	q.Set("access_token", c.accessToken)
	q.Set("code", code)
	query := q.Encode()
	resp, err := http.Get("https://oapi.dingtalk.com/user/getuserinfo?" + query)
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	resJson := gjson.ParseBytes(res)
	// TODO 错误处理
	staffId = resJson.Get("userid").String()
	//name = resJson.Get("name").String()
	//manager = resJson.Get("is_sys").Bool()
	//level = int(resJson.Get("sys_level").Int())
	return
}

type GetUserResponse struct {
	StaffId        string
	UnionId        string
	Name           string
	Avatar         string
	ManagerStaffId string
	JobNumber      string
	Title          string
	LoginId        string
	DeptIdList     []int
	Senior         bool              // 高管
	Admin          bool              // 管理员
	Boss           bool              // 老板
	LeaderInfoList []*DeptLeaderInfo // 部门中是否领导
}
type DeptLeaderInfo struct {
	DeptId int
	Leader bool
}

func (c *client) getUser(staffId string) (string, error) {
	if err := c.checkAccessToken(); err != nil {
		return "", err
	}
	resp, err := c.client.Post(
		c.baseUrl+"/topapi/v2/user/get?access_token="+c.accessToken,
		"application/json",
		strings.NewReader(fmt.Sprintf("{\"userid\":\"%s\"}", staffId)),
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

	return string(res), err

	//resJson := gjson.ParseBytes(res)
	//// TODO 错误处理
	//uj := resJson.Get("result")
	//u := &GetUserResponse{
	//	StaffId:        uj.Get("userid").String(),
	//	UnionId:        uj.Get("unionid").String(),
	//	Name:           uj.Get("name").String(),
	//	Avatar:         uj.Get("avatar").String(),
	//	ManagerStaffId: uj.Get("manager_userid").String(),
	//	JobNumber:      uj.Get("job_number").String(),
	//	Title:          uj.Get("title").String(),
	//	LoginId:        uj.Get("login_id").String(),
	//	//DeptIdList:     uj.Get(),
	//	Senior: uj.Get("senior").Bool(),
	//	Admin:  uj.Get("admin").Bool(),
	//	Boss:   uj.Get("boss").Bool(),
	//	//LeaderInfoList: DeptLeaderInfo{},
	//}
	//deptIdsJson := uj.Get("dept_id_list").Array()
	//for _, did := range deptIdsJson {
	//	u.DeptIdList = append(u.DeptIdList, int(did.Int()))
	//}
	//
	//deptLeadersJson := uj.Get("leader_in_dept").Array()
	//for _, dl := range deptLeadersJson {
	//	u.LeaderInfoList = append(
	//		u.LeaderInfoList,
	//		&DeptLeaderInfo{
	//			DeptId: int(dl.Get("dept_id").Int()),
	//			Leader: dl.Get("leader").Bool(),
	//		},
	//	)
	//}
	//
	//return u, nil
}

func (c *client) getUserIdsInDepartment(deptId int) ([]string, error) {
	if err := c.checkAccessToken(); err != nil {
		return nil, err
	}

	resp, err := c.client.Post(
		c.baseUrl+"/topapi/user/listid?access_token="+c.accessToken,
		"application/json",
		strings.NewReader(fmt.Sprintf("{\"dept_id\":%d}", deptId)),
	)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resJson := gjson.ParseBytes(res)
	// TODO 错误处理
	rj := resJson.Get("result.userid_list").Array()

	var ul []string
	for _, r := range rj {
		ul = append(ul, r.String())
	}
	return ul, nil
}

func GetUser(source *model.ThirdSource, staffId string) (string, error) {
	c := checkClient(source)
	if c == nil {
		return "", nil
	}
	return c.getUser(staffId)
}

func GetStaffIdByCode(source *model.ThirdSource, authCode string) (string, error) {
	c := checkClient(source)
	if c == nil {
		return "", nil
	}
	staffId, err := c.getUserInfoByCode(authCode)
	if err != nil {
		logger.Errorln(err)
		return "", err
	}
	return staffId, nil
}
