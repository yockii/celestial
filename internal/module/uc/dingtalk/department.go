package dingtalk

import (
	"errors"
	"fmt"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/module/uc/model"
	"io"
	"strings"

	"github.com/tidwall/gjson"
)

type GetDepartmentResponse struct {
	DeptId   string
	Name     string
	ParentId string

	OrderNum    int
	HasChildren bool

	OriginalJson string
}

func (c *client) GetChildrenDepartments(parentId string) ([]*GetDepartmentResponse, error) {
	if err := c.checkAccessToken(); err != nil {
		return nil, err
	}
	body := "{}"
	if parentId != "" {
		body = fmt.Sprintf("{\"dept_id\":%s}", parentId)
	}
	resp, err := c.client.Post(
		c.baseUrl+"/topapi/v2/department/listsub?access_token="+c.accessToken,
		"application/json",
		strings.NewReader(body),
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
	// 错误处理
	if resJson.Get("errcode").Int() != 0 {
		err = errors.New(resJson.Get("errmsg").String())
		logger.Error(err)
		return nil, err
	}
	rj := resJson.Get("result").Array()
	var result []*GetDepartmentResponse
	for _, j := range rj {
		dept, err := c.GetDepartment(j.Get("dept_id").String())
		if err != nil {
			return nil, err
		}
		result = append(result, dept)
	}
	return result, nil
}

func (c *client) GetDepartment(deptId string) (*GetDepartmentResponse, error) {
	if err := c.checkAccessToken(); err != nil {
		return nil, err
	}
	if deptId == "" {
		return nil, nil
	}
	body := fmt.Sprintf("{\"dept_id\":%s}", deptId)
	resp, err := c.client.Post(
		c.baseUrl+"/topapi/v2/department/get?access_token="+c.accessToken,
		"application/json",
		strings.NewReader(body),
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
	// 处理错误
	if resJson.Get("errcode").Int() != 0 {
		err = errors.New(resJson.Get("errmsg").String())
		logger.Error(err)
		return nil, err
	}

	rj := resJson.Get("result")
	dept := &GetDepartmentResponse{
		DeptId:       rj.Get("dept_id").String(),
		Name:         rj.Get("name").String(),
		OrderNum:     int(rj.Get("order").Int()),
		ParentId:     rj.Get("parent_id").String(),
		HasChildren:  rj.Get("group_contain_sub_dept").Bool(),
		OriginalJson: string(res),
	}

	return dept, nil
}

func GetChildrenDepartments(source *model.ThirdSource, parentId string) ([]*GetDepartmentResponse, error) {
	c := checkClient(source)
	if c == nil {
		return nil, nil
	}
	return c.GetChildrenDepartments(parentId)
}

func GetDepartment(source *model.ThirdSource, externalDeptId string) (*GetDepartmentResponse, error) {
	c := checkClient(source)
	if c == nil {
		return nil, nil
	}
	return c.GetDepartment(externalDeptId)
}
