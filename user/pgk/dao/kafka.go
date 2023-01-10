package dao

import (
	"fmt"
	"github.com/Shopify/sarama"
	"go.uber.org/zap"
	"sync"
)

type Client struct {
	sarama.Client
	producer sarama.SyncProducer
	config   *sarama.Config
	addrs    []string
	topic    string
	group    string
}

func Consumer(addrs []string, topic string) {
	var wg sync.WaitGroup
	consumer, err := sarama.NewConsumer(addrs, nil)
	if err != nil {
		fmt.Println("Failed to start consumer: %s", err)
		return
	}
	partitionList, err := consumer.Partitions("task-status-data") // 通过topic获取到所有的分区
	if err != nil {
		fmt.Println("Failed to get the list of partition: ", err)
		return
	}
	fmt.Println(partitionList)

	for partition := range partitionList { // 遍历所有的分区
		pc, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest) // 针对每个分区创建一个分区消费者
		if err != nil {
			fmt.Println("Failed to start consumer for partition %d: %s\n", partition, err)
		}
		wg.Add(1)
		go func(sarama.PartitionConsumer) { // 为每个分区开一个go协程取值
			for msg := range pc.Messages() { // 阻塞直到有值发送过来，然后再继续等待
				fmt.Printf("Partition:%d, Offset:%d, key:%s, value:%s\n", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
			}
			defer pc.AsyncClose()
			wg.Done()
		}(pc)
	}
	wg.Wait()
	consumer.Close()
}
func NewClient(addrs []string, topic, group string) (c *Client, err error) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_1_0_0
	config.Producer.Idempotent = true
	config.Producer.Return.Errors = true
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	config.Producer.Transaction.Retry.Backoff = 10

	config.Net.MaxOpenRequests = 1
	producer, err := sarama.NewSyncProducer([]string{"192.168.3.116:9092", "192.168.3.116:9092", "192.168.3.116:9092"}, config)
	if err != nil {
		fmt.Println("create producer failed, ", err.Error())
		return nil, err
	}
	c = &Client{
		config:   config,
		addrs:    addrs,
		topic:    topic,
		group:    group,
		producer: producer,
	}
	c.Client, err = sarama.NewClient(addrs, config)
	return
}
func (c *Client) SendSync(v string) error {
	msg := &sarama.ProducerMessage{}
	msg.Topic = c.topic

	msg.Value = sarama.StringEncoder(v)

	defer c.Close()

	// 发送消息
	pid, offset, err := c.producer.SendMessage(msg)
	if err != nil {
		zap.L().Info("send msg failed, err:", zap.Error(err))
		return err
	}
	fmt.Printf("pid:%v offset:%v\n", pid, offset)
	return nil
}
func (c *Client) Close() {
	c.Client.Close()
}

func (c *Client) Topics() []string {
	topics, err := c.Client.Topics()
	if err != nil {
		panic(err)
	}
	return topics
}

func (c *Client) Partitions() []int32 {
	ps, err := c.Client.Partitions(c.topic)
	if err != nil {
		panic(err)
	}
	return ps
}

func (c *Client) OffsetNew() (info map[int32]int64, err error) {
	var (
		offset int64
	)
	ps, err := c.Client.Partitions(c.topic)
	if err != nil {
		return
	}
	info = make(map[int32]int64)
	for _, p := range ps {
		offset, err = c.Client.GetOffset(c.topic, p, sarama.OffsetNewest)
		if err != nil {
			return
		}
		info[p] = offset
	}
	return
}

func (c *Client) OffsetOld() (info map[int32]int64, err error) {
	var (
		offset int64
	)
	ps, err := c.Client.Partitions(c.topic)
	if err != nil {
		return
	}
	info = make(map[int32]int64)
	for _, p := range ps {
		offset, err = c.Client.GetOffset(c.topic, p, sarama.OffsetOldest)
		if err != nil {
			return
		}
		info[p] = offset
	}
	return
}

// pool of producers that ensure transactional-id is unique.
type producerProvider struct {
	transactionIdGenerator int32

	producersLock sync.Mutex
	producers     []sarama.AsyncProducer

	producerProvider func() sarama.AsyncProducer
}

func newProducerProvider(brokers []string, producerConfigurationProvider func() *sarama.Config) *producerProvider {
	provider := &producerProvider{}
	provider.producerProvider = func() sarama.AsyncProducer {
		config := producerConfigurationProvider()
		suffix := provider.transactionIdGenerator
		// Append transactionIdGenerator to current config.Producer.Transaction.ID to ensure transaction-id uniqueness.
		if config.Producer.Transaction.ID != "" {
			provider.transactionIdGenerator++
			config.Producer.Transaction.ID = config.Producer.Transaction.ID + "-" + fmt.Sprint(suffix)
		}
		producer, err := sarama.NewAsyncProducer(brokers, config)
		if err != nil {
			return nil
		}
		return producer
	}
	return provider
}

func (p *producerProvider) borrow() (producer sarama.AsyncProducer) {
	p.producersLock.Lock()
	defer p.producersLock.Unlock()

	if len(p.producers) == 0 {
		for {
			producer = p.producerProvider()
			if producer != nil {
				return
			}
		}
	}

	index := len(p.producers) - 1
	producer = p.producers[index]
	p.producers = p.producers[:index]
	return
}

func (p *producerProvider) release(producer sarama.AsyncProducer) {
	p.producersLock.Lock()
	defer p.producersLock.Unlock()

	// If released producer is erroneous close it and don't return it to the producer pool.
	if producer.TxnStatus()&sarama.ProducerTxnFlagInError != 0 {
		// Try to close it
		_ = producer.Close()
		return
	}
	p.producers = append(p.producers, producer)
}

func (p *producerProvider) clear() {
	p.producersLock.Lock()
	defer p.producersLock.Unlock()

	for _, producer := range p.producers {
		producer.Close()
	}
	p.producers = p.producers[:0]
}
