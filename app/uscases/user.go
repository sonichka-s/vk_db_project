package uscases

import (
	"vk_db_project/app/models"
	"vk_db_project/app/repositories"
)

type IUserUsecase interface {
	GetUser(string) (models.UserModel, error)
	CreateUser(models.UserModel) (interface{}, error)
	UpdateUserData(models.UserModel) (models.UserModel, error)
	GetServerStatus() models.Status
	Clear()
}

type UserUsecaseImpl struct {
	userRepo repositories.IUserRepo
}

func NewUserUsecaseImpl(uRepo repositories.UserRepoImpl) UserUsecaseImpl {
	return UserUsecaseImpl{userRepo: uRepo}
}

func (UserUC UserUsecaseImpl) GetUser(nickname string) (models.UserModel, error) {
	return UserUC.userRepo.GetUserData(nickname)
}

func (UserUC UserUsecaseImpl) CreateUser(newUser models.UserModel) (interface{}, error) {
	answerData, err := UserUC.userRepo.CreateNewUser(newUser)
	if err != nil {
		return answerData, err
	}

	return answerData[0], err
}

func (UserUC UserUsecaseImpl) UpdateUserData(newUserData models.UserModel) (models.UserModel, error) {
	return UserUC.userRepo.UpdateUserData(newUserData)
}

func (UserUC UserUsecaseImpl) GetServerStatus() models.Status {
	return UserUC.userRepo.Status()
}

func (UserUC UserUsecaseImpl) Clear() {
	UserUC.userRepo.Clear()
}
