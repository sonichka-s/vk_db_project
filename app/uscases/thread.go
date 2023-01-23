package uscases

import (
	"strconv"
	"time"

	"vk_db_project/app/models"
	"vk_db_project/app/repositories"
)

type IThreadUsecase interface {
	CreatePosts(string, []models.Post) ([]models.Post, error)
	VoteThread(string, string, int) (models.Thread, error)
	GetThread(string) (models.Thread, error)
	GetPosts(string, int, int, string, bool) ([]models.Post, error)
	UpdateThread(string, models.Thread) (models.Thread, error)
}

type ThreadsUsecaseImpl struct {
	threadRepo repositories.IThreadRepository
}

func NewThreadsUsecaseImpl(tRepo repositories.ThreadRepoImpl) ThreadsUsecaseImpl {
	return ThreadsUsecaseImpl{threadRepo: tRepo}
}

func (ThreadUC ThreadsUsecaseImpl) CreatePosts(slugOrId string, posts []models.Post) ([]models.Post, error) {

	id, err := strconv.Atoi(slugOrId)

	if err != nil {
		id = 0
	} else {
		slugOrId = ""
	}

	t := time.Now()

	return ThreadUC.threadRepo.CreatePost(t, slugOrId, id, posts)
}

func (ThreadUC ThreadsUsecaseImpl) VoteThread(slug, nickname string, voice int) (models.Thread, error) {

	threadId, err := strconv.Atoi(slug)

	if err != nil {
		threadId = 0
	} else {
		slug = ""
	}

	return ThreadUC.threadRepo.VoteThread(nickname, voice, threadId, models.Thread{Slug: slug})
}

func (ThreadUC ThreadsUsecaseImpl) GetThread(slug string) (models.Thread, error) {

	threadId, err := strconv.Atoi(slug)

	if err != nil {
		threadId = 0
	} else {
		slug = ""
	}

	return ThreadUC.threadRepo.GetThread(threadId, models.Thread{Slug: slug})
}

func (ThreadUC ThreadsUsecaseImpl) GetPosts(slugOrId string, limit int, since int, sortType string, desc bool) ([]models.Post, error) {

	threadId, err := strconv.Atoi(slugOrId)

	if err != nil {
		threadId = 0
	} else {
		slugOrId = ""
	}

	data, err := ThreadUC.threadRepo.GetPostsSorted(slugOrId, threadId, limit, since, sortType, desc)
	return data, err
}

func (ThreadUC ThreadsUsecaseImpl) UpdateThread(slugOrId string, newThreadData models.Thread) (models.Thread, error) {
	threadId, err := strconv.Atoi(slugOrId)

	if err != nil {
		threadId = 0
	} else {
		slugOrId = ""
	}

	return ThreadUC.threadRepo.UpdateThread(slugOrId, threadId, newThreadData)
}
