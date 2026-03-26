package main

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
)

func main() {
	fmt.Println("Hello Go")
	producer := NewKafkaProducer()
	Publish("transferiu", "teste", producer, nil)
	// Você precisa de Flush porque o Producer do Kafka em Go é assíncrono: ele coloca as mensagens em um buffer interno e envia em background. Flush força o envio de todas as mensagens pendentes para o broker antes do programa terminar, garantindo que não se percam.
	producer.Flush(1000)
}

func NewKafkaProducer() *kafka.Producer {
	configMap := &kafka.ConfigMap{
		// Você só precisa passar um ou poucos brokers em "bootstrap.servers"; o client conecta nesses, descobre os outros no cluster automaticamente e usa todos conforme necessário.
		// usar docker ps para descobrir o nome do service broker
		// docker por default gera o nome, baseado no nome do projeto
		"bootstrap.servers":   "fc3-apache-kafka-kafka-1:9092",
	}
	p, err := kafka.NewProducer(configMap)
	if err != nil {
		log.Println(err.Error())
	}
	return p
}

func Publish(msg string, topic string, producer *kafka.Producer, key []byte) error {
	message := &kafka.Message{
		// []byte - converter de string para slice de bytes
		// Array: tamanho fixo, faz parte do tipo ([3]int).
		// Slice: tamanho dinâmico, referência para um array, tipo []int.
		// eu forcei string, mas é um slice de bytes no value, poderia passar um arquivo aqui
		Value:          []byte(msg),
		// enviando pra qualquer partição
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            key,
	}
	err := producer.Produce(message, nil)
	if err != nil {
		return err
	}
	return nil
}