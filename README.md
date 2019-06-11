# Kafka Training

时间：2019-06-11

版本信息：
- mac osx: 10.13.3
- go: 1.11.5 darwin/amd64
- jdk: openjdk 11.0.1
- kafka: 2.1.1
- zookeeper: 3.4.13

## Kafka安装

### mac osx下安装

```bash
$ brew install kafka
==> Installing dependencies for kafka: zookeeper
==> Installing kafka dependency: zookeeper
==> Downloading https://homebrew.bintray.com/bottles/zookeeper-3.4.13.high_sierra.bottle.tar.gz
Already downloaded: /Users/xxxx/Library/Caches/Homebrew/downloads/a890511ccb1774e42433d4bbe67e8505b40e29d792a55faa8f40e172c0a4d8f7--zookeeper-3.4.13.high_sierra.bottle.tar.gz
==> Pouring zookeeper-3.4.13.high_sierra.bottle.tar.gz
==> Caveats
To have launchd start zookeeper now and restart at login:
  brew services start zookeeper
Or, if you don't want/need a background service you can just run:
  zkServer start
==> Summary
🍺  /usr/local/Cellar/zookeeper/3.4.13: 244 files, 33.4MB
==> Installing kafka
==> Downloading https://homebrew.bintray.com/bottles/kafka-2.1.1.high_sierra.bottle.tar.gz
Already downloaded: /Users/xxxx/Library/Caches/Homebrew/downloads/b7faa699991c953d8f5cc421ddc7336b3bc25a4374dc0c73530985f24616df22--kafka-2.1.1.high_sierra.bottle.tar.gz
==> Pouring kafka-2.1.1.high_sierra.bottle.tar.gz
==> Caveats
To have launchd start kafka now and restart at login:
  brew services start kafka
Or, if you don't want/need a background service you can just run:
  zookeeper-server-start /usr/local/etc/kafka/zookeeper.properties & kafka-server-start /usr/local/etc/kafka/server.properties
==> Summary
🍺  /usr/local/Cellar/kafka/2.1.1: 162 files, 52.7MB
==> Caveats
==> zookeeper
To have launchd start zookeeper now and restart at login:
  brew services start zookeeper
Or, if you don't want/need a background service you can just run:
  zkServer start
==> kafka
To have launchd start kafka now and restart at login:
  brew services start kafka
Or, if you don't want/need a background service you can just run:
  zookeeper-server-start /usr/local/etc/kafka/zookeeper.properties & kafka-server-start /usr/local/etc/kafka/server.properties
```

## golang相关的包

1. github.com/Shopify/sarama
2. github.com/bsm/sarama-cluster

## kafka相关工具

### bin/kafka-topics

操作topic相关的命令，该命令所必须的参数有：

1. --list, --describe, --create, --alter, --delete这几个任意一个动作参数
2. --zookeeper host，zookeeper的配置信息

#### 创建主题

`bin/kafka-topics --create`

参数列表：

1. --partitions: 指定topic的分区数
2. --replication-factor: 指定副本数，注意，副本数一定要小于等于broker数量，否则报错
3. --topic topicname: 指定topic名字

例如：`bin/kafka-topics --create --zookeeper 127.0.0.1 --partitions 1 --replication-factor 1 --topic test_logs`

#### 查看已创建的主题列表

`bin/kafka-topics --list`

例如：`bin/kafka-topics --list --zookeeper 127.0.0.1`

#### 查看主题的详情

`bin/kafka-topics --describe --zookeeper 127.0.0.1 --topic test_logs`

```bash
Topic:test_logs	PartitionCount:1	ReplicationFactor:1	Configs:
	Topic: test_logs	Partition: 0	Leader: 0	Replicas: 0	Isr: 0
```

通过上面信息，可以看到`test_logs`这个主题的`Partition`在0号broker上，`Leader`也在0号broker上，`Replicas`副本也是在0号broker上的，`Isr`是选举要用到的，Leader和副本同步的策略就是，生产者只给Leader写数据，副本自动同步，当Leader宕掉之后，就看剩下的副本，谁同步的数据多，谁就被选举为下一任的Leader，`Isr`可能有多个值，例如：`Isr: 0, 2`，`Isr`值的顺序就是已经排好了，谁的数据多，谁在后面，假设此时`Leader`挂掉了，根据`Isr`提前排好的顺序，2号broker立刻就成为了`Leader`。

#### 删除已创建的主题

`bin/kafka-topics --delete`

例如：`bin/kafka-topics --delete --zookeeper 127.0.0.1 --topic test_logs`

删除之后可能会出现这样的提示：

```
## 如果delete.topic.enable不是true，则topic只是被标记为删除，并不会真实的删除 

Topic test_logs is marked for deletion.
Note: This will have no impact if delete.topic.enable is not set to true.
```

### bin/kafka-console-producer

生产者相关的命令，该命令所需要的参数有：

-- broker-list host: 指定kafka节点的主机
-- topic topicname: 指定要给哪个topic写数据

例如：`bin/kafka-console-producer --broker-list 127.0.0.1:9092 --topic test_logs`

### bin/kafka-console-consumer

消费者相关的命令，该命令需要的参数有：

- --zookeeper host: 指定zookeeper的主机地址（已废弃），高版本的kafka将offset维护在了本地，不交给zookeeper维护了，把offset维护在了leader中了，提高了效率
- --bootstrap-server host: 指定kafka的集群 
- --topic topicname: 指定要消费的topic
- --from-beginning: 默认消费者读取的是最新的数据，最大的offset数据，如果想让消费者从最早的数据进行消费，可以指定该参数

`bin/kafka-console-consumer --bootstrap-server 127.0.0.1:9092 --topic test_logs --from-beginning`

上面提到了，之前数据的标识offset维护在zk中，新版本直接放到了kafka中，那offset放在哪个地方，是在`__consumer_offsets`的topic里。

## kafka数据文件存放位置

- 默认数据存放位置：/usr/local/var/lib/kafka-logs
- 默认日志存放位置：/usr/local/var/log/kafka/kafka_output.log

## kafka其他的杂乱的知识点

1. kafka的broker id是整数类型
2. kafka指定topic的副本的时候，一定要小于等于broker数量

## docker-compose安装zk和kafka集群
