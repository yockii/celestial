package service

import (
	"encoding/json"
	"errors"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/module/test/domain"
	"github.com/yockii/celestial/internal/module/test/model"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-core/util"
	"sync"
	"time"
)

var ProjectTestService = new(projectTestService)

type projectTestService struct{}

func (s projectTestService) Add(instance *model.ProjectTest) (success bool, err error) {
	if instance.ProjectID == 0 {
		err = errors.New("projectId is required")
		return
	}
	// 查出最大轮次+1
	var maxRound int
	if err = database.DB.Model(&model.ProjectTest{}).Where(&model.ProjectTest{
		ProjectID: instance.ProjectID,
	}).Select("COALESCE(max(round),0)").Scan(&maxRound).Error; err != nil {
		logger.Errorln(err)
		return
	}

	instance.Round = maxRound + 1
	instance.StartTime = time.Now().UnixMilli()
	instance.ID = util.SnowflakeId()

	if err = database.DB.Create(instance).Error; err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

func (s projectTestService) Delete(id uint64) (success bool, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}
	if err = database.DB.Delete(&model.ProjectTest{ID: id}).Error; err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

func (s projectTestService) Update(instance *model.ProjectTest) (success bool, err error) {
	if instance.ID == 0 {
		err = errors.New("id is required")
		return
	}

	err = database.DB.Where(&model.ProjectTest{ID: instance.ID}).Updates(&model.ProjectTest{
		Remark:    instance.Remark,
		CreatorID: instance.CreatorID,
	}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

func (s projectTestService) Close(id, closerID uint64) (success bool, err error) {
	// 封版把所有记录全部取出并转为json存储到测试记录种
	if id == 0 || closerID == 0 {
		err = errors.New("id & closer Id is required")
		return
	}
	// 确认id的记录存在
	var instance model.ProjectTest
	if err = database.DB.Where(&model.ProjectTest{ID: id}).First(&instance).Error; err != nil {
		logger.Errorln(err)
		return
	}
	// 确认id的记录未封版
	if instance.EndTime != 0 {
		err = errors.New("this test has been closed")
		return
	}
	// 封版操作
	var result []*domain.ProjectTestCaseWithItemsWithSteps
	{ // 1、取出所有测试用例
		var testCaseList []*model.ProjectTestCase
		if err = database.DB.Where(&model.ProjectTestCase{
			ProjectID: instance.ProjectID,
		}).Find(&testCaseList).Error; err != nil {
			logger.Errorln(err)
			return
		}
		// 2、异步取出所有测试用例的测试项
		var wg sync.WaitGroup
		for _, testCase := range testCaseList {
			tcwi := &domain.ProjectTestCaseWithItemsWithSteps{
				ProjectTestCase: *testCase,
			}
			result = append(result, tcwi)
			wg.Add(1)
			// 异步取出所有测试项
			go func(caseWithItems *domain.ProjectTestCaseWithItemsWithSteps) {
				defer wg.Done()
				var items []*model.ProjectTestCaseItem
				if err = database.DB.Where(&model.ProjectTestCaseItem{
					TestCaseID: caseWithItems.ID,
				}).Find(&items).Error; err != nil {
					logger.Errorln(err)
					return
				}
				// 再次异步取出所有测试项的测试步骤
				var wg2 sync.WaitGroup
				for _, item := range items {
					itemWithSteps := &domain.ProjectTestCaseItemWithSteps{
						ProjectTestCaseItem: *item,
					}
					caseWithItems.Items = append(caseWithItems.Items, itemWithSteps)
					wg2.Add(1)
					go func(itemWithSteps *domain.ProjectTestCaseItemWithSteps) {
						defer wg2.Done()
						var steps []*model.ProjectTestCaseItemStep
						if err = database.DB.Where(&model.ProjectTestCaseItemStep{
							CaseItemID: itemWithSteps.ID,
						}).Find(&steps).Error; err != nil {
							logger.Errorln(err)
							return
						}
						itemWithSteps.Steps = steps
					}(itemWithSteps)
				}
				wg2.Wait()
			}(tcwi)
		}
		wg.Wait()
	}

	// 将测试记录转为json
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		logger.Errorln(err)
		return
	}
	// 将json存入测试记录
	if err = database.DB.Model(&model.ProjectTest{}).Where(&model.ProjectTest{ID: id}).
		Updates(&model.ProjectTest{
			TestRecord: string(jsonBytes),
			EndTime:    time.Now().UnixMilli(),
			CloserID:   instance.CloserID,
		}).Error; err != nil {
		logger.Errorln(err)
		return
	}

	success = true
	return
}

func (s projectTestService) List(instance *model.ProjectTest) (list []*model.ProjectTest, err error) {
	if instance.ProjectID == 0 {
		err = errors.New("projectId is required")
		return
	}
	err = database.DB.Omit("test_record").Where(&model.ProjectTest{
		ProjectID: instance.ProjectID,
	}).Find(&list).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	return
}

func (s projectTestService) Instance(id uint64) (instance *model.ProjectTest, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}
	instance = new(model.ProjectTest)
	err = database.DB.Where(&model.ProjectTest{ID: id}).First(&instance).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	return
}
