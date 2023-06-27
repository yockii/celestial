package mq

import (
	"github.com/panjf2000/ants/v2"
	"time"
)

var (
	topicMap = make(map[string][]func(*Message))
)

// RegisterTopic registers a topic and its handler.
func RegisterTopic(topic string, handler func(*Message)) {
	if _, ok := topicMap[topic]; !ok {
		topicMap[topic] = make([]func(*Message), 0)
	}
	topicMap[topic] = append(topicMap[topic], handler)
}

// Publish publishes a message to a topic.
func Publish(topic string, msg *Message, delay ...time.Duration) {
	if len(delay) > 0 {
		time.AfterFunc(delay[0], func() {
			publish(topic, msg)
		})
	} else {
		_ = ants.Submit(func() {
			publish(topic, msg)
		})
	}
}

func publish(topic string, msg *Message) {
	if handlers, ok := topicMap[topic]; ok {
		for _, handler := range handlers {
			handler(msg)
		}
	}
}
