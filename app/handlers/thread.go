package handlers

import (
	"net/http"
	"strconv"

	"github.com/jackc/pgx"
	"github.com/labstack/echo"
	"vk_db_project/app/models"
	"vk_db_project/app/uscases"
)

type ThreadHandler struct {
	threadLogic uscases.IThreadUsecase
}

func NewThreadHandler(tLogic uscases.ThreadsUsecaseImpl) ThreadHandler {
	return ThreadHandler{threadLogic: tLogic}
}

func (Thread ThreadHandler) CreatePosts(rwContext echo.Context) error {

	slugOrId := rwContext.Param("slug_or_id")

	posts := []models.Post{}
	rwContext.Bind(&posts)

	posts, err := Thread.threadLogic.CreatePosts(slugOrId, posts)
	if err == pgx.ErrNoRows {
		return rwContext.JSON(http.StatusNotFound, models.Error{Message: "can't find thread by slug_or_id: " + slugOrId})
	}

	if err != nil {
		if err.Error() == "no user" {
			return rwContext.JSON(http.StatusNotFound, models.Error{Message: "Can't find post author by nickname: "})
		}

		if err.Error() == "No parent message!" {
			return rwContext.JSON(http.StatusConflict, models.Error{Message: err.Error()})
		}

		if err.Error() == "Parent post was created in another thread" {
			return rwContext.JSON(http.StatusConflict, models.Error{Message: err.Error()})
		}
	}

	if err != nil {
		return rwContext.JSON(http.StatusConflict, models.Error{Message: "no parent message"})
	}

	return rwContext.JSON(http.StatusCreated, posts)
}

func (Thread ThreadHandler) VoteThread(rwContext echo.Context) error {
	slugOrId := rwContext.Param("slug_or_id")

	vote := new(models.Vote)
	rwContext.Bind(&vote)

	thread, err := Thread.threadLogic.VoteThread(slugOrId, vote.Nickname, vote.Voice)

	if err != nil {
		return rwContext.JSON(http.StatusNotFound, models.Error{Message: "can't vote by slug_or_id:" + slugOrId})
	}

	return rwContext.JSON(http.StatusOK, thread)
}

func (Thread ThreadHandler) GetThread(rwContext echo.Context) error {
	slugOrId := rwContext.Param("slug_or_id")

	thread, err := Thread.threadLogic.GetThread(slugOrId)

	if err != nil {
		return rwContext.JSON(http.StatusNotFound, models.Error{Message: "can't get thread by slug_or_id:" + slugOrId})
	}

	return rwContext.JSON(http.StatusOK, thread)
}

func (Thread ThreadHandler) GetPosts(rwContext echo.Context) error {
	slugOrId := rwContext.Param("slug_or_id")
	limit, _ := strconv.Atoi(rwContext.QueryParam("limit"))
	since, _ := strconv.Atoi(rwContext.QueryParam("since"))
	sortType := rwContext.QueryParam("sort")
	desc, err := strconv.ParseBool(rwContext.QueryParam("desc"))

	if err != nil {
		desc = false
	}

	if sortType == "" {
		sortType = "flat"
	}

	posts, err := Thread.threadLogic.GetPosts(slugOrId, limit, since, sortType, desc)
	if err != nil {
		return rwContext.JSON(http.StatusNotFound, models.Error{Message: "can't get thread by slug_or_id:" + slugOrId})
	}

	return rwContext.JSON(http.StatusOK, posts)
}

func (Thread ThreadHandler) UpdateThread(rwContext echo.Context) error {
	slugOrId := rwContext.Param("slug_or_id")
	newThread := new(models.Thread)
	rwContext.Bind(newThread)

	thread, err := Thread.threadLogic.UpdateThread(slugOrId, *newThread)

	if err != nil {
		return rwContext.JSON(http.StatusNotFound, models.Error{"Can't find thread by slug_or_id: " + slugOrId})
	}

	return rwContext.JSON(http.StatusOK, thread)
}

func (Thread ThreadHandler) SetupHandlers(server *echo.Echo) {
	server.POST("/api/thread/:slug_or_id/create", Thread.CreatePosts)
	server.POST("/api/thread/:slug_or_id/vote", Thread.VoteThread)
	server.POST("/api/thread/:slug_or_id/details", Thread.UpdateThread)
	server.GET("/api/thread/:slug_or_id/details", Thread.GetThread)
	server.GET("/api/thread/:slug_or_id/posts", Thread.GetPosts)
}