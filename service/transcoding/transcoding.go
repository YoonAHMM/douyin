package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"

	"douyin/service/transcoding/internal/config"
	"douyin/service/transcoding/internal/logic"
	"douyin/service/transcoding/internal/svc"
	"github.com/segmentio/kafka-go"
)

var configFile = flag.String("f", "etc/transcoding.yaml", "the config file")

func NewKafkaReader(kafkaURL, topic, groupID string, minBytes, maxBytes int) *kafka.Reader {
	brokers := strings.Split(kafkaURL, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: minBytes,
		MaxBytes: maxBytes,
	})
}

func main() {
	var c config.Config
	config.MustLoad(*configFile, &c)

	reader := NewKafkaReader(c.KafkaConfig.Host, c.KafkaConfig.Topic, c.KafkaConfig.GroupId, c.KafkaConfig.MinBytes, c.KafkaConfig.MaxBytes)
	defer reader.Close()

	svcctx := svc.NewServiceContext(c)
	l := logic.NewTranscodingLogic(context.Background(), svcctx)
	
	fmt.Println("TransCoding Service Start...")
	fmt.Println("start consuming ...")
	err := reader.SetOffset(kafka.LastOffset)
	if err != nil {
		return
	}

	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("message at topic:%v partition:%v offset:%v	%s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
		err = l.TransCoding(string(m.Key), m.Value)
		if err != nil {
			log.Println(err)
		}

		fmt.Println("TransCoding completed...")
	}
}
