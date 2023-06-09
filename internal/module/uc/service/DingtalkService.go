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

func (s *dingtalkService) SyncDingUserByThirdSourceOutsideDingtalk(source *model.ThirdSource, code string) (user *model.User, err error) {
	staffId, _, err := dingtalk.GetUserOutsideDingtalk(source, code)
	if err != nil {
		return nil, err
	}

	return s.SyncDingUserByThirdSource(source, staffId)
}

func (s *dingtalkService) SyncDingUserByThirdSource(source *model.ThirdSource, staffId string) (matchingUser *model.User, err error) {
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
			user.Info = dingUserResp
			user.UserID = matchingUser.ID

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
					if err = tx.Model(&model.ThirdUser{}).Where(&model.ThirdUser{ID: user.ID}).Updates(user).Error; err != nil {
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
				if err = database.DB.Transaction(func(tx *gorm.DB) error {
					if err = tx.Create(matchingUser).Error; err != nil {
						logger.Error(err)
						return err
					}

					if err = tx.Model(&model.ThirdUser{}).Where(&model.ThirdUser{ID: user.ID}).Updates(user).Error; err != nil {
						logger.Error(err)
						return err
					}

					return nil
				}); err != nil {
					return nil, err
				}

			}
		} else {
			logger.Error(err)
			return nil, err
		}
	}

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

				if err = database.DB.Transaction(func(tx *gorm.DB) error {
					if err = tx.Create(&dept).Error; err != nil {
						logger.Error(err)
						return err
					}
					if dept.ParentID != 0 {
						err = tx.Model(&model.Department{ID: dept.ParentID}).Update("child_count", gorm.Expr("child_count + ?", 1)).Error
						if err != nil {
							logger.Errorln(err)
							return err
						}
					}
					return nil
				}); err != nil {
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
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				logger.Error(err)
				return err
			}
		}
		// 获取原有部门ID
		var oldDeptIdList []uint64
		for _, userDept := range userDeptList {
			oldDeptIdList = append(oldDeptIdList, userDept.DepartmentID)
		}
		// 获取需要添加的部门ID
		var addDeptIdList []uint64
		// 获取需要删除的部门ID
		var deleteDeptIdList []uint64

		for i := len(deptIdList) - 1; i >= 0; i-- {
			var exists bool
			for k := len(oldDeptIdList) - 1; k >= 0; k-- {
				if deptIdList[i] == oldDeptIdList[k] {
					exists = true
					if k == len(oldDeptIdList)-1 {
						oldDeptIdList = oldDeptIdList[:k]
					} else {
						oldDeptIdList = append(oldDeptIdList[:k], oldDeptIdList[k+1:]...)
					}
					break
				}
			}
			if !exists {
				addDeptIdList = append(addDeptIdList, deptIdList[i])
				if i == len(deptIdList)-1 {
					deptIdList = deptIdList[:i]
				} else {
					deptIdList = append(deptIdList[:i], deptIdList[i+1:]...)
				}
			}
		}

		for i := len(oldDeptIdList) - 1; i >= 0; i-- {
			var exists bool
			for k := len(deptIdList) - 1; k >= 0; k-- {
				if oldDeptIdList[i] == deptIdList[k] {
					exists = true
					break
				}
			}
			if !exists {
				deleteDeptIdList = append(deleteDeptIdList, oldDeptIdList[i])
			}
		}

		// 删除需要删除的部门ID
		if len(deleteDeptIdList) > 0 {
			if err = tx.Where(&model.UserDepartment{UserID: matchingUser.ID}).Where("department_id in ?", deleteDeptIdList).Delete(&model.UserDepartment{}).Error; err != nil {
				logger.Error(err)
				return err
			}
		}
		// 添加新的部门关系
		var ud []*model.UserDepartment
		for _, deptId := range addDeptIdList {
			// 获取dept信息
			var dept *model.Department
			dept, err = DepartmentService.Instance(deptId)
			if err != nil {
				return err
			}
			ud = append(ud, &model.UserDepartment{
				ID:             util.SnowflakeId(),
				UserID:         matchingUser.ID,
				DepartmentID:   deptId,
				DepartmentPath: dept.FullPath,
			})
		}
		if len(ud) > 0 {
			if err = tx.Create(ud).Error; err != nil {
				logger.Error(err)
				return err
			}
		}
		return nil
	})

	return matchingUser, nil
}

func (s *dingtalkService) SyncDingUser(corpId string, staffId string) (*model.User, error) {
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
	return s.SyncDingUserByThirdSource(source, staffId)
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
				if err = database.DB.Where(parentDept).First(parentDept).Error; err != nil {
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

			if err = database.DB.Transaction(func(tx *gorm.DB) error {
				if err = tx.Create(&dept).Error; err != nil {
					logger.Error(err)
					return err
				}
				if dept.ParentID != 0 {
					err = tx.Model(&model.Department{ID: dept.ParentID}).Update("child_count", gorm.Expr("child_count + ?", 1)).Error
					if err != nil {
						logger.Errorln(err)
						return err
					}
				}
				return nil
			}); err != nil {
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

func (s *dingtalkService) SyncChildrenDepartmentsWithStaff(source *model.ThirdSource, parentId string) error {
	// 同步所有部门
	deptList, err := dingtalk.GetChildrenDepartments(source, parentId)
	if err != nil {
		return err
	}
	for _, dept := range deptList {
		d := &model.Department{ExternalID: dept.DeptId}
		if err = database.DB.Where(d).First(d).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				var dingDept *dingtalk.GetDepartmentResponse
				dingDept, err = dingtalk.GetDepartment(source, dept.DeptId)
				if err != nil {
					return err
				}
				d.OrderNum = dingDept.OrderNum
				d.Name = dingDept.Name
				// 获取父级部门
				if dingDept.ParentId != "" {
					parentDept := &model.Department{
						ExternalID: dingDept.ParentId,
					}
					if err = database.DB.Where(parentDept).First(parentDept).Error; err != nil {
						if errors.Is(err, gorm.ErrRecordNotFound) {
							// 没有找到父级部门，同步父级部门
							parentDept, err = s.SyncDingDept(source.CorpId, dingDept.ParentId, false)
							if err != nil {
								logger.Errorln(err)
								return err
							}
						} else {
							logger.Errorln(err)
							return err
						}
					}
					d.ParentID = parentDept.ID
					d.FullPath = parentDept.FullPath + "/" + d.Name
				} else {
					d.FullPath = d.Name
				}
				d.ID = util.SnowflakeId()
				d.ExternalJson = dept.OriginalJson

				if err = database.DB.Transaction(func(tx *gorm.DB) error {
					if err = tx.Create(d).Error; err != nil {
						logger.Errorln(err)
						return err
					}
					// 如果存在父级，则更新父级的子数量
					if d.ParentID != 0 {
						err = tx.Model(&model.Department{ID: d.ParentID}).Update("child_count", gorm.Expr("child_count + ?", 1)).Error
						if err != nil {
							logger.Errorln(err)
							return err
						}
					}

					return nil
				}); err != nil {
					return err
				}

			} else {
				logger.Errorln(err)
				return err
			}
		}

		// 同步该部门下的用户
		err = s.syncStaffInDepartment(source, dept)
		if err != nil {
			return err
		}

		// 同步子部门
		if err = s.SyncChildrenDepartmentsWithStaff(source, dept.DeptId); err != nil {
			return err
		}
	}

	return nil
}

// SyncAll 同步所有钉钉用户和部门相关信息
func (s *dingtalkService) SyncAll(source *model.ThirdSource) error {
	// 检查source
	if source == nil || source.ID == 0 {
		return errors.New("source不能为空")
	}

	return s.SyncChildrenDepartmentsWithStaff(source, "")
}

func (s *dingtalkService) syncStaffInDepartment(source *model.ThirdSource, dept *dingtalk.GetDepartmentResponse) (err error) {
	var staffIdList []string
	staffIdList, err = dingtalk.GetStaffIdListInDepartment(source, dept.DeptId)
	if err != nil {
		return
	}
	for _, staffId := range staffIdList {
		_, err = s.SyncDingUserByThirdSource(source, staffId)
		if err != nil {
			return
		}
	}
	return
}
