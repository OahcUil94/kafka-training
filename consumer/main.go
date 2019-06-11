package main

import (
	"fmt"
	"github.com/bsm/sarama-cluster"
	"log"
	"os"
	"os/signal"
)

var Address = []string{"127.0.0.1:9092"}
var Topic = []string{"test01"}

func main() {
	partitionConsume()
}

// 多路复用模式，多个分区，多个主题都通过一个channel获取数据
func multiplexedConsume() {
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true

	consumer, err := cluster.NewConsumer(Address, "group01", Topic, config)
	if err != nil {
		panic(err)
	}

	defer func () {
		if err := consumer.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	go func () {
		for err := range consumer.Errors() {
			log.Printf("Error: %s\n", err.Error())
		}
	}()

	go func () {
		for ntf := range consumer.Notifications() {
			log.Printf("Rebalanced: %+v\n", ntf)
		}
	}()

	for {
		select {
		case msg, ok := <-consumer.Messages():
			if ok {
				fmt.Fprintf(os.Stdout, "%s/%d/%d\t%s\t%s\n", msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
				consumer.MarkOffset(msg, "")	// mark message as processed
			}
		case <-signals:
			return
		}
	}
}

func partitionConsume() {
	config := cluster.NewConfig()
	config.Group.Mode = cluster.ConsumerModePartitions

	consumer, err := cluster.NewConsumer(Address, "group02", Topic, config)
	if err != nil {
		panic(err)
	}

	defer func () {
		if err := consumer.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	for {
		select {
		case part, ok := <- consumer.Partitions():
			if !ok {
				return
			}

			go func (pc cluster.PartitionConsumer) {
				for msg := range pc.Messages() {
					fmt.Fprintf(os.Stdout, "%s/%d/%d\t%s\t%s\n", msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
					consumer.MarkOffset(msg, "")	// mark message as processed
				}
			}(part)

		case <-signals:
			return
		}
	}
}