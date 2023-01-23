package uscases

import (
	"vk_db_project/app/models"
	"vk_db_project/app/repositories"
)

type IPostUsecase interface {
	GetPostData(int, []string) (models.FullPost, error)
	UpdatePost(int64, string) (models.Post, error)
}

type PostUsecaseImpl struct {
	postRepo repositories.PostRepoImpl
}

func NewPostUsecaseImpl(pRepo repositories.PostRepoImpl) PostUsecaseImpl {
	return PostUsecaseImpl{postRepo: pRepo}
}

func (PostUC PostUsecaseImpl) GetPostData(id int, flags []string) (models.FullPost, error) {
	return PostUC.postRepo.GetPost(id, flags)
}

func (PostUC PostUsecaseImpl) UpdatePost(id int64, message string) (models.Post, error) {
	return PostUC.postRepo.UpdatePost(models.Post{Id: id, Message: message})
}
