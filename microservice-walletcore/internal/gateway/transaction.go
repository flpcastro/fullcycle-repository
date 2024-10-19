package gateway

import "github.com/flpcastro/fullcycle-repository/microservice-walletcore/internal/entity"

type TransactionGateway interface {
	Create(transaction *entity.Transaction) error
}
