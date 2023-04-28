package dingtalk

import (
	"fmt"
	"github.com/yockii/celestial/internal/module/uc/model"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/tidwall/gjson"
)

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
	defer resp.Body.Close()
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
