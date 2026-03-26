package main

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
)

func main() {
	// Sim, um canal é uma região de memória usada para comunicação entre goroutines. Goroutines podem enviar e receber dados por ele; elas podem "ficar esperando" por dados no canal, permitindo sincronização.
	deliveryChan := make(chan kafka.Event)
	producer := NewKafkaProducer()
	Publish("transferiu", "teste", producer, nil,deliveryChan)

	// # trabalhando de forma sincrona
	
	// Sim, ao ler o deliveryChan, você acessa a região de memória compartilhada e pega o valor enviado por outra goroutine (no caso, o status da entrega do Kafka).
	// Ao ler de um channel (<-deliveryChan), você retira o valor da fila do channel (não é uma cópia; ele sai do channel). Para só ver sem retirar, não é possível com channels em Go.
	// Você só lê um valor do canal com <-deliveryChan; para ler todas as mensagens, precisa chamar várias vezes (por exemplo, em um loop).
	// pegar o resultado do que aconteceu com o envio da mensagem
	// receber o resultado do channel
	e := <-deliveryChan
	// type assertion
	msg := e.(*kafka.Message)
	// schema message
	// Dentro de msg (kafka.Message) há os campos: TopicPartition, Value, Key, Timestamp, Headers e outros metadados da mensagem.
	if msg.TopicPartition.Error != nil {
		fmt.Println("Erro ao enviar")
		} else {
			fmt.Println("Mensagem enviada:", msg.TopicPartition)
			//schema da msg.TopicPartition:
			// error, partition, topic, metadata, offset, string
			// result: Mensagem enviada: teste[0]@1
			// 0: partição, 1: offset
	}
	
	// Você precisa de Flush porque o Producer do Kafka em Go é assíncrono: ele coloca as mensagens em um buffer interno e envia em background. 
	// Ele força o envio de todas as mensagens pendentes, bloqueando a execução até que todas sejam enviadas ou até expirar o tempo (em ms) passado como parâmetro. Serve para garantir que as mensagens não fiquem no buffer antes de encerrar o programa.
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

func Publish(msg string, topic string, producer *kafka.Producer, key []byte, deliveryChan chan kafka.Event) error {
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
	// O Kafka Producer roda internamente goroutines para envio e confirmação, e você usa o channel apenas para receber os resultados dessas operações assíncronas.
	err := producer.Produce(message, deliveryChan)
	if err != nil {
		return err
	}
	return nil
}