package kafka

import (
	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type Producer struct {
	producer sarama.SyncProducer
	logger   *zap.Logger
}

func NewProducer(brokers []string, logger *zap.Logger) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Producer.Partitioner = sarama.NewHashPartitioner

	prod, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &Producer{producer: prod, logger: logger}, nil
}

func (p *Producer) SendMessage(topic, key string, message []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(message),
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		p.logger.Error("failed to send message", zap.Error(err))
		return err
	}

	p.logger.Info("message sent",
		zap.Int32("partition", partition),
		zap.Int64("offset", offset))

	return nil
}

func (p *Producer) Close() error {
	return p.producer.Close()
}
