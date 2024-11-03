package handler

import (
	"fmt"
	"sync"

	"github.com/flpcastro/fullcycle-repository/microservice-walletcore/pkg/events"
	"github.com/flpcastro/fullcycle-repository/microservice-walletcore/pkg/kafka"
)

type TransactionCreatedKafkaHandler struct {
	Kafka *kafka.Producer
}

func NewTransactionCreatedKafkaHandler(kafka *kafka.Producer) *TransactionCreatedKafkaHandler {
	return &TransactionCreatedKafkaHandler{
		Kafka: kafka,
	}
}

func (h *TransactionCreatedKafkaHandler) Handle(event events.EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()
	h.Kafka.Publish(event, nil, "transactions")
	fmt.Println("TransactionCreatedKafkaHandler: ", event)
}
