package scrum21stcenturypoker

import (
	"os"
	"fmt"
)

type Error interface {
	os.Error
	HasUserMessage() bool
	UserMessage() string
}

type ErrorData struct {
	InternalMessage string
	ExternalMessage string
}

func (this *ErrorData) String() string {
	return this.InternalMessage
}

func (this *ErrorData) HasUserMessage() bool {
	return this.ExternalMessage != ""
}

func (this *ErrorData) UserMessage() string {
	return this.ExternalMessage
}

func (this *ErrorData) Format(v ...interface{}) (e *ErrorData) {
	e = this.FormatInternalMessage(v...)
	e = e.FormatExternalMessage(v...)
	return
}

func (this *ErrorData) FormatInternalMessage(v ...interface{}) (e *ErrorData) {
	e = new(ErrorData)
	*e = *this
	e.InternalMessage = fmt.Sprintf(this.InternalMessage, v...)
	return
}

func (this *ErrorData) FormatExternalMessage(v ...interface{}) (e *ErrorData) {
	e = new(ErrorData)
	*e = *this
	e.ExternalMessage = fmt.Sprintf(this.ExternalMessage, v...)
	return e
}

func FromError(prefix string, e os.Error) Error {
	if e == nil {
		return nil
	}
	return &ErrorData{
		InternalMessage: prefix + ": " + e.String(),
	}
}
