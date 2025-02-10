package uhttp

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/jacobbrewer1/uhttp/common"
)

// NewMessage creates a new Message.
func NewMessage(message string, args ...any) *common.Message {
	var msg string
	if len(args) > 0 {
		msg = fmt.Sprintf(message, args...)
	} else {
		msg = message
	}
	return &common.Message{
		Message: msg,
	}
}

func SendMessageWithStatus(w http.ResponseWriter, status int, message string, args ...any) {
	msg := NewMessage(message, args...)
	err := EncodeJSON(w, status, msg)
	if err != nil {
		slog.Error("Error encoding message", slog.String(loggingKeyError, err.Error()))
	}
}

func SendMessage(w http.ResponseWriter, message string, args ...any) {
	SendMessageWithStatus(w, http.StatusOK, message, args...)
}
