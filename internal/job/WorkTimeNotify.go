package job

import (
	"fmt"
	logger "github.com/sirupsen/logrus"
	logModel "github.com/yockii/celestial/internal/module/log/model"
	"github.com/yockii/celestial/internal/module/uc/dingtalk"
	ucModel "github.com/yockii/celestial/internal/module/uc/model"
	"github.com/yockii/celestial/pkg/task"
	"github.com/yockii/ruomu-core/database"
	"time"
)

func init() {
	// 每周五下午5点执行一次检查
	_, err := task.AddFunc("0 0 17 * * 5", CheckWorkTime)
	if err != nil {
		logger.Errorln("每周工时检查通知初始化失败! ", err)
		return
	}
}

func getThisMondayTimestamp(now time.Time) time.Time {
	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}
	weekMonday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
	fmt.Println(weekMonday.Format(time.DateTime))
	return weekMonday
}

func CheckWorkTime() {
	// 检查本周工时
	now := time.Now()
	monday := getThisMondayTimestamp(now)
	// 查询所有结束时间大于monday的工时记录
	var workTimeList []*logModel.WorkTime
	err := database.DB.Where("end_date > ?", monday.UnixMilli()).Find(&workTimeList).Error
	if err != nil {
		logger.Errorln("每周工时检查通知失败! ", err)
		return
	}
	// 按照用户分组
	userMap := make(map[uint64][]*logModel.WorkTime)
	for _, w := range workTimeList {
		userMap[w.UserID] = append(userMap[w.UserID], w)
	}
	var noNeedNotifyUserIDList []uint64
	// 检查每个用户的工时
	for userID, wtList := range userMap {
		// 检查本周工时是否超过40小时
		var totalWorkTime int64
		for _, w := range wtList {
			t := w.WorkTime
			if w.StartDate < monday.UnixMilli() {
				// 根据天数判断占比
				duration := monday.Sub(time.UnixMilli(w.StartDate))
				sub := int64(duration.Hours() / 24)
				total := int64(time.UnixMilli(w.EndDate).Sub(time.UnixMilli(w.StartDate)).Hours()/24) + 1
				t = t * sub / total
			}
			totalWorkTime += t
		}
		if totalWorkTime >= 40*3600 {
			noNeedNotifyUserIDList = append(noNeedNotifyUserIDList, userID)
		}
	}
	// 找出所有需要发送通知的用户对应的staffId
	// 先找出所有需要填报工时的用户
	var userList []*ucModel.User
	tx := database.DB.Select("id")
	if len(noNeedNotifyUserIDList) > 0 {
		tx = tx.Where("id not int (?)", noNeedNotifyUserIDList)
	}
	err = tx.Find(&userList, &ucModel.User{
		ExtType: 1,
	}).Error
	if err != nil {
		logger.Errorln("获取待通知工时填报用户失败! ", err)
		return
	}
	var notifyUserIdList []uint64
	for _, u := range userList {
		notifyUserIdList = append(notifyUserIdList, u.ID)
	}
	// 从thirdUser中获取所有notifyUserIdList对应的openId
	var staffList []*ucModel.ThirdUser
	err = database.DB.Model(&ucModel.ThirdUser{}).Where("user_id in (?)", notifyUserIdList).Find(&staffList).Error
	if err != nil {
		logger.Errorln("获取待通知工时填报用户三方信息失败! ", err)
		return
	}
	// 发送通知
	var sourceUserMap = make(map[uint64][]string)
	for _, s := range staffList {
		sourceUserMap[s.SourceID] = append(sourceUserMap[s.SourceID], s.OpenID)
	}
	for sourceID, staffIDList := range sourceUserMap {
		source := &ucModel.ThirdSource{ID: sourceID}
		err = database.DB.First(source).Error
		if err != nil {
			logger.Errorln("获取三方信息失败! ", err)
			continue
		}
		message := &dingtalk.Message{
			MsgType: "text",
		}
		message.Text.Content = "本周工时填报提醒: 本周工时未足40小时，请抓紧时间进入系统进行填报!"
		_, _ = dingtalk.SendMessage(source, staffIDList, message)
	}
}
