package queue

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
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

	// ping kafka to not lose connection
	go func() {
		t := time.NewTicker(5 * time.Minute)
		for _ = range t.C {
			_, err := conn.ReadPartitions()
			if err != nil {
				t.Stop()
				log.Fatalf("can't ping kafka: %s", err.Error())
			}
		}
	}()
	return conn, nil
}
