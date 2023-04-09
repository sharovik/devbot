package mock

import (
	"errors"
	"fmt"
)

type Service struct {
}

var (
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

func (s Service) BeforeWSConnectionStart() error {
	return nil
}

func (s Service) ProcessMessage(message interface{}) error {
	return nil
}
