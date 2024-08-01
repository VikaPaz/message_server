package queue

import (
	"context"
	"github.com/segmentio/kafka-go"
	"time"
)

type Config struct {
	Topic     string
	Partition int
	Host      string
	Network   string
}

// to produce messages
func Connection(conf Config) (*kafka.Conn, error) {
	dialer := &kafka.Dialer{
		KeepAlive: 10 * time.Second,
	}
	conn, err := dialer.DialLeader(context.Background(), conf.Network, conf.Host, conf.Topic, conf.Partition)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
