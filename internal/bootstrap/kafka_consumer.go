// Package bootstrap
package bootstrap

import (
	"strings"

	"sharefood/internal/appctx"
	"sharefood/pkg/kafka"
)

func RegistryKafkaConsumer(cfg *appctx.Config) kafka.Consumer {
	return kafka.NewConsumerGroup(&kafka.Config{
		Consumer: kafka.ConsumerConfig{
			SessionTimeoutSecond: cfg.Kafka.Consumer.SessionTimeoutSecond,
			HeartbeatInterval:    cfg.Kafka.Consumer.HeartbeatIntervalMS,
			RebalanceStrategy:    cfg.Kafka.Consumer.RebalanceStrategy,
			OffsetInitial:        cfg.Kafka.Consumer.OffsetInitial,
		},
		Version:  cfg.Kafka.Version,
		Brokers:  strings.Split(cfg.Kafka.Brokers, ","),
		ClientID: cfg.Kafka.ClientID,
		SASL: kafka.SASL{
			Enable:    cfg.Kafka.SASL.Enable,
			User:      cfg.Kafka.SASL.User,
			Password:  cfg.Kafka.SASL.Password,
			Mechanism: cfg.Kafka.SASL.Mechanism,
			Version:   cfg.Kafka.SASL.Version,
			Handshake: cfg.Kafka.SASL.Handshake,
		},
		TLS: kafka.TLS{
			Enable:     cfg.Kafka.TLS.Enable,
			CaFile:     cfg.Kafka.TLS.CaFile,
			CertFile:   cfg.Kafka.TLS.CertFile,
			KeyFile:    cfg.Kafka.TLS.KeyFile,
			SkipVerify: cfg.Kafka.TLS.SkipVerify,
		},

		ChannelBufferSize: cfg.Kafka.ChannelBufferSize,
	})
}
