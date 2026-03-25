# Apache kafka

#### objetivo
- colocar o conhecimento em pratica

#### como executar o ambiente
- instalar o docker
- executar "docker compose up" no terminal

#### comandos e workflow

```bash
# ver conteudo da instancia
docker logs kafka

# rodar ao atualizar dockerfile
docker compose up --build

# entrar no container como root
docker exec -u 0 -it fc3-apache-kafka-kafka-1 bash
  # listar comandos
  ls /usr/bin | grep '^kafka-'
  # executar um utilitario
  kafka-topics

#DENTRO DO KAFKA - WORKFLOW DE TESTE

#qualquer coisa que eu for fazer no kafka vou precisar passar o bootstrap-server

- criar topic
kafka-topics --create --topic teste --bootstrap-server localhost:9092 --partitions 3

- comando: listar topics
- comando: ver descrição do topic

- umas das flags ao conectar consumidor
# cli para isso: kafka-console-consumer
# opcinal: escolher partição especifica
# opcinal: escolher consumer-group
# opcinal: --from-beginning ; por padrão começa o consumidor le das mensagens mais novas, mas com essa opção consigo começar no offset 0

cenario 1: 2 consumer lendo o mesmo topico em grupos separados via cli
- consumer 1: ler topicos em tempo real via cli
kafka-console-consumer --topic test --bootstrap-server localhost:9092
- consumer 2: ler topicos em tempo real via cli
kafka-console-consumer --topic test --bootstrap-server localhost:9092
- producer 1: produzindo mensagem
kafka-console-producer --topic test --bootstrap-server localhost:9092

cenario 2: 2 consumer lendo o mesmo topico no mesmo grupo só que em partições separadas
- consumer 1: ler topicos em tempo real via cli
kafka-console-consumer --bootstrap-server localhost:9092 --topic test --group=x
- consumer 1: ler topicos em tempo real via cli
kafka-console-consumer --bootstrap-server localhost:9092 --topic test --group=x
- producer 1: produzindo mensagem
kafka-console-producer --topic test --bootstrap-server localhost:9092


- ver os consumidores de um grupo
kafka-consumer-groups --bootstrap-server localhost:9092 --group x --describe

- acessar control-center no navegador
localhost:9021

```

## Kafka basics and commands
Prerequisite: make sure the Kafka broker is running and reachable at localhost:9092.

## How to read the diagram (kafka.png)
![Kafka overview](./docs/images/kafka.png)

## Delivery
![Kafka delivery overview](./docs/images/delivery.png)

# Commands

## Topics
- Create a topic (3 partitions):
```shell
kafka-topics --create --topic test --bootstrap-server localhost:9092 --partitions 3
```
- List all topics:
```shell
kafka-topics --list --bootstrap-server localhost:9092
```
- Describe a topic (partitions/replicas/leader):
```shell
kafka-topics --describe --topic test --bootstrap-server localhost:9092
```
- Delete a topic:
```shell
kafka-topics --delete --topic test --bootstrap-server localhost:9092
```

## Console Producer
- Write messages from stdin to a topic. Each line you type is a message:
```shell
kafka-console-producer --topic test --bootstrap-server localhost:9092
```

## Console Consumer
- Read new messages from the tip of the log:
```shell
kafka-console-consumer --topic test --bootstrap-server localhost:9092
```
- Read everything from the beginning of the topic:
```shell
kafka-console-consumer --topic test --bootstrap-server localhost:9092 --from-beginning
```
- Join a consumer group (enables scaling and offset tracking):
```shell
kafka-console-consumer --topic test --bootstrap-server localhost:9092 --group group-x
```

## Consumer Groups
- Describe a group (lags, offsets, assignments). Start at least one consumer in the group first:
```shell
kafka-consumer-groups --bootstrap-server localhost:9092 --group group-x --describe
```