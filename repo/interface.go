package repo

import "github.com/go-backend-test/model"

type IUserRepo interface {
	AddUser(req model.User) error
	CheckUserExist(req model.User) (bool, error)
	AddUserMessage(model.LineMessage) error
	ListUserMessageByUserId(userLineId string) ([]model.LineMessage, error)
}
