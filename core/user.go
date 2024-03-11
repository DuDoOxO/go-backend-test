package core

import (
	"fmt"

	"github.com/go-backend-test/model"
	"github.com/go-backend-test/repo"
)

type UserCore struct {
	userRepo repo.IUserRepo
}

func NewUserCore() IUserCore {
	uc := &UserCore{}
	uc.userRepo = repo.NewUserRepo()
	return uc
}

func (c *UserCore) AddUser(req model.User) error {
	has, err := c.userRepo.CheckUserExist(req)
	if err != nil {
		return err
	}

	if has {
		er := fmt.Errorf("The user already exists")
		return er
	}

	if err := c.userRepo.AddUser(req); err != nil {
		return err
	}

	return nil
}

func (c *UserCore) AddUserMessage(req model.LineMessage) error {
	return c.userRepo.AddUserMessage(req)
}

func (c *UserCore) ListUserMessageByUserId(userLineId string) ([]model.LineMessage, error) {
	return c.userRepo.ListUserMessageByUserId(userLineId)
}
