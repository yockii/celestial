package domain

import (
	"github.com/tidwall/gjson"
	projectModel "github.com/yockii/celestial/internal/module/project/model"
	"github.com/yockii/celestial/internal/module/task/model"
	"github.com/yockii/ruomu-core/server"
)

type ProjectTaskListTask struct {
	model.ProjectTask
	OnlyParent               bool                  `json:"onlyParent"`
	StartTimeCondition       *server.TimeCondition `json:"startTimeCondition"`
	EndTimeCondition         *server.TimeCondition `json:"endTimeCondition"`
	ActualStartTimeCondition *server.TimeCondition `json:"actualStartTimeCondition"`
	ActualEndTimeCondition   *server.TimeCondition `json:"actualEndTimeCondition"`
	CreateTimeCondition      *server.TimeCondition `json:"createTimeCondition"`
	UpdateTimeCondition      *server.TimeCondition `json:"updateTimeCondition"`
	OrderBy                  string                `json:"orderBy"`
}

type ProjectTaskWorkTimeStatistics struct {
	ProjectID uint64 `json:"projectId,string"`
	// 任务总数
	TaskCount int `json:"taskCount"`
	// 任务预计总工时
	EstimateDuration int64 `json:"estimateDuration"`
	// 任务实际总工时
	ActualDuration int64 `json:"actualDuration"`
}

type ProjectTaskWithMembers struct {
	model.ProjectTask
	Owner   *ProjectTaskMemberWithRealName   `json:"owner"`
	Members []*ProjectTaskMemberWithRealName `json:"members"`
	Issue   *projectModel.ProjectIssue       `json:"issue"`
}

func (ptwm *ProjectTaskWithMembers) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	ptwm.ID = j.Get("id").Uint()
	ptwm.ProjectID = j.Get("projectId").Uint()
	ptwm.RequirementID = j.Get("requirementId").Uint()
	ptwm.Priority = int(j.Get("priority").Int())
	ptwm.StageID = j.Get("stageId").Uint()
	ptwm.ModuleID = j.Get("moduleId").Uint()
	ptwm.OwnerID = j.Get("ownerId").Uint()
	ptwm.ParentID = j.Get("parentId").Uint()
	ptwm.Name = j.Get("name").String()
	ptwm.StartTime = j.Get("startTime").Int()
	ptwm.EndTime = j.Get("endTime").Int()
	ptwm.TaskDesc = j.Get("taskDesc").String()
	ptwm.ActualStartTime = j.Get("actualStartTime").Int()
	ptwm.ActualEndTime = j.Get("actualEndTime").Int()
	ptwm.EstimateDuration = j.Get("estimateDuration").Int()
	ptwm.ActualDuration = j.Get("actualDuration").Int()
	ptwm.Status = int(j.Get("status").Int())
	ptwm.CreatorID = j.Get("creatorId").Uint()
	ptwm.FullPath = j.Get("fullPath").String()
	ptwm.CreateTime = j.Get("createTime").Int()
	ptwm.UpdateTime = j.Get("updateTime").Int()
	ptwm.Members = make([]*ProjectTaskMemberWithRealName, 0)
	for _, m := range j.Get("members").Array() {
		member := &model.ProjectTaskMember{}
		member.UnmarshalJSON([]byte(m.Raw))
		ptwm.Members = append(ptwm.Members, &ProjectTaskMemberWithRealName{
			ProjectTaskMember: *member,
		})
	}

	return nil
}
