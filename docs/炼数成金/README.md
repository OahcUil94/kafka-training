# 炼数成金学习目录

1. 

## 1. kafka简介

查看zookeeper的2181端口是否被监听：lsof -i:2181
查看kafka的9092端口是否被监听：lsof -i:9092

查看kafka的进程是否被停掉：`ps -ef|grep -i kafka`
查看zookeeper的进程是否被停掉：`ps -ef|grep -i zookeeper`

根据zookeeper的dockerfile，构建一个镜像：

`docker build -t oahcuil94/zookeeper:3.4.6 -f zookeeper.Dockerfile .`

`docker run -itd --name zookeeper -h zookeeper -p 2181:2181 oahcuil94/zookeeper:3.4.6 bash`

根据kafka的dockerfile，构建一个镜像：

`docker build -t oahcuil94/kafka:0.8.2.2 -f kafka.Dockerfile .`

`docker run -itd --name kafka -h kafka -p 9092:9092 --link zookeeper oahcuil94/kafka:0.8.2.2 bash`

将一个topic中的数据，同步到另一个topic中：

`bin/kafka-replay-log-producer.sh --broker-list localhost:9092 --zookeeper zookeeper:2181 --inputtopic test1 --outputtopic test2`

## 2. kafka架构

### Partition

- 物理概念，一个Partition只分布于一个Broker上（不考虑备份）
- 一个Partition物理上对应一个文件夹
- 一个Partition包含多个Segment（Segment对用户透明）
- 一个Segment对应一个文件
- Segment由一个个不可变记录组成
- 记录只会被append到Segment中，不会单独删除或者修改
- 清除过期日志时，直接删除一个或多个Segment

清除日志，可以设定时间，默认清除168小时的日志，也可以设定大小，比如达到了1GB（1024M）就可以删除了，由于Segment也是可以设定大小的，假设一个Segment大小是100MB，那在创建了第11个Segment之后，Kafka的后台进程会进行整个扫描，然后删除第一个Segment，删除第一个Segment之后，大小又小于了1GB了，又可以继续写数据了，当然扫描也是有时间的，假如这次扫描的时候，已经创建了15个Segment了，那就会一次性删除5个Segment。

