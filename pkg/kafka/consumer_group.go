// Package kafka broker message
// @author Daud Valentino
package kafka

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"sharefood/pkg/logger"

	"github.com/Shopify/sarama"
)

type consumerGroup struct {
	config     *sarama.Config
	brokers    []string
	autoCommit bool
}

// NewConsumer return consumer message broker
func NewConsumerGroup(cfg *Config) Consumer {
	m := &consumerGroup{}

	/**
	 * Construct a new Sarama configuration.
	 * The Kafka cluster version has to be defined before the consumer/producer is initialized.
	 */
	config := sarama.NewConfig()

	if cfg.Version == "" {
		cfg.Version = defaultVersion
	}

	version, err := sarama.ParseKafkaVersion(cfg.Version)
	if err != nil {
		logger.Fatal(fmt.Sprintf("parse kafka version got: %v", err))
	}

	if cfg.SASL.Enable {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = cfg.SASL.User
		config.Net.SASL.Password = cfg.SASL.Password
		config.Net.SASL.Version = sarama.SASLHandshakeV0
		config.Net.SASL.Handshake = true
		config.Net.SASL.Mechanism = sarama.SASLMechanism(cfg.SASL.Mechanism)
		config.Net.TLS.Enable = true
	}

	// The TLS configuration to use for secure connections if
	// enabled (defaults to nil).
	if config.Net.TLS.Enable || cfg.TLS.Enable {
		config.Net.TLS.Config = createTlsConfig(cfg.TLS)
	}

	config.Version = version

	config.Consumer.Offsets.Initial = cfg.Consumer.OffsetInitial
	config.Consumer.Return.Errors = true
	config.Consumer.Group.Session.Timeout = time.Duration(cfg.Consumer.SessionTimeoutSecond) * time.Second
	config.Consumer.Group.Heartbeat.Interval = time.Duration(cfg.Consumer.HeartbeatInterval) * time.Millisecond

	if len(strings.Trim(cfg.Consumer.RebalanceStrategy, " ")) == 0 {
		cfg.Consumer.RebalanceStrategy = sarama.RangeBalanceStrategyName
	}

	st, ok := balanceStrategies[cfg.Consumer.RebalanceStrategy]

	if !ok {
		logger.Fatal(fmt.Sprintf(
			`rebalance strateggy only available : "%s", "%s", "%s",   on setting value : "%s"`,
			sarama.RoundRobinBalanceStrategyName,
			sarama.RangeBalanceStrategyName,
			sarama.StickyBalanceStrategyName,
			cfg.Consumer.RebalanceStrategy,
		))
	}

	if cfg.ChannelBufferSize > 0 {
		config.ChannelBufferSize = cfg.ChannelBufferSize
	}

	config.Consumer.IsolationLevel = sarama.IsolationLevel(cfg.Consumer.IsolationLevel)

	config.Consumer.Group.Rebalance.Strategy = st

	m.brokers = cfg.Brokers
	m.config = config
	m.autoCommit = cfg.Consumer.AutoCommit
	return m
}

// Subscribe message
func (k *consumerGroup) Subscribe(ctx *ConsumerContext) {

	fields := []logger.Field{
		logger.Any("topics", ctx.Topics),
	}

	client, err := sarama.NewConsumerGroup(k.brokers, ctx.GroupID, k.config)

	if err != nil {
		logger.Fatal(err.Error(), fields...)
	}

	handler := NewConsumerHandler(ctx.Handler, k.autoCommit)

	// kafka consumer client
	nCtx, cancel := context.WithCancel(ctx.Context)

	defer func() {
		_ = client.Close()
	}()

	// subscriber errors
	go func() {
		for err := range client.Errors() {
			logger.Error(fmt.Sprintf("[consumer] error %s", err.Error()), fields...)
		}
	}()

	go func() {
		for {
			select {
			case <-nCtx.Done():
				logger.Warn(logger.MessageFormat("[consumer] stopped consume topics %v", ctx.Topics))
				return
			default:
				err := client.Consume(nCtx, ctx.Topics, handler)
				if err != nil {
					logger.Error(logger.MessageFormat("[consumer] topic %v consume message error %s", ctx.Topics, err.Error()))
				}
			}
		}
	}()

	logger.Info(fmt.Sprintf("[consumer] sarama consumer up and running!... group %s, queue %v", ctx.Context, ctx.Topics), fields...)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	<-sigterm // Await a sigterm signal before safely closing the consumer

	cancel()
	logger.Info("[consumer] Cancelled message without marking offsets", fields...)

}
