package mq

type Message struct {
	Topic string
	Data  interface{}
}

type TaskMemberAddedMessage struct {
	TaskId       uint64
	MemberIdList []uint64
}

type IssueAssignedMessage struct {
	IssueId    uint64
	OperatorId uint64
	AssigneeId uint64
}
