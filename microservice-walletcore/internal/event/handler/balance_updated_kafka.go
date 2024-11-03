package handler

import (
	"fmt"
	"sync"

	"github.com/flpcastro/fullcycle-repository/microservice-walletcore/pkg/events"
	"github.com/flpcastro/fullcycle-repository/microservice-walletcore/pkg/kafka"
)

type UpdatedBalanceKafkaHandler struct {
	Kafka *kafka.Producer
}

func NewUpdatedBalanceKafkaHandler(kafka *kafka.Producer) *UpdatedBalanceKafkaHandler {
	return &UpdatedBalanceKafkaHandler{
		Kafka: kafka,
	}
}

func (h *UpdatedBalanceKafkaHandler) Handle(message events.EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()
	h.Kafka.Publish(message, nil, "balances")
	fmt.Println("Balance updated event published to Kafka")
}
