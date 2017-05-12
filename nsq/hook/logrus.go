package hook

import (
	"encoding/json"

	"fmt"

	"github.com/Sirupsen/logrus"
	nsq "github.com/go-nsq"
)

type Logrushooker struct {
	name     string
	producer *nsq.Producer
	topic    string
	levels   []logrus.Level
}

func NewLogrushooker(name string, producer *nsq.Producer, topic string, lv int) *Logrushooker {
	hook := &Logrushooker{
		name:     name,
		producer: producer,
		topic:    topic,
		levels:   make([]logrus.Level, 0),
	}

	lvs := logrus.Level(lv)
	for i := logrus.PanicLevel; i <= logrus.DebugLevel; i++ {
		if i <= lvs {
			hook.levels = append(hook.levels, i)
		}
	}

	return hook
}

func (h *Logrushooker) Fire(entry *logrus.Entry) error {
	data := make(map[string]interface{})
	data["Name"] = h.name
	data["Level"] = entry.Level.String()
	data["Time"] = entry.Time
	data["Message"] = entry.Message
	if errData, ok := entry.Data[logrus.ErrorKey]; ok {
		if err, ok := errData.(error); ok && entry.Data[logrus.ErrorKey] != nil {
			data[logrus.ErrorKey] = err.Error()
		}
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("Failed to send log entry to nsq: %s", err)
	}

	h.producer.Publish(h.topic, bytes)
	return nil
}

func (h *Logrushooker) Levels() []logrus.Level {
	return h.levels
}
