package main

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
)

func main() {
	// canal usado pelo producer e pela goroutine delivery report
	deliveryChan := make(chan kafka.Event)
	producer := NewKafkaProducer()
	Publish("transferiu", "teste", producer, nil,deliveryChan)
	//DeliveryReport(deliveryChan) // prossamento sincrono
	
	// com esse código só há uma goroutine lendo do canal, processando todas as mensagens sequencialmente conforme chegam.
	// O programa só termina após main finalizar. Se não houver código bloqueante depois do go DeliveryReport(deliveryChan), ele pode encerrar antes de processar todas as mensagens.

	// ao usar keyword "go" digo que é processamento assincrono e o programa continua
	// async que depende do processo principal estiver ativo
	go DeliveryReport(deliveryChan) // async
	fmt.Println("executei antes do código assincrono, mesmo sendo declarado depois")

	// # processando delivery report de forma sincrona
	
	// Sim, ao ler o deliveryChan, você acessa a região de memória compartilhada e pega o valor enviado por outra goroutine (no caso, o status da entrega do Kafka).
	// Ao ler de um channel (<-deliveryChan), você retira o valor da fila do channel (não é uma cópia; ele sai do channel). Para só ver sem retirar, não é possível com channels em Go.
	// Você só lê um valor do canal com <-deliveryChan; para ler todas as mensagens, precisa chamar várias vezes (por exemplo, em um loop).
	// pegar o resultado do que aconteceu com o envio da mensagem
	// receber o resultado do channel
	// e := <-deliveryChan
	// type assertion
	// msg := e.(*kafka.Message)
	// schema message
	// Dentro de msg (kafka.Message) há os campos: TopicPartition, Value, Key, Timestamp, Headers e outros metadados da mensagem.
	// if msg.TopicPartition.Error != nil {
	// 	fmt.Println("Erro ao enviar")
	// 	} else {
	// 		fmt.Println("Mensagem enviada:", msg.TopicPartition)
	// 		//schema da msg.TopicPartition:
	// 		// error, partition, topic, metadata, offset, string
	// 		// result: Mensagem enviada: teste[0]@1
	// 		// 0: partição, 1: offset
	// }
	
	// Você precisa de Flush porque o Producer do Kafka em Go é assíncrono: ele coloca as mensagens em um buffer interno e envia em background. 
	// Ele força o envio de todas as mensagens pendentes, bloqueando a execução até que todas sejam enviadas ou até expirar o tempo (em ms) passado como parâmetro. Serve para garantir que as mensagens não fiquem no buffer antes de encerrar o programa.
	// Não é obrigado usar Flush, mas é recomendado ao encerrar. O Kafka Producer em Go envia mensagens automaticamente para o canal, porém mantém um buffer interno. Se o programa termina antes desse buffer ser enviado, você pode perder mensagens. O Flush garante que tudo é processado antes de sair.
	producer.Flush(1000)
}

func NewKafkaProducer() *kafka.Producer {
	configMap := &kafka.ConfigMap{
		// Você só precisa passar um ou poucos brokers em "bootstrap.servers"; o client conecta nesses, descobre os outros no cluster automaticamente e usa todos conforme necessário.
		// usar docker ps para descobrir o nome do service broker
		// docker por default gera o nome, baseado no nome do projeto
			"bootstrap.servers":   "fc3-apache-kafka-kafka-1:9092",
			// "delivery.timeout.ms": "0" desativa o timeout de entrega do Kafka Producer, ou seja, não haverá limite de tempo para tentar entregar uma mensagem antes de considerar falha. Útil para garantir entrega garantida, mas pode causar espera indefinida se o broker estiver indisponível.
			"delivery.timeout.ms": "3000",
			/*
				"acks": "all" faz o producer aguardar confirmação de todos os brokers (líder + réplicas) antes de considerar a mensagem entregue, aumentando a durabilidade e segurança.
				"enable.idempotence": "true" garante que mensagens duplicadas não sejam geradas em caso de retransmissão, tornando o envio "exatamente uma vez".
				// kafka dentro dele deve criar um hash da mensagem + clientid + sequence, entao ele sabe que a mensagem esta duplicada
				Juntos, essas opções maximizam a confiabilidade e previnem duplicidades no Kafka Producer.
			*/
			"acks":                "all",
			"enable.idempotence":  "true",
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

func DeliveryReport(deliveryChan chan kafka.Event) {
	// range em Go é usado para iterar sobre elementos de uma coleção (como arrays, slices, maps, strings ou channels). No caso de for e := range deliveryChan, ele lê continuamente os valores enviados para o canal deliveryChan até que este seja fechado.
	for e := range deliveryChan {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				fmt.Println("Erro ao enviar")
			} else {
				fmt.Println("Mensagem enviada:", ev.TopicPartition)
				// anotar no banco de dados que a mensagem foi processado.
				// ex: confirma que uma transferencia bancaria ocorreu.
			}
		}
	}
}