package createtransaction

import (
	"context"
	"testing"

	"github.com/flpcastro/fullcycle-repository/microservice-walletcore/internal/entity"
	"github.com/flpcastro/fullcycle-repository/microservice-walletcore/internal/event"
	"github.com/flpcastro/fullcycle-repository/microservice-walletcore/internal/usecase/mocks"
	"github.com/flpcastro/fullcycle-repository/microservice-walletcore/pkg/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateTransactionUseCase_Execute(t *testing.T) {
	client1, _ := entity.NewClient("John Doe", "j@j")
	account1 := entity.NewAccount(client1)
	account1.Credit(1000)

	client2, _ := entity.NewClient("Jane Doe", "j@j2")
	account2 := entity.NewAccount(client2)
	account2.Credit(1000)

	mockUow := &mocks.UowMock{}
	mockUow.On("Do", mock.Anything, mock.Anything).Return(nil)

	dispatcher := events.NewEventDispatcher()
	eventTransaction := event.NewTransactionCreated()
	eventBalance := event.NewBalanceUpdated()
	ctx := context.Background()

	inputDTO := CreateTransactionInputDTO{
		AccountIDFrom: account1.ID,
		AccountIDTo:   account2.ID,
		Amount:        100,
	}

	uc := NewCreateTransactionUseCase(mockUow, dispatcher, eventTransaction, eventBalance)
	output, err := uc.Execute(ctx, inputDTO)
	assert.Nil(t, err)
	assert.NotNil(t, output.ID)
	mockUow.AssertExpectations(t)
	mockUow.AssertNumberOfCalls(t, "Do", 1)
}
