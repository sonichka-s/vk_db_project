package handlers

import (
	"net/http"
	"strconv"

	"github.com/jackc/pgx"
	"github.com/labstack/echo"
	"vk_db_project/app/models"
	"vk_db_project/app/uscases"
)

type ForumHandler struct {
	ForumLogic uscases.IForumUsecase
}

func NewForumHandler(fLogic uscases.ForumUsecaseImpl) ForumHandler {
	return ForumHandler{ForumLogic: fLogic}
}

func (ForumHandler ForumHandler) CreateForum(rwContext echo.Context) error {
	newForumData := new(models.Forum)

	rwContext.Bind(newForumData)
	answer, err := ForumHandler.ForumLogic.CreateForum(*newForumData)
	if err == pgx.ErrNoRows {
		return rwContext.JSON(http.StatusNotFound, &models.Error{Message: "Can't find user by nickname: " + newForumData.User})
	}

	if err != nil {
		return rwContext.JSON(http.StatusConflict, answer)
	}

	return rwContext.JSON(http.StatusCreated, answer)
}

func (ForumHandler ForumHandler) GetForumUsers(rwContext echo.Context) error {
	slug := rwContext.Param("slug")
	limit, _ := strconv.Atoi(rwContext.QueryParam("limit"))
	since := rwContext.QueryParam("since")
	desc, _ := strconv.ParseBool(rwContext.QueryParam("desc"))

	data, err := ForumHandler.ForumLogic.GetForumUsers(slug, limit, since, desc)
	if err != nil {
		return rwContext.JSON(http.StatusNotFound, &models.Error{Message: "can't find forum by slug: " + slug})
	}

	return rwContext.JSON(http.StatusOK, data)
}

func (ForumHandler ForumHandler) GetForum(rwContext echo.Context) error {
	slug := rwContext.Param("slug")

	answer, err := ForumHandler.ForumLogic.GetForumData(slug)

	if err == pgx.ErrNoRows {
		return rwContext.JSON(http.StatusNotFound, &models.Error{Message: "Can't find forum by slug: " + slug})
	}

	return rwContext.JSON(http.StatusOK, answer)
}

func (ForumHandler ForumHandler) CreateThread(rwContext echo.Context) error {
	slug := rwContext.Param("slug")

	threadReq := new(models.Thread)
	rwContext.Bind(threadReq)

	thread, err := ForumHandler.ForumLogic.CreateThread(slug, *threadReq)
	if err == pgx.ErrNoRows {
		return rwContext.JSON(http.StatusNotFound, &models.Error{Message: "Can't create thread by slug: " + slug})
	}

	if err != nil {
		return rwContext.JSON(http.StatusConflict, thread)
	}

	return rwContext.JSON(http.StatusCreated, thread)
}

func (ForumHandler ForumHandler) GetSortedThreads(rwContext echo.Context) error {
	slug := rwContext.Param("slug")
	limit, _ := strconv.Atoi(rwContext.QueryParam("limit"))
	since := rwContext.QueryParam("since")
	desc, _ := strconv.ParseBool(rwContext.QueryParam("desc"))

	threads, err := ForumHandler.ForumLogic.GetThreads(slug, limit, since, desc)

	if err != nil {
		return rwContext.JSON(http.StatusNotFound, &models.Error{Message: "Can't create thread by slug: " + slug})
	}

	return rwContext.JSON(http.StatusOK, threads)
}

func (ForumHandler ForumHandler) SetupHandlers(server *echo.Echo) {
	server.POST("/api/forum/create", ForumHandler.CreateForum)
	server.GET("/api/forum/:slug/details", ForumHandler.GetForum)
	server.POST("/api/forum/:slug/create", ForumHandler.CreateThread)
	server.GET("/api/forum/:slug/threads", ForumHandler.GetSortedThreads)
	server.GET("/api/forum/:slug/users", ForumHandler.GetForumUsers)
}
