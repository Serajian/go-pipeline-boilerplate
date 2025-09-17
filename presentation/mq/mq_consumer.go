package mq

import (
	"fmt"

	"go-pipeline/pkg/logger"

	"github.com/IBM/sarama"
)

// ConsumerHandler for handling consume message
type ConsumerHandler struct{}

func NewConsumerHandler() *ConsumerHandler {
	return &ConsumerHandler{}
}

func (h *ConsumerHandler) Setup(sarama.ConsumerGroupSession) error {
	logger.GetLogger().Info(&logger.Log{
		Event: "start consumer",
	})
	return nil
}

func (h *ConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	logger.GetLogger().Info(&logger.Log{
		Event: "cleanup consumer",
		Additional: map[string]interface{}{
			"msg": "🧹 Kafka consumer cleanup done",
		},
	})
	return nil
}

func (h *ConsumerHandler) ConsumeClaim(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) error {
	for msg := range claim.Messages() {
		m := fmt.Sprintf("📩 Consumed: topic=%s partition=%d offset=%d key=%s value=%s",
			msg.Topic, msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
		logger.GetLogger().Info(&logger.Log{
			Event: "consume",
			Additional: map[string]interface{}{
				"msg": m,
			},
		})

		session.MarkMessage(msg, "done")
	}
	return nil
}
