package uscases

import (
	"vk_db_project/app/models"
	"vk_db_project/app/repositories"
)

type IForumUsecase interface {
	CreateForum(models.Forum) (models.Forum, error)
	GetForumData(string) (models.Forum, error)
	CreateThread(string, models.Thread) (models.Thread, error)
	GetThreads(string, int, string, bool) ([]models.Thread, error)
	GetForumUsers(string, int, string, bool) ([]models.UserModel, error)
}

type ForumUsecaseImpl struct {
	ForumRepo repositories.IForumRepository
}

func NewForumUsecaseImpl(fRepo repositories.ForumRepoImpl) ForumUsecaseImpl {
	return ForumUsecaseImpl{ForumRepo: fRepo}
}

func (ForumUC ForumUsecaseImpl) CreateForum(forum models.Forum) (models.Forum, error) {
	return ForumUC.ForumRepo.CreateNewForum(forum)
}

func (ForumUC ForumUsecaseImpl) GetForumData(slug string) (models.Forum, error) {
	return ForumUC.ForumRepo.GetForum(slug)
}

func (ForumUC ForumUsecaseImpl) CreateThread(slug string, thread models.Thread) (models.Thread, error) {
	thread.Forum = slug
	return ForumUC.ForumRepo.CreateThread(thread)
}

func (ForumUC ForumUsecaseImpl) GetThreads(slug string, limit int, since string, sort bool) ([]models.Thread, error) {
	forum := models.Forum{Slug: slug}

	return ForumUC.ForumRepo.GetThreads(forum, limit, since, sort)
}

func (ForumUC ForumUsecaseImpl) GetForumUsers(slug string, limit int, since string, desc bool) ([]models.UserModel, error) {
	return ForumUC.ForumRepo.GetForumUsers(slug, limit, since, desc)
}