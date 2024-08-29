package http

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/Jacobbrewer1/f1-data/pkg/codegen/apis/common"
	"github.com/Jacobbrewer1/f1-data/pkg/utils"
	"github.com/stretchr/testify/suite"
)

type MessageSuite struct {
	suite.Suite

	w *httptest.ResponseRecorder
}

func TestMessageSuite(t *testing.T) {
	suite.Run(t, new(MessageSuite))
}

func (s *MessageSuite) SetupTest() {
	s.w = httptest.NewRecorder()
}

func (s *MessageSuite) TestSendMessage() {
	SendMessage(s.w, "test")

	s.Equal(200, s.w.Code)

	expectedResponse := "{\"message\":\"test\"}\n"
	s.Equal(expectedResponse, s.w.Body.String())
}

func (s *MessageSuite) TestSendMessageWithStatus() {
	SendMessageWithStatus(s.w, 400, "test")

	s.Equal(400, s.w.Code)

	expectedResponse := "{\"message\":\"test\"}\n"
	s.Equal(expectedResponse, s.w.Body.String())
}

type ErrorMessageSuite struct {
	suite.Suite

	w *httptest.ResponseRecorder
}

func TestErrorMessageSuite(t *testing.T) {
	suite.Run(t, new(ErrorMessageSuite))
}

func (s *ErrorMessageSuite) SetupTest() {
	s.w = httptest.NewRecorder()
}

func (s *ErrorMessageSuite) TestSendErrorMessage() {
	SendErrorMessage(s.w, "test message", errors.New("test error"))

	s.Equal(500, s.w.Code)

	expectedJson := "{\"error\":\"test error\",\"message\":\"test message\"}\n"
	s.JSONEq(expectedJson, s.w.Body.String())
}

func (s *ErrorMessageSuite) TestSendErrorMessageWithStatus() {
	SendErrorMessageWithStatus(s.w, 400, "test message", errors.New("test error"))

	s.Equal(400, s.w.Code)

	expectedJson := "{\"error\":\"test error\",\"message\":\"test message\"}\n"
	s.JSONEq(expectedJson, s.w.Body.String())
}

func (s *ErrorMessageSuite) TestNewErrorMessage() {
	msg := NewErrorMessage("test message", errors.New("test error"))

	expected := &common.ErrorMessage{
		Message: utils.Ptr("test message"),
		Error:   utils.Ptr("test error"),
	}
	s.Equal(expected, msg)
}

func (s *ErrorMessageSuite) TestNewErrorMessageNoError() {
	msg := NewErrorMessage("test message", nil)

	expected := &common.ErrorMessage{
		Message: utils.Ptr("test message"),
		Error:   utils.Ptr(""),
	}
	s.Equal(expected, msg)
}

func (s *ErrorMessageSuite) TestSendErrorMessageWithStatusNoError() {
	SendErrorMessageWithStatus(s.w, 400, "test message", nil)

	s.Equal(400, s.w.Code)

	expectedJson := "{\"error\":\"\",\"message\":\"test message\"}\n"
	s.JSONEq(expectedJson, s.w.Body.String())
}
