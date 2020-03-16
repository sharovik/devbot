package mock

import (
	"errors"
	"fmt"

	"github.com/sharovik/devbot/internal/service/base"
)

type Service struct {
}

var (
	MockedSlackService               base.ServiceInterface = Service{}
	ErrorInitWebSocketReceiver       map[int]error
	numberInitWebSocketReceiverCalls = 0
)

func (s Service) InitWebSocketReceiver() error {
	err := errors.New("Default error ")
	if ErrorInitWebSocketReceiver[numberInitWebSocketReceiverCalls] != nil {
		err = ErrorInitWebSocketReceiver[numberInitWebSocketReceiverCalls]
	}

	fmt.Println(fmt.Sprintf("Retry: %d", numberInitWebSocketReceiverCalls))
	numberInitWebSocketReceiverCalls++
	return err
}
