// Package appctx
package appctx

import (
	"fmt"
	"sync"

	"sharefood/internal/consts"
	"sharefood/pkg/file"
	"sharefood/pkg/logger"
)

var (
	once sync.Once
	_cfg *Config
)

// NewConfig initialize config object
func NewConfig() *Config {
	fpath := []string{consts.ConfigPath}
	once.Do(func() {
		c, err := readCfg("app.yaml", fpath...)
		if err != nil {
			logger.Fatal(err)
		}

		_cfg = c
	})

	return _cfg
}

// Config object contract
//
//go:generate easytags $GOFILE yaml,json
type Config struct {
	App     *Common      `yaml:"app" json:"app"`
	Logger  Logging      `yaml:"logger" json:"logger"`
	WriteDB *Database    `yaml:"db_write" json:"db_write"`
	ReadDB  *Database    `yaml:"db_read" json:"read_db"`
	Redis   *RedisConf   `yaml:"redis" json:"redis"`
	AWS     AWS          `yaml:"aws" json:"aws"`
	Kafka   *KafkaConfig `yaml:"kafka" json:"kafka"`
	APM     APM          `yaml:"apm" json:"apm"`
	Pubsub  PubSub       `yaml:"pubsub" json:"pubsub"`
	GCS     GCS          `yaml:"gcs" json:"gcs"`
}

// Common general config object contract
type Common struct {
	AppName            string `yaml:"name" json:"name"`
	ApiKey             string `yaml:"key" json:"api_key"`
	Debug              bool   `yaml:"debug" json:"debug"`
	Maintenance        bool   `yaml:"maintenance" json:"maintenance"`
	Timezone           string `yaml:"timezone" json:"timezone"`
	Env                string `yaml:"env" json:"env"`
	Port               int    `yaml:"port" json:"port"`
	ReadTimeoutSecond  int    `yaml:"read_timeout_second" json:"read_timeout_second"`
	WriteTimeoutSecond int    `yaml:"write_timeout_second" json:"write_timeout_second"`
	DefaultLang        string `yaml:"default_lang" json:"default_lang"`
	EncryptKey         string `yaml:"encrypt_key" json:"encrypt_key"`
	JWTSecret          string `yaml:"jwt_secret" json:"jwt_secret"`
}

// Database configuration structure
type Database struct {
	Name          string `yaml:"name" json:"name"`
	User          string `yaml:"user" json:"user"`
	Pass          string `yaml:"pass" json:"pass"`
	Host          string `yaml:"host" json:"host"`
	Port          int    `yaml:"port" json:"port"`
	MaxOpen       int    `yaml:"max_open" json:"max_open"`
	MaxIdle       int    `yaml:"max_idle" json:"max_idle"`
	TimeoutSecond int    `yaml:"timeout_second" json:"timeout_second"`
	MaxLifeTimeMS int    `yaml:"life_time_ms" json:"max_life_time_ms"`
	Charset       string `yaml:"charset" json:"charset"`
	Driver        string `yaml:"driver" json:"driver"`
	Timezone      string `yaml:"timezone" json:"timezone"`
}

// RedisConf general config redis
type RedisConf struct {
	Hosts              string `yaml:"host"`
	DB                 int    `yaml:"db"`
	ReadTimeoutSecond  int    `yaml:"read_timeout_second"`
	WriteTimeoutSecond int    `yaml:"write_timeout_second"`
	PoolSize           int    `yaml:"pool_size"`
	PoolTimeoutSecond  int    `yaml:"pool_timeout_second"`
	MinIdleConn        int    `yaml:"min_idle_conn"`
	IdleTimeoutSecond  int    `yaml:"idle_timeout_second"`
	RouteByLatency     bool   `yaml:"route_by_latency"`
	IdleFrequencyCheck int    `yaml:"idle_frequency_check"`
	Password           string `yaml:"password"`
	ReadOnly           bool   `yaml:"read_only"`
	RouteRandomly      bool   `yaml:"route_randomly"`
	MaxRedirect        int    `yaml:"max_redirect"`
	ClusterMode        bool   `yaml:"cluster_mode"`
	TLSEnable          bool   `yaml:"tls_enable"`
	InsecureSkipVerify bool   `yaml:"insecure_skip_verify"`
}

// AWS config object for aws account
type AWS struct {
	AccessKey    string `yaml:"access_key" json:"access_key"`
	AccessSecret string `yaml:"access_secret" json:"access_secret"`
	Profile      string `yaml:"profile" json:"profile"`
	Region       string `yaml:"region" json:"region"`
	StackName    string `yaml:"stack_name" json:"stack_name"`
}

// SQS config contract for aws sqs
type SQS struct {
	QueueName      string `yaml:"queue_name" json:"queue_name"`
	QueueURL       string `yaml:"queue_url" json:"queue_url"`
	MaxMessage     int    `yaml:"max_message" json:"max_message"`
	WaitTimeSecond int    `yaml:"wait_time_second" json:"wait_time_second"`
}

// readCfg reads the configuration from file
// args:
//
//	fname: filename
//	ps: full path of possible configuration files
//
// returns:
//
//	*config.Configuration: configuration ptr object
//	error: error operation
func readCfg(fname string, ps ...string) (*Config, error) {
	var cfg *Config
	var errs []error

	for _, p := range ps {
		f := fmt.Sprint(p, fname)

		err := file.ReadFromYAML(f, &cfg)
		if err != nil {
			errs = append(errs, fmt.Errorf("file %s error %s", f, err.Error()))
			continue
		}
		break
	}

	if cfg == nil {
		return nil, fmt.Errorf("file config parse error %v", errs)
	}

	return cfg, nil
}

// Client is a config contract for third party  service provider
type Client struct {
	URL           string `yaml:"url" json:"url"`
	ApiKey        string `yaml:"api_key" json:"api_key"`
	ApiSecret     string `yaml:"api_secret" json:"api_secret"`
	Version       string `yaml:"version" json:"version"`
	TimeoutSecond int    `yaml:"timeout_second" json:"timeout_second"`
	VendorID      int    `yaml:"vendor_id" json:"vendor_id"`
}

// Logging config
type Logging struct {
	Name  string `yaml:"name" json:"name"`
	Level string `yaml:"level" json:"level"`
}

// Config entity of kafka broker
type KafkaConfig struct {
	// Brokers list of brokers connection hostname or ip address
	Brokers string `yaml:"brokers" json:"brokers"`
	SASL    SASL   `yaml:"sasl" json:"sasl"`
	// kafka broker Version
	Version  string        `yaml:"version" json:"version"`
	ClientID string        `yaml:"client_id" json:"client_id"`
	Producer KafkaProducer `yaml:"producer" json:"producer"`
	Consumer KafkaConsumer `yaml:"consumer" json:"consumer"`
	TLS      TLS           `yaml:"tls" json:"tls"`
	// The number of events to buffer in internal and external channels. This
	// permits the producer and consumer to continue processing some messages
	// in the background while user code is working, greatly improving throughput.
	// Defaults to 256.
	ChannelBufferSize int `json:"channel_buffer_size" yaml:"channel_buffer_size"`
}

// KafkaProducer config
type KafkaProducer struct {
	// The maximum duration the broker will wait the receipt of the number of
	// RequiredAcks (defaults to 10 seconds). This is only relevant when
	// RequiredAcks is set to WaitForAll or a number > 1. Only supports
	// millisecond resolution, nanoseconds will be truncated. Equivalent to
	// the JVM producer's `request.timeout.ms` setting.
	TimeoutSecond int `yaml:"timeout_second" json:"timeout_second"`
	// RequireACK
	// 0 = NoResponse doesn't send any response, the TCP ACK is all you get.
	// 1 =  WaitForLocal waits for only the local commit to succeed before responding.
	// - 1 =  WaitForAll
	// WaitForAll waits for all in-sync replicas to commit before responding.
	// The minimum number of in-sync replicas is configured on the broker via
	// the `min.insync.replicas` configuration key.
	RequireACK int16 `yaml:"ack" json:"require_ack"`
	// If enabled, the producer will ensure that exactly one copy of each message is
	// written.
	IdemPotent bool `yaml:"idem_potent" json:"idem_potent"`

	// Generates partitioners for choosing the partition to send messages to
	// (defaults to hashing the message key). Similar to the `partitioner.class`
	// setting for the JVM producer.
	PartitionStrategy string `yaml:"partition_strategy" json:"partition_strategy"`
}

// KafkaConsumer config
type KafkaConsumer struct {
	// Minimum is 10s
	SessionTimeoutSecond int    `yaml:"session_timeout_second" json:"session_timeout_second"`
	OffsetInitial        int64  `yaml:"offset_initial" json:"offset_initial"`
	HeartbeatIntervalMS  int    `yaml:"heartbeat_interval_ms" json:"heartbeat_interval_ms"`
	RebalanceStrategy    string `yaml:"rebalance_strategy" json:"rebalance_strategy"`
	AutoCommit           bool   `yaml:"auto_commit" json:"auto_commit"`
	IsolationLevel       int8   `json:"isolation_level" yaml:"isolation_level"`
}

// SALS secure connection config
type SASL struct {
	// Whether or not to use SASL authentication when connecting to the broker
	// (defaults to false).
	Enable bool `yaml:"enable" json:"enable"`
	// SASLMechanism is the name of the enabled SASL mechanism.
	// Possible values: OAUTHBEARER, PLAIN (defaults to PLAIN).
	Mechanism string `yaml:"mechanism" json:"mechanism"`
	// Version is the SASL Protocol Version to use
	// Kafka > 1.x should use V1, except on Azure EventHub which use V0
	Version int16 `yaml:"version" json:"version"`
	// Whether or not to send the Kafka SASL handshake first if enabled
	// (defaults to true). You should only set this to false if you're using
	// a non-Kafka SASL proxy.
	Handshake bool `yaml:"handshake" json:"handshake"`
	// User is the authentication identity (authcid) to present for
	// SASL/PLAIN or SASL/SCRAM authentication
	User string `yaml:"user" json:"user"`
	// Password for SASL/PLAIN authentication
	Password string `yaml:"password" json:"password"`
}

// TLS config
type TLS struct {
	Enable     bool   `yaml:"enable" json:"enable"`
	CaFile     string `yaml:"ca_file" json:"ca_file"`
	KeyFile    string `yaml:"key_file" json:"key_file"`
	CertFile   string `yaml:"cert_file" json:"cert_file"`
	SkipVerify bool   `yaml:"skip_verify" json:"skip_verify"`
}

// APM config
type APM struct {
	Address string `yaml:"address" json:"address"`
	Enable  bool   `yaml:"enable" json:"enable"`
	Name    string `yaml:"name" json:"name"`
}

type (
	PubSub struct {
		AccountPath string `yaml:"account_path" json:"account_path"`
		ProjectID   string `yaml:"project_id" json:"project_id"`
		Topic       string `yaml:"topic" json:"topic"`
	}

	GCS struct {
		AccountPath string `yaml:"account_path" json:"account_path"`
		Bucket      string `yaml:"bucket" json:"bucket"`
		Prefix      string `yaml:"prefix"`
	}
)
