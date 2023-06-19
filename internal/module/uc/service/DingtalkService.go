package service

import (
	"errors"
	logger "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/yockii/celestial/internal/module/uc/dingtalk"
	"github.com/yockii/celestial/internal/module/uc/model"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-core/util"
	"gorm.io/gorm"
)

var DingtalkService = new(dingtalkService)

type dingtalkService struct{}

func (s *dingtalkService) SyncDingUserByThirdSourceOutsideDingtalk(source *model.ThirdSource, code string, withDept bool) (user *model.User, err error) {
	staffId, _, err := dingtalk.GetUserOutsideDingtalk(source, code)
	if err != nil {
		return nil, err
	}

	return s.SyncDingUserByThirdSource(source, staffId, withDept)
}

func (s *dingtalkService) SyncDingUserByThirdSource(source *model.ThirdSource, staffId string, withDept bool) (matchingUser *model.User, err error) {
	user := &model.ThirdUser{
		SourceID: source.ID,
		OpenID:   staffId,
	}
	if err = database.DB.Where(user).Attrs(model.ThirdUser{
		ID:         util.SnowflakeId(),
		SourceCode: source.Code,
		Status:     model.UserStatusNormal,
	}).FirstOrCreate(user).Error; err != nil {
		logger.Error(err)
		return nil, err
	}

	dingUserResp, err := dingtalk.GetUser(source, staffId)
	if err != nil {
		return nil, err
	}
	if dingUserResp == "" {
		logger.Warn("未能获取到钉钉用户信息")
		return nil, errors.New("未能获取到钉钉用户信息")
	}
	dingUserRespJson := gjson.Parse(dingUserResp)
	if errCode := dingUserRespJson.Get("errcode").Int(); errCode != 0 {
		err = errors.New(dingUserRespJson.Get("errmsg").String())
		logger.Error(err)
		return nil, err
	}

	// 检查匹配的用户
	matchingUser = new(model.User)
	matchConf := gjson.Parse(source.MatchConfig)
	matchFields := matchConf.Get("match").Array()
	hasMatchField := false
	for _, matchField := range matchFields {
		switch matchField.String() {
		case "realName":
			matchingUser.RealName = dingUserRespJson.Get(matchConf.Get("realName").String()).String()
			hasMatchField = true
		case "mobile":
			matchingUser.Mobile = dingUserRespJson.Get(matchConf.Get("mobile").String()).String()
			hasMatchField = true
		case "email":
			matchingUser.Email = dingUserRespJson.Get(matchConf.Get("email").String()).String()
			hasMatchField = true
			// 以后继续扩展其他可用字段
		}
	}
	if !hasMatchField {
		logger.Warn("未配置匹配字段")
		return nil, errors.New("未配置匹配字段")
	}

	if err = database.DB.Where(matchingUser).First(matchingUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			matchingUser.ID = util.SnowflakeId()
			matchingUser.Status = model.UserStatusNormal
			if usernameField := matchConf.Get("username").String(); usernameField != "" {
				matchingUser.Username = dingUserRespJson.Get(usernameField).String()
			} else {
				matchingUser.Username = matchingUser.RealName
			}

			// 获取默认角色
			defaultRole := &model.Role{
				DefaultRole: 1,
			}
			if err = database.DB.Where(defaultRole).First(defaultRole).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					logger.Warn("未配置默认角色")
				} else {
					logger.Error(err)
					return nil, err
				}
			}
			if defaultRole != nil && defaultRole.ID > 0 {
				// 添加用户的同时添加默认角色
				if err = database.DB.Transaction(func(tx *gorm.DB) error {
					if err = tx.Create(matchingUser).Error; err != nil {
						logger.Error(err)
						return err
					}
					userRole := &model.UserRole{
						ID:     util.SnowflakeId(),
						UserID: matchingUser.ID,
						RoleID: defaultRole.ID,
					}
					if err = tx.Create(userRole).Error; err != nil {
						logger.Error(err)
						return err
					}
					return nil
				}); err != nil {
					return nil, err
				}
			} else {
				if err = database.DB.Create(matchingUser).Error; err != nil {
					logger.Error(err)
					return nil, err
				}
			}
		} else {
			logger.Error(err)
			return nil, err
		}
	}

	// 只更新第三方用户信息的userID
	if err = database.DB.Model(user).Update("user_id", matchingUser.ID).Error; err != nil {
		logger.Error(err)
		return nil, err
	}

	if withDept {
		// 同步部门信息
		deptIdListJson := dingUserRespJson.Get("result.dept_id_list").Array()
		var deptIdList []uint64
		for _, deptIdJson := range deptIdListJson {
			dept := model.Department{
				ExternalID: deptIdJson.String(),
			}
			if err = database.DB.Where(&dept).First(&dept).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					var dingDept *dingtalk.GetDepartmentResponse
					dingDept, err = dingtalk.GetDepartment(source, deptIdJson.String())
					if err != nil {
						return nil, err
					}
					dept.OrderNum = dingDept.OrderNum
					dept.Name = dingDept.Name
					// 获取父级部门
					if dingDept.ParentId != "" {
						parentDept := &model.Department{
							ExternalID: dingDept.ParentId,
						}
						if err = database.DB.Where(&parentDept).First(&parentDept).Error; err != nil {
							if errors.Is(err, gorm.ErrRecordNotFound) {
								parentDept, err = s.SyncDingDept(source.CorpId, dingDept.ParentId, false)
								if err != nil {
									return nil, err
								}
							} else {
								logger.Error(err)
								return nil, err
							}
						}
						dept.ParentID = parentDept.ID
						dept.FullPath = parentDept.FullPath + "/" + dept.Name
					} else {
						dept.FullPath = dept.Name
					}
					dept.ID = util.SnowflakeId()

					if err = database.DB.Create(&dept).Error; err != nil {
						logger.Error(err)
						return nil, err
					}
				} else {
					logger.Error(err)
					return nil, err
				}
			}
			deptIdList = append(deptIdList, dept.ID)
		}
		// 更新用户部门关系
		err = database.DB.Transaction(func(tx *gorm.DB) error {
			// 查询原有部门信息，并更新用户部门关系
			var userDeptList []*model.UserDepartment
			if err = tx.Where(&model.UserDepartment{
				UserID: matchingUser.ID,
			}).Find(&userDeptList).Error; err != nil {
				logger.Error(err)
				return err
			}
			// 获取原有部门ID
			var oldDeptIdList []uint64
			for _, userDept := range userDeptList {
				oldDeptIdList = append(oldDeptIdList, userDept.DepartmentID)
			}
			// 获取需要删除的部门ID
			var deleteDeptIdList []uint64
			for _, oldDeptId := range oldDeptIdList {
				for _, deptId := range deptIdList {
					if oldDeptId == deptId {
						break
					}
					deleteDeptIdList = append(deleteDeptIdList, oldDeptId)
				}
			}
			// 删除需要删除的部门ID
			if len(deleteDeptIdList) > 0 {
				if err = tx.Where(&model.UserDepartment{UserID: matchingUser.ID}).Where("department_id in ?", deleteDeptIdList).Delete(&model.UserDepartment{}).Error; err != nil {
					logger.Error(err)
					return err
				}
			}
			// 获取需要添加的部门ID
			var addDeptIdList []uint64
			for _, deptId := range deptIdList {
				for _, oldDeptId := range oldDeptIdList {
					if deptId == oldDeptId {
						break
					}
					addDeptIdList = append(addDeptIdList, deptId)
				}
			}
			// 添加新的部门关系
			for _, deptId := range addDeptIdList {
				if err = tx.Create(&model.UserDepartment{
					ID:           util.SnowflakeId(),
					UserID:       matchingUser.ID,
					DepartmentID: deptId,
				}).Error; err != nil {
					logger.Error(err)
					return err
				}
			}
			return nil
		})
	}

	return matchingUser, nil
}

func (s *dingtalkService) SyncDingUser(corpId string, staffId string, withDept bool) (*model.User, error) {
	source, err := ThirdSourceService.Instance(&model.ThirdSource{
		CorpId: corpId,
	})
	if err != nil {
		return nil, err
	}
	if source == nil {
		logger.Warn("未找到对应的第三方源配置")
		return nil, nil
	}
	return s.SyncDingUserByThirdSource(source, staffId, withDept)
}

func (s *dingtalkService) SyncDingDept(corpId string, deptId string, withChildren bool) (*model.Department, error) {
	source, err := ThirdSourceService.Instance(&model.ThirdSource{
		CorpId: corpId,
	})
	if err != nil {
		return nil, err
	}
	if source == nil {
		logger.Warn("未找到对应的第三方源配置")
		return nil, nil
	}

	dept := &model.Department{
		ExternalID: deptId,
	}
	if err = database.DB.Where(dept).First(dept).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			var dingDept *dingtalk.GetDepartmentResponse
			dingDept, err = dingtalk.GetDepartment(source, deptId)
			if err != nil {
				return nil, err
			}
			dept.OrderNum = dingDept.OrderNum
			dept.Name = dingDept.Name
			// 获取父级部门
			if dingDept.ParentId != "" {
				parentDept := &model.Department{
					ExternalID: dingDept.ParentId,
				}
				if err = database.DB.Where(&parentDept).First(&parentDept).Error; err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						parentDept, err = s.SyncDingDept(corpId, dingDept.ParentId, false)
						if err != nil {
							return nil, err
						}
					} else {
						logger.Error(err)
						return nil, err
					}
				}
				dept.ParentID = parentDept.ID
				dept.FullPath = parentDept.FullPath + "/" + dept.Name
			} else {
				dept.FullPath = dept.Name
			}
			dept.ID = util.SnowflakeId()

			if err = database.DB.Create(&dept).Error; err != nil {
				logger.Error(err)
				return nil, err
			}
		} else {
			logger.Error(err)
			return nil, err
		}
	}

	if withChildren {
		// 同步子部门
		dingDeptList, err := dingtalk.GetChildrenDepartments(source, deptId)
		if err != nil {
			return nil, err
		}
		for _, dingDept := range dingDeptList {
			_, err = s.SyncDingDept(corpId, dingDept.DeptId, true)
			if err != nil {
				return nil, err
			}
		}
	}

	return dept, nil
}
