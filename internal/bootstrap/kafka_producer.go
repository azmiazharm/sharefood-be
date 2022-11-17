// Package bootstrap
package bootstrap

import (
	"strings"

	"sharefood/internal/appctx"
	"sharefood/pkg/kafka"
)

func RegistryKafkaProducer(cfg *appctx.Config) kafka.Producer {
	return kafka.NewProducer(&kafka.Config{
		Producer: kafka.ProducerConfig{
			TimeoutSecond:     cfg.Kafka.Producer.TimeoutSecond,
			RequireACK:        cfg.Kafka.Producer.RequireACK,
			IdemPotent:        cfg.Kafka.Producer.IdemPotent,
			PartitionStrategy: cfg.Kafka.Producer.PartitionStrategy,
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
