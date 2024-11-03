package main

import (
	"context"
	"database/sql"
	"fmt"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/flpcastro/fullcycle-repository/microservice-walletcore/internal/database"
	"github.com/flpcastro/fullcycle-repository/microservice-walletcore/internal/event"
	"github.com/flpcastro/fullcycle-repository/microservice-walletcore/internal/event/handler"
	createaccount "github.com/flpcastro/fullcycle-repository/microservice-walletcore/internal/usecase/create_account"
	createclient "github.com/flpcastro/fullcycle-repository/microservice-walletcore/internal/usecase/create_client"
	createtransaction "github.com/flpcastro/fullcycle-repository/microservice-walletcore/internal/usecase/create_transaction"
	"github.com/flpcastro/fullcycle-repository/microservice-walletcore/internal/web"
	"github.com/flpcastro/fullcycle-repository/microservice-walletcore/internal/web/webserver"
	"github.com/flpcastro/fullcycle-repository/microservice-walletcore/pkg/events"
	"github.com/flpcastro/fullcycle-repository/microservice-walletcore/pkg/kafka"
	"github.com/flpcastro/fullcycle-repository/microservice-walletcore/pkg/uow"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", "root", "root", "mysql", "3306", "wallet"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	configMap := ckafka.ConfigMap{
		"bootstrap.servers": "kafka:29092",
		"group.id":          "wallet",
	}

	kafkaProducer := kafka.NewKafkaProducer(&configMap)

	eventDispatcher := events.NewEventDispatcher()
	transactionCreatedEvent := event.NewTransactionCreated()
	eventDispatcher.Register("TransactionCreated", handler.NewTransactionCreatedKafkaHandler(kafkaProducer))
	balanceUpdatedEvent := event.NewBalanceUpdated()
	eventDispatcher.Register("BalanceUpdated", handler.NewUpdatedBalanceKafkaHandler(kafkaProducer))

	clientDb := database.NewClientDB(db)
	accountDb := database.NewAccountDB(db)

	ctx := context.Background()
	uow := uow.NewUow(ctx, db)

	uow.Register("AccountDB", func(tx *sql.Tx) interface{} {
		return database.NewAccountDB(db)
	})

	uow.Register("TransactionDB", func(tx *sql.Tx) interface{} {
		return database.NewTransactionDB(db)
	})

	createClientUseCase := createclient.NewCreateClientUseCase(clientDb)
	createAccountUseCase := createaccount.NewCreateAccountUseCase(accountDb, clientDb)
	createTransactionUseCase := createtransaction.NewCreateTransactionUseCase(uow, eventDispatcher, transactionCreatedEvent, balanceUpdatedEvent)

	webserver := webserver.NewWebServer(":8080")

	clientHandler := web.NewWebClientHandler(*createClientUseCase)
	accountHandler := web.NewWebAccountHandler(*createAccountUseCase)
	transactionHandler := web.NewWebTransactionHandler(*createTransactionUseCase)

	webserver.AddHandler("/clients", clientHandler.CreateClient)
	webserver.AddHandler("/accounts", accountHandler.CreateAccount)
	webserver.AddHandler("/transactions", transactionHandler.CreateTransaction)

	fmt.Println("Server is running")
	webserver.Start()
}
