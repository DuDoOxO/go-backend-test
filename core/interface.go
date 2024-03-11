package core

import (
	"github.com/gin-gonic/gin"
	"github.com/go-backend-test/model"
)

type IUserCore interface {
	AddUser(req model.User) error
	AddUserMessage(req model.LineMessage) error
	ListUserMessageByUserId(userLineId string) ([]model.LineMessage, error)
}

type ILine interface {
	// Callback(w http.ResponseWriter, req *http.Request)
	HandleCallback() gin.HandlerFunc
	HandleSendMessageByAPI() gin.HandlerFunc
}
