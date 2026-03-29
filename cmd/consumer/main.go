package main

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	configMap := &kafka.ConfigMap{
		"bootstrap.servers":   "fc3-apache-kafka-kafka-1:9092",
		"client.id":         "goapp-consumer",
		"group.id":          "goapp-group",
		// começa a ler do inicio
		// em grupo tem comportamento diferente
		"auto.offset.reset": "earliest",
	}
	c, err := kafka.NewConsumer(configMap)
	if err != nil {
		fmt.Println("error consumer", err.Error())
	}
	topics := []string{"teste"}
	// SubscribeTopics espera receber um slice de strings com os nomes dos tópicos (no caso, topics) e um ponteiro para função callback opcional (nil ou uma função para lidar com rebalanceamento de partições).
	c.SubscribeTopics(topics, nil)
	for {
		// -1 = espera para sempre ate a mensagem chegar, processa, depois vai para a proxima
		msg, err := c.ReadMessage(-1)
		if err == nil {
			fmt.Println(string(msg.Value), msg.TopicPartition)
		}
	}
}