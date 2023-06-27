package mq

type Message struct {
	Topic string
	Data  interface{}
}

type TaskMemberAddedMessage struct {
	TaskId       uint64
	MemberIdList []uint64
}
