package hook

import (
	"encoding/json"
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/Sirupsen/logrus"
)

type Logrushooker struct {
	name     string
	producer sarama.AsyncProducer
	topic    string
	levels   []logrus.Level
}

func NewLogrushooker(name string, producer sarama.AsyncProducer, topic string, lv int) *Logrushooker {
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

	msg := &sarama.ProducerMessage{
		Topic: h.topic,
		Key:   sarama.StringEncoder(fmt.Sprintf("log_%s_%s", h.name, entry.Level.String())),
		Value: sarama.ByteEncoder(bytes),
	}

	h.producer.Input() <- msg
	return nil
}

func (h *Logrushooker) Levels() []logrus.Level {
	return h.levels
}
