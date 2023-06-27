package message

import (
	"fmt"
	"github.com/yockii/celestial/internal/core/mq"
	taskService "github.com/yockii/celestial/internal/module/task/service"
	"github.com/yockii/celestial/internal/module/uc/dingtalk"
	ucModel "github.com/yockii/celestial/internal/module/uc/model"
	ucService "github.com/yockii/celestial/internal/module/uc/service"
	"github.com/yockii/ruomu-core/config"
)

var DingtalkAdapter = new(dingtalkMessageAdapter)

type dingtalkMessageAdapter struct{}

func (a dingtalkMessageAdapter) TaskMemberAdded(message *mq.Message) {
	if taskMemberAddedMessage, ok := message.Data.(*mq.TaskMemberAddedMessage); ok {
		// 获取任务信息
		task, err := taskService.ProjectTaskService.Instance(taskMemberAddedMessage.TaskId)
		if err != nil {
			return
		}
		if task == nil {
			return
		}

		// 任务负责人owner
		var owner *ucModel.User
		owner, err = ucService.UserService.Instance(&ucModel.User{ID: task.OwnerID})
		if err != nil {
			return
		}
		if owner == nil {
			return
		}

		// 获取当前可用的三方用户信息
		thirdUserList, err := ucService.ThirdUserService.ListByUserIDListAndSourceCode(taskMemberAddedMessage.MemberIdList, ucModel.ThirdSourceCodeDingtalk)
		sourceStaffIdsMap := make(map[uint64][]string)
		for _, thirdUser := range thirdUserList {
			if thirdUser.OpenID != "" {
				sourceStaffIdsMap[thirdUser.SourceID] = append(sourceStaffIdsMap[thirdUser.UserID], thirdUser.OpenID)
			}
		}

		// 每个任务成员发送钉消息
		for sourceID, staffIdList := range sourceStaffIdsMap {
			// 获取三方源信息
			var source *ucModel.ThirdSource
			source, err = ucService.ThirdSourceService.Instance(&ucModel.ThirdSource{ID: sourceID})
			if err != nil {
				continue
			}
			if source == nil {
				continue
			}
			msg := &dingtalk.Message{
				MsgType: dingtalk.MessageTypeActionCard,
			}
			msg.ActionCard.Title = task.Name
			msg.ActionCard.Markdown = "### " + task.Name + "\n" +
				"**任务描述**：" + task.TaskDesc + "\n\n" +
				"**任务负责人**：" + owner.RealName + "\n\n"
			msg.ActionCard.SingleTitle = "查看任务详情"
			msg.ActionCard.SingleUrl = fmt.Sprintf(

				"%s/project/detail/%d/task?id=%d",
				config.GetString("server.baseUrl"),
				task.ProjectID,
				task.ID,
			)

			_, _ = dingtalk.SendMessage(source, staffIdList, msg)
		}
	}
}

func InitDingtalkMessageAdapter() {
	mq.RegisterTopic(mq.TopicTaskMemberAdded, DingtalkAdapter.TaskMemberAdded)
}
