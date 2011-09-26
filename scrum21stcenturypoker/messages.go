package scrum21stcenturypoker

import ()

const (
	MSG_SUCCESS = iota
	MSG_ERROR
)

type Message struct {
	Code        int
	Description string
	Data        interface{}
}

func NewSuccessMessage(descr string, data interface{}) *Message {
	return &Message{
		Code:        MSG_SUCCESS,
		Description: descr,
		Data:        data,
	}
}

func NewErrorMessage(descr string) *Message {
	return &Message{
		Code:        MSG_ERROR,
		Description: descr,
		Data:        nil,
	}
}
