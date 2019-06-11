package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"time"
)

var Address = []string{"127.0.0.1:9092"}

func main() {
	asyncProduce()
}

func syncProduce() {
	config := sarama.NewConfig()
	// 设置使用的kafka版本,如果低于V0_10_0_0版本,消息中的timestrap没有作用.需要消费和生产同时配置
	// 注意，版本设置不对的话，kafka会返回很奇怪的错误，并且无法成功发送消息
	config.Version = sarama.V2_1_0_0

	// 等待所有副本都保存成功后再响应
	config.Producer.RequiredAcks = sarama.WaitForAll

	// 随机向Partition发送消息
	config.Producer.Partitioner = sarama.NewRandomPartitioner

	//是否等待成功和失败后的响应,只有上面的RequireAcks设置不是NoReponse这里才有用.
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	// 超时时间
	config.Producer.Timeout = 5 * time.Second

	// 创建一个同步的生产者
	p, err := sarama.NewSyncProducer(Address, config)
	if err != nil {
		log.Printf("sarama.NewSyncProducer err, message=%s \n", err)
		return
	}

	defer func () {
		if err := p.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	topic := "test01"
	srcValue := "sync: this is a message. index = %d"
	for i := 0; i < 10; i++ {
		value := fmt.Sprintf(srcValue, i)

		// 组建要写入的数据
		msg := &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.StringEncoder(value),
			Timestamp: time.Now(),
		}

		// 生产者往kafka写入数据
		part, offset, err := p.SendMessage(msg)
		if err != nil {
			log.Printf("send message(%s) err=%s \n", value, err)
		} else {
			fmt.Printf(value + "发送成功，partition=%d, offset=%d \n", part, offset)
		}

		time.Sleep(2 * time.Second)
	}
}

func asyncProduce() {
	config := sarama.NewConfig()
	config.Version = sarama.V2_1_0_0

	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	p, e := sarama.NewAsyncProducer(Address, config)
	if e != nil {
		log.Println(e)
		return
	}

	defer p.AsyncClose()

	go func (p sarama.AsyncProducer) {
		for {
			select {
			case msg := <-p.Successes():	// 异步处理指的是返回的信息是异步的
				fmt.Println("success: ", msg)
			case fail := <-p.Errors():
				fmt.Println("err: ", fail.Err)
			}
		}
	}(p)

	var value string
	for i := 0; ; i++ {
		time.Sleep(500 * time.Millisecond)
		value = "this is a message hello " + time.Now().Format("15:04:05")

		msg := &sarama.ProducerMessage{
			Topic: "test01",
			Timestamp: time.Now(),
		}
		msg.Value = sarama.StringEncoder(value)
		p.Input() <- msg
	}
}
